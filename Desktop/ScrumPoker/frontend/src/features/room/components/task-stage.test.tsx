import { render, screen } from "@testing-library/react";

import type { RoomSnapshot } from "../types";
import { TaskStage } from "./task-stage";

const baseRoom: RoomSnapshot = {
  room: {
    id: 1,
    code: "ABC123",
    hostMemberId: 1,
    status: "lobby",
    createdAt: new Date().toISOString()
  },
  me: {
    id: 1,
    displayName: "مرصاد",
    isHost: true,
    isActive: true
  },
  members: [],
  history: []
};

describe("TaskStage", () => {
  it("renders voting state", () => {
    render(
      <TaskStage
        room={{
          ...baseRoom,
          currentTask: {
            id: 10,
            title: "پیاده‌سازی صفحه ورود",
            status: "voting"
          }
        }}
        onCreateTask={() => undefined}
        onSubmitVote={() => undefined}
      />
    );

    expect(screen.getByText("رأی‌گیری فعال")).toBeInTheDocument();
    expect(screen.getByText("ثبت رأی")).toBeInTheDocument();
  });

  it("renders waiting state", () => {
    render(
      <TaskStage
        room={{
          ...baseRoom,
          me: { ...baseRoom.me, isHost: false },
          currentTask: {
            id: 10,
            title: "پیاده‌سازی صفحه ورود",
            status: "voting",
            myVote: "5"
          },
          votingProgress: {
            votedCount: 3,
            totalCount: 5
          }
        }}
        onCreateTask={() => undefined}
        onSubmitVote={() => undefined}
      />
    );

    expect(screen.getByText("رأی شما ثبت شد")).toBeInTheDocument();
    expect(screen.getByText(/۳ از ۵/)).toBeInTheDocument();
  });

  it("renders revealed state", () => {
    render(
      <TaskStage
        room={{
          ...baseRoom,
          currentTask: {
            id: 10,
            title: "پیاده‌سازی صفحه ورود",
            status: "revealed",
            averageLabel: "4",
            votes: [
              { memberId: 1, displayName: "مرصاد", value: "3" },
              { memberId: 2, displayName: "ندا", value: "5" }
            ]
          }
        }}
        onCreateTask={() => undefined}
        onSubmitVote={() => undefined}
      />
    );

    expect(screen.getByText("میانگین: 4")).toBeInTheDocument();
    expect(screen.getByText("ندا")).toBeInTheDocument();
  });
});
