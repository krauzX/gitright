import { Routes, Route, Navigate } from "react-router-dom";
import { useAuthStore } from "./store/authStore";
import LandingPage from "./app/pages/LandingPage";
import { DashboardPage } from "./app/pages/DashboardPage";
import { AuthCallbackPage } from "./app/pages/AuthCallbackPage";
import { ProfileBuilderPage } from "./app/pages/ProfileBuilderPage";
import { NotFoundPage } from "./app/pages/NotFoundPage";

function App() {
  const { isAuthenticated } = useAuthStore();

  return (
    <Routes>
      <Route path="/" element={<LandingPage />} />
      <Route path="/auth/callback" element={<AuthCallbackPage />} />

      {/* Protected routes */}
      <Route
        path="/dashboard"
        element={
          isAuthenticated ? <DashboardPage /> : <Navigate to="/" replace />
        }
      />
      <Route
        path="/profile-builder"
        element={
          isAuthenticated ? <ProfileBuilderPage /> : <Navigate to="/" replace />
        }
      />

      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}

export default App;
