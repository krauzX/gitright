import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import {
  ArrowLeft,
  ArrowRight,
  Check,
  Download,
  GitBranch,
  Loader2,
  Sparkles,
  Zap,
  User,
  MapPin,
  Building2,
  Globe,
  Mail,
  X,
  Image,
  Copy,
} from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import { useAuthStore } from "../../store/authStore";
import { profileAPI } from "../../lib/api";
import type {
  RepositoryAnalysis,
  ContentGenerationResponse,
  AutoImportResponse,
} from "../../types";

type Step = "import" | "configure" | "generate" | "preview";

export function ProfileBuilderPage() {
  const user = useAuthStore((state) => state.user);
  const [step, setStep] = useState<Step>("import");
  const [analyses, setAnalyses] = useState<Record<string, RepositoryAnalysis>>({});
  const [importedData, setImportedData] = useState<AutoImportResponse | null>(null);
  const [generatedContent, setGeneratedContent] = useState<ContentGenerationResponse | null>(null);
  const [userApiKey, setUserApiKey] = useState("");
  const [targetRole, setTargetRole] = useState("");
  const [tone, setTone] = useState<"professional" | "friendly" | "technical">("professional");
  const [skills, setSkills] = useState<string[]>([]);
  const [newSkill, setNewSkill] = useState("");
  const [contact, setContact] = useState({ email: "", linkedin: "", website: "", twitter: "", personal_website: "" });
  const [loading, setLoading] = useState<"import" | "generate" | "deploy" | null>(null);
  const [bannerSvg, setBannerSvg] = useState<string | null>(null);
  const [timelineSvg, setTimelineSvg] = useState<string | null>(null);
  const [loadingBanner, setLoadingBanner] = useState(false);
  const [loadingTimeline, setLoadingTimeline] = useState(false);
  const [toast, setToast] = useState<{ message: string; type: "success" | "error" | "info" } | null>(null);

  const showToast = (message: string, type: "success" | "error" | "info" = "info") => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 4000);
  };

  useEffect(() => {
    const stored = sessionStorage.getItem("analyses");
    if (stored) {
      try {
        const parsed = JSON.parse(stored);
        const valid: Record<string, RepositoryAnalysis> = {};
        Object.entries(parsed).forEach(([k, v]) => {
          if (v && typeof v === "object" && (v as RepositoryAnalysis).repository) {
            valid[k] = v as RepositoryAnalysis;
          }
        });
        if (Object.keys(valid).length > 0) setAnalyses(valid);
      } catch {}
    }
    const key = localStorage.getItem("gemini_api_key");
    if (key) setUserApiKey(key);
  }, []);

  const hasAnalyses = Object.keys(analyses).length > 0;

  const handleAutoImport = async () => {
    setLoading("import");
    try {
      const data = await profileAPI.autoImport(userApiKey || "");
      setImportedData(data);
      setSkills(data.skills || []);
      setContact({
        email: data.email || user?.email || "",
        linkedin: data.social_links?.find((s) => s.provider === "LINKEDIN")?.url || "",
        website: data.website_url || "",
        twitter: data.twitter_username ? `@${data.twitter_username}` : "",
        personal_website: data.website_url || "",
      });
      setStep("configure");
    } catch {
      showToast("Failed to import from GitHub.", "error");
    } finally {
      setLoading(null);
    }
  };

  const handleGenerate = async () => {
    if (!userApiKey.trim()) {
      showToast("Gemini API key required. Get one free at aistudio.google.com/app/apikey", "error");
      return;
    }
    if (!hasAnalyses) {
      showToast("No repositories found. Go back to Dashboard to select repos.", "error");
      return;
    }
    setLoading("generate");
    try {
      const projects = Object.values(analyses).filter((a) => a?.repository);
      const response = await profileAPI.generate({
        target_role: targetRole || "Software Engineer",
        emphasized_skills: skills,
        tone_of_voice: tone,
        contact_prefs: { ...contact, preferred_order: ["email"] },
        projects,
        user_api_key: userApiKey,
      });
      setGeneratedContent(response);
      setStep("preview");
    } catch (e: any) {
      showToast(e.response?.data?.message || "Generation failed", "error");
    } finally {
      setLoading(null);
    }
  };

  const handleDeploy = async () => {
    if (!generatedContent) return;
    setLoading("deploy");
    try {
      const projects = Object.values(analyses).filter((a) => a?.repository);
      await profileAPI.deploy({
        target_role: targetRole || "Software Engineer",
        emphasized_skills: skills,
        tone_of_voice: tone,
        contact_prefs: { ...contact, preferred_order: ["email"] },
        projects,
        user_api_key: userApiKey,
      });
      showToast("Deployed! Check your GitHub profile.", "success");
      window.open(`https://github.com/${user?.username}`, "_blank");
    } catch (e: any) {
      showToast(e.response?.data?.message || "Deploy failed", "error");
    } finally {
      setLoading(null);
    }
  };

  const addSkill = () => {
    const s = newSkill.trim();
    if (s && !skills.includes(s)) {
      setSkills([...skills, s]);
      setNewSkill("");
    }
  };

  const fetchBanner = async () => {
    setLoadingBanner(true);
    try {
      const svg = await profileAPI.getBanner("dark");
      setBannerSvg(svg);
    } catch {
      showToast("Failed to generate banner.", "error");
    } finally {
      setLoadingBanner(false);
    }
  };

  const fetchTimeline = async () => {
    setLoadingTimeline(true);
    try {
      const svg = await profileAPI.getTimeline("dark");
      setTimelineSvg(svg);
    } catch {
      showToast("Failed to generate timeline.", "error");
    } finally {
      setLoadingTimeline(false);
    }
  };

  const downloadSvg = (svgContent: string, filename: string) => {
    const blob = new Blob([svgContent], { type: "image/svg+xml" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
  };

  const copyEmbedCode = (type: "banner" | "timeline") => {
    const username = user?.username || "username";
    const code = `![${type}](https://raw.githubusercontent.com/${username}/${username}/main/${type}.svg)`;
    navigator.clipboard.writeText(code);
    showToast("Embed code copied!", "success");
  };

  const steps: { id: Step; label: string; num: string }[] = [
    { id: "import", label: "Import", num: "01" },
    { id: "configure", label: "Configure", num: "02" },
    { id: "generate", label: "Generate", num: "03" },
    { id: "preview", label: "Preview", num: "04" },
  ];

  const stepIndex = steps.findIndex((s) => s.id === step);

  return (
    <div className="min-h-screen bg-[#030712] text-white">
      {toast && (
        <div className={`fixed top-4 right-4 z-50 px-4 py-3 rounded-lg text-sm font-medium shadow-lg transition-all ${
          toast.type === "success" ? "bg-emerald-600 text-white" :
          toast.type === "error" ? "bg-red-600 text-white" :
          "bg-violet-600 text-white"
        }`}>
          {toast.message}
        </div>
      )}
      <header className="border-b border-white/5">
        <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
          <Link to="/dashboard" className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors text-sm">
            <ArrowLeft className="w-4 h-4" />
            Dashboard
          </Link>
          <div className="flex items-center gap-2">
            <GitBranch className="w-4 h-4 text-violet-400" />
            <span className="text-sm font-medium">Profile Builder</span>
          </div>
          <div className="w-16" />
        </div>
      </header>

      <div className="max-w-6xl mx-auto px-6 py-8">
        {/* Step Indicator */}
        <div className="flex items-center justify-center gap-1 mb-10">
          {steps.map((s, i) => (
            <div key={s.id} className="flex items-center">
              <button
                onClick={() => {
                  if (i < stepIndex || (s.id === "configure" && importedData) || (s.id === "preview" && generatedContent)) {
                    setStep(s.id);
                  }
                }}
                className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium transition-all ${
                  step === s.id
                    ? "bg-violet-600 text-white"
                    : i < stepIndex
                    ? "bg-emerald-600/20 text-emerald-400"
                    : "bg-white/5 text-gray-500"
                }`}
              >
                {i < stepIndex ? (
                  <Check className="w-3 h-3" />
                ) : (
                  <span>{s.num}</span>
                )}
                <span className="hidden sm:inline">{s.label}</span>
              </button>
              {i < steps.length - 1 && (
                <div className={`w-8 h-px mx-1 ${i < stepIndex ? "bg-emerald-500/50" : "bg-white/10"}`} />
              )}
            </div>
          ))}
        </div>

        {/* Step Content */}
        <div className="grid grid-cols-1 lg:grid-cols-5 gap-8">
          <div className="lg:col-span-2 space-y-4">
            {step === "import" && (
              <ImportStep
                userApiKey={userApiKey}
                setUserApiKey={setUserApiKey}
                hasAnalyses={hasAnalyses}
                onImport={handleAutoImport}
                loading={loading === "import"}
              />
            )}
            {step === "configure" && importedData && (
              <ConfigureStep
                targetRole={targetRole}
                setTargetRole={setTargetRole}
                tone={tone}
                setTone={setTone}
                skills={skills}
                setSkills={setSkills}
                newSkill={newSkill}
                setNewSkill={setNewSkill}
                addSkill={addSkill}
                contact={contact}
                setContact={setContact}
                userApiKey={userApiKey}
                setUserApiKey={setUserApiKey}
                onNext={() => setStep("generate")}
              />
            )}
            {(step === "generate" || step === "preview") && (
              <GenerateStep
                targetRole={targetRole}
                tone={tone}
                skills={skills}
                userApiKey={userApiKey}
                onGenerate={handleGenerate}
                onDeploy={handleDeploy}
                loading={loading}
                hasContent={!!generatedContent}
              />
            )}
          </div>

          <div className="lg:col-span-3">
            {step === "import" && !importedData && (
              <div className="bg-white/[0.02] border border-white/[0.06] rounded-xl p-8 text-center">
                <div className="w-16 h-16 mx-auto rounded-2xl bg-gradient-to-br from-violet-500/20 to-cyan-500/20 flex items-center justify-center mb-4">
                  <User className="w-8 h-8 text-violet-400" />
                </div>
                <h3 className="text-lg font-semibold mb-2">Your Profile Data</h3>
                <p className="text-sm text-gray-400 max-w-sm mx-auto">
                  Import your GitHub data to auto-fill skills, stats, and contact info.
                  Or skip and configure manually.
                </p>
              </div>
            )}

            {step === "import" && importedData && (
              <ImportedDataPreview data={importedData} />
            )}

            {step === "import" && importedData && (
              <div className="bg-white/[0.02] border border-white/[0.06] rounded-xl p-6 space-y-4">
                <h3 className="text-base font-semibold">SVG Assets</h3>
                <p className="text-xs text-gray-400">Generate premium banners for your profile README.</p>
                <div className="grid grid-cols-2 gap-3">
                  <button
                    onClick={fetchBanner}
                    disabled={loadingBanner}
                    className="px-4 py-3 rounded-lg bg-white/5 border border-white/10 text-sm font-medium text-gray-300 hover:bg-white/10 transition-all flex items-center justify-center gap-2"
                  >
                    {loadingBanner ? <Loader2 className="w-4 h-4 animate-spin" /> : <Image className="w-4 h-4" />}
                    Banner
                  </button>
                  <button
                    onClick={fetchTimeline}
                    disabled={loadingTimeline}
                    className="px-4 py-3 rounded-lg bg-white/5 border border-white/10 text-sm font-medium text-gray-300 hover:bg-white/10 transition-all flex items-center justify-center gap-2"
                  >
                    {loadingTimeline ? <Loader2 className="w-4 h-4 animate-spin" /> : <Image className="w-4 h-4" />}
                    Timeline
                  </button>
                </div>
                {bannerSvg && (
                  <div className="space-y-2">
                    <div className="rounded-lg overflow-hidden border border-white/10">
                      <div dangerouslySetInnerHTML={{ __html: bannerSvg }} />
                    </div>
                    <div className="flex gap-2">
                      <button onClick={() => downloadSvg(bannerSvg, "banner.svg")} className="flex-1 px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-[10px] text-gray-400 hover:bg-white/10 flex items-center justify-center gap-1">
                        <Download className="w-3 h-3" /> Download
                      </button>
                      <button onClick={() => copyEmbedCode("banner")} className="flex-1 px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-[10px] text-gray-400 hover:bg-white/10 flex items-center justify-center gap-1">
                        <Copy className="w-3 h-3" /> Copy Embed
                      </button>
                    </div>
                  </div>
                )}
                {timelineSvg && (
                  <div className="space-y-2">
                    <div className="rounded-lg overflow-hidden border border-white/10">
                      <div dangerouslySetInnerHTML={{ __html: timelineSvg }} />
                    </div>
                    <div className="flex gap-2">
                      <button onClick={() => downloadSvg(timelineSvg, "timeline.svg")} className="flex-1 px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-[10px] text-gray-400 hover:bg-white/10 flex items-center justify-center gap-1">
                        <Download className="w-3 h-3" /> Download
                      </button>
                      <button onClick={() => copyEmbedCode("timeline")} className="flex-1 px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-[10px] text-gray-400 hover:bg-white/10 flex items-center justify-center gap-1">
                        <Copy className="w-3 h-3" /> Copy Embed
                      </button>
                    </div>
                  </div>
                )}
              </div>
            )}

            {step === "configure" && importedData && (
              <ConfigPreview data={importedData} skills={skills} contact={contact} />
            )}

            {(step === "generate" || step === "preview") && generatedContent && (
              <PreviewPanel content={generatedContent} />
            )}

            {(step === "generate" || step === "preview") && !generatedContent && (
              <div className="bg-white/[0.02] border border-white/[0.06] rounded-xl p-8 text-center">
                <div className="w-16 h-16 mx-auto rounded-2xl bg-gradient-to-br from-violet-500/20 to-cyan-500/20 flex items-center justify-center mb-4">
                  <Sparkles className="w-8 h-8 text-violet-400" />
                </div>
                <h3 className="text-lg font-semibold mb-2">Generate Your Profile</h3>
                <p className="text-sm text-gray-400 max-w-sm mx-auto">
                  Click Generate to have AI write your GitHub profile README
                  based on your real repositories and skills.
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

function ImportStep({
  userApiKey,
  setUserApiKey,
  hasAnalyses,
  onImport,
  loading,
}: {
  userApiKey: string;
  setUserApiKey: (v: string) => void;
  hasAnalyses: boolean;
  onImport: () => void;
  loading: boolean;
}) {
  return (
    <div className="bg-white/[0.02] border border-white/[0.06] rounded-xl p-6 space-y-5">
      <div>
        <h3 className="text-base font-semibold mb-1">Import from GitHub</h3>
        <p className="text-xs text-gray-400">Auto-fill skills, stats, and contact info from your profile.</p>
      </div>

      <div>
        <label className="block text-xs font-medium text-gray-400 mb-1.5">
          Gemini API Key <span className="text-red-400">*</span>
        </label>
        <input
          type="password"
          value={userApiKey}
          onChange={(e) => setUserApiKey(e.target.value)}
          placeholder="AIza..."
          className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-sm text-white placeholder-gray-500 focus:outline-none focus:border-violet-500/50"
        />
        <a href="https://aistudio.google.com/app/apikey" target="_blank" rel="noopener noreferrer"
          className="text-[10px] text-violet-400 hover:text-violet-300 mt-1 inline-block">
          Get free key at Google AI Studio
        </a>
      </div>

      <button
        onClick={onImport}
        disabled={loading}
        className="w-full px-4 py-2.5 rounded-lg bg-gradient-to-r from-violet-600 to-cyan-600 text-sm font-medium text-white disabled:opacity-50 transition-all flex items-center justify-center gap-2"
      >
        {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Zap className="w-4 h-4" />}
        {loading ? "Importing..." : "Auto-Import from GitHub"}
      </button>

      {!hasAnalyses && (
        <p className="text-[10px] text-amber-400/80 text-center">
          No repos analyzed yet.{" "}
          <Link to="/dashboard" className="underline hover:text-amber-300">Go to Dashboard</Link> to select repos first.
        </p>
      )}
    </div>
  );
}

function ImportedDataPreview({ data }: { data: AutoImportResponse }) {
  return (
    <div className="bg-white/[0.02] border border-white/[0.06] rounded-xl p-6 space-y-4">
      <h3 className="text-base font-semibold">Imported Data</h3>

      <div className="grid grid-cols-3 gap-3">
        {[
          { label: "Contributions", value: data.total_contributions },
          { label: "Followers", value: data.followers },
          { label: "Repos", value: data.top_repositories?.length || 0 },
        ].map((s) => (
          <div key={s.label} className="text-center p-3 rounded-lg bg-white/[0.03]">
            <div className="text-lg font-bold text-white">{s.value.toLocaleString()}</div>
            <div className="text-[10px] text-gray-500">{s.label}</div>
          </div>
        ))}
      </div>

      <div className="space-y-2 text-xs">
        {data.name && (
          <div className="flex items-center gap-2 text-gray-300">
            <User className="w-3.5 h-3.5 text-gray-500" />
            {data.name}
          </div>
        )}
        {data.location && (
          <div className="flex items-center gap-2 text-gray-300">
            <MapPin className="w-3.5 h-3.5 text-gray-500" />
            {data.location}
          </div>
        )}
        {data.company && (
          <div className="flex items-center gap-2 text-gray-300">
            <Building2 className="w-3.5 h-3.5 text-gray-500" />
            {data.company}
          </div>
        )}
        {data.email && (
          <div className="flex items-center gap-2 text-gray-300">
            <Mail className="w-3.5 h-3.5 text-gray-500" />
            {data.email}
          </div>
        )}
      </div>

      {data.skills?.length > 0 && (
        <div>
          <div className="text-[10px] text-gray-500 mb-1.5">SKILLS</div>
          <div className="flex flex-wrap gap-1">
            {data.skills.map((s) => (
              <span key={s} className="px-2 py-0.5 rounded text-[10px] bg-violet-500/10 text-violet-300 border border-violet-500/20">
                {s}
              </span>
            ))}
          </div>
        </div>
      )}

      {data.top_repositories?.length > 0 && (
        <div>
          <div className="text-[10px] text-gray-500 mb-1.5">TOP REPOS</div>
          <div className="space-y-1">
            {data.top_repositories.slice(0, 5).map((r) => (
              <div key={r.name} className="flex items-center justify-between text-xs">
                <span className="text-gray-300 truncate">{r.name}</span>
                <span className="text-gray-500 shrink-0 ml-2">{r.language || "—"}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function ConfigureStep({
  targetRole, setTargetRole,
  tone, setTone,
  skills, setSkills, newSkill, setNewSkill, addSkill,
  contact, setContact,
  userApiKey, setUserApiKey,
  onNext,
}: {
  targetRole: string;
  setTargetRole: (v: string) => void;
  tone: "professional" | "friendly" | "technical";
  setTone: (v: "professional" | "friendly" | "technical") => void;
  skills: string[];
  setSkills: (v: string[]) => void;
  newSkill: string;
  setNewSkill: (v: string) => void;
  addSkill: () => void;
  contact: { email: string; linkedin: string; website: string; twitter: string; personal_website: string };
  setContact: (v: { email: string; linkedin: string; website: string; twitter: string; personal_website: string }) => void;
  userApiKey: string;
  setUserApiKey: (v: string) => void;
  onNext: () => void;
}) {
  return (
    <div className="bg-white/[0.02] border border-white/[0.06] rounded-xl p-6 space-y-5">
      <h3 className="text-base font-semibold">Configure Profile</h3>

      <div>
        <label className="block text-xs font-medium text-gray-400 mb-1.5">Target Role</label>
        <input
          type="text"
          value={targetRole}
          onChange={(e) => setTargetRole(e.target.value)}
          placeholder="e.g., Senior Backend Engineer"
          className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-sm text-white placeholder-gray-500 focus:outline-none focus:border-violet-500/50"
        />
      </div>

      <div>
        <label className="block text-xs font-medium text-gray-400 mb-1.5">Tone</label>
        <div className="grid grid-cols-3 gap-1.5">
          {(["professional", "friendly", "technical"] as const).map((t) => (
            <button
              key={t}
              onClick={() => setTone(t)}
              className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all ${
                tone === t
                  ? "bg-violet-600 text-white"
                  : "bg-white/5 text-gray-400 hover:bg-white/10"
              }`}
            >
              {t.charAt(0).toUpperCase() + t.slice(1)}
            </button>
          ))}
        </div>
      </div>

      <div>
        <label className="block text-xs font-medium text-gray-400 mb-1.5">Skills</label>
        <div className="flex gap-1.5 mb-2">
          <input
            type="text"
            value={newSkill}
            onChange={(e) => setNewSkill(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && (e.preventDefault(), addSkill())}
            placeholder="Add skill..."
            className="flex-1 px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-violet-500/50"
          />
          <button onClick={addSkill} className="px-3 py-1.5 rounded-lg bg-white/10 text-xs text-gray-300 hover:bg-white/15">
            Add
          </button>
        </div>
        <div className="flex flex-wrap gap-1">
          {skills.map((s, i) => (
            <span key={`${s}-${i}`} className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] bg-violet-500/10 text-violet-300 border border-violet-500/20">
              {s}
              <button onClick={() => setSkills(skills.filter((_, j) => j !== i))} className="hover:text-white">
                <X className="w-2.5 h-2.5" />
              </button>
            </span>
          ))}
        </div>
      </div>

      <div>
        <label className="block text-xs font-medium text-gray-400 mb-1.5">Contact</label>
        <div className="space-y-1.5">
          {[
            { key: "email" as const, icon: Mail, placeholder: "Email" },
            { key: "linkedin" as const, icon: Globe, placeholder: "LinkedIn URL" },
            { key: "website" as const, icon: Globe, placeholder: "Website URL" },
            { key: "twitter" as const, icon: Globe, placeholder: "@twitter" },
          ].map(({ key, placeholder }) => (
            <input
              key={key}
              type="text"
              value={contact[key]}
              onChange={(e) => setContact({ ...contact, [key]: e.target.value })}
              placeholder={placeholder}
              className="w-full px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-violet-500/50"
            />
          ))}
        </div>
      </div>

      <div>
        <label className="block text-xs font-medium text-gray-400 mb-1.5">
          Gemini API Key <span className="text-red-400">*</span>
        </label>
        <input
          type="password"
          value={userApiKey}
          onChange={(e) => setUserApiKey(e.target.value)}
          placeholder="AIza..."
          className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-sm text-white placeholder-gray-500 focus:outline-none focus:border-violet-500/50"
        />
      </div>

      <button
        onClick={onNext}
        disabled={!userApiKey}
        className="w-full px-4 py-2.5 rounded-lg bg-gradient-to-r from-violet-600 to-cyan-600 text-sm font-medium text-white disabled:opacity-50 transition-all flex items-center justify-center gap-2"
      >
        Next
        <ArrowRight className="w-4 h-4" />
      </button>
    </div>
  );
}

