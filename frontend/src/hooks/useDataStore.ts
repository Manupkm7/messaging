import { useMutation } from "@tanstack/react-query";

import api from "../Lib/api";
import { Profile } from "../Lib/types";

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
  });

  return {
    login: login.mutateAsync,
  };
};
