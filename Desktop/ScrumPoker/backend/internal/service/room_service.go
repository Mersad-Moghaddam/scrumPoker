package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"
	"net/http"
	"strings"
	"time"

	"scrum-poker/backend/internal/cache"
	"scrum-poker/backend/internal/models"
	"scrum-poker/backend/internal/store"
)

type RoomService struct {
	repo  *store.Repository
	cache *cache.RedisCache
}

func NewRoomService(repo *store.Repository, cache *cache.RedisCache) *RoomService {
	return &RoomService{repo: repo, cache: cache}
}

func (s *RoomService) CreateRoom(ctx context.Context, displayName string) (*models.AuthSession, error) {
	displayName = strings.TrimSpace(displayName)
	if len([]rune(displayName)) < 2 || len([]rune(displayName)) > 30 {
		return nil, newAppError(http.StatusBadRequest, "نام نمایشی باید بین ۲ تا ۳۰ کاراکتر باشد.")
	}

	roomCode, err := s.generateUniqueRoomCode(ctx)
	if err != nil {
		return nil, err
	}
	token := randomToken(32)
	room, member, err := s.repo.CreateRoom(ctx, roomCode, displayName, token)
	if err != nil {
		return nil, err
	}
	_ = s.repo.TouchMember(ctx, member.ID)
	_ = s.cache.TouchPresence(ctx, room.Code, member.ID)
	_ = s.publish(ctx, room.Code, "room.created", map[string]any{"memberId": member.ID})
	return &models.AuthSession{RoomCode: room.Code, Token: token, MemberID: member.ID}, nil
}

func (s *RoomService) JoinRoom(ctx context.Context, roomCode, displayName string) (*models.AuthSession, error) {
	roomCode = normalizeRoomCode(roomCode)
	displayName = strings.TrimSpace(displayName)
	if !isValidRoomCode(roomCode) {
		return nil, newAppError(http.StatusBadRequest, "کد اتاق نامعتبر است.")
	}
	if len([]rune(displayName)) < 2 || len([]rune(displayName)) > 30 {
		return nil, newAppError(http.StatusBadRequest, "نام نمایشی باید بین ۲ تا ۳۰ کاراکتر باشد.")
	}

	room, err := s.repo.GetRoomByCode(ctx, roomCode)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, newAppError(http.StatusNotFound, "اتاق پیدا نشد.")
	}
	if err != nil {
		return nil, err
	}

	token := randomToken(32)
	member, err := s.repo.CreateMember(ctx, room.ID, displayName, token)
	if err != nil {
		return nil, err
	}

	_ = s.cache.TouchPresence(ctx, room.Code, member.ID)
	_ = s.publish(ctx, room.Code, "member.joined", map[string]any{"memberId": member.ID})
	return &models.AuthSession{RoomCode: room.Code, Token: token, MemberID: member.ID}, nil
}

func (s *RoomService) GetRoomSnapshot(ctx context.Context, roomCode, token string) (*models.RoomSnapshot, error) {
	room, member, err := s.authorizeMember(ctx, roomCode, token)
	if err != nil {
		return nil, err
	}
	_ = s.repo.TouchMember(ctx, member.ID)
	_ = s.cache.TouchPresence(ctx, room.Code, member.ID)
	return s.buildSnapshot(ctx, room, member)
}

func (s *RoomService) GetRoomHistory(ctx context.Context, roomCode, token string) ([]models.HistoryTask, error) {
	room, member, err := s.authorizeMember(ctx, roomCode, token)
	if err != nil {
		return nil, err
	}
	return s.buildHistory(ctx, room.ID, member.ID)
}

func (s *RoomService) CreateTask(ctx context.Context, roomCode, token, title string) (*models.RoomSnapshot, error) {
	room, member, err := s.authorizeMember(ctx, roomCode, token)
	if err != nil {
		return nil, err
	}
	if !member.IsHost {
		return nil, newAppError(http.StatusForbidden, "فقط میزبان می‌تواند تسک جدید بسازد.")
	}
	title = strings.TrimSpace(title)
	if len([]rune(title)) < 1 || len([]rune(title)) > 120 {
		return nil, newAppError(http.StatusBadRequest, "عنوان تسک باید بین ۱ تا ۱۲۰ کاراکتر باشد.")
	}

	currentTask, err := s.repo.GetCurrentTask(ctx, room.ID)
	if err != nil {
		return nil, err
	}
	if currentTask != nil && currentTask.Status == models.TaskStatusVoting {
		return nil, newAppError(http.StatusConflict, "در هر اتاق فقط یک تسک فعال می‌تواند وجود داشته باشد.")
	}

	if err := s.repo.ArchiveRevealedTasks(ctx, room.ID); err != nil {
		return nil, err
	}
	task, err := s.repo.CreateTask(ctx, room.ID, title)
	if err != nil {
		return nil, err
	}

	activeMembers, _ := s.activeMemberIDs(ctx, room)
	_ = s.cache.SetCurrentTask(ctx, room.Code, task.ID, task.Status)
	_ = s.cache.SetVoteProgress(ctx, room.Code, task.ID, 0, len(activeMembers))
	_ = s.publish(ctx, room.Code, "task.created", map[string]any{"taskId": task.ID})
	return s.buildSnapshot(ctx, room, member)
}

