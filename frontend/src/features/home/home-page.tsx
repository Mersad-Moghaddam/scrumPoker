import { useMutation } from "@tanstack/react-query";
import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";

import { AppShell } from "../../components/layout/app-shell";
import { Button } from "../../components/ui/button";
import { Card } from "../../components/ui/card";
import { createRoom, joinRoom } from "../../lib/api";
import { saveRoomSession } from "../../lib/session";

export function HomePage() {
  const navigate = useNavigate();
  const [displayName, setDisplayName] = useState("");
  const [roomCode, setRoomCode] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  const createMutation = useMutation({
    mutationFn: () => createRoom(displayName),
    onSuccess: (session) => {
      saveRoomSession(session.roomCode, { token: session.token, memberId: session.memberId });
      navigate(`/room/${session.roomCode}`);
    },
    onError: (error: Error) => setErrorMessage(error.message)
  });

  const joinMutation = useMutation({
    mutationFn: () => joinRoom(displayName, roomCode),
    onSuccess: (session) => {
      saveRoomSession(session.roomCode, { token: session.token, memberId: session.memberId });
      navigate(`/room/${session.roomCode}`);
    },
    onError: (error: Error) => setErrorMessage(error.message)
  });

  const hasName = displayName.trim().length >= 2;

  const handleCreate = (event: FormEvent) => {
    event.preventDefault();
    setErrorMessage("");
    createMutation.mutate();
  };

  const handleJoin = (event: FormEvent) => {
    event.preventDefault();
    setErrorMessage("");
    joinMutation.mutate();
  };

  return (
    <AppShell>
      <div className="mx-auto max-w-sm">
        <Card className="space-y-4 p-4">
          <label className="block space-y-2">
            <span className="text-sm text-app-muted">نام</span>
            <input
              value={displayName}
              onChange={(event) => setDisplayName(event.target.value)}
              className="w-full rounded-xl border border-white/10 bg-surface-2 px-3 py-2 outline-none transition focus:border-brand-400"
            />
          </label>

          <label className="block space-y-2">
            <span className="text-sm text-app-muted">کد</span>
            <input
              value={roomCode}
              onChange={(event) => setRoomCode(event.target.value.toUpperCase())}
              className="w-full rounded-xl border border-white/10 bg-surface-2 px-3 py-2 font-mono text-left tracking-[0.25em] outline-none transition focus:border-brand-400"
              dir="ltr"
              maxLength={6}
            />
          </label>

          {errorMessage ? <p className="text-sm text-rose-300">{errorMessage}</p> : null}

          <div className="grid grid-cols-2 gap-2">
            <form onSubmit={handleCreate}>
              <Button type="submit" className="w-full" disabled={!hasName || createMutation.isPending}>
                {createMutation.isPending ? "..." : "ساخت"}
              </Button>
            </form>
            <form onSubmit={handleJoin}>
              <Button
                type="submit"
                variant="secondary"
                className="w-full"
                disabled={!hasName || roomCode.trim().length !== 6 || joinMutation.isPending}
              >
                {joinMutation.isPending ? "..." : "ورود"}
              </Button>
            </form>
          </div>
        </Card>
      </div>
    </AppShell>
  );
}
