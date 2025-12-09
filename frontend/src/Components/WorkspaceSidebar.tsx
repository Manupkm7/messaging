"use client";

import Button from "../Components/CommonButton";
import { ScrollArea } from "../Components/ScrollArea";
import { Separator } from "../Components/Separator";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../Components/DropdownMenu";
import { Avatar, AvatarFallback } from "../Components/Avatar";
import {
  Hash,
  Lock,
  Plus,
  ChevronDown,
  Settings,
  LogOut,
  Users,
  Shield,
} from "lucide-react";

interface WorkspaceSidebarProps {
  workspace: Workspace;
  channels: Channel[];
  profile: Profile;
  onCreateChannel: () => void;
  onLogout: () => void;
}

export function WorkspaceSidebar({
  workspace,
  channels,
  profile,
  onCreateChannel,
  onLogout,
}: WorkspaceSidebarProps) {
  const pathname = window.location.pathname;

  const publicChannels = channels.filter((c) => !c.is_private);
  const privateChannels = channels.filter((c) => c.is_private);

  const showAdminOptions = profile.is_admin;

  return (
    <div className="flex h-full w-60 flex-col bg-secondary">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button className="h-12 w-full justify-between px-4 font-semibold hover:bg-accent">
            <span className="truncate">{workspace.name}</span>
            <ChevronDown className="h-4 w-4 opacity-50" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start" className="w-56">
          <DropdownMenuLabel>Workspace</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem>
            <Button href={`/app/workspace/${workspace.slug}/settings`}>
              <Settings className="mr-2 h-4 w-4" />
              Settings
            </Button>
          </DropdownMenuItem>
          {showAdminOptions && (
            <DropdownMenuItem>
              <Button href={`/app/workspace/${workspace.slug}/members`}>
                <Users className="mr-2 h-4 w-4" />
                Manage Members
              </Button>
            </DropdownMenuItem>
          )}
          <DropdownMenuSeparator />
          <DropdownMenuItem onClick={onLogout}>
            <LogOut className="mr-2 h-4 w-4" />
            Log out
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <Separator />

      <ScrollArea className="flex-1 px-2 py-4">
        <div className="space-y-4">
          <div>
            <div className="mb-2 flex items-center justify-between px-2">
              <span className="text-xs font-semibold uppercase text-muted-foreground">
                Channels
              </span>
              <Button className="h-5 w-5" onClick={onCreateChannel}>
                <Plus className="h-4 w-4" />
              </Button>
            </div>
            <div className="space-y-0.5">
              {publicChannels.map((channel) => {
                const isActive = pathname.includes(`/channel/${channel.id}`);
                return (
                  <Button
                    key={channel.id}
                    href={`/app/workspace/${workspace.slug}/channel/${channel.id}`}
                    className={`w-full justify-start px-2 ${
                      isActive && "bg-accent text-accent-foreground"
                    }`}
                  >
                    <Hash className="mr-2 h-4 w-4" />
                    <span className="truncate">{channel.name}</span>
                  </Button>
                );
              })}
            </div>
          </div>

          {privateChannels.length > 0 && (
            <div>
              <div className="mb-2 px-2">
                <span className="text-xs font-semibold uppercase text-muted-foreground">
                  Private
                </span>
              </div>
              <div className="space-y-0.5">
                {privateChannels.map((channel) => {
                  const isActive = pathname.includes(`/channel/${channel.id}`);
                  return (
                    <Button
                      key={channel.id}
                      className={`
                          w-full justify-start px-2
                          ${isActive && "bg-accent text-accent-foreground"}
                        `}
                      href={`/app/workspace/${workspace.slug}/channel/${channel.id}`}
                    >
                      <Lock className="mr-2 h-4 w-4" />
                      <span className="truncate">{channel.name}</span>
                    </Button>
                  );
                })}
              </div>
            </div>
          )}
        </div>
      </ScrollArea>

      <Separator />

      <div className="p-2">
        <div className="flex items-center gap-2 rounded-md p-2 hover:bg-accent">
          <Avatar className="h-8 w-8">
            <AvatarFallback>
              {profile.display_name?.[0]?.toUpperCase() || "U"}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1 overflow-hidden">
            <div className="flex items-center gap-1">
              <p className="truncate text-sm font-medium">
                {profile.display_name || "User"}
              </p>
              {profile.is_admin && <Shield className="h-3 w-3 text-primary" />}
            </div>
            <p className="truncate text-xs text-muted-foreground">
              {profile.status}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
