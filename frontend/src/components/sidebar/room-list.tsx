"use client";

import { useState } from "react";
import { Plus, Hash, Users } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Room, User } from "@/lib/types";
import { cn } from "@/lib/utils";

interface RoomListProps {
  rooms: Room[];
  currentRoomId: string | null;
  onlineUsers: User[];
  onSelectRoom: (room: Room) => void;
  onCreateRoom: (name: string) => void;
}

export function RoomList({
  rooms,
  currentRoomId,
  onlineUsers,
  onSelectRoom,
  onCreateRoom,
}: RoomListProps) {
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [newRoomName, setNewRoomName] = useState("");

  const handleCreate = () => {
    if (newRoomName.trim()) {
      onCreateRoom(newRoomName.trim());
      setNewRoomName("");
      setIsCreateOpen(false);
    }
  };

  return (
    <div className="flex flex-col h-full">
      <div className="p-4 border-b border-border">
        <Button
          onClick={() => setIsCreateOpen(true)}
          className="w-full cursor-pointer"
          size="sm"
        >
          <Plus className="mr-2 h-4 w-4" />
          Create Room
        </Button>
      </div>

      <ScrollArea className="flex-1 custom-scrollbar">
        <div className="p-2">
          <div className="text-xs font-semibold text-muted-foreground px-2 py-1 uppercase">
            Rooms
          </div>
          {rooms.map((room) => (
            <button
              key={room.id}
              onClick={() => onSelectRoom(room)}
              className={cn(
                "w-full flex items-center gap-2 px-2 py-2 rounded-md text-sm transition-colors cursor-pointer",
                currentRoomId === room.id
                  ? "bg-primary text-primary-foreground"
                  : "hover:bg-accent text-foreground"
              )}
            >
              <Hash className="h-4 w-4 shrink-0" />
              <span className="truncate">{room.name}</span>
            </button>
          ))}
          {rooms.length === 0 && (
            <p className="text-sm text-muted-foreground px-2 py-4 text-center">
              No rooms yet. Create one!
            </p>
          )}
        </div>
      </ScrollArea>

      {onlineUsers.length > 0 && (
        <div className="p-4 border-t border-border">
          <div className="flex items-center gap-2 text-xs font-semibold text-muted-foreground uppercase mb-2">
            <Users className="h-3 w-3" />
            Online ({onlineUsers.length})
          </div>
          <div className="space-y-1">
            {onlineUsers.slice(0, 5).map((user) => (
              <div key={user.id} className="flex items-center gap-2 text-sm">
                <div className="relative">
                  <div className="w-2 h-2 rounded-full bg-green-500 online-pulse" />
                </div>
                <span className="truncate">{user.username}</span>
              </div>
            ))}
            {onlineUsers.length > 5 && (
              <p className="text-xs text-muted-foreground">
                +{onlineUsers.length - 5} more
              </p>
            )}
          </div>
        </div>
      )}

      <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
        <DialogContent className="sm:max-w-[350px]">
          <DialogHeader>
            <DialogTitle>Create Room</DialogTitle>
          </DialogHeader>
          <div className="py-4">
            <Input
              value={newRoomName}
              onChange={(e) => setNewRoomName(e.target.value)}
              placeholder="Room name"
              onKeyDown={(e) => e.key === "Enter" && handleCreate()}
              autoFocus
            />
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setIsCreateOpen(false)}
              className="cursor-pointer"
            >
              Cancel
            </Button>
            <Button
              onClick={handleCreate}
              disabled={!newRoomName.trim()}
              className="cursor-pointer"
            >
              Create
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
