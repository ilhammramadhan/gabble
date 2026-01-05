export interface User {
  id: string;
  github_id: string;
  username: string;
  avatar_url: string;
  created_at: string;
}

export interface Room {
  id: string;
  name: string;
  created_by: string;
  created_at: string;
  member_count?: number;
}

export interface Message {
  id: string;
  room_id: string;
  user_id: string;
  content: string;
  created_at: string;
  user?: User;
}

export type WSEventType =
  | "join_room"
  | "leave_room"
  | "send_message"
  | "typing"
  | "message"
  | "user_joined"
  | "user_left"
  | "online_users"
  | "error";

export interface WSMessage {
  type: WSEventType;
  payload: unknown;
}

export interface JoinRoomPayload {
  room_id: string;
}

export interface SendMessagePayload {
  room_id: string;
  content: string;
}

export interface TypingPayload {
  room_id: string;
  is_typing: boolean;
}

export interface MessagePayload {
  id: string;
  room_id: string;
  content: string;
  user: User;
  created_at: string;
}

export interface UserEventPayload {
  room_id: string;
  user: User;
}

export interface TypingEventPayload {
  room_id: string;
  user: User;
  is_typing: boolean;
}

export interface OnlineUsersPayload {
  room_id: string;
  users: User[];
}

export interface ErrorPayload {
  message: string;
}
