import { useEffect, useState, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { motion } from "framer-motion";
import { Loader2, CheckCircle, XCircle } from "lucide-react";
import { useAuthStore } from "../../store/authStore";
import { api } from "../../lib/api";

export function AuthCallbackPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState<"loading" | "success" | "error">(
    "loading"
  );
  const [errorMessage, setErrorMessage] = useState("");
  const setUser = useAuthStore((state) => state.setUser);
  const setToken = useAuthStore((state) => state.setToken);
  const hasCalledRef = useRef(false);

  useEffect(() => {
    // Prevent duplicate calls in React StrictMode
    if (hasCalledRef.current) return;
    hasCalledRef.current = true;

    handleCallback();
  }, []);

  const handleCallback = async () => {
    const code = searchParams.get("code");
    const state = searchParams.get("state");
    const error = searchParams.get("error");

    // Check for OAuth error
    if (error) {
      setStatus("error");
      setErrorMessage(`Authentication failed: ${error}`);
      setTimeout(() => navigate("/"), 3000);
      return;
    }

    // Validate required parameters
    if (!code || !state) {
      setStatus("error");
      setErrorMessage("Missing authorization code or state");
      setTimeout(() => navigate("/"), 3000);
      return;
    }

    try {
      // Exchange code for token
      // Backend will validate the state parameter
      const response = await api.get(
        `/auth/callback?code=${code}&state=${state}`
      );

      const { user, token } = response.data;

      // Store user and token
      setUser(user);
      setToken(token);

      // Clear stored state
      sessionStorage.removeItem("oauth_state");

      setStatus("success");

      // Redirect to dashboard after short delay
      setTimeout(() => {
        navigate("/dashboard");
      }, 1500);
    } catch (error: any) {
      console.error("Authentication failed:", error);
      setStatus("error");
      setErrorMessage(
        error.response?.data?.message ||
          "Failed to complete authentication. Please try again."
      );
      setTimeout(() => navigate("/"), 3000);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950 flex items-center justify-center px-4 relative overflow-hidden">
      {/* Animated background effects */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-20 left-10 w-72 h-72 bg-purple-500/20 rounded-full blur-3xl animate-pulse"></div>
        <div
          className="absolute bottom-20 right-20 w-96 h-96 bg-blue-500/20 rounded-full blur-3xl animate-pulse"
          style={{ animationDelay: "1s" }}
        ></div>
      </div>

      <motion.div
        initial={{ opacity: 0, scale: 0.9, y: 20 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        transition={{ duration: 0.5, ease: "easeOut" }}
        className="relative bg-gradient-to-br from-slate-900/90 via-slate-800/90 to-slate-900/90 backdrop-blur-xl rounded-3xl p-10 border border-slate-700/50 shadow-2xl max-w-md w-full text-center"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-purple-500/5 rounded-3xl"></div>

        <div className="relative">
          {status === "loading" && (
            <>
              <motion.div
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                className="mb-6"
              >
                <Loader2 className="w-20 h-20 text-primary mx-auto drop-shadow-lg" />
              </motion.div>
              <h2 className="text-3xl font-bold text-white mb-3 bg-gradient-to-r from-white to-slate-300 bg-clip-text text-transparent">
                Authenticating...
              </h2>
              <p className="text-slate-400 text-lg">
                Please wait while we complete your sign in
              </p>
              <div className="mt-6 flex justify-center gap-2">
                <div className="w-2 h-2 bg-primary rounded-full animate-bounce"></div>
                <div
                  className="w-2 h-2 bg-primary rounded-full animate-bounce"
                  style={{ animationDelay: "0.1s" }}
                ></div>
                <div
                  className="w-2 h-2 bg-primary rounded-full animate-bounce"
                  style={{ animationDelay: "0.2s" }}
                ></div>
              </div>
            </>
          )}

          {status === "success" && (
            <>
              <motion.div
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ type: "spring", stiffness: 200, damping: 10 }}
                className="mb-6"
              >
                <div className="relative inline-block">
                  <div className="absolute inset-0 bg-green-400/30 rounded-full blur-xl animate-pulse"></div>
                  <CheckCircle className="relative w-20 h-20 text-green-400 mx-auto drop-shadow-lg" />
                </div>
              </motion.div>
              <h2 className="text-3xl font-bold text-white mb-3 bg-gradient-to-r from-green-400 to-emerald-400 bg-clip-text text-transparent">
                Success!
              </h2>
              <p className="text-slate-400 text-lg">
                You've been successfully authenticated. Redirecting to
                dashboard...
              </p>
            </>
          )}

          {status === "error" && (
            <>
              <motion.div
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ type: "spring", stiffness: 200, damping: 10 }}
                className="mb-6"
              >
                <div className="relative inline-block">
                  <div className="absolute inset-0 bg-red-400/30 rounded-full blur-xl animate-pulse"></div>
                  <XCircle className="relative w-20 h-20 text-red-400 mx-auto drop-shadow-lg" />
                </div>
              </motion.div>
              <h2 className="text-3xl font-bold text-white mb-3 bg-gradient-to-r from-red-400 to-rose-400 bg-clip-text text-transparent">
                Authentication Failed
              </h2>
              <p className="text-slate-300 text-base mb-6 leading-relaxed px-4">
                {errorMessage}
              </p>
              <button
                onClick={() => navigate("/")}
                className="px-8 py-3 rounded-xl bg-gradient-to-r from-primary via-purple-500 to-pink-500 text-white font-semibold hover:shadow-lg hover:shadow-primary/50 transition-all duration-300 hover:scale-105"
              >
                Return to Home
              </button>
            </>
          )}
        </div>
      </motion.div>
    </div>
  );
}
