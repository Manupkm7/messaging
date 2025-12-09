import { Download, Play, Pause } from "lucide-react";
import { useState, useRef } from "react";
import { Avatar, AvatarFallback } from "./Avatar";
import Button from "./CommonButton";

interface MessageItemProps {
  message: Message;
  isOwn: boolean;
}

export function MessageItem({ message, isOwn }: MessageItemProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const audioRef = useRef<HTMLAudioElement>(null);

  const handlePlayPause = () => {
    if (!audioRef.current) return;

    if (isPlaying) {
      audioRef.current.pause();
    } else {
      audioRef.current.play();
    }
    setIsPlaying(!isPlaying);
  };

  const formatTime = (date: string) => {
    return new Date(date).toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <div
      className={`
        flex gap-3 px-4 py-2 hover:bg-accent/50 rounded-lg
        ${isOwn && "bg-accent/30"}
      `}
    >
      <Avatar className="h-10 w-10">
        <AvatarFallback>
          {message.profile?.display_name?.[0]?.toUpperCase() || "U"}
        </AvatarFallback>
      </Avatar>
      <div className="flex-1 space-y-1">
        <div className="flex items-baseline gap-2">
          <span className="font-semibold text-sm">
            {message.profile?.display_name || "Unknown User"}
          </span>
          <span className="text-xs text-muted-foreground">
            {formatTime(message.created_at)}
          </span>
        </div>

        {message.message_type === "text" && (
          <p className="text-sm leading-relaxed whitespace-pre-wrap wrap-break-word">
            {message.content}
          </p>
        )}

        {message.message_type === "image" && message.file_url && (
          <div className="mt-2">
            <img
              src={message.file_url || "/placeholder.svg"}
              alt="Shared image"
              className="max-w-sm rounded-lg"
            />
            {message.content && (
              <p className="mt-2 text-sm leading-relaxed">{message.content}</p>
            )}
          </div>
        )}

        {message.message_type === "audio" && message.file_url && (
          <div className="flex items-center gap-2 rounded-lg bg-secondary p-3 max-w-sm">
            <Button onClick={handlePlayPause} className="h-8 w-8">
              {isPlaying ? (
                <Pause className="h-4 w-4" />
              ) : (
                <Play className="h-4 w-4" />
              )}
            </Button>
            <audio
              ref={audioRef}
              src={message.file_url}
              onEnded={() => setIsPlaying(false)}
              className="hidden"
            />
            <div className="flex-1">
              <div className="h-1 rounded-full bg-muted">
                <div className="h-1 w-1/3 rounded-full bg-primary" />
              </div>
            </div>
            <a href={message.file_url} download>
              <Button className="h-8 w-8">
                <Download className="h-4 w-4" />
              </Button>
            </a>
          </div>
        )}

        {message.message_type === "file" && message.file_url && (
          <div className="mt-2">
            <a
              href={message.file_url}
              download={message.file_name}
              className="flex items-center gap-2 rounded-lg bg-secondary p-3 hover:bg-secondary/80 max-w-sm"
            >
              <Download className="h-5 w-5" />
              <div className="flex-1 overflow-hidden">
                <p className="truncate text-sm font-medium">
                  {message.file_name}
                </p>
                {message.file_size && (
                  <p className="text-xs text-muted-foreground">
                    {(message.file_size / 1024).toFixed(1)} KB
                  </p>
                )}
              </div>
            </a>
            {message.content && (
              <p className="mt-2 text-sm leading-relaxed">{message.content}</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
