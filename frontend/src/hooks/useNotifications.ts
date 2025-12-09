import { useState } from "react";

export function useNotifications() {
  const [permission, setPermission] = useState<NotificationPermission>(
    typeof window !== "undefined" && "Notification" in window
      ? Notification.permission
      : "default"
  );

  const requestPermission = async () => {
    if (typeof window !== "undefined" && "Notification" in window) {
      const result = await Notification.requestPermission();
      setPermission(result);
      return result;
    }
    return "denied";
  };

  const showNotification = (title: string, options?: NotificationOptions) => {
    if (permission === "granted" && "Notification" in window) {
      new Notification(title, {
        icon: "/icon.svg",
        badge: "/icon.svg",
        ...options,
      });
    }
  };

  return { permission, requestPermission, showNotification };
}
