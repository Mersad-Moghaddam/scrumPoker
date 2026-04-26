import { CheckCircle2, Coffee, Infinity as InfinityIcon, HelpCircle, Hourglass } from "lucide-react";
import { useEffect, useState } from "react";

import { Button } from "../../../components/ui/button";
import { Card } from "../../../components/ui/card";
import { cn, toPersianDigits } from "../../../lib/utils";
import type { RoomSnapshot } from "../types";
import { CreateTaskCard } from "./create-task-card";

const voteOptions = ["0.5", "1", "2", "3", "5", "8", "13", "21", "34", "?", "∞", "Coffee"];

type Props = {
  room: RoomSnapshot;
  creatingTask?: boolean;
  submittingVote?: boolean;
  onCreateTask: (title: string) => void;
  onSubmitVote: (value: string) => void;
};

export function TaskStage({ room, creatingTask, submittingVote, onCreateTask, onSubmitVote }: Props) {
  const [selectedVote, setSelectedVote] = useState<string | null>(room.currentTask?.myVote ?? null);

  useEffect(() => {
    setSelectedVote(room.currentTask?.myVote ?? null);
  }, [room.currentTask?.id, room.currentTask?.myVote]);

  if (!room.currentTask) {
    return room.me.isHost ? (
      <CreateTaskCard onSubmit={onCreateTask} pending={creatingTask} />
    ) : (
      <IdleParticipantCard />
    );
  }

  if (room.currentTask.status === "voting" && !room.currentTask.myVote) {
    return (
      <Card className="space-y-4 p-4">
        <h2 className="text-lg font-bold">رأی‌گیری فعال</h2>
        <p className="text-sm text-app-muted">{room.currentTask.title}</p>

        <div className="grid grid-cols-3 gap-2 sm:grid-cols-4 xl:grid-cols-6">
          {voteOptions.map((value) => (
            <button
              key={value}
              type="button"
              onClick={() => setSelectedVote(value)}
              className={cn(
                "flex min-h-20 items-center justify-center rounded-xl border border-white/10 bg-surface-2 p-3 text-center transition hover:border-brand-400/60",
                selectedVote === value && "border-brand-400 bg-brand-500/15"
              )}
            >
              <span className="text-xl font-bold">{voteLabel(value)}</span>
            </button>
          ))}
        </div>

        <Button
          className="w-full sm:w-auto"
          disabled={!selectedVote || submittingVote}
          onClick={() => selectedVote && onSubmitVote(selectedVote)}
        >
          {submittingVote ? "..." : "ثبت رأی"}
        </Button>
      </Card>
    );
  }

  if (room.currentTask.status === "voting" && room.currentTask.myVote) {
    return (
      <Card className="space-y-4 p-4 text-center">
        <div className="mx-auto inline-flex rounded-full bg-emerald-500/15 p-3 text-emerald-300">
          <CheckCircle2 size={22} />
        </div>
        <h2 className="text-lg font-bold">رأی شما ثبت شد</h2>
        <p className="text-sm text-app-muted">{room.currentTask.title}</p>

        <div className="rounded-xl border border-white/10 bg-surface-2/80 p-4">
          <p className="text-xs text-app-muted">
            {toPersianDigits(room.votingProgress?.votedCount ?? 0)} از{" "}
            {toPersianDigits(room.votingProgress?.totalCount ?? 0)}
          </p>
          <div className="mt-2 h-2 overflow-hidden rounded-full bg-white/10">
            <div
              className="h-full rounded-full bg-brand-500 transition-all duration-300"
              style={{
                width: `${Math.min(
                  100,
                  ((room.votingProgress?.votedCount ?? 0) / Math.max(room.votingProgress?.totalCount ?? 1, 1)) * 100
                )}%`
              }}
            />
          </div>
        </div>
      </Card>
    );
  }

  return (
    <Card className="space-y-4 p-4">
      <div className="flex items-center justify-between gap-2">
        <h2 className="text-lg font-bold">نتیجه</h2>
        <span className="rounded-full bg-brand-500/15 px-3 py-1 text-xs font-semibold text-brand-200">
          میانگین: {room.currentTask.averageLabel ?? "-"}
        </span>
      </div>
      <p className="text-sm text-app-muted">{room.currentTask.title}</p>

      <div className="grid gap-2 md:grid-cols-2">
        {room.currentTask.votes?.map((vote) => (
          <div
            key={vote.memberId}
            className="flex items-center justify-between rounded-xl border border-white/10 bg-surface-2/80 px-3 py-2"
          >
            <span className="text-sm font-semibold">{vote.displayName}</span>
            <span className="rounded-lg bg-white/10 px-2 py-1 font-mono text-sm text-brand-100">{vote.value}</span>
          </div>
        ))}
      </div>

      {room.me.isHost ? (
        <CreateTaskCard onSubmit={onCreateTask} pending={creatingTask} />
      ) : (
        <div className="rounded-xl border border-dashed border-white/10 bg-surface-2/60 p-3 text-sm text-app-muted">
          منتظر میزبان
        </div>
      )}
    </Card>
  );
}

function IdleParticipantCard() {
  return (
    <Card className="space-y-3 p-4 text-center">
      <div className="mx-auto inline-flex rounded-full bg-white/10 p-3 text-app-text">
        <Hourglass size={22} />
      </div>
      <h2 className="text-lg font-bold">منتظر میزبان</h2>
    </Card>
  );
}

function voteLabel(value: string) {
  if (value === "Coffee") {
    return <Coffee size={22} />;
  }
  if (value === "∞") {
    return <InfinityIcon size={22} />;
  }
  if (value === "?") {
    return <HelpCircle size={22} />;
  }
  return value;
}
