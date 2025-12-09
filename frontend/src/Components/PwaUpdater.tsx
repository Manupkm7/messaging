import { useEffect, useState } from "react";
import { RefreshCw } from "lucide-react";
import { Card } from "./Card";
import Button from "./CommonButton";

export function PWAUpdatePrompt() {
  const [showPrompt, setShowPrompt] = useState(false);
  const [waitingWorker, setWaitingWorker] = useState<ServiceWorker | null>(
    null
  );

  useEffect(() => {
    if ("serviceWorker" in navigator) {
      navigator.serviceWorker.ready.then((registration) => {
        registration.addEventListener("updatefound", () => {
          const newWorker = registration.installing;
          if (newWorker) {
            newWorker.addEventListener("statechange", () => {
              if (
                newWorker.state === "installed" &&
                navigator.serviceWorker.controller
              ) {
                setWaitingWorker(newWorker);
                setShowPrompt(true);
              }
            });
          }
        });
      });

      // Check for updates every hour
      const interval = setInterval(() => {
        navigator.serviceWorker.ready.then((registration) => {
          registration.update();
        });
      }, 60 * 60 * 1000);

      return () => clearInterval(interval);
    }
  }, []);

  const handleUpdate = () => {
    if (waitingWorker) {
      waitingWorker.postMessage({ type: "SKIP_WAITING" });
      setShowPrompt(false);
      window.location.reload();
    }
  };

  if (!showPrompt) return null;

  return (
    <div className="fixed bottom-4 left-1/2 z-50 w-full max-w-sm -translate-x-1/2 px-4">
      <Card className="border-primary shadow-lg">
        <div className="pb-3">
          <h2 className="text-lg">Update Available</h2>
          <span>A new version of ChatSpace is available</span>
        </div>
        <div className="pt-0">
          <Button onClick={handleUpdate} className="w-full">
            <RefreshCw className="mr-2 h-4 w-4" />
            Update Now
          </Button>
        </div>
      </Card>
    </div>
  );
}
