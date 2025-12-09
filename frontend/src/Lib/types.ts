export interface Profile {
  //	token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDY4MjlkYTItYzVmMy00YzQ2LTk2NTktMGNmOTFkMjczNmE3Iiwicm9sZSI6ImNsaWVudCIsImV4cCI6MTc2NDg2MDE2OSwibmJmIjoxNzY0NzczNzY5LCJpYXQiOjE3NjQ3NzM3Njl9.VZkOyWB5sgOYjpHyUGvpPCMVpojvgqZN0eZTZM-pS64",
  token: string;
  user: {
    id: string;
    username: string;
    role: string;
    created_at: string;
    updated_at: string;
  };
  session_token: string;
}

export interface Workspace {
  id: string;
  name: string;
  slug: string;
  icon_url: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface WorkspaceMember {
  id: string;
  workspace_id: string;
  user_id: string;
  role: "admin" | "member";
  joined_at: string;
}

export interface Channel {
  id: string;
  workspace_id: string;
  name: string;
  description: string | null;
  is_private: boolean;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface Message {
  id: string;
  channel_id: string;
  user_id: string;
  content: string | null;
  message_type: "text" | "image" | "audio" | "file";
  file_url: string | null;
  file_name: string | null;
  file_size: number | null;
  reply_to: string | null;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  profile?: Profile;
}

export interface DirectMessage {
  id: string;
  workspace_id: string;
  created_at: string;
}

export interface DMMessage {
  id: string;
  dm_id: string;
  user_id: string;
  content: string | null;
  message_type: "text" | "image" | "audio" | "file";
  file_url: string | null;
  file_name: string | null;
  file_size: number | null;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  profile?: Profile;
}

export interface ResponseCreateWorkspace {
  id: string;
  name: string;
  slug: string;
  created_at: string;
  updated_at: string;
}
