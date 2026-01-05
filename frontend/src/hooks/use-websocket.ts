"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { getWebSocketURL } from "@/lib/api";
import {
  WSMessage,
  MessagePayload,
  UserEventPayload,
  TypingEventPayload,
  OnlineUsersPayload,
  User,
} from "@/lib/types";

interface UseWebSocketOptions {
  onMessage?: (payload: MessagePayload) => void;
  onUserJoined?: (payload: UserEventPayload) => void;
  onUserLeft?: (payload: UserEventPayload) => void;
  onTyping?: (payload: TypingEventPayload) => void;
  onOnlineUsers?: (payload: OnlineUsersPayload) => void;
  onError?: (message: string) => void;
}

export function useWebSocket(options: UseWebSocketOptions = {}) {
  const [isConnected, setIsConnected] = useState(false);
  const [onlineUsers, setOnlineUsers] = useState<User[]>([]);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const connect = useCallback(() => {
    const token = localStorage.getItem("token");
    if (!token) return;

    try {
      const ws = new WebSocket(getWebSocketURL());

      ws.onopen = () => {
        setIsConnected(true);
      };

      ws.onclose = () => {
        setIsConnected(false);
        reconnectTimeoutRef.current = setTimeout(connect, 3000);
      };

      ws.onerror = () => {
        ws.close();
      };

      ws.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data);

          switch (message.type) {
            case "message":
              options.onMessage?.(message.payload as MessagePayload);
              break;
            case "user_joined":
              options.onUserJoined?.(message.payload as UserEventPayload);
              break;
            case "user_left":
              options.onUserLeft?.(message.payload as UserEventPayload);
              break;
            case "typing":
              options.onTyping?.(message.payload as TypingEventPayload);
              break;
            case "online_users":
              const payload = message.payload as OnlineUsersPayload;
              setOnlineUsers(payload.users);
              options.onOnlineUsers?.(payload);
              break;
            case "error":
              options.onError?.((message.payload as { message: string }).message);
              break;
          }
        } catch (err) {
          console.error("Failed to parse WebSocket message:", err);
        }
      };

      wsRef.current = ws;
    } catch (err) {
      console.error("Failed to connect WebSocket:", err);
    }
  }, [options]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
    }
    wsRef.current?.close();
    wsRef.current = null;
  }, []);

  const send = useCallback((message: WSMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    }
  }, []);

  const joinRoom = useCallback(
    (roomId: string) => {
      send({ type: "join_room", payload: { room_id: roomId } });
    },
    [send]
  );

  const leaveRoom = useCallback(
    (roomId: string) => {
      send({ type: "leave_room", payload: { room_id: roomId } });
    },
    [send]
  );

  const sendMessage = useCallback(
    (roomId: string, content: string) => {
      send({
        type: "send_message",
        payload: { room_id: roomId, content },
      });
    },
    [send]
  );

  const sendTyping = useCallback(
    (roomId: string, isTyping: boolean) => {
      send({
        type: "typing",
        payload: { room_id: roomId, is_typing: isTyping },
      });
    },
    [send]
  );

  useEffect(() => {
    return () => {
      disconnect();
    };
  }, [disconnect]);

  return {
    isConnected,
    onlineUsers,
    connect,
    disconnect,
    joinRoom,
    leaveRoom,
    sendMessage,
    sendTyping,
  };
}
