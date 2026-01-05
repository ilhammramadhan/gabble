"use client";

import { useState, useEffect, useCallback } from "react";
import { Header } from "@/components/header";
import { RoomList } from "@/components/sidebar/room-list";
import { ChatRoom } from "@/components/chat/chat-room";
import { useAuth } from "@/hooks/use-auth";
import { useWebSocket } from "@/hooks/use-websocket";
import { api } from "@/lib/api";
import { Room, Message, User, MessagePayload, TypingEventPayload } from "@/lib/types";
import { MessageCircle } from "lucide-react";

export default function Home() {
  const { user, loading, logout } = useAuth();
  const [rooms, setRooms] = useState<Room[]>([]);
  const [currentRoom, setCurrentRoom] = useState<Room | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [typingUsers, setTypingUsers] = useState<User[]>([]);

  const handleMessage = useCallback((payload: MessagePayload) => {
    if (payload.room_id === currentRoom?.id) {
      setMessages((prev) => [
        ...prev,
        {
          id: payload.id,
          room_id: payload.room_id,
          user_id: payload.user.id,
          content: payload.content,
          created_at: payload.created_at,
          user: payload.user,
        },
      ]);
    }
  }, [currentRoom?.id]);

  const handleTyping = useCallback((payload: TypingEventPayload) => {
    if (payload.room_id === currentRoom?.id && payload.user.id !== user?.id) {
      setTypingUsers((prev) => {
        if (payload.is_typing) {
          if (!prev.some((u) => u.id === payload.user.id)) {
            return [...prev, payload.user];
          }
        } else {
          return prev.filter((u) => u.id !== payload.user.id);
        }
        return prev;
      });
    }
  }, [currentRoom?.id, user?.id]);

  const {
    isConnected,
    onlineUsers,
    connect,
    joinRoom,
    leaveRoom,
    sendMessage,
    sendTyping,
  } = useWebSocket({
    onMessage: handleMessage,
    onTyping: handleTyping,
  });

  // Load rooms
  useEffect(() => {
    const loadRooms = async () => {
      try {
        const data = await api.getRooms();
        setRooms(data);
      } catch (error) {
        console.error("Failed to load rooms:", error);
      }
    };
    loadRooms();
  }, []);

  // Connect WebSocket when user is authenticated
  useEffect(() => {
    if (user && !isConnected) {
      connect();
    }
  }, [user, isConnected, connect]);

  // Load messages when room changes
  useEffect(() => {
    if (currentRoom) {
      const loadMessages = async () => {
        try {
          const data = await api.getMessages(currentRoom.id);
          setMessages(data);
        } catch (error) {
          console.error("Failed to load messages:", error);
        }
      };
      loadMessages();
      joinRoom(currentRoom.id);
      setTypingUsers([]);
    }

    return () => {
      if (currentRoom) {
        leaveRoom(currentRoom.id);
      }
    };
  }, [currentRoom, joinRoom, leaveRoom]);

  const handleSelectRoom = (room: Room) => {
    if (currentRoom?.id !== room.id) {
      if (currentRoom) {
        leaveRoom(currentRoom.id);
      }
      setCurrentRoom(room);
    }
  };

  const handleCreateRoom = async (name: string) => {
    try {
      const room = await api.createRoom(name);
      setRooms((prev) => [room, ...prev]);
      setCurrentRoom(room);
    } catch (error) {
      console.error("Failed to create room:", error);
    }
  };

  const handleSendMessage = (content: string) => {
    if (currentRoom) {
      sendMessage(currentRoom.id, content);
    }
  };

  const handleTypingChange = (isTyping: boolean) => {
    if (currentRoom) {
      sendTyping(currentRoom.id, isTyping);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col bg-background">
      <Header user={user} onLogout={logout} />

      <div className="flex-1 flex overflow-hidden">
        {/* Sidebar */}
        <aside className="w-60 border-r border-border bg-card shrink-0">
          <RoomList
            rooms={rooms}
            currentRoomId={currentRoom?.id || null}
            onlineUsers={onlineUsers}
            onSelectRoom={handleSelectRoom}
            onCreateRoom={handleCreateRoom}
          />
        </aside>

        {/* Main Content */}
        <main className="flex-1 flex flex-col overflow-hidden">
          {currentRoom && user ? (
            <ChatRoom
              room={currentRoom}
              messages={messages}
              currentUser={user}
              typingUsers={typingUsers}
              onSendMessage={handleSendMessage}
              onTyping={handleTypingChange}
            />
          ) : (
            <div className="flex-1 flex flex-col items-center justify-center text-center p-8">
              <MessageCircle className="h-16 w-16 text-muted-foreground mb-4" />
              <h2 className="text-xl font-semibold mb-2">Welcome to Gabble</h2>
              <p className="text-muted-foreground max-w-md">
                {user
                  ? "Select a room from the sidebar or create a new one to start chatting."
                  : "Login with GitHub to join the conversation."}
              </p>
            </div>
          )}
        </main>
      </div>
    </div>
  );
}
