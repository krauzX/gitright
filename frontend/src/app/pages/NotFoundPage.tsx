import { Link } from "react-router-dom";
import { Home, GitBranch } from "lucide-react";

export function NotFoundPage() {
  return (
    <div className="min-h-screen bg-[#030712] flex items-center justify-center px-4">
      <div className="text-center max-w-md">
        <div className="text-[120px] font-bold leading-none bg-gradient-to-r from-violet-400 to-cyan-400 bg-clip-text text-transparent mb-4">
          404
        </div>
        <h1 className="text-2xl font-semibold text-white mb-2">Page Not Found</h1>
        <p className="text-sm text-gray-400 mb-8">
          The page you're looking for doesn't exist or has been moved.
        </p>
        <div className="flex gap-3 justify-center">
          <Link
            to="/"
            className="inline-flex items-center gap-2 px-6 py-2.5 rounded-lg bg-gradient-to-r from-violet-600 to-cyan-600 text-sm font-medium text-white"
          >
            <Home className="w-4 h-4" />
            Go Home
          </Link>
          <Link
            to="/dashboard"
            className="inline-flex items-center gap-2 px-6 py-2.5 rounded-lg bg-white/5 border border-white/10 text-sm font-medium text-gray-300 hover:bg-white/10"
          >
            <GitBranch className="w-4 h-4" />
            Dashboard
          </Link>
        </div>
      </div>
    </div>
  );
}
