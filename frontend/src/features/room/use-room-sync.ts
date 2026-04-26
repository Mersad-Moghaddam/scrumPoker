import { useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";

import { wsUrl } from "../../lib/api";
import type { RoomEvent } from "./types";

export function useRoomSync(roomCode: string | undefined, token: string | undefined) {
  const queryClient = useQueryClient();

  useEffect(() => {
    if (!roomCode || !token) {
      return;
    }

    const socket = new WebSocket(wsUrl(roomCode, token));

    socket.addEventListener("message", (event) => {
      const payload = JSON.parse(event.data) as RoomEvent;
      if (payload.type !== "system.connected") {
        queryClient.invalidateQueries({ queryKey: ["room", roomCode] });
      }
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