func (s *RoomService) SubmitVote(ctx context.Context, roomCode string, taskID int64, token, value string) (*models.RoomSnapshot, error) {
	room, member, err := s.authorizeMember(ctx, roomCode, token)
	if err != nil {
		return nil, err
	}
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, newAppError(http.StatusNotFound, "تسک فعال پیدا نشد.")
	}
	if err != nil {
		return nil, err
	}
	if task.RoomID != room.ID {
		return nil, newAppError(http.StatusBadRequest, "این تسک به این اتاق تعلق ندارد.")
	}
	if task.Status != models.TaskStatusVoting {
		return nil, newAppError(http.StatusConflict, "رأی‌گیری این تسک تمام شده است.")
	}
	if !IsValidVoteValue(value) {
		return nil, newAppError(http.StatusBadRequest, "مقدار رأی نامعتبر است.")
	}
	existingVote, err := s.repo.GetTaskVoteByMember(ctx, taskID, member.ID)
	if err != nil {
		return nil, err
	}
	if existingVote != nil {
		return nil, newAppError(http.StatusConflict, "شما قبلا رأی خود را ثبت کرده‌اید.")
	}

	numericValue, _ := NumericVoteValue(value)
	if _, err := s.repo.CreateVote(ctx, taskID, member.ID, value, numericValue); err != nil {
		return nil, err
	}

	return s.syncRevealState(ctx, room, member)
}

func (s *RoomService) HandlePresenceConnected(ctx context.Context, roomCode, token string) error {
	room, member, err := s.authorizeMember(ctx, roomCode, token)
	if err != nil {
		return err
	}
	_ = s.repo.TouchMember(ctx, member.ID)
	_ = s.cache.TouchPresence(ctx, room.Code, member.ID)
	return s.syncPresence(ctx, room, member.ID, "member.joined")
}

func (s *RoomService) HandlePresenceHeartbeat(ctx context.Context, roomCode, token string) error {
	room, member, err := s.authorizeMember(ctx, roomCode, token)
	if err != nil {
		return err
	}
	_ = s.repo.TouchMember(ctx, member.ID)
	return s.cache.TouchPresence(ctx, room.Code, member.ID)
}

func (s *RoomService) HandlePresenceDisconnected(ctx context.Context, roomCode string, memberID int64) error {
	room, err := s.repo.GetRoomByCode(ctx, roomCode)
	if err != nil {
		return err
	}
	_ = s.cache.RemovePresence(ctx, roomCode, memberID)
	_ = s.repo.SetMemberInactive(ctx, memberID)
	return s.syncPresence(ctx, room, memberID, "member.left")
}

func (s *RoomService) syncPresence(ctx context.Context, room *models.Room, memberID int64, eventType string) error {
	currentTask, err := s.repo.GetCurrentTask(ctx, room.ID)
	if err != nil {
		return err
	}
	if currentTask != nil && currentTask.Status == models.TaskStatusVoting {
		votes, err := s.repo.ListVotesByTask(ctx, currentTask.ID)
		if err != nil {
			return err
		}
		activeMemberIDs, err := s.activeMemberIDs(ctx, room)
		if err != nil {
			return err
		}
		_ = s.cache.SetVoteProgress(ctx, room.Code, currentTask.ID, len(votes), len(activeMemberIDs))
		if ShouldReveal(activeMemberIDs, votes) {
			averageLabel, averageNumeric := CalculateAverageLabel(votes)
			if err := s.repo.RevealTask(ctx, currentTask.ID, averageLabel, averageNumeric); err != nil {
				return err
			}
			_ = s.cache.SetCurrentTask(ctx, room.Code, currentTask.ID, models.TaskStatusRevealed)
			_ = s.publish(ctx, room.Code, "results.revealed", map[string]any{"taskId": currentTask.ID})
			_ = s.publish(ctx, room.Code, "history.updated", map[string]any{"taskId": currentTask.ID})
			return nil
		}
	}
	return s.publish(ctx, room.Code, eventType, map[string]any{"memberId": memberID})
}

