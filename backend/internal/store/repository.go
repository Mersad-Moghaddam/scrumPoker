package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"scrum-poker/backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateRoom(ctx context.Context, roomCode, displayName, sessionToken string) (*models.Room, *models.Member, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO rooms (code, status)
		VALUES (?, ?)
	`, roomCode, models.RoomStatusLobby)
	if err != nil {
		_ = tx.Rollback()
		return nil, nil, err
	}
	roomID, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, nil, err
	}

	now := time.Now().UTC()
	memberRes, err := tx.ExecContext(ctx, `
		INSERT INTO members (room_id, display_name, is_host, is_active, session_token, last_seen_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, roomID, displayName, true, true, sessionToken, now)
	if err != nil {
		_ = tx.Rollback()
		return nil, nil, err
	}
	memberID, err := memberRes.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, nil, err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE rooms SET host_member_id = ? WHERE id = ?`, memberID, roomID); err != nil {
		_ = tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	room, err := r.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, nil, err
	}
	member, err := r.GetMemberByID(ctx, memberID)
	if err != nil {
		return nil, nil, err
	}
	return room, member, nil
}

func (r *Repository) CreateMember(ctx context.Context, roomID int64, displayName, sessionToken string) (*models.Member, error) {
	now := time.Now().UTC()
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO members (room_id, display_name, is_host, is_active, session_token, last_seen_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, roomID, displayName, false, true, sessionToken, now)
	if err != nil {
		return nil, err
	}
	memberID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetMemberByID(ctx, memberID)
}

func (r *Repository) GetRoomByID(ctx context.Context, roomID int64) (*models.Room, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, code, host_member_id, status, created_at
		FROM rooms
		WHERE id = ?
	`, roomID)
	return scanRoom(row)
}

func (r *Repository) GetRoomByCode(ctx context.Context, roomCode string) (*models.Room, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, code, host_member_id, status, created_at
		FROM rooms
		WHERE code = ?
	`, roomCode)
	return scanRoom(row)
}

func (r *Repository) GetMemberByID(ctx context.Context, memberID int64) (*models.Member, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, room_id, display_name, is_host, is_active, session_token, joined_at, last_seen_at
		FROM members
		WHERE id = ?
	`, memberID)
	return scanMember(row)
}

func (r *Repository) GetMemberByToken(ctx context.Context, token string) (*models.Member, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, room_id, display_name, is_host, is_active, session_token, joined_at, last_seen_at
		FROM members
		WHERE session_token = ?
	`, token)
	return scanMember(row)
}

func (r *Repository) TouchMember(ctx context.Context, memberID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE members
		SET is_active = TRUE, last_seen_at = ?
		WHERE id = ?
	`, time.Now().UTC(), memberID)
	return err
}

func (r *Repository) SetMemberInactive(ctx context.Context, memberID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE members
		SET is_active = FALSE, last_seen_at = ?
		WHERE id = ?
	`, time.Now().UTC(), memberID)
	return err
}

func (r *Repository) ListMembers(ctx context.Context, roomID int64) ([]models.Member, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, room_id, display_name, is_host, is_active, session_token, joined_at, last_seen_at
		FROM members
		WHERE room_id = ?
		ORDER BY joined_at ASC
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		member, err := scanMember(rows)
		if err != nil {
			return nil, err
		}
		members = append(members, *member)
	}
	return members, rows.Err()
}

func (r *Repository) GetCurrentTask(ctx context.Context, roomID int64) (*models.Task, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, room_id, title, status, average_label, average_numeric, created_at, revealed_at, completed_at
		FROM tasks
		WHERE room_id = ? AND status IN (?, ?)
		ORDER BY created_at DESC
		LIMIT 1
	`, roomID, models.TaskStatusVoting, models.TaskStatusRevealed)
	task, err := scanTask(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return task, err
}

func (r *Repository) ArchiveRevealedTasks(ctx context.Context, roomID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tasks
		SET status = ?
		WHERE room_id = ? AND status = ?
	`, models.TaskStatusArchived, roomID, models.TaskStatusRevealed)
	return err
}

func (r *Repository) CreateTask(ctx context.Context, roomID int64, title string) (*models.Task, error) {
	now := time.Now().UTC()
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO tasks (room_id, title, status, created_at)
		VALUES (?, ?, ?, ?)
	`, roomID, title, models.TaskStatusVoting, now)
	if err != nil {
		return nil, err
	}
	taskID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	_, err = r.db.ExecContext(ctx, `UPDATE rooms SET status = ? WHERE id = ?`, models.TaskStatusVoting, roomID)
	if err != nil {
		return nil, err
	}
	return r.GetTaskByID(ctx, taskID)
}

