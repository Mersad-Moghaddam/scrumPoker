type RoomSession = {
  token: string;
  memberId: number;
};

const STORAGE_KEY = "scrum-poker-sessions";

function readAll(): Record<string, RoomSession> {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? (JSON.parse(raw) as Record<string, RoomSession>) : {};
  } catch {
    return {};
  }
}

export function saveRoomSession(roomCode: string, session: RoomSession) {
  const sessions = readAll();
  sessions[roomCode.toUpperCase()] = session;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(sessions));
}

export function getRoomSession(roomCode: string) {
  return readAll()[roomCode.toUpperCase()] ?? null;
}

export function clearRoomSession(roomCode: string) {
  const sessions = readAll();
  delete sessions[roomCode.toUpperCase()];
  localStorage.setItem(STORAGE_KEY, JSON.stringify(sessions));
}
