import { useEffect, useRef } from "react";
import { Link } from "react-router-dom";
import { gsap } from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { api } from "@/lib/api";
import {
  Github,
  Sparkles,
  Zap,
  FileText,
  Rocket,
  Code2,
  Star,
  ArrowRight,
  CheckCircle2,
} from "lucide-react";

gsap.registerPlugin(ScrollTrigger);

export default function LandingPage() {
  const heroRef = useRef<HTMLDivElement>(null);
  const featuresRef = useRef<HTMLDivElement>(null);
  const ctaRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const ctx = gsap.context(() => {
      // Hero animations
      gsap.from(".hero-title", {
        opacity: 0,
        y: 50,
        duration: 1,
        ease: "power4.out",
      });

      gsap.from(".hero-subtitle", {
        opacity: 0,
        y: 30,
        duration: 1,
        delay: 0.3,
        ease: "power4.out",
      });

      gsap.from(".hero-buttons", {
        opacity: 0,
        y: 30,
        duration: 1,
        delay: 0.6,
        ease: "power4.out",
      });

      // Floating animation for decorative elements
      gsap.to(".float", {
        y: -20,
        duration: 2,
        repeat: -1,
        yoyo: true,
        ease: "power1.inOut",
        stagger: 0.2,
      });

      // Feature cards animation on scroll
      gsap.from(".feature-card", {
        scrollTrigger: {
          trigger: featuresRef.current,
          start: "top 80%",
        },
        opacity: 1,
        y: 50,
        duration: 0.8,
        stagger: 0.2,
        ease: "power3.out",
      });

      // CTA section animation
      gsap.from(".cta-content", {
        scrollTrigger: {
          trigger: ctaRef.current,
          start: "top 80%",
        },
        opacity: 1,
        scale: 0.9,
        duration: 1,
        ease: "back.out(1.7)",
      });
    });

    return () => ctx.revert();
  }, []);

  const handleGetStarted = async () => {
    try {
      // Call backend to get OAuth URL with server-generated state
      const response = await api.get("/auth/login");
      const { auth_url, state } = response.data;

      // Store state for validation on callback
      sessionStorage.setItem("oauth_state", state);

      // Redirect to GitHub OAuth
      window.location.href = auth_url;
    } catch (error) {
      console.error("Failed to initiate OAuth:", error);
      alert("Failed to start authentication. Please try again.");
    }
  };

  const features = [
    {
      id: "ai-powered",
      icon: <Sparkles className="w-8 h-8" />,
      title: "AI-Powered Generation",
      description:
        "Leverages Google Gemini Pro to create compelling, personalized README profiles that showcase your best work.",
      gradient: "from-purple-500 to-pink-500",
    },
    {
      id: "smart-analysis",
      icon: <Code2 className="w-8 h-8" />,
      title: "Smart Analysis",
      description:
        "Automatically analyzes your repositories, extracting languages, frameworks, and architectural patterns.",
      gradient: "from-blue-500 to-cyan-500",
    },
    {
      id: "instant-deploy",
      icon: <Zap className="w-8 h-8" />,
      title: "Instant Deployment",
      description:
        "One-click deployment to your GitHub profile. No manual copying, no hassle—just results.",
      gradient: "from-orange-500 to-red-500",
    },
    {
      id: "live-preview",
      icon: <FileText className="w-8 h-8" />,
      title: "Live Preview",
      description:
        "See your profile in real-time as you configure tone, role, and skills emphasis before deploying.",
      gradient: "from-green-500 to-emerald-500",
    },
    {
      id: "performance",
      icon: <Rocket className="w-8 h-8" />,
      title: "Performance Optimized",
      description:
        "Built with Go backend and React 19 frontend for blazing-fast performance and seamless UX.",
      gradient: "from-violet-500 to-purple-500",
    },
    {
      id: "production-ready",
      icon: <Star className="w-8 h-8" />,
      title: "Production Ready",
      description:
        "Enterprise-grade architecture with PostgreSQL, Redis caching, and OAuth2 security.",
      gradient: "from-yellow-500 to-orange-500",
    },
  ];

  const benefits = [
    "Generate professional README in under 60 seconds",
    "Customize tone and style for different roles",
    "Automatic skill extraction from your code",
    "Beautiful badges and visual elements",
    "SEO-optimized markdown formatting",
    "Secure GitHub OAuth integration",
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950 text-white overflow-hidden">
      {/* Animated background effects */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-20 left-10 w-72 h-72 bg-purple-500/20 rounded-full blur-3xl float"></div>
        <div className="absolute top-40 right-20 w-96 h-96 bg-blue-500/20 rounded-full blur-3xl float"></div>
        <div className="absolute bottom-20 left-1/3 w-80 h-80 bg-pink-500/20 rounded-full blur-3xl float"></div>
      </div>

      {/* Navigation */}
      <nav className="relative z-10 container mx-auto px-6 py-6 flex items-center justify-between">
        <div className="flex items-center space-x-3 group cursor-pointer">
          <div className="relative">
            <div className="absolute inset-0 bg-gradient-to-r from-primary to-purple-500 blur-lg opacity-50 group-hover:opacity-75 transition-opacity"></div>
            <Github className="relative w-8 h-8 text-white" />
          </div>
          <span className="text-2xl font-bold bg-gradient-to-r from-primary via-purple-400 to-pink-400 bg-clip-text text-transparent">
            GitRight
          </span>
        </div>
        <div className="flex items-center space-x-4">
          <Link to="/dashboard">
            <Button
              variant="ghost"
              className="text-slate-300 hover:text-white hover:bg-white/10 transition-all duration-300"
            >
              Dashboard
            </Button>
          </Link>
          <Button
            onClick={handleGetStarted}
            className="relative overflow-hidden bg-gradient-to-r from-primary via-purple-500 to-pink-500 hover:shadow-lg hover:shadow-primary/50 transition-all duration-300 hover:scale-105"
          >
            <Github className="w-4 h-4 mr-2" />
            Get Started
          </Button>
        </div>
      </nav>

      {/* Hero Section */}
      <section
        ref={heroRef}
        className="relative z-10 container mx-auto px-6 py-20 md:py-32"
      >
        <div className="max-w-5xl mx-auto text-center space-y-8">
          <Badge className="hero-title bg-gradient-to-r from-primary/20 to-purple-500/20 text-primary border-primary/40 px-5 py-2 backdrop-blur-sm shadow-lg shadow-primary/10">
            <Sparkles className="w-4 h-4 mr-2 inline animate-pulse" />
            AI-Powered Profile Generation
          </Badge>

          <h1 className="hero-title text-5xl md:text-7xl font-extrabold leading-tight">
            Transform Your GitHub Profile
            <br />
            <span className="bg-gradient-to-r from-primary via-purple-400 to-pink-400 bg-clip-text text-transparent">
              In Seconds
            </span>
          </h1>

          <p className="hero-subtitle text-xl md:text-2xl text-slate-300 max-w-3xl mx-auto">
            GitRight uses advanced AI to analyze your repositories and generate
            stunning, professional README profiles that make you stand out.
          </p>

          <div className="hero-buttons flex flex-col sm:flex-row gap-4 justify-center items-center pt-6">
            <Button
              onClick={handleGetStarted}
              size="lg"
              className="bg-gradient-to-r from-primary to-purple-500 hover:opacity-90 text-lg px-8 py-6"
            >
              <Github className="w-5 h-5 mr-2" />
              Create Your Profile
              <ArrowRight className="w-5 h-5 ml-2" />
            </Button>
            <Button
              size="lg"
              variant="outline"
              className=" from-primary to-purple-500 bg-gradient-to-r border-slate-700  hover: text-white bg-inherit-800 text-lg px-8 py-6"
            >
              Watch Demo
            </Button>
          </div>

          <div className="flex flex-wrap items-center justify-center gap-8 pt-8 text-sm text-slate-400">
            <div className="flex items-center gap-2">
              <CheckCircle2 className="w-5 h-5 text-green-400" />
              Free & Open Source
            </div>
            <div className="flex items-center gap-2">
              <CheckCircle2 className="w-5 h-5 text-green-400" />
              No Credit Card
            </div>
            <div className="flex items-center gap-2">
              <CheckCircle2 className="w-5 h-5 text-green-400" />
              60s Setup
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section
        ref={featuresRef}
        className="relative z-10 container mx-auto px-6 py-20"
      >
        <div className="max-w-6xl mx-auto space-y-12">
          <div className="text-center space-y-4">
            <h2 className="text-4xl md:text-5xl font-bold">
              Powerful Features for
              <span className="bg-gradient-to-r from-primary to-purple-400 bg-clip-text text-transparent">
                {" "}
                Developers
              </span>
            </h2>
            <p className="text-xl text-slate-400 max-w-2xl mx-auto">
              Everything you need to create an impressive GitHub profile that
              attracts recruiters and collaborators.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((feature) => (
              <Card
                key={feature.id}
                className="feature-card group relative text-white  bg-neutral-950/80 border border-slate-800 hover:border-fuchsia-500/50 transition-all duration-300 hover:scale-[1.02] hover:shadow-xl hover:shadow-purple-500/30 backdrop-blur-md p-6 space-y-4 overflow-hidden"
              >
                {/* Gradient Overlay for Purple Light Source */}
                <div
                  className="absolute inset-0 bg-gradient-to-tr from-fuchsia-900/10 to-transparent 
    group-hover:from-fuchsia-800/20 group-hover:to-transparent 
    transition-opacity duration-300"
                ></div>

                <div className="relative">
                  <div
                    // Bright Purple Icon
                    className={`w-16 h-16 rounded-2xl bg-gradient-to-br from-fuchsia-500 to-violet-400 
      flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300`}
                  >
                    {feature.icon}
                  </div>
                  <h3
                    // Title Brightens to match the aura
                    className="text-xl font-semibold mt-4 group-hover:text-fuchsia-400 transition-colors duration-300"
                  >
                    {feature.title}
                  </h3>
                  <p
                    // Text color adjusted for better visibility on near-black
                    className="text-slate-400 group-hover:text-slate-200 transition-colors duration-300"
                  >
                    {feature.description}
                  </p>
                </div>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Benefits Section */}
      <section className="relative z-10 container mx-auto px-6 py-20">
        <div className="max-w-4xl mx-auto">
          <Card className="relative overflow-hidden bg-gradient-to-br from-slate-900/95 via-slate-800/95 to-slate-900/95 border-slate-700/50 backdrop-blur-xl p-8 md:p-12 shadow-2xl">
            <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-purple-500/5"></div>
            <div className="relative">
              <h2 className="text-3xl md:text-4xl font-bold mb-8 text-center bg-gradient-to-r from-white to-slate-300 bg-clip-text text-transparent">
                Why Choose GitRight?
              </h2>
              <div className="grid md:grid-cols-2 gap-4">
                {benefits.map((benefit) => (
                  <div
                    key={benefit}
                    className="flex items-start gap-3 group hover:translate-x-2 transition-transform duration-300"
                  >
                    <CheckCircle2 className="w-6 h-6 text-green-400 flex-shrink-0 mt-0.5 group-hover:scale-110 transition-transform duration-300" />
                    <span className="text-slate-300 group-hover:text-white transition-colors duration-300">
                      {benefit}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </Card>
        </div>
      </section>

      {/* CTA Section */}
      <section
        ref={ctaRef}
        className="relative z-10 container mx-auto px-6 py-20 pb-32"
      >
        <div className="cta-content max-w-4xl mx-auto text-center space-y-8">
          <div className="relative">
            <div className="absolute inset-0 bg-gradient-to-r from-primary/30 via-purple-500/30 to-pink-500/30 blur-3xl animate-pulse"></div>
            <Card className="relative overflow-hidden bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 border-primary/30 p-12 shadow-2xl shadow-primary/20">
              <div className="absolute inset-0 bg-gradient-to-br from-primary/10 via-transparent to-purple-500/10"></div>
              <div className="relative">
                <div className="inline-block mb-4">
                  <Rocket className="w-16 h-16 text-primary animate-bounce" />
                </div>
                <h2 className="text-4xl md:text-5xl font-bold mb-6 bg-gradient-to-r from-white via-slate-100 to-slate-300 bg-clip-text text-transparent">
                  Ready to Stand Out?
                </h2>
                <p className="text-xl text-slate-300 mb-8 max-w-2xl mx-auto leading-relaxed">
                  Join thousands of developers who have transformed their GitHub
                  profiles with GitRight. Start building your professional
                  presence today.
                </p>
                <Button
                  onClick={handleGetStarted}
                  size="lg"
                  className="group relative overflow-hidden bg-gradient-to-r from-primary via-purple-500 to-pink-500 hover:shadow-2xl hover:shadow-primary/50 transition-all duration-500 text-lg px-12 py-7 hover:scale-110"
                >
                  <span className="absolute inset-0 bg-gradient-to-r from-pink-500 via-purple-500 to-primary opacity-0 group-hover:opacity-100 transition-opacity duration-500"></span>
                  <span className="relative flex items-center">
                    <Github className="w-5 h-5 mr-2" />
                    Create Your Profile Now
                    <ArrowRight className="w-5 h-5 ml-2 group-hover:translate-x-2 transition-transform duration-300" />
                  </span>
                </Button>
              </div>
            </Card>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="relative z-10 border-t border-slate-800/50 bg-slate-950/50 backdrop-blur-sm py-8">
        <div className="container mx-auto px-6 text-center text-slate-400">
          <p className="text-base">
            Built with <span className="text-red-400 animate-pulse">❤️</span>{" "}
            using Go, React 19, PostgreSQL, and Google Gemini
          </p>
          <p className="mt-2 text-sm text-slate-500">
            © 2025 GitRight. Open Source under MIT License.
          </p>
        </div>
      </footer>
    </div>
  );
}
