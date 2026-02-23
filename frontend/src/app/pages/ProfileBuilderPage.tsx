import { useState, useEffect, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import { Sparkles, Download, Eye, Loader2 } from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import { useAuthStore } from "../../store/authStore";
import { profileAPI } from "../../lib/api";
import type {
  RepositoryAnalysis,
  ProfileConfig,
  ContentGenerationResponse,
} from "../../types";

export function ProfileBuilderPage() {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const [analyses, setAnalyses] = useState<Record<string, RepositoryAnalysis>>(
    {}
  );
  const [userApiKey, setUserApiKey] = useState<string>(""); // User's Gemini API key
  const [config, setConfig] = useState<Partial<ProfileConfig>>({
    target_role: "",
    skills_emphasis: [],
    tone_of_voice: "professional",
    template_id: "technical_deep_dive",
    contact_prefs: {
      email: user?.email || "",
      linkedin: "",
      personal_website: "",
      twitter: "",
      preferred_order: ["email"],
    },
    show_private_repos: false,
  });
  const [generatedContent, setGeneratedContent] =
    useState<ContentGenerationResponse | null>(null);
  const [isGenerating, setIsGenerating] = useState(false);
  const [isDeploying, setIsDeploying] = useState(false);

  useEffect(() => {
    // Load analyses from session storage
    const storedAnalyses = sessionStorage.getItem("analyses");
    if (!storedAnalyses) {
      navigate("/dashboard");
      return;
    }

    try {
      const parsed = JSON.parse(storedAnalyses);
      // Filter out null analyses
      const validAnalyses: Record<string, RepositoryAnalysis> = {};
      Object.entries(parsed).forEach(([key, value]) => {
        if (value !== null && typeof value === "object") {
          validAnalyses[key] = value as RepositoryAnalysis;
        }
      });

      if (Object.keys(validAnalyses).length === 0) {
        alert(
          "No valid repository analyses found. Please select repositories and try again."
        );
        navigate("/dashboard");
        return;
      }

      setAnalyses(validAnalyses);
    } catch (error) {
      console.error("Failed to parse analyses:", error);
      navigate("/dashboard");
    }
  }, [navigate]);

  const handleGenerate = async () => {
    // Validate API key
    if (!userApiKey || userApiKey.trim() === "") {
      alert(
        "Please provide your Gemini API Key. Get one for free at: https://aistudio.google.com/app/apikey"
      );
      return;
    }

    setIsGenerating(true);
    try {
      // Get ALL valid analyses
      const analysesArray = Object.values(analyses).filter(
        (a) => a !== null && a.repository
      );

      if (analysesArray.length === 0) {
        alert(
          "No valid repository analyses found. Please go back to the dashboard."
        );
        setIsGenerating(false);
        return;
      }

      // Use the unified /generate endpoint with projects array AND user API key
      const response = await profileAPI.generate({
        target_role: config.target_role || "Software Engineer",
        emphasized_skills: config.skills_emphasis || [],
        tone_of_voice: config.tone_of_voice || "professional",
        contact_prefs: config.contact_prefs || {
          linkedin: "",
          personal_website: "",
          email: "",
          twitter: "",
          preferred_order: [],
        },
        projects: analysesArray, // Send ALL analyses
        user_api_key: userApiKey, // Include user's API key
      });

      setGeneratedContent(response);
    } catch (error: any) {
      console.error("Failed to generate profile:", error);
      const errorMessage =
        error.response?.data?.message ||
        error.message ||
        "Unknown error occurred";
      alert(`Failed to generate profile: ${errorMessage}`);
    } finally {
      setIsGenerating(false);
    }
  };

  const handleDeploy = async () => {
    if (!generatedContent) return;

    // Validate API key
    if (!userApiKey || userApiKey.trim() === "") {
      alert("API key is required for deployment");
      return;
    }

    setIsDeploying(true);
    try {
      // Get the first valid analysis to use as the main project
      const analysesArray = Object.values(analyses).filter(
        (a) => a !== null && a.repository
      );

      if (analysesArray.length === 0) {
        alert("No valid repository analyses found.");
        setIsDeploying(false);
        return;
      }

      // Use same request as generate with API key
      await profileAPI.deploy({
        target_role: config.target_role || "Software Engineer",
        emphasized_skills: config.skills_emphasis || [],
        tone_of_voice: config.tone_of_voice || "professional",
        contact_prefs: config.contact_prefs || {
          linkedin: "",
          personal_website: "",
          email: "",
          twitter: "",
          preferred_order: [],
        },
        projects: analysesArray,
        user_api_key: userApiKey, // Include user's API key
      });

      // Show success and redirect
      alert("Profile deployed successfully! Check your GitHub profile.");
      window.open(`https://github.com/${user?.username}`, "_blank");
    } catch (error: any) {
      console.error("Failed to deploy profile:", error);
      const errorMessage =
        error.response?.data?.message ||
        error.message ||
        "Unknown error occurred";
      alert(`Failed to deploy profile: ${errorMessage}`);
    } finally {
      setIsDeploying(false);
    }
  };

  const toneOptions = useMemo(
    () => [
      { value: "professional", label: "Professional" },
      { value: "casual", label: "Casual" },
      { value: "technical", label: "Technical" },
    ],
    []
  );

  const templateOptions = useMemo(
    () => [
      { value: "technical_deep_dive", label: "Technical Deep Dive" },
      { value: "hiring_manager_scan", label: "Hiring Manager Scan" },
      { value: "community_contributor", label: "Community Contributor" },
    ],
    []
  );

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Configuration Panel */}
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            className="space-y-6"
          >
            <div className="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
              <h2 className="text-2xl font-bold text-white mb-6 flex items-center gap-2">
                <Sparkles className="w-6 h-6 text-purple-400" />
                Configure Your Profile
              </h2>

              {/* Target Role */}
              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Target Role
                </label>
                <input
                  type="text"
                  value={config.target_role}
                  onChange={(e) =>
                    setConfig({ ...config, target_role: e.target.value })
                  }
                  placeholder="e.g., Senior Backend Engineer"
                  className="w-full px-4 py-3 rounded-lg bg-white/10 border border-white/20 text-white placeholder-gray-500 focus:outline-none focus:border-purple-400"
                />
              </div>

              {/* Gemini API Key (REQUIRED) */}
              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  🔑 Your Gemini API Key{" "}
                  <span className="text-red-400">*Required</span>
                </label>
                <input
                  type="password"
                  value={userApiKey}
                  onChange={(e) => setUserApiKey(e.target.value)}
                  placeholder="AIza..."
                  className="w-full px-4 py-3 rounded-lg bg-white/10 border border-white/20 text-white placeholder-gray-500 focus:outline-none focus:border-purple-400"
                />
                <p className="text-xs text-gray-400 mt-2">
                  Get your free API key:{" "}
                  <a
                    href="https://aistudio.google.com/app/apikey"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-purple-400 hover:text-purple-300 underline"
                  >
                    Google AI Studio
                  </a>
                  <br />
                  Your key is never stored on our servers - it's only used for
                  this session.
                </p>
              </div>

              {/* Tone of Voice */}
              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Tone of Voice
                </label>
                <div className="grid grid-cols-3 gap-2">
                  {toneOptions.map((option) => (
                    <button
                      key={option.value}
                      onClick={() =>
                        setConfig({
                          ...config,
                          tone_of_voice: option.value as any,
                        })
                      }
                      className={`
                        px-4 py-2 rounded-lg font-medium transition-all
                        ${
                          config.tone_of_voice === option.value
                            ? "bg-purple-600 text-white"
                            : "bg-white/10 text-gray-300 hover:bg-white/20"
                        }
                      `}
                    >
                      {option.label}
                    </button>
                  ))}
                </div>
              </div>

              {/* Template */}
              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Template Style
                </label>
                <select
                  value={config.template_id}
                  onChange={(e) =>
                    setConfig({ ...config, template_id: e.target.value as any })
                  }
                  className="w-full px-4 py-3 rounded-lg bg-white/10 border border-white/20 text-white focus:outline-none focus:border-purple-400"
                >
                  {templateOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>

              {/* Skills Emphasis */}
              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Skills to Emphasize (press Enter to add)
                </label>
                <div className="space-y-3">
                  <input
                    type="text"
                    placeholder="e.g., Go, Kubernetes, PostgreSQL"
                    className="w-full px-4 py-3 rounded-lg bg-white/10 border border-white/20 text-white placeholder-gray-500 focus:outline-none focus:border-purple-400"
                    onKeyDown={(e) => {
                      if (e.key === "Enter" && e.currentTarget.value.trim()) {
                        e.preventDefault();
                        const newSkill = e.currentTarget.value.trim();
                        if (!config.skills_emphasis?.includes(newSkill)) {
                          setConfig({
                            ...config,
                            skills_emphasis: [
                              ...(config.skills_emphasis || []),
                              newSkill,
                            ],
                          });
                        }
                        e.currentTarget.value = "";
                      }
                    }}
                  />
                  {config.skills_emphasis &&
                    config.skills_emphasis.length > 0 && (
                      <div className="flex flex-wrap gap-2">
                        {config.skills_emphasis.map((skill, index) => (
                          <span
                            key={`${skill}-${index}`}
                            className="px-3 py-1 rounded-full bg-purple-600/20 text-purple-300 text-sm flex items-center gap-2"
                          >
                            {skill}
                            <button
                              onClick={() =>
                                setConfig({
                                  ...config,
                                  skills_emphasis:
                                    config.skills_emphasis?.filter(
                                      (_, i) => i !== index
                                    ),
                                })
                              }
                              className="text-purple-300 hover:text-white"
                            >
                              ×
                            </button>
                          </span>
                        ))}
                      </div>
                    )}
                </div>
              </div>

              {/* Contact Preferences */}
              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Contact Information
                </label>
                <div className="space-y-3">
                  <div>
                    <label className="block text-sm text-gray-400 mb-2">
                      Email
                    </label>
                    <input
                      type="email"
                      value={config.contact_prefs?.email || ""}
                      onChange={(e) =>
                        setConfig({
                          ...config,
                          contact_prefs: {
                            ...config.contact_prefs!,
                            email: e.target.value,
                          },
                        })
                      }
                      placeholder="your@email.com"
                      className="w-full px-4 py-2 rounded-lg bg-white/10 border border-white/20 text-white"
                    />
                  </div>
                  <div>
                    <label className="block text-sm text-gray-400 mb-2">
                      LinkedIn
                    </label>
                    <input
                      type="text"
                      value={config.contact_prefs?.linkedin || ""}
                      onChange={(e) =>
                        setConfig({
                          ...config,
                          contact_prefs: {
                            ...config.contact_prefs!,
                            linkedin: e.target.value,
                          },
                        })
                      }
                      placeholder="linkedin.com/in/yourusername"
                      className="w-full px-4 py-2 rounded-lg bg-white/10 border border-white/20 text-white"
                    />
                  </div>
                  <div>
                    <label className="block text-sm text-gray-400 mb-2">
                      Personal Website
                    </label>
                    <input
                      type="text"
                      value={config.contact_prefs?.personal_website || ""}
                      onChange={(e) =>
                        setConfig({
                          ...config,
                          contact_prefs: {
                            ...config.contact_prefs!,
                            personal_website: e.target.value,
                          },
                        })
                      }
                      placeholder="yourwebsite.com"
                      className="w-full px-4 py-2 rounded-lg bg-white/10 border border-white/20 text-white"
                    />
                  </div>
                  <div>
                    <label className="block text-sm text-gray-400 mb-2">
                      Twitter
                    </label>
                    <input
                      type="text"
                      value={config.contact_prefs?.twitter || ""}
                      onChange={(e) =>
                        setConfig({
                          ...config,
                          contact_prefs: {
                            ...config.contact_prefs!,
                            twitter: e.target.value,
                          },
                        })
                      }
                      placeholder="@yourusername"
                      className="w-full px-4 py-2 rounded-lg bg-white/10 border border-white/20 text-white"
                    />
                  </div>
                </div>
              </div>

              {/* Action Buttons */}
              <div className="flex gap-4">
                <button
                  onClick={handleGenerate}
                  disabled={isGenerating || !config.target_role}
                  className="flex-1 px-6 py-3 rounded-lg bg-gradient-to-r from-purple-600 to-blue-600 text-white font-semibold hover:shadow-lg hover:shadow-purple-500/50 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
                >
                  {isGenerating ? (
                    <span className="flex items-center justify-center gap-2">
                      <Loader2 className="w-5 h-5 animate-spin" />
                      Generating...
                    </span>
                  ) : (
                    <span className="flex items-center justify-center gap-2">
                      <Sparkles className="w-5 h-5" />
                      Generate
                    </span>
                  )}
                </button>
              </div>
            </div>
          </motion.div>

          {/* Preview Panel */}
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            className="space-y-6"
          >
            <div className="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-2xl font-bold text-white flex items-center gap-2">
                  <Eye className="w-6 h-6 text-purple-400" />
                  Preview
                </h2>
                {generatedContent && (
                  <button
                    onClick={handleDeploy}
                    disabled={isDeploying}
                    className="px-4 py-2 rounded-lg bg-green-600 text-white font-semibold hover:bg-green-700 disabled:opacity-50 transition-all flex items-center gap-2"
                  >
                    {isDeploying ? (
                      <>
                        <Loader2 className="w-4 h-4 animate-spin" />
                        Deploying...
                      </>
                    ) : (
                      <>
                        <Download className="w-4 h-4" />
                        Deploy to GitHub
                      </>
                    )}
                  </button>
                )}
              </div>

              {generatedContent ? (
                <div className="bg-[#0d1117] rounded-lg p-8 overflow-auto max-h-[700px] border border-[#30363d]">
                  {/* GitHub-style README preview */}
                  <div className="markdown-body">
                    <style>{`
                      .markdown-body {
                        color: #c9d1d9;
                        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
                        font-size: 16px;
                        line-height: 1.6;
                      }
                      .markdown-body h1, .markdown-body h2, .markdown-body h3 {
                        color: #c9d1d9;
                        font-weight: 600;
                        padding-bottom: 0.3em;
                        margin-top: 24px;
                        margin-bottom: 16px;
                      }
                      .markdown-body h1 {
                        font-size: 2em;
                        border-bottom: 1px solid #21262d;
                      }
                      .markdown-body h2 {
                        font-size: 1.5em;
                        border-bottom: 1px solid #21262d;
                      }
                      .markdown-body h3 {
                        font-size: 1.25em;
                      }
                      .markdown-body a {
                        color: #58a6ff;
                        text-decoration: none;
                      }
                      .markdown-body a:hover {
                        text-decoration: underline;
                      }
                      .markdown-body img {
                        max-width: 100%;
                        background-color: transparent;
                      }
                      .markdown-body code {
                        background-color: rgba(110, 118, 129, 0.4);
                        border-radius: 6px;
                        padding: 0.2em 0.4em;
                        font-family: ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, Liberation Mono, monospace;
                        font-size: 85%;
                      }
                      .markdown-body pre {
                        background-color: #161b22;
                        border-radius: 6px;
                        padding: 16px;
                        overflow: auto;
                      }
                      .markdown-body blockquote {
                        border-left: 0.25em solid #3b434b;
                        padding: 0 1em;
                        color: #8b949e;
                      }
                      .markdown-body hr {
                        height: 0.25em;
                        padding: 0;
                        margin: 24px 0;
                        background-color: #21262d;
                        border: 0;
                      }
                      .markdown-body ul, .markdown-body ol {
                        padding-left: 2em;
                      }
                      .markdown-body li {
                        margin-top: 0.25em;
                      }
                      .markdown-body table {
                        border-collapse: collapse;
                        margin-top: 16px;
                        margin-bottom: 16px;
                      }
                      .markdown-body th, .markdown-body td {
                        border: 1px solid #30363d;
                        padding: 6px 13px;
                      }
                      .markdown-body th {
                        background-color: #161b22;
                        font-weight: 600;
                      }
                      .markdown-body p {
                        margin-top: 0;
                        margin-bottom: 16px;
                      }
                      /* Center aligned divs */
                      .markdown-body div[align="center"] {
                        text-align: center;
                      }
                      /* Badge styling */
                      .markdown-body img[src*="shields.io"],
                      .markdown-body img[src*="badge"] {
                        display: inline-block;
                        margin: 4px;
                        vertical-align: middle;
                      }
                    `}</style>
                    <ReactMarkdown
                      remarkPlugins={[remarkGfm]}
                      rehypePlugins={[rehypeRaw]}
                      components={{
                        // Render images properly
                        img: ({ node, ...props }) => (
                          <img
                            {...props}
                            style={{ display: "inline-block", margin: "4px" }}
                          />
                        ),
                      }}
                    >
                      {generatedContent.markdown || generatedContent.summary}
                    </ReactMarkdown>
                  </div>
                </div>
              ) : (
                <div className="text-center py-20 text-gray-400 bg-[#0d1117] rounded-lg border border-[#30363d]">
                  <Eye className="w-16 h-16 mx-auto mb-4 opacity-30" />
                  <p>
                    Configure your profile and click Generate to see the preview
                  </p>
                </div>
              )}
            </div>
          </motion.div>
        </div>
      </div>
    </div>
  );
}
