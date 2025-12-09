import { Paperclip, Loader2 } from "lucide-react";
import { ChangeEvent, useRef, useState } from "react";
import Button from "./CommonButton";

interface FileUploadButtonProps {
  onFileSelect: (file: File, type: "image" | "audio" | "file") => Promise<void>;
  disabled?: boolean;
}

export function FileUploadButton({
  onFileSelect,
  disabled,
}: FileUploadButtonProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isUploading, setIsUploading] = useState(false);

  const handleFileChange = async (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setIsUploading(true);
    try {
      let type: "image" | "audio" | "file" = "file";
      if (file.type.startsWith("image/")) {
        type = "image";
      } else if (file.type.startsWith("audio/")) {
        type = "audio";
      }

      await onFileSelect(file, type);
    } finally {
      setIsUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  return (
    <>
      <input
        ref={fileInputRef}
        type="file"
        className="hidden"
        onChange={handleFileChange}
        accept="image/*,audio/*,.pdf,.doc,.docx,.txt"
      />
      <Button
        type="button"
        onClick={() => fileInputRef.current?.click()}
        disabled={disabled || isUploading}
      >
        {isUploading ? (
          <Loader2 className="h-5 w-5 animate-spin" />
        ) : (
          <Paperclip className="h-5 w-5" />
        )}
      </Button>
    </>
  );
}
