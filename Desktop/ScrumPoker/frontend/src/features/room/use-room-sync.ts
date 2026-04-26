import { useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";

import { wsUrl } from "../../lib/api";
import type { RoomEvent } from "./types";

const REFRESH_EVENTS = new Set([
  "room.created",
  "member.joined",
  "member.left",
  "task.created",
  "vote.progress",
  "results.revealed",
  "history.updated"
]);

export function useRoomSync(roomCode: string | undefined, token: string | undefined) {
  const queryClient = useQueryClient();

  useEffect(() => {
    if (!roomCode || !token) {
      return;
    }

    const socket = new WebSocket(wsUrl(roomCode, token));

    socket.addEventListener("message", (event) => {
      try {
        const payload = JSON.parse(event.data) as RoomEvent;
        if (REFRESH_EVENTS.has(payload.type)) {
          queryClient.invalidateQueries({ queryKey: ["room", roomCode] });
        }
      } catch {
        // ignore malformed message
      }
    });

    socket.addEventListener("error", () => {
      // no-op: keep UI stable on transient websocket errors
    });

    socket.addEventListener("close", () => {
      // no-op: room page fetches latest state with HTTP when needed
    });

    const interval = window.setInterval(() => {
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: "heartbeat" }));
      }
    }, 20_000);

    return () => {
      window.clearInterval(interval);
      socket.close();
    };
  }, [queryClient, roomCode, token]);
}
