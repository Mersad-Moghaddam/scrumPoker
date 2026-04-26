import { useMutation } from "@tanstack/react-query";
import { ArrowLeft, Sparkles, UsersRound } from "lucide-react";
import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";

import { AppShell } from "../../components/layout/app-shell";
import { Card } from "../../components/ui/card";
import { Button } from "../../components/ui/button";
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
      <div className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <Card className="relative overflow-hidden p-8 sm:p-10">
          <div className="pointer-events-none absolute inset-0">
            <div className="absolute -right-16 top-0 h-48 w-48 rounded-full bg-brand-500/20 blur-3xl" />
            <div className="absolute -left-10 bottom-0 h-56 w-56 rounded-full bg-accent-500/20 blur-3xl" />
          </div>
          <div className="relative space-y-6">
            <span className="inline-flex items-center gap-2 rounded-full bg-white/10 px-4 py-2 text-sm text-brand-100">
              <Sparkles size={16} />
              برای پلنینگ سریع، شفاف و خوش‌حال‌کننده
            </span>
            <div className="space-y-4">
              <h2 className="max-w-2xl text-4xl font-black leading-tight sm:text-5xl">
                برآورد تیمی را از حالت خشک و خسته‌کننده خارج کنید.
              </h2>
              <p className="max-w-xl text-base leading-8 text-app-muted">
                یک اتاق بسازید، تیم را با کد کوتاه دعوت کنید، برای هر تسک فقط عنوان وارد کنید و
                بگذارید سیستم درست در لحظه مناسب نتایج را رو کند.
              </p>
            </div>
            <div className="grid gap-4 sm:grid-cols-3">
              <Metric title="فارسـی و RTL" value="از روز اول" />
              <Metric title="رأی مخفی" value="تا رأی آخر" />
              <Metric title="تجربه سریع" value="کمتر از ۳۰ ثانیه" />
            </div>
          </div>
        </Card>

        <div className="space-y-6">
          <Card className="space-y-4">
            <div className="space-y-1">
              <h3 className="text-xl font-bold">اول اسم خودت را وارد کن</h3>
              <p className="text-sm text-app-muted">همین اسم در اتاق و نتایج نمایش داده می‌شود.</p>
            </div>
            <label className="block space-y-2">
              <span className="text-sm text-app-muted">نام نمایشی</span>
              <input
                value={displayName}
                onChange={(event) => setDisplayName(event.target.value)}
                className="w-full rounded-2xl border border-white/10 bg-surface-2 px-4 py-3 text-base outline-none transition focus:border-brand-400"
                placeholder="مثلا مرصاد"
              />
            </label>
            {errorMessage ? <p className="text-sm text-rose-300">{errorMessage}</p> : null}
          </Card>

          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-1 xl:grid-cols-2">
            <Card className="space-y-4">
              <div className="flex items-center gap-3">
                <span className="rounded-2xl bg-brand-500/15 p-3 text-brand-200">
                  <UsersRound size={18} />
                </span>
                <div>
                  <h3 className="font-bold">ساخت اتاق جدید</h3>
                  <p className="text-sm text-app-muted">میزبان شوید و رأی‌گیری را شروع کنید.</p>
                </div>
              </div>
              <form onSubmit={handleCreate} className="space-y-4">
                <Button type="submit" className="w-full gap-2" disabled={!hasName || createMutation.isPending}>
                  <span>{createMutation.isPending ? "در حال ساخت..." : "ساخت اتاق"}</span>
                  <ArrowLeft size={16} />
                </Button>
              </form>
            </Card>

            <Card className="space-y-4">
              <div>
                <h3 className="font-bold">پیوستن به اتاق</h3>
                <p className="text-sm text-app-muted">کد اتاق را وارد کنید و مستقیم وارد رأی‌گیری شوید.</p>
              </div>
              <form onSubmit={handleJoin} className="space-y-4">
                <label className="block space-y-2">
                  <span className="text-sm text-app-muted">کد اتاق</span>
                  <input
                    value={roomCode}
                    onChange={(event) => setRoomCode(event.target.value.toUpperCase())}
                    className="w-full rounded-2xl border border-white/10 bg-surface-2 px-4 py-3 font-mono text-left tracking-[0.35em] outline-none transition focus:border-brand-400"
                    dir="ltr"
                    maxLength={6}
                    placeholder="A1B2C3"
                  />
                </label>
                <Button
                  type="submit"
                  variant="secondary"
                  className="w-full"
                  disabled={!hasName || roomCode.trim().length !== 6 || joinMutation.isPending}
                >
                  {joinMutation.isPending ? "در حال ورود..." : "ورود به اتاق"}
                </Button>
              </form>
            </Card>
          </div>
        </div>
      </div>
    </AppShell>
  );
}

function Metric({ title, value }: { title: string; value: string }) {
  return (
    <div className="rounded-3xl border border-white/10 bg-white/5 p-4 backdrop-blur">
      <p className="text-sm text-app-muted">{title}</p>
      <p className="mt-2 text-xl font-bold">{value}</p>
    </div>
  );
}
