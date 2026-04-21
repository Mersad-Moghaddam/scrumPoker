import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ArrowRight, Copy, DoorOpen, LoaderCircle } from "lucide-react";
import { useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { AppShell } from "../../components/layout/app-shell";
import { Button } from "../../components/ui/button";
import { Card } from "../../components/ui/card";
import { createTask, fetchRoom, submitVote } from "../../lib/api";
import { clearRoomSession, getRoomSession } from "../../lib/session";
import { toPersianDigits } from "../../lib/utils";
import { HistoryPanel } from "./components/history-panel";
import { MembersPanel } from "./components/members-panel";
import { TaskStage } from "./components/task-stage";
import { useRoomSync } from "./use-room-sync";

export function RoomPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const params = useParams();
  const roomCode = params.code?.toUpperCase();
  const session = useMemo(() => (roomCode ? getRoomSession(roomCode) : null), [roomCode]);
  const [actionError, setActionError] = useState("");

  useRoomSync(roomCode, session?.token);

  const roomQuery = useQuery({
    queryKey: ["room", roomCode],
    queryFn: () => fetchRoom(roomCode!, session!.token),
    enabled: Boolean(roomCode && session?.token)
  });

  const createTaskMutation = useMutation({
    mutationFn: (title: string) => createTask(roomCode!, title, session!.token),
    onSuccess: (snapshot) => {
      setActionError("");
      queryClient.setQueryData(["room", roomCode], snapshot);
    },
    onError: (error: Error) => setActionError(error.message)
  });

  const voteMutation = useMutation({
    mutationFn: (value: string) => submitVote(roomCode!, roomQuery.data!.currentTask!.id, value, session!.token),
    onSuccess: (snapshot) => {
      setActionError("");
      queryClient.setQueryData(["room", roomCode], snapshot);
    },
    onError: (error: Error) => setActionError(error.message)
  });

  if (!session || !roomCode) {
    return (
      <AppShell>
        <Card className="space-y-4 text-center">
          <h2 className="text-2xl font-black">جلسه این اتاق پیدا نشد</h2>
          <p className="text-app-muted">برای ورود دوباره، از صفحه اصلی وارد اتاق شوید.</p>
          <Link to="/">
            <Button className="mx-auto">بازگشت به خانه</Button>
          </Link>
        </Card>
      </AppShell>
    );
  }

  if (roomQuery.isLoading) {
    return (
      <AppShell>
        <Card className="flex items-center justify-center gap-3 py-16 text-app-muted">
          <LoaderCircle className="animate-spin" />
          در حال بارگذاری وضعیت اتاق...
        </Card>
      </AppShell>
    );
  }

  if (roomQuery.isError || !roomQuery.data) {
    return (
      <AppShell>
        <Card className="space-y-4 text-center">
          <h2 className="text-2xl font-black">ورود به اتاق ممکن نشد</h2>
          <p className="text-app-muted">{(roomQuery.error as Error)?.message ?? "لطفا دوباره تلاش کنید."}</p>
          <div className="flex justify-center gap-3">
            <Link to="/">
              <Button>بازگشت به خانه</Button>
            </Link>
            <Button variant="ghost" onClick={() => roomQuery.refetch()}>
              تلاش دوباره
            </Button>
          </div>
        </Card>
      </AppShell>
    );
  }

  const room = roomQuery.data;

  return (
    <AppShell>
      <div className="space-y-6">
        <Card className="overflow-hidden">
          <div className="grid gap-6 lg:grid-cols-[1fr_auto] lg:items-center">
            <div className="space-y-3">
              <p className="text-sm text-app-muted">کد اتاق</p>
              <div className="flex flex-wrap items-center gap-3">
                <span className="rounded-2xl bg-surface-2 px-4 py-3 font-mono text-2xl tracking-[0.35em]" dir="ltr">
                  {room.room.code}
                </span>
                <button
                  type="button"
                  onClick={() => navigator.clipboard.writeText(room.room.code)}
                  className="inline-flex items-center gap-2 rounded-2xl bg-white/5 px-3 py-2 text-sm text-app-muted transition hover:bg-white/10"
                >
                  <Copy size={16} />
                  کپی
                </button>
              </div>
              <p className="text-sm text-app-muted">
                اعضای فعال: {toPersianDigits(room.members.filter((member) => member.isActive).length)} نفر
              </p>
            </div>

            <div className="flex flex-wrap gap-3">
              <Button
                variant="ghost"
                className="gap-2"
                onClick={() => {
                  clearRoomSession(roomCode);
                  navigate("/");
                }}
              >
                <DoorOpen size={16} />
                خروج از اتاق
              </Button>
              <Link to="/" className="inline-flex">
                <Button variant="secondary" className="gap-2">
                  <ArrowRight size={16} />
                  خانه
                </Button>
              </Link>
            </div>
          </div>
        </Card>

        <div className="grid gap-6 xl:grid-cols-[280px_minmax(0,1fr)_320px]">
          <MembersPanel members={room.members} meId={room.me.id} />

          <div className="space-y-4">
            {actionError ? (
              <Card className="border border-rose-400/30 bg-rose-500/10 text-sm text-rose-100">
                {actionError}
              </Card>
            ) : null}
            <TaskStage
              room={room}
              creatingTask={createTaskMutation.isPending}
              submittingVote={voteMutation.isPending}
              onCreateTask={(title) => createTaskMutation.mutate(title)}
              onSubmitVote={(value) => voteMutation.mutate(value)}
            />
          </div>

          <HistoryPanel items={room.history} />
        </div>
      </div>
    </AppShell>
  );
}
