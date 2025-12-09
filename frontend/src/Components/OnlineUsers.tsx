import { Avatar, AvatarFallback } from "./Avatar";
import { ScrollArea } from "./ScrollArea";

interface OnlineUsersProps {
  users: OnlineUser[];
}

export function OnlineUsers({ users }: OnlineUsersProps) {
  if (users.length === 0) return null;

  return (
    <div className="flex h-full w-60 flex-col border-l bg-secondary/30">
      <div className="flex h-14 items-center px-4 border-b">
        <h3 className="font-semibold text-sm">Online ({users.length})</h3>
      </div>
      <ScrollArea className="flex-1 p-2">
        <div className="space-y-1">
          {users.map((user) => (
            <div
              key={user.presence_ref}
              className="flex items-center gap-2 rounded-md px-2 py-1.5 hover:bg-accent"
            >
              <div className="relative">
                <Avatar className="h-8 w-8">
                  <AvatarFallback className="text-xs">
                    {user.display_name[0]?.toUpperCase() || "U"}
                  </AvatarFallback>
                </Avatar>
                <div className="absolute bottom-0 right-0 h-2.5 w-2.5 rounded-full border-2 border-background bg-green-500" />
              </div>
              <span className="text-sm truncate">{user.display_name}</span>
            </div>
          ))}
        </div>
      </ScrollArea>
    </div>
  );
}
