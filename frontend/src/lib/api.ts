import { Room, Message, User } from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function fetchAPI<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(token && { Authorization: `Bearer ${token}` }),
    ...options?.headers,
  };

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }

  return response.json();
}

export const api = {
  // Auth
  getCurrentUser: () => fetchAPI<User>("/api/auth/me"),

  // Rooms
  getRooms: () => fetchAPI<Room[]>("/api/rooms"),
  createRoom: (name: string) =>
    fetchAPI<Room>("/api/rooms", {
      method: "POST",
      body: JSON.stringify({ name }),
    }),
  getRoom: (id: string) => fetchAPI<Room>(`/api/rooms/${id}`),
  deleteRoom: (id: string) =>
    fetchAPI<void>(`/api/rooms/${id}`, { method: "DELETE" }),

  // Messages
  getMessages: (roomId: string) =>
    fetchAPI<Message[]>(`/api/rooms/${roomId}/messages`),
};

export function getWebSocketURL(): string {
  const wsProtocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const apiHost = API_URL.replace(/^https?:\/\//, "");
  const token = localStorage.getItem("token");
  return `${wsProtocol}//${apiHost}/ws?token=${token}`;
}

export function getGithubAuthURL(): string {
  return `${API_URL}/auth/github`;
}