func (s *RoomService) syncRevealState(ctx context.Context, room *models.Room, member *models.Member) (*models.RoomSnapshot, error) {
	task, err := s.repo.GetCurrentTask(ctx, room.ID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return s.buildSnapshot(ctx, room, member)
	}

	votes, err := s.repo.ListVotesByTask(ctx, task.ID)
	if err != nil {
		return nil, err
	}
	activeMemberIDs, err := s.activeMemberIDs(ctx, room)
	if err != nil {
		return nil, err
	}

	_ = s.cache.SetVoteProgress(ctx, room.Code, task.ID, len(votes), len(activeMemberIDs))
	_ = s.publish(ctx, room.Code, "vote.progress", map[string]any{
		"taskId":     task.ID,
		"votedCount": len(votes),
		"totalCount": len(activeMemberIDs),
	})

	if ShouldReveal(activeMemberIDs, votes) {
		averageLabel, averageNumeric := CalculateAverageLabel(votes)
		if err := s.repo.RevealTask(ctx, task.ID, averageLabel, averageNumeric); err != nil {
			return nil, err
		}
		_ = s.cache.SetCurrentTask(ctx, room.Code, task.ID, models.TaskStatusRevealed)
		_ = s.publish(ctx, room.Code, "results.revealed", map[string]any{"taskId": task.ID})
		_ = s.publish(ctx, room.Code, "history.updated", map[string]any{"taskId": task.ID})
	}
	return s.buildSnapshot(ctx, room, member)
}

func (s *RoomService) authorizeMember(ctx context.Context, roomCode, token string) (*models.Room, *models.Member, error) {
	roomCode = normalizeRoomCode(roomCode)
	if !isValidRoomCode(roomCode) {
		return nil, nil, newAppError(http.StatusBadRequest, "کد اتاق نامعتبر است.")
	}
	if token == "" {
		return nil, nil, newAppError(http.StatusUnauthorized, "توکن جلسه ارسال نشده است.")
	}

	room, err := s.repo.GetRoomByCode(ctx, roomCode)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, newAppError(http.StatusNotFound, "اتاق پیدا نشد.")
	}
	if err != nil {
		return nil, nil, err
	}

	member, err := s.repo.GetMemberByToken(ctx, token)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, newAppError(http.StatusUnauthorized, "جلسه معتبر نیست.")
	}
	if err != nil {
		return nil, nil, err
	}
	if member.RoomID != room.ID {
		return nil, nil, newAppError(http.StatusForbidden, "این جلسه به این اتاق تعلق ندارد.")
	}

	return room, member, nil
}

func (s *RoomService) buildSnapshot(ctx context.Context, room *models.Room, viewer *models.Member) (*models.RoomSnapshot, error) {
	members, err := s.repo.ListMembers(ctx, room.ID)
	if err != nil {
		return nil, err
	}
	currentTask, err := s.repo.GetCurrentTask(ctx, room.ID)
	if err != nil {
		return nil, err
	}
	history, err := s.buildHistory(ctx, room.ID, viewer.ID)
	if err != nil {
		return nil, err
	}

	snapshot := &models.RoomSnapshot{
		Room:    *room,
		Me:      models.MemberView{ID: viewer.ID, DisplayName: viewer.DisplayName, IsHost: viewer.IsHost, IsActive: viewer.IsActive},
		History: history,
	}

	memberViews := make([]models.MemberView, 0, len(members))
	for _, member := range members {
		memberViews = append(memberViews, models.MemberView{
			ID:          member.ID,
			DisplayName: member.DisplayName,
			IsHost:      member.IsHost,
			IsActive:    member.IsActive,
		})
	}
	snapshot.Members = memberViews

	if currentTask == nil {
		return snapshot, nil
	}

	taskState := &models.TaskState{
		ID:           currentTask.ID,
		Title:        currentTask.Title,
		Status:       currentTask.Status,
		AverageLabel: currentTask.AverageLabel,
		RevealedAt:   currentTask.RevealedAt,
	}

	votes, err := s.repo.ListVotesByTask(ctx, currentTask.ID)
	if err != nil {
		return nil, err
	}
	var myVote *string
	for _, vote := range votes {
		if vote.MemberID == viewer.ID {
			myVote = &vote.Value
			break
		}
	}
	taskState.MyVote = myVote
	snapshot.CurrentTask = taskState

	progress, err := s.cache.GetVoteProgress(ctx, room.Code, currentTask.ID)
	if err == nil && progress != nil {
		snapshot.VotingProgress = progress
	} else {
		activeMemberIDs, activeErr := s.activeMemberIDs(ctx, room)
		if activeErr == nil {
			snapshot.VotingProgress = &models.VotingProgress{
				VotedCount: len(votes),
				TotalCount: len(activeMemberIDs),
			}
		}
	}

	if currentTask.Status == models.TaskStatusRevealed {
		taskState.Votes = s.joinMemberVotes(members, votes)
	}
	return snapshot, nil
}

