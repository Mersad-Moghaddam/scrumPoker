import { CheckCircle2, Coffee, Infinity as InfinityIcon, HelpCircle, Hourglass, Sparkles } from "lucide-react";
import { useEffect, useState } from "react";

import { Button } from "../../../components/ui/button";
import { Card } from "../../../components/ui/card";
import { cn, toPersianDigits } from "../../../lib/utils";
import type { RoomSnapshot } from "../types";
import { CreateTaskCard } from "./create-task-card";

const voteOptions = [
  "0.5",
  "1",
  "2",
  "3",
  "5",
  "8",
  "13",
  "21",
  "34",
  "?",
  "∞",
  "Coffee"
];

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
      <Card className="space-y-6">
        <div className="space-y-2">
          <span className="inline-flex items-center gap-2 rounded-full bg-brand-500/15 px-3 py-1 text-sm text-brand-200">
            <Sparkles size={15} />
            رأی‌گیری فعال
          </span>
          <h2 className="text-2xl font-black sm:text-3xl">{room.currentTask.title}</h2>
          <p className="text-sm leading-7 text-app-muted">
            یکی از کارت‌ها را انتخاب کنید. تا وقتی همه افراد فعال رأی ندهند، هیچ رأیی نمایش داده
            نمی‌شود.
          </p>
        </div>

        <div className="grid grid-cols-3 gap-3 sm:grid-cols-4 xl:grid-cols-6">
          {voteOptions.map((value) => (
            <button
              key={value}
              type="button"
              onClick={() => setSelectedVote(value)}
              className={cn(
                "group flex min-h-28 flex-col items-center justify-center rounded-[26px] border border-white/10 bg-surface-2 p-4 text-center transition-all duration-200 hover:-translate-y-1 hover:border-brand-400/60",
                selectedVote === value && "border-brand-400 bg-brand-500/15 shadow-lg shadow-brand-500/20"
              )}
            >
              <span className="text-2xl font-black">{voteLabel(value)}</span>
              <span className="mt-2 text-xs text-app-muted">{voteHint(value)}</span>
            </button>
          ))}
        </div>

        <Button
          className="w-full sm:w-auto"
          disabled={!selectedVote || submittingVote}
          onClick={() => selectedVote && onSubmitVote(selectedVote)}
        >
          {submittingVote ? "در حال ثبت رأی..." : "ثبت رأی"}
        </Button>
      </Card>
    );
  }

  if (room.currentTask.status === "voting" && room.currentTask.myVote) {
    return (
      <Card className="space-y-6 text-center">
        <div className="mx-auto inline-flex rounded-full bg-emerald-500/15 p-4 text-emerald-300">
          <CheckCircle2 size={26} />
        </div>
        <div className="space-y-2">
          <h2 className="text-2xl font-black">رأی شما ثبت شد</h2>
          <p className="text-app-muted">
            برای تسک <span className="font-semibold text-app-text">{room.currentTask.title}</span> رأی
            شما ثبت شد. تا رأی آخر منتظر می‌مانیم.
          </p>
        </div>

        <div className="rounded-3xl border border-white/10 bg-surface-2/80 p-5">
          <p className="text-sm text-app-muted">پیشرفت رأی‌گیری</p>
          <p className="mt-2 text-3xl font-black">
            {toPersianDigits(room.votingProgress?.votedCount ?? 0)} از{" "}
            {toPersianDigits(room.votingProgress?.totalCount ?? 0)}
          </p>
          <div className="mt-4 h-3 overflow-hidden rounded-full bg-white/10">
            <div
              className="h-full rounded-full bg-brand-500 transition-all duration-300"
              style={{
                width: `${Math.min(
                  100,
                  ((room.votingProgress?.votedCount ?? 0) /
                    Math.max(room.votingProgress?.totalCount ?? 1, 1)) *
                    100
                )}%`
              }}
            />
          </div>
        </div>
      </Card>
    );
  }

  return (
    <Card className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <p className="text-sm text-app-muted">نتیجه نهایی برای تسک</p>
          <h2 className="text-2xl font-black sm:text-3xl">{room.currentTask.title}</h2>
        </div>
        <span className="rounded-full bg-brand-500/15 px-4 py-2 text-sm font-semibold text-brand-200">
          میانگین: {room.currentTask.averageLabel ?? "بدون برآورد عددی"}
        </span>
      </div>

      <div className="grid gap-3 md:grid-cols-2">
        {room.currentTask.votes?.map((vote) => (
          <div
            key={vote.memberId}
            className="flex items-center justify-between rounded-3xl border border-white/10 bg-surface-2/80 px-4 py-4"
          >
            <span className="font-semibold">{vote.displayName}</span>
            <span className="rounded-2xl bg-white/10 px-3 py-2 font-mono text-brand-100">{vote.value}</span>
          </div>
        ))}
      </div>

      {room.me.isHost ? (
        <CreateTaskCard
          onSubmit={onCreateTask}
          pending={creatingTask}
          title="آماده برای تسک بعدی؟"
          description="به‌محض شروع تسک جدید، رأی‌گیری قبلی به تاریخچه منتقل می‌شود."
        />
      ) : (
        <div className="rounded-3xl border border-dashed border-white/10 bg-surface-2/60 p-5 text-sm text-app-muted">
          میزبان تسک بعدی را شروع می‌کند. فعلا نتیجه این برآورد را با تیم جمع‌بندی کنید.
        </div>
      )}
    </Card>
  );
}

function IdleParticipantCard() {
  return (
    <Card className="space-y-4 text-center">
      <div className="mx-auto inline-flex rounded-full bg-white/10 p-4 text-app-text">
        <Hourglass size={26} />
      </div>
      <div>
        <h2 className="text-2xl font-black">منتظر میزبان هستیم</h2>
        <p className="mt-2 text-app-muted">
          به‌محض این‌که میزبان عنوان تسک را وارد کند، کارت‌های رأی برای همه نمایش داده می‌شود.
        </p>
      </div>
    </Card>
  );
}

function voteLabel(value: string) {
  if (value === "Coffee") {
    return <Coffee size={26} />;
  }
  if (value === "∞") {
    return <InfinityIcon size={26} />;
  }
  if (value === "?") {
    return <HelpCircle size={26} />;
  }
  return value;
}

function voteHint(value: string) {
  switch (value) {
    case "Coffee":
      return "خیلی کوچک";
    case "∞":
      return "اطلاعات کافی نیست";
    case "?":
      return "تسک مبهم است";
    default:
      return "برآورد عددی";
  }
}
