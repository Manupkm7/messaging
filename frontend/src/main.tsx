import React from "react";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Toaster } from "react-hot-toast";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

import { Provider } from "jotai";
import { createRoot } from "react-dom/client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import "./index.css";
import LoginPage from "./Pages/Login";
import { ThemeProvider } from "./Context/ThemeContext";
import CreateWorkspacePage from "./Pages/CreateWorkspace";
import Workspace from "./Pages/Workspace";

export function AppRoutes() {
  return (
    <Routes>
      <Route path="*" element={<LoginPage />} />
      <Route path="/create-workspace" element={<CreateWorkspacePage />} />
      <Route path="/workspace/:slug/:channelId" element={<Workspace />} />
    </Routes>
  );
}

export function App() {
  return (
    <Provider>
      <Router>
        <ThemeProvider>
          <AppRoutes />
          <Toaster
            toastOptions={{
              duration: 5000,
            }}
          />
        </ThemeProvider>
      </Router>
    </Provider>
  );
}

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: true,
    },
  },
});

// Montar la aplicación
const container = document.getElementById("root");
if (container) {
  const root = createRoot(container);
  root.render(
    <React.StrictMode>
      <QueryClientProvider client={queryClient}>
        <App />
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </React.StrictMode>
  );
}
