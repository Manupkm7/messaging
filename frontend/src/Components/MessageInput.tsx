import { Send } from "lucide-react";
import { useState, useRef } from "react";
import { AudioRecorder } from "./AudioRecorder";
import { Input } from "./Input";
import { EmojiPicker } from "./EmojiPicker";
import Button from "./CommonButton";
import { FileUploadButton } from "./FileUpload";

interface MessageInputProps {
  onSendMessage: (content: string, type: "text") => Promise<void>;
  onSendFile: (
    file: File,
    type: "image" | "audio" | "file",
    caption?: string
  ) => Promise<void>;
  disabled?: boolean;
}

export function MessageInput({
  onSendMessage,
  onSendFile,
  disabled,
}: MessageInputProps) {
  const [message, setMessage] = useState("");
  const [isSending, setIsSending] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!message.trim() || isSending) return;

    setIsSending(true);
    try {
      await onSendMessage(message, "text");
      setMessage("");
    } finally {
      setIsSending(false);
    }
  };

  const handleEmojiSelect = (emoji: string) => {
    setMessage((prev) => prev + emoji);
    inputRef.current?.focus();
  };

  const handleFileSelect = async (
    file: File,
    type: "image" | "audio" | "file"
  ) => {
    await onSendFile(file, type);
  };

  const handleRecordingComplete = async (audioBlob: Blob) => {
    const file = new File([audioBlob], `audio-${Date.now()}.webm`, {
      type: "audio/webm",
    });
    await onSendFile(file, "audio");
  };

  return (
    <form onSubmit={handleSubmit} className="border-t bg-background p-4">
      <div className="flex items-end gap-2">
        <FileUploadButton
          onFileSelect={handleFileSelect}
          disabled={disabled || isSending}
        />
        <AudioRecorder
          onRecordingComplete={handleRecordingComplete}
          disabled={disabled || isSending}
        />

        <div className="flex-1 flex items-end gap-2">
          <Input
            ref={inputRef}
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Type a message..."
            disabled={disabled || isSending}
            className="flex-1"
          />
          <EmojiPicker onEmojiSelect={handleEmojiSelect} />
        </div>

        <Button
          type="submit"
          disabled={!message.trim() || disabled || isSending}
        >
          <Send className="h-5 w-5" />
        </Button>
      </div>
    </form>
  );
}
