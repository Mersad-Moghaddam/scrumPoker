import type { AuthSession, HistoryTask, RoomSnapshot } from "../features/room/types";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

async function request<T>(path: string, options: RequestInit = {}, token?: string): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(options.headers ?? {})
    }
  });

  if (!response.ok) {
    const payload = (await response.json().catch(() => null)) as { message?: string } | null;
    throw new Error(payload?.message ?? "خطای ارتباط با سرور");
  }

  return response.json() as Promise<T>;
}

export function createRoom(displayName: string) {
  return request<AuthSession>("/api/rooms", {
    method: "POST",
    body: JSON.stringify({ displayName })
  });
}

export function joinRoom(displayName: string, code: string) {
  return request<AuthSession>("/api/rooms/join", {
    method: "POST",
    body: JSON.stringify({ displayName, code })
  });
}

export function fetchRoom(code: string, token: string) {
  return request<RoomSnapshot>(`/api/rooms/${code}`, undefined, token);
}

export function createTask(code: string, title: string, token: string) {
  return request<RoomSnapshot>(
    `/api/rooms/${code}/tasks`,
    {
      method: "POST",
      body: JSON.stringify({ title })
    },
    token
  );
}

export function submitVote(code: string, taskId: number, value: string, token: string) {
  return request<RoomSnapshot>(
    `/api/rooms/${code}/tasks/${taskId}/votes`,
    {
      method: "POST",
      body: JSON.stringify({ value })
    },
    token
  );
}

export async function fetchHistory(code: string, token: string) {
  const response = await request<{ items: HistoryTask[] }>(`/api/rooms/${code}/history`, undefined, token);
  return response.items;
}

export function wsUrl(roomCode: string, token: string) {
  const base = (import.meta.env.VITE_WS_URL as string | undefined) ?? API_URL.replace("http", "ws");
  return `${base}/ws/rooms/${roomCode}?token=${encodeURIComponent(token)}`;
}