function ConfigPreview({
  data,
  skills,
  contact,
}: {
  data: AutoImportResponse;
  skills: string[];
  contact: { email: string; linkedin: string; website: string; twitter: string; personal_website: string };
}) {
  return (
    <div className="bg-white/[0.02] border border-white/[0.06] rounded-xl p-6 space-y-4">
      <h3 className="text-base font-semibold">Configuration Summary</h3>
      <div className="space-y-2 text-xs">
        <div className="flex justify-between">
          <span className="text-gray-500">Name</span>
          <span className="text-gray-300">{data.name || data.login}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-500">Skills</span>
          <span className="text-gray-300">{skills.length} selected</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-500">Contributions</span>
          <span className="text-gray-300">{data.total_contributions}</span>
        </div>
        {contact.email && (
          <div className="flex justify-between">
            <span className="text-gray-500">Email</span>
            <span className="text-gray-300">{contact.email}</span>
          </div>
        )}
      </div>
    </div>
  );
}

function GenerateStep({
  targetRole,
  tone,
  skills,
  userApiKey,
  onGenerate,
  onDeploy,
  loading,
  hasContent,
}: {
  targetRole: string;
  tone: string;
  skills: string[];
  userApiKey: string;
  onGenerate: () => void;
  onDeploy: () => void;
  loading: "import" | "generate" | "deploy" | null;
  hasContent: boolean;
}) {
  return (
    <div className="bg-white/[0.02] border border-white/[0.06] rounded-xl p-6 space-y-4">
      <h3 className="text-base font-semibold">Generate & Deploy</h3>
      <div className="space-y-2 text-xs text-gray-400">
        <div className="flex justify-between"><span>Role</span><span className="text-gray-300">{targetRole || "Software Engineer"}</span></div>
        <div className="flex justify-between"><span>Tone</span><span className="text-gray-300">{tone}</span></div>
        <div className="flex justify-between"><span>Skills</span><span className="text-gray-300">{skills.length}</span></div>
      </div>
      <div className="flex gap-2">
        <button
          onClick={onGenerate}
          disabled={loading !== null || !userApiKey}
          className="flex-1 px-4 py-2.5 rounded-lg bg-gradient-to-r from-violet-600 to-cyan-600 text-sm font-medium text-white disabled:opacity-50 transition-all flex items-center justify-center gap-2"
        >
          {loading === "generate" ? <Loader2 className="w-4 h-4 animate-spin" /> : <Sparkles className="w-4 h-4" />}
          {loading === "generate" ? "Generating..." : "Generate"}
        </button>
        {hasContent && (
          <button
            onClick={onDeploy}
            disabled={loading !== null}
            className="px-4 py-2.5 rounded-lg bg-emerald-600 text-sm font-medium text-white disabled:opacity-50 transition-all flex items-center justify-center gap-2"
          >
            {loading === "deploy" ? <Loader2 className="w-4 h-4 animate-spin" /> : <Download className="w-4 h-4" />}
            {loading === "deploy" ? "Deploying..." : "Deploy"}
          </button>
        )}
      </div>
    </div>
  );
}

