import axios, { AxiosInstance, InternalAxiosRequestConfig } from "axios";
import { useAuthStore } from "@store/authStore";

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

const apiClient: AxiosInstance = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  timeout: 300000,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = useAuthStore.getState().token;
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

export const profileAPI = {
  generate: async (data: import("@/types").ContentGenerationRequest) => {
    const response = await apiClient.post<
      import("@/types").ContentGenerationResponse
    >("/profile/generate", data);
    return response.data;
  },
  deploy: async (data: import("@/types").ContentGenerationRequest) => {
    const response = await apiClient.post("/profile/deploy", data);
    return response.data;
  },
  autoImport: async (geminiApiKey: string) => {
    const response = await apiClient.post<import("@/types").AutoImportResponse>(
      "/profile/auto-import",
      { gemini_api_key: geminiApiKey }
    );
    return response.data;
  },
  getBanner: async (theme: "dark" | "light" = "dark") => {
    const response = await apiClient.get("/profile/banner", {
      params: { theme },
      responseType: "text",
    });
    return response.data as string;
  },
  getTimeline: async (theme: "dark" | "light" = "dark") => {
    const response = await apiClient.get("/profile/timeline", {
      params: { theme },
      responseType: "text",
    });
    return response.data as string;
  },
};

export { apiClient as api };
