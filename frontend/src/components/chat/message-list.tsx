"use client";

import { useEffect, useRef } from "react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { MessageBubble } from "./message-bubble";
import { TypingIndicator } from "./typing-indicator";
import { Message, User } from "@/lib/types";

interface MessageListProps {
  messages: Message[];
  currentUser: User;
  typingUsers: User[];
}

export function MessageList({
  messages,
  currentUser,
  typingUsers,
}: MessageListProps) {
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, typingUsers]);

  return (
    <ScrollArea className="flex-1 custom-scrollbar">
      <div className="py-4">
        {messages.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full py-20 text-center">
            <p className="text-muted-foreground">No messages yet</p>
            <p className="text-sm text-muted-foreground">
              Be the first to say something!
            </p>
          </div>
        ) : (
          messages.map((message, index) => (
            <MessageBubble
              key={message.id}
              message={message}
              currentUser={currentUser}
              isNew={index === messages.length - 1}
            />
          ))
        )}

        <TypingIndicator users={typingUsers} />
        <div ref={bottomRef} />
      </div>
    </ScrollArea>
  );
}
