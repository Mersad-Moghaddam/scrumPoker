package models

import "time"

const (
	RoomStatusLobby    = "lobby"
	TaskStatusVoting   = "voting"
	TaskStatusRevealed = "revealed"
	TaskStatusArchived = "archived"
)

var ValidVoteValues = []string{"0.5", "1", "2", "3", "5", "8", "13", "21", "34", "?", "∞", "Coffee"}

type Room struct {
	ID           int64     `json:"id"`
	Code         string    `json:"code"`
	HostMemberID int64     `json:"hostMemberId"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Member struct {
	ID           int64      `json:"id"`
	RoomID       int64      `json:"roomId"`
	DisplayName  string     `json:"displayName"`
	IsHost       bool       `json:"isHost"`
	IsActive     bool       `json:"isActive"`
	SessionToken string     `json:"-"`
	JoinedAt     time.Time  `json:"joinedAt"`
	LastSeenAt   *time.Time `json:"lastSeenAt,omitempty"`
}

type Task struct {
	ID             int64      `json:"id"`
	RoomID         int64      `json:"roomId"`
	Title          string     `json:"title"`
	Status         string     `json:"status"`
	AverageLabel   *string    `json:"averageLabel,omitempty"`
	AverageNumeric *float64   `json:"averageNumeric,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	RevealedAt     *time.Time `json:"revealedAt,omitempty"`
	CompletedAt    *time.Time `json:"completedAt,omitempty"`
}

type Vote struct {
	ID           int64     `json:"id"`
	TaskID       int64     `json:"taskId"`
	MemberID     int64     `json:"memberId"`
	Value        string    `json:"value"`
	NumericValue *float64  `json:"numericValue,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type AuthSession struct {
	RoomCode string `json:"roomCode"`
	Token    string `json:"token"`
	MemberID int64  `json:"memberId"`
}

type MemberView struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"displayName"`
	IsHost      bool   `json:"isHost"`
	IsActive    bool   `json:"isActive"`
}

type MemberVoteView struct {
	MemberID     int64    `json:"memberId"`
	DisplayName  string   `json:"displayName"`
	Value        string   `json:"value"`
	NumericValue *float64 `json:"numericValue,omitempty"`
}

type TaskState struct {
	ID           int64            `json:"id"`
	Title        string           `json:"title"`
	Status       string           `json:"status"`
	MyVote       *string          `json:"myVote,omitempty"`
	AverageLabel *string          `json:"averageLabel,omitempty"`
	RevealedAt   *time.Time       `json:"revealedAt,omitempty"`
	Votes        []MemberVoteView `json:"votes,omitempty"`
}

type VotingProgress struct {
	VotedCount int `json:"votedCount"`
	TotalCount int `json:"totalCount"`
}

type HistoryTask struct {
	ID           int64            `json:"id"`
	Title        string           `json:"title"`
	Status       string           `json:"status"`
	AverageLabel *string          `json:"averageLabel,omitempty"`
	CreatedAt    time.Time        `json:"createdAt"`
	CompletedAt  *time.Time       `json:"completedAt,omitempty"`
	Votes        []MemberVoteView `json:"votes"`
}

type RoomSnapshot struct {
	Room           Room            `json:"room"`
	Me             MemberView      `json:"me"`
	Members        []MemberView    `json:"members"`
	CurrentTask    *TaskState      `json:"currentTask,omitempty"`
	VotingProgress *VotingProgress `json:"votingProgress,omitempty"`
	History        []HistoryTask   `json:"history"`
}

type RoomEvent struct {
	Type      string         `json:"type"`
	RoomCode  string         `json:"roomCode"`
	Timestamp time.Time      `json:"timestamp"`
	Payload   map[string]any `json:"payload,omitempty"`
}
