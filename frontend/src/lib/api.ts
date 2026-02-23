import axios, { AxiosInstance, InternalAxiosRequestConfig } from "axios";
import { useAuthStore } from "@store/authStore";

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:5173";

// Create axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  timeout: 300000, // 5 minutes for profile generation
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor - add auth token
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
  }
);

// API helper functions
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
};

export default apiClient;
export { apiClient as api };
