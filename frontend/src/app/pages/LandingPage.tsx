import { useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { api } from "@/lib/api";
import {
  GitBranch,
  Sparkles,
  Zap,
  FileText,
  Rocket,
  Code2,
  Star,
  ArrowRight,
  CheckCircle2,
  Shield,
  Globe,
} from "lucide-react";

const features = [
  {
    icon: <Sparkles className="w-6 h-6" />,
    title: "AI-Powered Generation",
    desc: "Gemini analyzes your code and writes a README that highlights your actual strengths.",
    color: "text-violet-400",
    bg: "bg-violet-500/10",
  },
  {
    icon: <Code2 className="w-6 h-6" />,
    title: "Deep Code Analysis",
    desc: "Scans repos for languages, frameworks, dependencies, and contribution patterns.",
    color: "text-cyan-400",
    bg: "bg-cyan-500/10",
  },
  {
    icon: <Zap className="w-6 h-6" />,
    title: "One-Click Deploy",
    desc: "Push your generated profile directly to GitHub. No copy-paste needed.",
    color: "text-amber-400",
    bg: "bg-amber-500/10",
  },
  {
    icon: <FileText className="w-6 h-6" />,
    title: "Live Preview",
    desc: "See your profile render in real-time before deploying.",
    color: "text-emerald-400",
    bg: "bg-emerald-500/10",
  },
  {
    icon: <Shield className="w-6 h-6" />,
    title: "Secure BYOK",
    desc: "Your Gemini API key stays in your browser. Never stored on our servers.",
    color: "text-rose-400",
    bg: "bg-rose-500/10",
  },
  {
    icon: <Globe className="w-6 h-6" />,
    title: "SVG Banner Generator",
    desc: "Premium animated hero banners with glassmorphism and typing effects.",
    color: "text-indigo-400",
    bg: "bg-indigo-500/10",
  },
];

const steps = [
  { num: "01", title: "Sign In", desc: "GitHub OAuth" },
  { num: "02", title: "Select Repos", desc: "Pick your best" },
  { num: "03", title: "Configure", desc: "Role, tone, skills" },
  { num: "04", title: "Generate", desc: "AI writes it" },
  { num: "05", title: "Deploy", desc: "One click" },
];

export default function LandingPage() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleGetStarted = async () => {
    setLoading(true);
    try {
      const { data } = await api.get("/auth/login");
      sessionStorage.setItem("oauth_state", data.state);
      window.location.href = data.auth_url;
    } catch {
      setLoading(false);
      setError("Failed to start authentication. Please try again.");
    }
  };

  return (
    <div className="min-h-screen bg-[#030712] text-white">
      {/* Background glows — pure CSS, no JS animation */}
      <div className="fixed inset-0 pointer-events-none overflow-hidden" aria-hidden="true">
        <div className="absolute -top-40 -left-40 w-[500px] h-[500px] rounded-full bg-violet-600/8 blur-[120px]" />
        <div className="absolute top-1/3 -right-32 w-[400px] h-[400px] rounded-full bg-cyan-500/6 blur-[100px]" />
        <div className="absolute -bottom-32 left-1/3 w-[350px] h-[350px] rounded-full bg-emerald-500/5 blur-[100px]" />
      </div>

      {/* Nav */}
      <header className="relative z-10 border-b border-white/5">
        <nav className="container mx-auto px-6 h-16 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2.5 group">
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-violet-500 to-cyan-500 blur-md opacity-30 group-hover:opacity-50 transition-opacity" />
              <GitBranch className="relative w-6 h-6 text-white" />
            </div>
            <span className="text-lg font-bold tracking-tight bg-gradient-to-r from-violet-400 via-cyan-400 to-emerald-400 bg-clip-text text-transparent">
              GitRight
            </span>
          </Link>
          <div className="flex items-center gap-3">
            <Link to="/dashboard">
              <Button variant="ghost" size="sm" className="text-gray-400 hover:text-white hover:bg-white/5">
                Dashboard
              </Button>
            </Link>
            <Button
              onClick={handleGetStarted}
              size="sm"
              disabled={loading}
              className="bg-gradient-to-r from-violet-600 to-cyan-600 hover:from-violet-500 hover:to-cyan-500 text-white font-medium"
            >
              {loading ? "Connecting..." : "Get Started"}
              {!loading && <ArrowRight className="w-3.5 h-3.5 ml-1" />}
            </Button>
          </div>
        </nav>
      </header>

      {error && (
        <div className="fixed top-4 right-4 z-50 px-4 py-3 rounded-lg text-sm font-medium bg-red-600 text-white shadow-lg">
          {error}
        </div>
      )}

      {/* Hero */}
      <section className="relative z-10 container mx-auto px-6 pt-24 pb-20 md:pt-36 md:pb-28">
        <div className="max-w-4xl mx-auto text-center">
          <Badge className="mb-6 bg-violet-500/10 text-violet-400 border-violet-500/20 px-3 py-1 text-xs font-medium">
            <Sparkles className="w-3 h-3 mr-1.5 inline" />
            AI-Powered Profile Generation
          </Badge>

          <h1 className="text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-bold tracking-tight leading-[1.1]">
            Your GitHub Profile,
            <br />
            <span className="bg-gradient-to-r from-violet-400 via-cyan-400 to-emerald-400 bg-clip-text text-transparent">
              Automated
            </span>
          </h1>

          <p className="mt-6 text-base md:text-lg text-gray-400 max-w-xl mx-auto leading-relaxed">
            Analyze your repositories, extract your skills, and generate a
            professional README — all in one click. BYOK, open source.
          </p>

          <div className="mt-8 flex flex-col sm:flex-row gap-3 justify-center">
            <Button
              onClick={handleGetStarted}
              disabled={loading}
              className="bg-gradient-to-r from-violet-600 to-cyan-600 hover:from-violet-500 hover:to-cyan-500 text-white font-medium px-8 py-6 text-base"
            >
              <GitBranch className="w-5 h-5 mr-2" />
              Create Your Profile
              <ArrowRight className="w-5 h-5 ml-2" />
            </Button>
            <a href="https://github.com/krauzX/gitright" target="_blank" rel="noopener noreferrer">
              <Button variant="outline" className="border-white/10 hover:bg-white/5 px-8 py-6 text-base">
                <Star className="w-5 h-5 mr-2" />
                Star on GitHub
              </Button>
            </a>
          </div>

          <div className="mt-10 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-sm text-gray-500">
            <span className="flex items-center gap-1.5">
              <CheckCircle2 className="w-3.5 h-3.5 text-emerald-500" />
              Free & Open Source
            </span>
            <span className="flex items-center gap-1.5">
              <CheckCircle2 className="w-3.5 h-3.5 text-emerald-500" />
              BYOK — Your Key
            </span>
            <span className="flex items-center gap-1.5">
              <CheckCircle2 className="w-3.5 h-3.5 text-emerald-500" />
              Deploy in 60s
            </span>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="relative z-10 container mx-auto px-6 py-16 md:py-24">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-12">
            <h2 className="text-2xl md:text-3xl font-bold tracking-tight">
              Built for{" "}
              <span className="bg-gradient-to-r from-violet-400 to-cyan-400 bg-clip-text text-transparent">
                Developers
              </span>
            </h2>
            <p className="mt-3 text-gray-400 max-w-lg mx-auto">
              Every feature designed to showcase your actual work, not generic templates.
            </p>
          </div>

          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {features.map((f, i) => (
              <Card
                key={i}
                className="group bg-white/[0.02] border-white/[0.06] hover:border-white/[0.12] transition-all duration-200 p-5"
              >
                <div className={`w-10 h-10 rounded-lg ${f.bg} flex items-center justify-center mb-3 ${f.color} group-hover:scale-110 transition-transform duration-200`}>
                  {f.icon}
                </div>
                <h3 className="text-sm font-semibold text-white/90 mb-1.5">{f.title}</h3>
                <p className="text-xs text-gray-400 leading-relaxed">{f.desc}</p>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section className="relative z-10 container mx-auto px-6 py-16 md:py-24">
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-12">
            <h2 className="text-2xl md:text-3xl font-bold tracking-tight">
              How It{" "}
              <span className="bg-gradient-to-r from-violet-400 to-cyan-400 bg-clip-text text-transparent">
                Works
              </span>
            </h2>
            <p className="mt-3 text-gray-400">Five steps to a professional profile</p>
          </div>
          <div className="grid grid-cols-2 sm:grid-cols-5 gap-6">
            {steps.map((s, i) => (
              <div key={i} className="text-center">
                <div className="w-11 h-11 mx-auto rounded-full bg-gradient-to-br from-violet-500/20 to-cyan-500/20 border border-white/10 flex items-center justify-center mb-3">
                  <span className="text-xs font-bold bg-gradient-to-r from-violet-400 to-cyan-400 bg-clip-text text-transparent">
                    {s.num}
                  </span>
                </div>
                <p className="text-sm font-semibold text-white/90">{s.title}</p>
                <p className="text-xs text-gray-500 mt-0.5">{s.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="relative z-10 container mx-auto px-6 py-16 md:py-24 pb-24">
        <div className="max-w-2xl mx-auto">
          <Card className="relative overflow-hidden bg-white/[0.02] border-white/[0.06] p-10 text-center">
            <div className="absolute inset-0 bg-gradient-to-br from-violet-500/5 via-transparent to-cyan-500/5 pointer-events-none" />
            <div className="relative">
              <div className="w-12 h-12 mx-auto rounded-2xl bg-gradient-to-br from-violet-500 to-cyan-500 flex items-center justify-center mb-5">
                <Rocket className="w-6 h-6 text-white" />
              </div>
              <h2 className="text-2xl md:text-3xl font-bold tracking-tight mb-3">
                Ready to Stand Out?
              </h2>
              <p className="text-gray-400 mb-6 max-w-md mx-auto">
                Join developers who use GitRight to create compelling profiles
                that attract recruiters and collaborators.
              </p>
              <Button
                onClick={handleGetStarted}
                disabled={loading}
                className="bg-gradient-to-r from-violet-600 to-cyan-600 hover:from-violet-500 hover:to-cyan-500 font-medium px-8 py-6 text-base"
              >
                <GitBranch className="w-5 h-5 mr-2" />
                Get Started Free
                <ArrowRight className="w-5 h-5 ml-2" />
              </Button>
            </div>
          </Card>
        </div>
      </section>

      {/* Footer */}
      <footer className="relative z-10 border-t border-white/5 py-6">
        <div className="container mx-auto px-6 flex flex-col sm:flex-row items-center justify-between text-xs text-gray-500 gap-2">
          <p>Built with Go, React, PostgreSQL & Gemini</p>
          <p>
            MIT License ·{" "}
            <a href="https://github.com/krauzX/gitright" className="text-gray-400 hover:text-white transition-colors">
              GitHub
            </a>
          </p>
        </div>
      </footer>
    </div>
  );
}
