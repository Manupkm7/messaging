import { useCallback, useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { Bell, BellOff, Hash } from "lucide-react";

import { MessageInput } from "../Components/MessageInput";
import { MessageItem } from "../Components/MessageItem";
import { OnlineUsers } from "../Components/OnlineUsers";
import { ScrollArea } from "../Components/ScrollArea";
import { WorkspaceSidebar } from "../Components/WorkspaceSidebar";
import Button from "../Components/CommonButton";
import type {
  Channel,
  Message,
  Profile,
  Workspace as WorkspaceType,
} from "../Lib/types";
import { useNotifications } from "../hooks/useNotifications";

const Workspace = () => {
  const scrollRef = useRef<HTMLDivElement>(null);
  const params = useParams<{ slug: string; channelId: string }>();
  const [channel, setChannel] = useState<Channel | null>(null);
  const [workspace, setWorkspace] = useState<WorkspaceType | null>(null);
  const [channels, setChannels] = useState<Channel[]>([]);
  const [messages, setMessages] = useState<Message[]>([]);
  const [profile, setProfile] = useState<Profile | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const { permission, requestPermission } = useNotifications();

  //  const onlineUsers = useRealtimePresence(workspace?.id || "", profile?.id || "", profile?.display_name || "")
  const onlineUsers: any[] = [];

  const scrollToBottom = () => {
    if (scrollRef.current) {
      scrollRef.current.scrollIntoView({ behavior: "smooth" });
    }
  };

  const handleUpdateMessage = useCallback((message: Message) => {
    setMessages((prev) => prev.map((m) => (m.id === message.id ? message : m)));
  }, []);

  const handleDeleteMessage = useCallback((messageId: string) => {
    setMessages((prev) => prev.filter((m) => m.id !== messageId));
  }, []);

  // useRealtimeMessages(params.channelId, handleNewMessage, handleUpdateMessage, handleDeleteMessage)

  useEffect(() => {
    console.log("Workspace component mounted");
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, []);

  const handleLogout = async () => {
    console.log("Logout clicked");
  };

  const handleToggleNotifications = () => {
    if (permission === "default") {
      requestPermission();
    }
  };

  const handleSendMessage = async (content: string, type: "text") => {
    if (!profile) return;

    /*    const { error } = await supabase.from("messages").insert({
      channel_id: params.channelId,
      user_id: profile.id,
      content,
      message_type: type,
    });
    if (error) {
      console.error("[v0] Error sending message:", error);
    }
      */
    console.log("Message sent:", content, type);
  };

  const handleSendFile = async (
    file: File,
    type: "image" | "audio" | "file",
    caption?: string
  ) => {
    if (!profile || !workspace) return;

    try {
      /*      const fileExt = file.name.split(".").pop();
      const fileName = `${workspace.id}/${
        params.channelId
      }/${Date.now()}.${fileExt}`;

      const { data: uploadData, error: uploadError } = await supabase.storage
        .from("chat-files")
        .upload(fileName, file);

      if (uploadError) throw uploadError;

      const {
        data: { publicUrl },
      } = supabase.storage.from("chat-files").getPublicUrl(fileName);

      const { error: insertError } = await supabase.from("messages").insert({
        channel_id: params.channelId,
        user_id: profile.id,
        content: caption || null,
        message_type: type,
        file_url: publicUrl,
        file_name: file.name,
        file_size: file.size,
      });

      if (insertError) throw insertError;
      */
      console.log("File uploaded:", file.name, type, caption);
    } catch (error) {
      console.error("[v0] Error uploading file:", error);
      alert("Failed to upload file. Please try again.");
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  if (!channel || !workspace || !profile) {
    return (
      <div className="flex h-screen items-center justify-center">
        <p className="text-muted-foreground">Channel not found</p>
      </div>
    );
  }

  return (
    <>
      <WorkspaceSidebar
        workspace={workspace}
        channels={channels}
        profile={profile}
        onCreateChannel={() => { }}
        onLogout={handleLogout}
      />

      <div className="flex flex-1 flex-col">
        <div className="flex h-14 items-center gap-3 border-b px-6">
          <Hash className="h-5 w-5 text-muted-foreground" />
          <div className="flex-1">
            <h2 className="font-semibold">{channel.name}</h2>
            {channel.description && (
              <p className="text-xs text-muted-foreground">
                {channel.description}
              </p>
            )}
          </div>
          <Button
            icon={
              permission === "granted" ? (
                <Bell className="h-5 w-5" />
              ) : (
                <BellOff className="h-5 w-5" />
              )
            }
            onClick={handleToggleNotifications}
          >
            {permission === "granted"
              ? "Notifications enabled"
              : "Enable notifications"}
          </Button>
        </div>

        <ScrollArea className="flex-1">
          <div className="flex flex-col">
            {messages.map((message) => (
              <MessageItem
                key={message.id}
                message={message}
                isOwn={message.user_id === profile.user?.id}
              />
            ))}
            <div ref={scrollRef} />
          </div>
        </ScrollArea>

        <MessageInput
          onSendMessage={handleSendMessage}
          onSendFile={handleSendFile}
        />
      </div>

      <OnlineUsers users={onlineUsers} />
    </>
  );
};

export default Workspace;
