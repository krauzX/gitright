import { useEffect, useState, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Loader2, CheckCircle, XCircle } from "lucide-react";
import { useAuthStore } from "../../store/authStore";
import { api } from "../../lib/api";

export function AuthCallbackPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState<"loading" | "success" | "error">("loading");
  const [errorMessage, setErrorMessage] = useState("");
  const setUser = useAuthStore((state) => state.setUser);
  const setToken = useAuthStore((state) => state.setToken);
  const hasCalledRef = useRef(false);

  useEffect(() => {
    if (hasCalledRef.current) return;
    hasCalledRef.current = true;
    handleCallback();
  }, []);

  const handleCallback = async () => {
    const code = searchParams.get("code");
    const state = searchParams.get("state");
    const error = searchParams.get("error");

    if (error) {
      setStatus("error");
      setErrorMessage(`Authentication failed: ${error}`);
      setTimeout(() => navigate("/"), 3000);
      return;
    }

    if (!code || !state) {
      setStatus("error");
      setErrorMessage("Missing authorization code or state");
      setTimeout(() => navigate("/"), 3000);
      return;
    }

    try {
      const response = await api.get(`/auth/callback?code=${code}&state=${state}`);
      const { user, token } = response.data;
      setUser(user);
      setToken(token);
      sessionStorage.removeItem("oauth_state");
      setStatus("success");
      setTimeout(() => navigate("/dashboard"), 1500);
    } catch (error: any) {
      setStatus("error");
      setErrorMessage(error.response?.data?.message || "Failed to complete authentication.");
      setTimeout(() => navigate("/"), 3000);
    }
  };

  return (
    <div className="min-h-screen bg-[#030712] flex items-center justify-center px-4">
      <div className="bg-white/[0.02] border border-white/[0.06] rounded-xl p-10 max-w-sm w-full text-center">
        {status === "loading" && (
          <>
            <Loader2 className="w-16 h-16 text-violet-400 mx-auto mb-4 animate-spin" />
            <h2 className="text-xl font-semibold text-white mb-2">Authenticating...</h2>
            <p className="text-sm text-gray-400">Please wait while we complete your sign in</p>
          </>
        )}

        {status === "success" && (
          <>
            <CheckCircle className="w-16 h-16 text-emerald-400 mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-white mb-2">Success!</h2>
            <p className="text-sm text-gray-400">Redirecting to dashboard...</p>
          </>
        )}

        {status === "error" && (
          <>
            <XCircle className="w-16 h-16 text-red-400 mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-white mb-2">Authentication Failed</h2>
            <p className="text-sm text-gray-400 mb-4">{errorMessage}</p>
            <button
              onClick={() => navigate("/")}
              className="px-6 py-2 rounded-lg bg-gradient-to-r from-violet-600 to-cyan-600 text-sm font-medium text-white"
            >
              Return Home
            </button>
          </>
        )}
      </div>
    </div>
  );
}