function PreviewPanel({ content }: { content: ContentGenerationResponse }) {
  return (
    <div className="bg-[#0d1117] rounded-xl border border-[#30363d] overflow-hidden">
      <div className="px-4 py-2 border-b border-[#30363d] flex items-center gap-2">
        <div className="flex gap-1.5">
          <div className="w-3 h-3 rounded-full bg-[#ff5f56]" />
          <div className="w-3 h-3 rounded-full bg-[#ffbd2e]" />
          <div className="w-3 h-3 rounded-full bg-[#27c93f]" />
        </div>
        <span className="text-[10px] text-gray-500 ml-2">README.md preview</span>
      </div>
      <div className="p-6 overflow-auto max-h-[600px]">
        <style>{`
          .preview-md { color: #c9d1d9; font-size: 15px; line-height: 1.6; }
          .preview-md h1 { font-size: 1.8em; border-bottom: 1px solid #21262d; padding-bottom: 0.3em; margin: 24px 0 16px; }
          .preview-md h2 { font-size: 1.4em; border-bottom: 1px solid #21262d; padding-bottom: 0.3em; margin: 20px 0 12px; }
          .preview-md h3 { font-size: 1.15em; margin: 16px 0 8px; }
          .preview-md a { color: #58a6ff; }
          .preview-md img { max-width: 100%; }
          .preview-md code { background: rgba(110,118,129,0.4); padding: 0.2em 0.4em; border-radius: 4px; font-size: 85%; }
          .preview-md pre { background: #161b22; border-radius: 6px; padding: 16px; overflow: auto; }
          .preview-md blockquote { border-left: 0.25em solid #3b434b; color: #8b949e; padding: 0 1em; }
          .preview-md ul, .preview-md ol { padding-left: 2em; }
          .preview-md li { margin-top: 0.25em; }
          .preview-md div[align="center"] { text-align: center; }
          .preview-md table { border-collapse: collapse; margin: 16px 0; }
          .preview-md th, .preview-md td { border: 1px solid #30363d; padding: 6px 13px; }
          .preview-md th { background: #161b22; }
        `}</style>
        <div className="preview-md">
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            rehypePlugins={[rehypeRaw]}
            components={{
              img: ({ ...props }) => <img {...props} style={{ display: "inline-block", margin: "4px" }} />,
            }}
          >
            {content.markdown}
          </ReactMarkdown>
        </div>
      </div>
    </div>
  );
}
