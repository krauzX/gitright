import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import {
  CheckCircle,
  Circle,
  RefreshCw,
  GitBranch,
  Star,
  GitFork,
} from "lucide-react";
import { useAuthStore } from "../../store/authStore";
import { api } from "../../lib/api";
import type { Repository } from "../../types";

export function DashboardPage() {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [selectedRepos, setSelectedRepos] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(true);
  const [analyzing, setAnalyzing] = useState(false);
  const [includePrivate, setIncludePrivate] = useState(false);

  useEffect(() => {
    fetchRepositories();
  }, [includePrivate]);

  const fetchRepositories = async () => {
    try {
      setLoading(true);
      const response = await api.get(
        `/github/repositories?include_private=${includePrivate}`
      );
      setRepositories(response.data.repositories || []);
    } catch (error) {
      console.error("Failed to fetch repositories:", error);
    } finally {
      setLoading(false);
    }
  };

  const toggleRepository = (fullName: string) => {
    const newSelected = new Set(selectedRepos);
    if (newSelected.has(fullName)) {
      newSelected.delete(fullName);
    } else {
      if (newSelected.size < 6) {
        newSelected.add(fullName);
      }
    }
    setSelectedRepos(newSelected);
  };

  const handleContinue = async () => {
    if (selectedRepos.size === 0) return;

    setAnalyzing(true);
    try {
      // Batch analyze selected repositories
      const response = await api.post("/github/repositories/batch-analyze", {
        repositories: Array.from(selectedRepos),
      });

      // Store analyses in session storage for next page
      sessionStorage.setItem(
        "analyses",
        JSON.stringify(response.data.analyses)
      );

      navigate("/profile-builder");
    } catch (error) {
      console.error("Failed to analyze repositories:", error);
    } finally {
      setAnalyzing(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex items-center justify-center">
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
        >
          <RefreshCw className="w-12 h-12 text-purple-400" />
        </motion.div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      {/* Header */}
      <header className="border-b border-white/10 bg-black/20 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <GitBranch className="w-8 h-8 text-purple-400" />
              <div>
                <h1 className="text-xl font-bold text-white">GitRight</h1>
                <p className="text-sm text-gray-400">
                  Select your best projects
                </p>
              </div>
            </div>
            <div className="flex items-center gap-4">
              <img
                src={user?.avatar_url}
                alt={user?.username}
                className="w-10 h-10 rounded-full border-2 border-purple-400"
              />
              <span className="text-white">{user?.username}</span>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Instructions */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white/5 backdrop-blur-sm rounded-xl p-6 mb-8 border border-white/10"
        >
          <h2 className="text-2xl font-bold text-white mb-2">
            Select 3-6 repositories to showcase
          </h2>
          <p className="text-gray-300">
            Choose your best projects that represent your skills and experience.
            We'll analyze them and generate a compelling profile.
          </p>
          <div className="mt-4 flex items-center gap-4">
            <label className="flex items-center gap-2 text-white cursor-pointer">
              <input
                type="checkbox"
                checked={includePrivate}
                onChange={(e) => setIncludePrivate(e.target.checked)}
                className="w-4 h-4 rounded border-gray-600 text-purple-600 focus:ring-purple-500"
              />
              Include private repositories
            </label>
            <span className="text-sm text-gray-400">
              {selectedRepos.size}/6 selected
            </span>
          </div>
        </motion.div>

        {/* Repository Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
          {repositories.map((repo, index) => {
            const isSelected = selectedRepos.has(repo.full_name);
            const canSelect = selectedRepos.size < 6 || isSelected;

            return (
              <motion.div
                key={repo.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.05 }}
                onClick={() => canSelect && toggleRepository(repo.full_name)}
                className={`
                  relative p-5 rounded-xl border-2 transition-all cursor-pointer
                  ${
                    isSelected
                      ? "bg-purple-500/20 border-purple-400"
                      : "bg-white/5 border-white/10 hover:border-purple-400/50"
                  }
                  ${!canSelect && "opacity-50 cursor-not-allowed"}
                `}
              >
                {/* Selection indicator */}
                <div className="absolute top-3 right-3">
                  {isSelected ? (
                    <CheckCircle className="w-6 h-6 text-purple-400" />
                  ) : (
                    <Circle className="w-6 h-6 text-gray-500" />
                  )}
                </div>

                {/* Repository info */}
                <div className="pr-8">
                  <h3 className="text-lg font-semibold text-white mb-1 truncate">
                    {repo.name}
                  </h3>
                  <p className="text-sm text-gray-400 mb-3 line-clamp-2 h-10">
                    {repo.description || "No description available"}
                  </p>

                  {/* Stats */}
                  <div className="flex items-center gap-4 text-sm text-gray-400 mb-3">
                    <div className="flex items-center gap-1">
                      <Star className="w-4 h-4" />
                      <span>{repo.stargazers_count}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <GitFork className="w-4 h-4" />
                      <span>{repo.forks_count}</span>
                    </div>
                    {repo.language && (
                      <div className="flex items-center gap-1">
                        <div className="w-3 h-3 rounded-full bg-blue-400" />
                        <span>{repo.language}</span>
                      </div>
                    )}
                  </div>

                  {/* Topics */}
                  {repo.topics && repo.topics.length > 0 && (
                    <div className="flex flex-wrap gap-1">
                      {repo.topics.slice(0, 3).map((topic) => (
                        <span
                          key={topic}
                          className="px-2 py-1 text-xs rounded-full bg-white/10 text-gray-300"
                        >
                          {topic}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
              </motion.div>
            );
          })}
        </div>

        {/* Continue Button */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="flex justify-center"
        >
          <button
            onClick={handleContinue}
            disabled={selectedRepos.size === 0 || analyzing}
            className={`
              px-8 py-4 rounded-xl font-semibold text-lg
              transition-all transform hover:scale-105
              ${
                selectedRepos.size === 0 || analyzing
                  ? "bg-gray-700 text-gray-500 cursor-not-allowed"
                  : "bg-gradient-to-r from-purple-600 to-blue-600 text-white hover:shadow-lg hover:shadow-purple-500/50"
              }
            `}
          >
            {analyzing ? (
              <span className="flex items-center gap-2">
                <RefreshCw className="w-5 h-5 animate-spin" />
                Analyzing repositories...
              </span>
            ) : (
              `Continue with ${selectedRepos.size} repositories`
            )}
          </button>
        </motion.div>
      </main>
    </div>
  );
}
