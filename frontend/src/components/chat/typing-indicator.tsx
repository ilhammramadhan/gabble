"use client";

import { User } from "@/lib/types";

interface TypingIndicatorProps {
  users: User[];
}

export function TypingIndicator({ users }: TypingIndicatorProps) {
  if (users.length === 0) return null;

  const names =
    users.length === 1
      ? users[0].username
      : users.length === 2
      ? `${users[0].username} and ${users[1].username}`
      : `${users[0].username} and ${users.length - 1} others`;

  return (
    <div className="flex items-center gap-2 px-4 py-2 text-sm text-muted-foreground">
      <div className="flex gap-1">
        <div className="w-2 h-2 rounded-full bg-muted-foreground typing-dot" />
        <div className="w-2 h-2 rounded-full bg-muted-foreground typing-dot" />
        <div className="w-2 h-2 rounded-full bg-muted-foreground typing-dot" />
      </div>
      <span>{names} {users.length === 1 ? "is" : "are"} typing...</span>
    </div>
  );
}
