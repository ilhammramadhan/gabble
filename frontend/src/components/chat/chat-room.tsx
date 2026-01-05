"use client";

import { Hash } from "lucide-react";
import { MessageList } from "./message-list";
import { MessageInput } from "./message-input";
import { Room, Message, User } from "@/lib/types";

interface ChatRoomProps {
  room: Room;
  messages: Message[];
  currentUser: User;
  typingUsers: User[];
  onSendMessage: (content: string) => void;
  onTyping: (isTyping: boolean) => void;
}

export function ChatRoom({
  room,
  messages,
  currentUser,
  typingUsers,
  onSendMessage,
  onTyping,
}: ChatRoomProps) {
  return (
    <div className="flex flex-col h-full">
      {/* Room Header */}
      <div className="h-14 border-b border-border bg-card px-4 flex items-center gap-2 shrink-0">
        <Hash className="h-5 w-5 text-muted-foreground" />
        <span className="font-semibold">{room.name}</span>
      </div>

      {/* Messages */}
      <MessageList
        messages={messages}
        currentUser={currentUser}
        typingUsers={typingUsers}
      />

      {/* Input */}
      <MessageInput onSend={onSendMessage} onTyping={onTyping} />
    </div>
  );
}
