export const API_CONFIGS = {
  development: {
    BASE_URL: import.meta.env.VITE_API_URL,
    WEBSOCKET_URL: import.meta.env.VITE_WEBSOCKET_URL || "",
    TIMEOUT: 30000,
    RETRY_ATTEMPTS: 3,
  },

  production: {
    BASE_URL: import.meta.env.VITE_API_URL,
    WEBSOCKET_URL: import.meta.env.VITE_WEBSOCKET_URL || "wss://",
    TIMEOUT: 30000,
    RETRY_ATTEMPTS: 3,
  },

  local: {
    BASE_URL: "http://localhost:8080/api",
    WEBSOCKET_URL: "ws://localhost:20001/go/ws",
    TIMEOUT: 30000,
    RETRY_ATTEMPTS: 3,
  },
};

export const getApiConfig = () => {
  console.log("import.meta.env.MODE", import.meta.env.MODE);
  console.log("API_CONFIGS", API_CONFIGS);
  const env = import.meta.env.MODE || "development";
  return (
    API_CONFIGS[env as keyof typeof API_CONFIGS] || API_CONFIGS.development
  );
};
