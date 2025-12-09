import { Plus } from "lucide-react";
import Button from "./CommonButton";
import { Workspace } from "../Lib/types";
import { Avatar, AvatarFallback } from "./Avatar";

interface WorkspaceSwitcherProps {
  workspaces: Workspace[];
  currentWorkspace: Workspace;
}

export function WorkspaceSwitcher({
  workspaces,
  currentWorkspace,
}: WorkspaceSwitcherProps) {
  const pathname = window.location.pathname;
  console.debug("Current pathname:", currentWorkspace);
  
  return (
    <div className="flex h-full w-16 flex-col items-center gap-2 bg-background py-3">
      <div className="w-full flex-1 overflow-auto">
        <div className="flex flex-col items-center gap-2 px-2">
          {workspaces.map((workspace) => {
            const isActive = pathname.includes(`/workspace/${workspace.slug}`);
            return (
              <a key={workspace.id} href={`/app/workspace/${workspace.slug}`}>
                <Button
                  className={`
                    h-12 w-12 rounded-lg
                    ${isActive ? "bg-primary text-primary-foreground" : ""}
                  `}
                >
                  <Avatar className="h-10 w-10">
                    <AvatarFallback>
                      {workspace.name[0]?.toUpperCase() || "W"}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </a>
            );
          })}
        </div>
      </div>
      <Button
        className="h-12 w-12 rounded-lg border-2 border-dashed border-muted-foreground/25 hover:border-primary"
        href="/app/create-workspace"
      >
        Crear espacio de trabajo
        <Plus className="h-6 w-6" />
      </Button>
    </div>
  );
}