func (s *RoomService) buildHistory(ctx context.Context, roomID, _ int64) ([]models.HistoryTask, error) {
	tasks, err := s.repo.ListHistory(ctx, roomID)
	if err != nil {
		return nil, err
	}
	members, err := s.repo.ListMembers(ctx, roomID)
	if err != nil {
		return nil, err
	}
	memberNames := make(map[int64]string, len(members))
	for _, member := range members {
		memberNames[member.ID] = member.DisplayName
	}

	history := make([]models.HistoryTask, 0, len(tasks))
	for _, task := range tasks {
		votes, err := s.repo.ListVotesByTask(ctx, task.ID)
		if err != nil {
			return nil, err
		}
		memberVotes := make([]models.MemberVoteView, 0, len(votes))
		for _, vote := range votes {
			memberVotes = append(memberVotes, models.MemberVoteView{
				MemberID:     vote.MemberID,
				DisplayName:  memberNames[vote.MemberID],
				Value:        vote.Value,
				NumericValue: vote.NumericValue,
			})
		}
		history = append(history, models.HistoryTask{
			ID:           task.ID,
			Title:        task.Title,
			Status:       task.Status,
			AverageLabel: task.AverageLabel,
			CreatedAt:    task.CreatedAt,
			CompletedAt:  task.CompletedAt,
			Votes:        memberVotes,
		})
	}
	return history, nil
}

func (s *RoomService) joinMemberVotes(members []models.Member, votes []models.Vote) []models.MemberVoteView {
	memberNames := make(map[int64]string, len(members))
	for _, member := range members {
		memberNames[member.ID] = member.DisplayName
	}
	memberVotes := make([]models.MemberVoteView, 0, len(votes))
	for _, vote := range votes {
		memberVotes = append(memberVotes, models.MemberVoteView{
			MemberID:     vote.MemberID,
			DisplayName:  memberNames[vote.MemberID],
			Value:        vote.Value,
			NumericValue: vote.NumericValue,
		})
	}
	return memberVotes
}

func (s *RoomService) activeMemberIDs(ctx context.Context, room *models.Room) ([]int64, error) {
	activeIDs, err := s.cache.ActiveMemberIDs(ctx, room.Code)
	if err == nil && len(activeIDs) > 0 {
		return activeIDs, nil
	}

	members, err := s.repo.ListMembers(ctx, room.ID)
	if err != nil {
		return nil, err
	}
	result := make([]int64, 0, len(members))
	for _, member := range members {
		if member.IsActive {
			result = append(result, member.ID)
		}
	}
	return result, nil
}

func (s *RoomService) publish(ctx context.Context, roomCode, eventType string, payload map[string]any) error {
	return s.cache.PublishRoomEvent(ctx, models.RoomEvent{
		Type:      eventType,
		RoomCode:  roomCode,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	})
}

func (s *RoomService) generateUniqueRoomCode(ctx context.Context) (string, error) {
	for range 10 {
		code := randomRoomCode()
		_, err := s.repo.GetRoomByCode(ctx, code)
		if errors.Is(err, sql.ErrNoRows) {
			return code, nil
		}
		if err != nil {
			return "", err
		}
	}
	return "", newAppError(http.StatusInternalServerError, "ساخت کد اتاق ممکن نشد.")
}

var roomAlphabet = []rune("ABCDEFGHJKLMNPQRSTUVWXYZ23456789")

func randomRoomCode() string {
	result := make([]rune, 6)
	for i := range result {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(roomAlphabet))))
		result[i] = roomAlphabet[index.Int64()]
	}
	return string(result)
}

func randomToken(length int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		result[i] = alphabet[index.Int64()]
	}
	return string(result)
}

func normalizeRoomCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func isValidRoomCode(code string) bool {
	if len(code) != 6 {
		return false
	}
	for _, char := range code {
		if !strings.ContainsRune("ABCDEFGHJKLMNPQRSTUVWXYZ23456789", char) {
			return false
		}
	}
	return true
}
