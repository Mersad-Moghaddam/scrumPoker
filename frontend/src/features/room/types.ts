export type MemberView = {
  id: number;
  displayName: string;
  isHost: boolean;
  isActive: boolean;
};

export type MemberVoteView = {
  memberId: number;
  displayName: string;
  value: string;
  numericValue?: number | null;
};

export type TaskState = {
  id: number;
  title: string;
  status: "voting" | "revealed" | "archived";
  myVote?: string | null;
  averageLabel?: string | null;
  revealedAt?: string | null;
  votes?: MemberVoteView[];
};

export type VotingProgress = {
  votedCount: number;
  totalCount: number;
};

export type HistoryTask = {
  id: number;
  title: string;
  status: string;
  averageLabel?: string | null;
  createdAt: string;
  completedAt?: string | null;
  votes: MemberVoteView[];
};

export type RoomSnapshot = {
  room: {
    id: number;
    code: string;
    hostMemberId: number;
    status: string;
    createdAt: string;
  };
  me: MemberView;
  members: MemberView[];
  currentTask?: TaskState | null;
  votingProgress?: VotingProgress | null;
  history: HistoryTask[];
};

export type AuthSession = {
  roomCode: string;
  token: string;
  memberId: number;
};

export type RoomEvent = {
  type: string;
  roomCode: string;
  timestamp: string;
  payload?: Record<string, unknown>;
};
