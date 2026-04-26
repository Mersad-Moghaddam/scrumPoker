package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"scrum-poker/backend/internal/models"
)

const roomEventsChannel = "room-events"

type RedisCache struct {
	client      *redis.Client
	presenceTTL time.Duration
}

func NewRedisCache(addr, password string, db int, presenceTTL time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RedisCache{client: client, presenceTTL: presenceTTL}, nil
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) TouchPresence(ctx context.Context, roomCode string, memberID int64) error {
	expireAt := float64(time.Now().Add(c.presenceTTL).Unix())
	key := presenceKey(roomCode)
	if err := c.client.ZAdd(ctx, key, redis.Z{Score: expireAt, Member: memberID}).Err(); err != nil {
		return err
	}
	return c.client.Expire(ctx, key, c.presenceTTL*2).Err()
}

func (c *RedisCache) RemovePresence(ctx context.Context, roomCode string, memberID int64) error {
	return c.client.ZRem(ctx, presenceKey(roomCode), memberID).Err()
}

func (c *RedisCache) ActiveMemberIDs(ctx context.Context, roomCode string) ([]int64, error) {
	key := presenceKey(roomCode)
	now := strconv.FormatInt(time.Now().Unix(), 10)
	if err := c.client.ZRemRangeByScore(ctx, key, "-inf", now).Err(); err != nil {
		return nil, err
	}
	entries, err := c.client.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(entries))
	for _, entry := range entries {
		id, err := strconv.ParseInt(entry, 10, 64)
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (c *RedisCache) SetVoteProgress(ctx context.Context, roomCode string, taskID int64, votedCount, totalCount int) error {
	key := voteProgressKey(roomCode, taskID)
	_, err := c.client.HSet(ctx, key,
		"voted_count", votedCount,
		"total_count", totalCount,
	).Result()
	if err != nil {
		return err
	}
	return c.client.Expire(ctx, key, 24*time.Hour).Err()
}

func (c *RedisCache) GetVoteProgress(ctx context.Context, roomCode string, taskID int64) (*models.VotingProgress, error) {
	values, err := c.client.HGetAll(ctx, voteProgressKey(roomCode, taskID)).Result()
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}
	votedCount, _ := strconv.Atoi(values["voted_count"])
	totalCount, _ := strconv.Atoi(values["total_count"])
	return &models.VotingProgress{VotedCount: votedCount, TotalCount: totalCount}, nil
}

func (c *RedisCache) SetCurrentTask(ctx context.Context, roomCode string, taskID int64, status string) error {
	key := roomStateKey(roomCode)
	_, err := c.client.HSet(ctx, key,
		"task_id", taskID,
		"status", status,
	).Result()
	if err != nil {
		return err
	}
	return c.client.Expire(ctx, key, 24*time.Hour).Err()
}

func (c *RedisCache) ClearCurrentTask(ctx context.Context, roomCode string) error {
	return c.client.Del(ctx, roomStateKey(roomCode)).Err()
}

func (c *RedisCache) PublishRoomEvent(ctx context.Context, event models.RoomEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return c.client.Publish(ctx, roomEventsChannel, payload).Err()
}

func (c *RedisCache) SubscribeRoomEvents(ctx context.Context) *redis.PubSub {
	return c.client.Subscribe(ctx, roomEventsChannel)
}

func presenceKey(roomCode string) string {
	return "room:" + roomCode + ":presence"
}

func voteProgressKey(roomCode string, taskID int64) string {
	return "room:" + roomCode + ":task:" + strconv.FormatInt(taskID, 10) + ":progress"
}

func roomStateKey(roomCode string) string {
	return "room:" + roomCode + ":state"
}
