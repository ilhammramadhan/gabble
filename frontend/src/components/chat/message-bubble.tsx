"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Message, User } from "@/lib/types";
import { cn } from "@/lib/utils";

interface MessageBubbleProps {
  message: Message;
  currentUser: User;
  isNew?: boolean;
}

export function MessageBubble({ message, currentUser, isNew }: MessageBubbleProps) {
  const isSent = message.user?.id === currentUser.id;
  const user = message.user;

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  };

  return (
    <div
      className={cn(
        "flex gap-3 px-4 py-1",
        isSent ? "flex-row-reverse" : "flex-row",
        isNew && "message-enter"
      )}
    >
      {!isSent && (
        <Avatar className="h-8 w-8 shrink-0">
          <AvatarImage src={user?.avatar_url} alt={user?.username} />
          <AvatarFallback>
            {user?.username?.[0]?.toUpperCase() || "?"}
          </AvatarFallback>
        </Avatar>
      )}

      <div
        className={cn(
          "flex flex-col max-w-[70%]",
          isSent ? "items-end" : "items-start"
        )}
      >
        {!isSent && (
          <span className="text-xs text-muted-foreground mb-1">
            {user?.username}
          </span>
        )}

        <div
          className={cn(
            "px-4 py-2 text-sm",
            isSent
              ? "bg-primary text-primary-foreground rounded-[16px_16px_4px_16px]"
              : "bg-secondary text-secondary-foreground rounded-[16px_16px_16px_4px]"
          )}
        >
          {message.content}
        </div>

        <span className="text-xs text-muted-foreground mt-1">
          {formatTime(message.created_at)}
        </span>
      </div>
    </div>
  );
}
