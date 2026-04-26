import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Copy, DoorOpen, LoaderCircle } from "lucide-react";
import { useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { AppShell } from "../../components/layout/app-shell";
import { Button } from "../../components/ui/button";
import { Card } from "../../components/ui/card";
import { createTask, fetchRoom, submitVote } from "../../lib/api";
import { clearRoomSession, getRoomSession } from "../../lib/session";
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
        <Card className="space-y-3 p-4 text-center">
          <h2 className="text-lg font-bold">جلسه نامعتبر</h2>
          <Link to="/">
            <Button className="mx-auto">خانه</Button>
          </Link>
        </Card>
      </AppShell>
    );
  }

  if (roomQuery.isLoading) {
    return (
      <AppShell>
        <Card className="flex items-center justify-center gap-2 py-8 text-app-muted">
          <LoaderCircle size={16} className="animate-spin" />
          در حال بارگذاری
        </Card>
      </AppShell>
    );
  }

  if (roomQuery.isError || !roomQuery.data) {
    return (
      <AppShell>
        <Card className="space-y-3 p-4 text-center">
          <h2 className="text-lg font-bold">خطا</h2>
          <p className="text-sm text-app-muted">{(roomQuery.error as Error)?.message ?? "دوباره تلاش کنید"}</p>
          <div className="flex justify-center gap-2">
            <Link to="/">
              <Button>خانه</Button>
            </Link>
            <Button variant="ghost" onClick={() => roomQuery.refetch()}>
              تلاش
            </Button>
          </div>
        </Card>
      </AppShell>
    );
  }

  const room = roomQuery.data;

  return (
    <AppShell>
      <div className="space-y-4">
        <Card className="p-4">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="flex items-center gap-2">
              <span className="rounded-xl bg-surface-2 px-3 py-2 font-mono tracking-[0.25em]" dir="ltr">
                {room.room.code}
              </span>
              <button
                type="button"
                onClick={() => navigator.clipboard.writeText(room.room.code)}
                className="inline-flex items-center gap-1 rounded-xl bg-white/5 px-2 py-2 text-sm text-app-muted transition hover:bg-white/10"
              >
                <Copy size={14} />
                کپی
              </button>
            </div>

            <Button
              variant="ghost"
              className="gap-1"
              onClick={() => {
                clearRoomSession(roomCode);
                navigate("/");
              }}
            >
              <DoorOpen size={14} />
              خروج
            </Button>
          </div>
        </Card>

        <div className="grid gap-4 lg:grid-cols-[240px_minmax(0,1fr)_280px]">
          <MembersPanel members={room.members} meId={room.me.id} />

          <div className="space-y-3">
            {actionError ? (
              <Card className="border border-rose-400/30 bg-rose-500/10 p-3 text-sm text-rose-100">{actionError}</Card>
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
