import { useMutation } from "@tanstack/react-query";

import api from "../Lib/api";
import { Profile, ResponseCreateWorkspace } from "../Lib/types";

export const useAuthStore = () => {
  const login = useMutation({
    mutationFn: async ({
      username,
      password,
      session_token,
    }: {
      username: string;
      password: string;
      session_token?: string;
    }) => {
      const response = await api.post<Profile>("/auth/login", {
        username,
        password,
        session_token,
      });
      return response.data;
    },
    onSuccess: (data) => {
      localStorage.setItem("profile", JSON.stringify(data));
      localStorage.setItem("session_token", data.session_token);
      localStorage.setItem("token", data.token);
    },
  });

  return {
    login: login.mutateAsync,
  };
};

/* ### Workspaces */

export const useWorkspaces = () => {
  const createWorkspace = useMutation({
    mutationFn: async (name: string) => {
      const response = await api.post<ResponseCreateWorkspace>("/workspaces", {
        name,
      });
      return response.data;
    },
  });

  const updateWorkspace = useMutation({
    mutationFn: async ({
      workspaceId,
      name,
    }: {
      workspaceId: string;
      name: string;
    }) => {
      const response = await api.put(`/workspaces/${workspaceId}`, {
        name,
      });
      return response.data;
    },
  });

  return {
    createWorkspace: createWorkspace.mutateAsync,
    updateWorkspace: updateWorkspace.mutateAsync,
  };
};


