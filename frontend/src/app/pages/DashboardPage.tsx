import { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import {
  CheckCircle,
  Circle,
  RefreshCw,
  GitBranch,
  Star,
  GitFork,
  ArrowLeft,
  ArrowRight,
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
      const response = await api.get(`/github/repositories?include_private=${includePrivate}`);
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
    } else if (newSelected.size < 6) {
      newSelected.add(fullName);
    }
    setSelectedRepos(newSelected);
  };

  const handleContinue = async () => {
    if (selectedRepos.size === 0) return;
    setAnalyzing(true);
    try {
      const response = await api.post("/github/repositories/batch-analyze", {
        repositories: Array.from(selectedRepos),
      });
      sessionStorage.setItem("analyses", JSON.stringify(response.data.analyses));
      navigate("/profile-builder");
    } catch (error) {
      console.error("Failed to analyze repositories:", error);
    } finally {
      setAnalyzing(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-[#030712] flex items-center justify-center">
        <RefreshCw className="w-8 h-8 text-violet-400 animate-spin" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#030712] text-white">
      <header className="border-b border-white/5">
        <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors text-sm">
            <ArrowLeft className="w-4 h-4" />
            Home
          </Link>
          <div className="flex items-center gap-2">
            <GitBranch className="w-4 h-4 text-violet-400" />
            <span className="text-sm font-medium">Select Repos</span>
          </div>
          <div className="flex items-center gap-3">
            <img src={user?.avatar_url} alt="" className="w-7 h-7 rounded-full" />
            <span className="text-sm text-gray-300">{user?.username}</span>
          </div>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-6 py-8">
        <div className="mb-8">
          <h1 className="text-2xl font-bold tracking-tight mb-2">
            Select 3-6 repositories to showcase
          </h1>
          <p className="text-sm text-gray-400 mb-4">
            Choose your best projects. We'll analyze them and generate a compelling profile.
          </p>
          <div className="flex items-center gap-4">
            <label className="flex items-center gap-2 text-sm text-gray-400 cursor-pointer">
              <input
                type="checkbox"
                checked={includePrivate}
                onChange={(e) => setIncludePrivate(e.target.checked)}
                className="w-4 h-4 rounded border-white/20 bg-white/5 text-violet-500 focus:ring-violet-500"
              />
              Include private repos
            </label>
            <span className="text-xs text-gray-500">{selectedRepos.size}/6 selected</span>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3 mb-8">
          {repositories.map((repo) => {
            const isSelected = selectedRepos.has(repo.full_name);
            const canSelect = selectedRepos.size < 6 || isSelected;

            return (
              <button
                key={repo.id}
                onClick={() => canSelect && toggleRepository(repo.full_name)}
                disabled={!canSelect}
                className={`text-left p-4 rounded-xl border transition-all ${
                  isSelected
                    ? "bg-violet-500/10 border-violet-500/30"
                    : "bg-white/[0.02] border-white/[0.06] hover:border-white/[0.12]"
                } ${!canSelect ? "opacity-40 cursor-not-allowed" : "cursor-pointer"}`}
              >
                <div className="flex items-start justify-between mb-2">
                  <h3 className="text-sm font-semibold text-white/90 truncate pr-6">{repo.name}</h3>
                  {isSelected ? (
                    <CheckCircle className="w-4 h-4 text-violet-400 shrink-0" />
                  ) : (
                    <Circle className="w-4 h-4 text-gray-600 shrink-0" />
                  )}
                </div>
                <p className="text-xs text-gray-500 line-clamp-2 mb-2 h-8">
                  {repo.description || "No description"}
                </p>
                <div className="flex items-center gap-3 text-[10px] text-gray-500">
                  <span className="flex items-center gap-1">
                    <Star className="w-3 h-3" />{repo.stargazers_count}
                  </span>
                  <span className="flex items-center gap-1">
                    <GitFork className="w-3 h-3" />{repo.forks_count}
                  </span>
                  {repo.language && (
                    <span className="flex items-center gap-1">
                      <div className="w-2 h-2 rounded-full bg-violet-400" />
                      {repo.language}
                    </span>
                  )}
                </div>
                {repo.topics?.length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-2">
                    {repo.topics.slice(0, 3).map((t) => (
                      <span key={t} className="px-1.5 py-0.5 text-[9px] rounded bg-white/5 text-gray-500">{t}</span>
                    ))}
                  </div>
                )}
              </button>
            );
          })}
        </div>

        {repositories.length === 0 && (
          <div className="text-center py-16 text-gray-500">
            <p className="text-sm">No repositories found. Check your GitHub account.</p>
          </div>
        )}

        <div className="flex justify-center">
          <button
            onClick={handleContinue}
            disabled={selectedRepos.size === 0 || analyzing}
            className="px-8 py-3 rounded-lg bg-gradient-to-r from-violet-600 to-cyan-600 text-sm font-medium text-white disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center gap-2"
          >
            {analyzing ? (
              <>
                <RefreshCw className="w-4 h-4 animate-spin" />
                Analyzing...
              </>
            ) : (
              <>
                Continue with {selectedRepos.size} repos
                <ArrowRight className="w-4 h-4" />
              </>
            )}
          </button>
        </div>
      </main>
    </div>
  );
}
