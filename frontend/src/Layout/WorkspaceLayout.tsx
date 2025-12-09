import { useEffect, useMemo, useState } from "react";

import { ResponseGetAllWorkspaces } from "../Lib/types";
import { useNavigate, useParams } from "react-router-dom";
import { GetAllWorkspaces } from "../hooks/useDataStore";
import { WorkspaceSwitcher } from "@/Components/WorkspaceSwitcher";

export default function WorkspaceLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const params = useParams<{ slug: string }>();
  const router = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const [workspaces, setWorkspaces] = useState<
    ResponseGetAllWorkspaces[] | undefined
  >(undefined);

  useEffect(() => {
    GetAllWorkspaces().then((result) => {
      setWorkspaces(result.data);
    });
  }, []);

  const currentWorkspace = useMemo(() => {
    if (!workspaces) return null;
    const workspace = workspaces.find(
      (ws: ResponseGetAllWorkspaces) => ws.slug === params.slug
    );
    if (!workspace) {
      router("/create-workspace");
      return null;
    }
    setIsLoading(false);
    return workspace;
  }, [params.slug, workspaces]);

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  if (!currentWorkspace) {
    return null;
  }

  return (
    <div className="flex h-screen overflow-hidden">
      <WorkspaceSwitcher
        workspaces={workspaces}
        currentWorkspace={currentWorkspace}
      />
      {children}
    </div>
  );
}