func (r *Repository) GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, room_id, title, status, average_label, average_numeric, created_at, revealed_at, completed_at
		FROM tasks
		WHERE id = ?
	`, taskID)
	return scanTask(row)
}

func (r *Repository) GetTaskVoteByMember(ctx context.Context, taskID, memberID int64) (*models.Vote, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, task_id, member_id, value, numeric_value, created_at
		FROM votes
		WHERE task_id = ? AND member_id = ?
	`, taskID, memberID)
	vote, err := scanVote(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return vote, err
}

func (r *Repository) CreateVote(ctx context.Context, taskID, memberID int64, value string, numericValue *float64) (*models.Vote, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO votes (task_id, member_id, value, numeric_value)
		VALUES (?, ?, ?, ?)
	`, taskID, memberID, value, numericValue)
	if err != nil {
		return nil, err
	}
	voteID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRowContext(ctx, `
		SELECT id, task_id, member_id, value, numeric_value, created_at
		FROM votes
		WHERE id = ?
	`, voteID)
	return scanVote(row)
}

func (r *Repository) ListVotesByTask(ctx context.Context, taskID int64) ([]models.Vote, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, task_id, member_id, value, numeric_value, created_at
		FROM votes
		WHERE task_id = ?
		ORDER BY created_at ASC
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []models.Vote
	for rows.Next() {
		vote, err := scanVote(rows)
		if err != nil {
			return nil, err
		}
		votes = append(votes, *vote)
	}
	return votes, rows.Err()
}

func (r *Repository) RevealTask(ctx context.Context, taskID int64, averageLabel *string, averageNumeric *float64) error {
	now := time.Now().UTC()
	task, err := r.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE tasks
		SET status = ?, average_label = ?, average_numeric = ?, revealed_at = ?, completed_at = ?
		WHERE id = ?
	`, models.TaskStatusRevealed, averageLabel, averageNumeric, now, now, taskID); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE rooms
		SET status = ?
		WHERE id = ?
	`, models.TaskStatusRevealed, task.RoomID); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) ListHistory(ctx context.Context, roomID int64) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, room_id, title, status, average_label, average_numeric, created_at, revealed_at, completed_at
		FROM tasks
		WHERE room_id = ? AND status IN (?, ?)
		ORDER BY created_at DESC
	`, roomID, models.TaskStatusRevealed, models.TaskStatusArchived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}
	return tasks, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRoom(row rowScanner) (*models.Room, error) {
	room := &models.Room{}
	err := row.Scan(&room.ID, &room.Code, &room.HostMemberID, &room.Status, &room.CreatedAt)
	return room, err
}

func scanMember(row rowScanner) (*models.Member, error) {
	member := &models.Member{}
	var lastSeen sql.NullTime
	err := row.Scan(
		&member.ID,
		&member.RoomID,
		&member.DisplayName,
		&member.IsHost,
		&member.IsActive,
		&member.SessionToken,
		&member.JoinedAt,
		&lastSeen,
	)
	if err != nil {
		return nil, err
	}
	if lastSeen.Valid {
		member.LastSeenAt = &lastSeen.Time
	}
	return member, nil
}

func scanTask(row rowScanner) (*models.Task, error) {
	task := &models.Task{}
	var averageLabel sql.NullString
	var averageNumeric sql.NullFloat64
	var revealedAt sql.NullTime
	var completedAt sql.NullTime
	err := row.Scan(
		&task.ID,
		&task.RoomID,
		&task.Title,
		&task.Status,
		&averageLabel,
		&averageNumeric,
		&task.CreatedAt,
		&revealedAt,
		&completedAt,
	)
	if err != nil {
		return nil, err
	}
	if averageLabel.Valid {
		task.AverageLabel = &averageLabel.String
	}
	if averageNumeric.Valid {
		task.AverageNumeric = &averageNumeric.Float64
	}
	if revealedAt.Valid {
		task.RevealedAt = &revealedAt.Time
	}
	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}
	return task, nil
}

func scanVote(row rowScanner) (*models.Vote, error) {
	vote := &models.Vote{}
	var numericValue sql.NullFloat64
	err := row.Scan(
		&vote.ID,
		&vote.TaskID,
		&vote.MemberID,
		&vote.Value,
		&numericValue,
		&vote.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if numericValue.Valid {
		vote.NumericValue = &numericValue.Float64
	}
	return vote, nil
}
