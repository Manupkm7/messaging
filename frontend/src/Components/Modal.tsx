import { ReactNode } from "react";

import {
  Dialog,
  DialogPanel,
  Transition,
  TransitionChild,
} from "@headlessui/react";

export default function Modal({
  children,
  show = false,
  maxWidth = "2xl",
  closeable = true,
  onClose = () => {},
  className = "",
}: {
  children: ReactNode;
  show: boolean;
  maxWidth?: string;
  closeable?: boolean;
  onClose?: () => void;
  className?: string;
}) {
  const close = () => {
    if (closeable) {
      onClose?.();
    }
  };

  const maxWidthClass = {
    sm: "sm:max-w-sm",
    md: "sm:max-w-md",
    lg: "sm:max-w-lg",
    xl: "sm:max-w-xl",
    "2xl": "sm:max-w-2xl",
    "3xl": "sm:max-w-3xl",
    "4xl": "sm:max-w-4xl",
    "5xl": "sm:max-w-5xl",
    "6xl": "sm:max-w-6xl",
  }[maxWidth];

  return (
    <Transition show={show} leave="duration-200">
      <Dialog
        as="div"
        id="modal"
        className=" fixed inset-0 z-50 transform px-4 py-6 transition-all sm:px-0 flex items-center justify-center"
        onClose={close}
        aria-modal="true"
        role="dialog"
        data-testid="modal"
      >
        <TransitionChild
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div
            onClick={close}
            className="absolute inset-0 bg-shark-500/25 dark:bg-shark-900/25"
          />
        </TransitionChild>

        <TransitionChild
          enter="ease-out duration-300"
          enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          enterTo="opacity-100 translate-y-0 sm:scale-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100 translate-y-0 sm:scale-100"
          leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
        >
          <DialogPanel
            className={`mb-6 transform max-h-screen overflow-y-auto overflow-x-hidden rounded-lg bg-white shadow-xl transition-all dark:bg-shark-800 sm:mx-auto sm:w-full ${maxWidthClass} ${className}`}
            style={{ maxHeight: "100vh" }}
            data-testid="modal-panel"
          >
            {children}
          </DialogPanel>
        </TransitionChild>
      </Dialog>
    </Transition>
  );
}
