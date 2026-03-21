"use client";

import { useState, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useAuthStore } from "@/stores/authStore";
import {
  Github,
  Code2,
  ArrowLeft,
  ArrowRight,
  Check,
  Link as LinkIcon,
} from "lucide-react";
import Input from "@/components/ui/Input";
import RepositoryList from "@/components/github/RepositoryList";
import { LANGUAGES, DIFFICULTIES } from "@/types/entity";
import { ProjectStep } from "@/types/entity";
import {
  Difficulty,
  Language,
  ProjectMode,
  Repository,
  SourceType,
} from "@/types/entity";
import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import axios from "axios";
import { authService } from "@/services/authService";
import { githubService } from "@/services/githubService";
import { projectService } from "@/services/projectService";
import { STEPS } from "@/types/entity";

export type GithubInputMethod = "repos" | "url";

const STEP_TITLES: Record<ProjectStep, string> = {
  source: "Where does your project come from?",
  details: "Tell us about your project.",
  mode: "How do you want to learn?",
};

export default function NewProjectPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const apiPrivate = useAxiosPrivate();
  const apiAuth = authService(apiPrivate);
  const apiGithub = githubService(apiPrivate);
  const apiProject = projectService(apiPrivate);
  const { refreshUser, _hasHydrated, account } = useAuthStore();

  const [step, setStep] = useState<ProjectStep>("source");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  // source
  const [sourceType, setSourceType] = useState<SourceType>();
  const [githubMethod, setGithubMethod] = useState<GithubInputMethod>("repos");
  const [sourceURL, setSourceURL] = useState("");
  const [selectedRepo, setSelectedRepo] = useState<Repository | null>(null);

  // details
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [language, setLanguage] = useState<Language>();
  const [difficulty, setDifficulty] = useState<Difficulty>();

  // mode
  const [mode, setMode] = useState<ProjectMode>();

  // github repos
  const [repos, setRepos] = useState<Repository[]>([]);
  const [loadingRepos, setLoadingRepos] = useState(false);
  const [reposError, setReposError] = useState("");

  useEffect(() => {
    if (!_hasHydrated) return;
    if (!account?.accessToken) return;
    if (searchParams.get("github_connected") !== "true") return;

    router.replace("/dashboard/projects/new"); //cleanup the url
    const refreshActiveUser = async () => {
      try {
        const { user } = await apiAuth.me();

        refreshUser(user);
      } catch (err) {
        console.error("error when hitting me endpoint: ", err);
      }
    };

    refreshActiveUser();
  }, [_hasHydrated, account?.accessToken]);

  async function fetchRepos() {
    setLoadingRepos(true);
    setReposError("");
    try {
      const { repositories } = await apiGithub.getRepositories();
      setRepos(repositories ?? []);
    } catch {
      setReposError(
        "Could not load repositories. Make sure you logged in with GitHub.",
      );
    } finally {
      setLoadingRepos(false);
    }
  }

  function handleSourceSelect(type: "scratch" | "github") {
    setSourceType(type);
    if (type === "github" && repos.length === 0) fetchRepos();
  }

  function handleGithubMethodSwitch(method: GithubInputMethod) {
    setGithubMethod(method);
    if (method === "repos" && repos.length === 0) fetchRepos();
  }

  function handleRepoSelect(repo: Repository) {
    setSelectedRepo(repo);
    setSourceURL(repo.html_url);
    setTitle(repo.name);
    if (repo.language) setLanguage(repo.language.toLowerCase() as Language);
  }

  function canProceed() {
    if (step === "source") {
      if (!sourceType) return false;
      if (sourceType === "github") {
        if (githubMethod === "repos") return !!selectedRepo;
        if (githubMethod === "url")
          return sourceURL.startsWith("https://github.com/");
      }
      return true;
    }
    if (step === "details") return !!title && !!language;
    if (step === "mode") return !!mode;
    return false;
  }

  const stepIndex = STEPS.indexOf(step);

  function goBack() {
    if (step === "source") router.back();
    else setStep(STEPS[stepIndex - 1]);
  }

  function goNext() {
    setStep(STEPS[stepIndex + 1]);
  }

  async function handleSubmit() {
    setIsSubmitting(true);
    setError("");
    try {
      const { project } = await apiProject.createProject({
        title,
        description,
        language: language as Language,
        mode: mode as ProjectMode,
        source_type: sourceType as SourceType,
        source_url: sourceType === "github" ? sourceURL : undefined,
        difficulty: difficulty || undefined,
      });

      console.log("project: ", project);
      router.push(`/dashboard/projects/${project.id}`);
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.message ?? "Failed to create project");
      } else {
        console.log("error: ", err);
        setError("Something went wrong");
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div>
      {/* Back */}
      <button
        onClick={goBack}
        className="flex items-center gap-1.5 text-sm text-gray-400 dark:text-zinc-500 hover:text-gray-700 dark:hover:text-zinc-300 transition-colors mb-8"
      >
        <ArrowLeft size={14} />
        {step === "source" ? "Back to projects" : "Back"}
      </button>

      {/* Header */}
      <div className="mb-8">
        <p className="text-xs font-mono text-[#dc503c] tracking-widest uppercase mb-2">
          new project
        </p>
        <h1 className="text-3xl font-serif font-normal text-gray-900 dark:text-zinc-100 tracking-tight">
          {STEP_TITLES[step]}
        </h1>
      </div>

      {/* Step indicator */}
      <div className="flex items-center gap-2 mb-10">
        {STEPS.map((s, i) => (
          <div key={s} className="flex items-center gap-2">
            <div
              className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-mono transition-colors ${
                i < stepIndex
                  ? "bg-[#dc503c] text-white"
                  : i === stepIndex
                    ? "border-2 border-[#dc503c] text-[#dc503c]"
                    : "border border-gray-200 dark:border-zinc-700 text-gray-400 dark:text-zinc-600"
              }`}
            >
              {i < stepIndex ? <Check size={12} /> : i + 1}
            </div>
            {i < STEPS.length - 1 && (
              <div
                className={`h-px w-8 transition-colors ${i < stepIndex ? "bg-[#dc503c]" : "bg-gray-200 dark:bg-zinc-700"}`}
              />
            )}
          </div>
        ))}
      </div>

      {/* ── Step 1: Source ── */}
      {step === "source" && (
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-3">
            {[
              {
                type: "scratch" as const,
                Icon: Code2,
                label: "From scratch",
                sub: "Start a new idea with Smarty",
              },
              {
                type: "github" as const,
                Icon: Github,
                label: "From GitHub",
                sub: "Import an existing repository",
              },
            ].map(({ type, Icon, label, sub }) => (
              <button
                key={type}
                onClick={() => handleSourceSelect(type)}
                className={`p-5 rounded-xl border text-left transition-all ${
                  sourceType === type
                    ? "border-[#dc503c] bg-[#dc503c]/5"
                    : "border-black/[0.06] dark:border-white/[0.06] hover:border-black/10 dark:hover:border-white/10"
                }`}
              >
                <Icon
                  size={20}
                  className={
                    sourceType === type
                      ? "text-[#dc503c]"
                      : "text-gray-400 dark:text-zinc-500"
                  }
                />
                <p className="text-sm font-medium text-gray-900 dark:text-zinc-100 mt-3">
                  {label}
                </p>
                <p className="text-xs text-gray-500 dark:text-zinc-500 mt-1">
                  {sub}
                </p>
              </button>
            ))}
          </div>

          {sourceType === "github" && (
            <div className="space-y-3">
              {/* Toggle */}
              <div className="flex gap-1 p-1 bg-gray-100 dark:bg-zinc-800 rounded-lg w-fit">
                {[
                  {
                    method: "repos" as const,
                    Icon: Github,
                    label: "My repositories",
                  },
                  {
                    method: "url" as const,
                    Icon: LinkIcon,
                    label: "Paste URL",
                  },
                ].map(({ method, Icon, label }) => (
                  <button
                    key={method}
                    onClick={() => handleGithubMethodSwitch(method)}
                    className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium transition-colors ${
                      githubMethod === method
                        ? "bg-white dark:bg-zinc-700 text-gray-900 dark:text-zinc-100 shadow-sm"
                        : "text-gray-500 dark:text-zinc-400 hover:text-gray-700 dark:hover:text-zinc-300"
                    }`}
                  >
                    <Icon size={12} /> {label}
                  </button>
                ))}
              </div>

              {githubMethod === "repos" && (
                <RepositoryList
                  repos={repos}
                  loading={loadingRepos}
                  error={reposError}
                  selectedId={selectedRepo?.id}
                  onSelect={handleRepoSelect}
                />
              )}

              {githubMethod === "url" && (
                <div className="space-y-1">
                  <Input
                    label="Repository URL"
                    placeholder="https://github.com/username/repo"
                    value={sourceURL}
                    onChange={(e) => {
                      setSourceURL(e.target.value);
                      const parts = e.target.value.split("/");
                      if (parts.length >= 5) setTitle(parts[parts.length - 1]);
                    }}
                  />
                  {sourceURL &&
                    !sourceURL.startsWith("https://github.com/") && (
                      <p className="text-xs text-[#dc503c]">
                        Must be a valid GitHub URL
                      </p>
                    )}
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* ── Step 2: Details ── */}
      {step === "details" && (
        <div className="space-y-5">
          <Input
            label="Project title"
            placeholder="e.g. REST API with Go"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />

          <div className="flex flex-col gap-1.5">
            <label className="text-xs font-medium text-gray-500 dark:text-zinc-400 tracking-wide">
              Description{" "}
              <span className="text-gray-300 dark:text-zinc-600">
                (optional)
              </span>
            </label>
            <textarea
              rows={3}
              placeholder="What does this project do?"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-3 py-2 rounded-lg border text-sm outline-none transition-colors resize-none
                bg-white dark:bg-zinc-900 text-gray-900 dark:text-zinc-100
                placeholder:text-gray-400 dark:placeholder:text-zinc-600
                border-gray-200 dark:border-zinc-700 focus:border-gray-400 dark:focus:border-zinc-500"
            />
          </div>

          {/* Language pills */}
          <div className="flex flex-col gap-2">
            <label className="text-xs font-medium text-gray-500 dark:text-zinc-400 tracking-wide">
              Language
            </label>
            <div className="flex flex-wrap gap-2">
              {LANGUAGES.map((lang) => (
                <button
                  key={lang}
                  onClick={() => setLanguage(lang)}
                  className={`px-3 py-1.5 rounded-lg text-sm font-mono transition-colors ${
                    language === lang
                      ? "bg-[#dc503c] text-white"
                      : "border border-gray-200 dark:border-zinc-700 text-gray-600 dark:text-zinc-400 hover:border-gray-400 dark:hover:border-zinc-500"
                  }`}
                >
                  {lang}
                </button>
              ))}
            </div>
          </div>

          {/* Difficulty pills */}
          <div className="flex flex-col gap-2">
            <label className="text-xs font-medium text-gray-500 dark:text-zinc-400 tracking-wide">
              Difficulty{" "}
              <span className="text-gray-300 dark:text-zinc-600">
                (optional — Smarty can decide)
              </span>
            </label>
            <div className="flex gap-2">
              {DIFFICULTIES.map((d) => (
                <button
                  key={d}
                  onClick={() =>
                    setDifficulty(difficulty === d ? undefined : d)
                  }
                  className={`px-3 py-1.5 rounded-lg text-sm transition-colors ${
                    difficulty === d
                      ? "bg-[#dc503c] text-white"
                      : "border border-gray-200 dark:border-zinc-700 text-gray-600 dark:text-zinc-400 hover:border-gray-400 dark:hover:border-zinc-500"
                  }`}
                >
                  {d}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* ── Step 3: Mode ── */}
      {step === "mode" && (
        <div className="space-y-3">
          {[
            {
              value: "learn" as const,
              label: "Learn from scratch",
              description:
                "Smarty breaks your project into levels. You build it step by step with guided checkpoints and hints. Best for learning how to build something new.",
            },
            {
              value: "help" as const,
              label: "Help me finish or fix",
              description:
                "Smarty analyzes your existing code, suggests what to do next, detects issues, and explains design decisions. Best for working on a real project with AI guidance.",
            },
          ].map(({ value, label, description }) => (
            <button
              key={value}
              onClick={() => setMode(value)}
              className={`w-full p-6 rounded-xl border text-left transition-all ${
                mode === value
                  ? "border-[#dc503c] bg-[#dc503c]/5"
                  : "border-black/[0.06] dark:border-white/[0.06] hover:border-black/10 dark:hover:border-white/10"
              }`}
            >
              <p
                className={`text-sm font-medium mb-1.5 ${mode === value ? "text-[#dc503c]" : "text-gray-900 dark:text-zinc-100"}`}
              >
                {label}
              </p>
              <p className="text-xs text-gray-500 dark:text-zinc-500 leading-relaxed">
                {description}
              </p>
            </button>
          ))}
        </div>
      )}

      {error && (
        <p className="text-xs text-red-500 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2 mt-4">
          {error}
        </p>
      )}

      {/* Actions */}
      <div className="flex justify-end mt-8">
        {step !== "mode" ? (
          <button
            onClick={goNext}
            disabled={!canProceed()}
            className="flex items-center gap-2 px-5 py-2.5 rounded-lg text-sm font-medium transition-opacity
              bg-gray-900 dark:bg-zinc-100 text-white dark:text-zinc-900
              hover:opacity-85 disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Continue <ArrowRight size={14} />
          </button>
        ) : (
          <button
            onClick={handleSubmit}
            disabled={!canProceed() || isSubmitting}
            className="flex items-center gap-2 px-5 py-2.5 rounded-lg text-sm font-medium transition-opacity
              bg-[#dc503c] text-white hover:opacity-85 disabled:opacity-40 disabled:cursor-not-allowed"
          >
            {isSubmitting ? "Creating..." : "Create project"}{" "}
            {!isSubmitting && <ArrowRight size={14} />}
          </button>
        )}
      </div>
    </div>
  );
}
