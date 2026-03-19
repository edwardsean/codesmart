import Link from "next/link";
import { HomeButton } from "@/components/ui/Button";
import LoginForm from "@/components/auth/LoginForm";
import SignInForm from "@/components/auth/SignInForm";
import GithubButton from "@/components/auth/GithubButton";

export default function AuthPage({
  greetings,
  login,
}: {
  greetings: string[];
  login: boolean;
}) {
  return (
    <div className="min-h-screen flex flex-col justify-center items-center px-4 py-8 bg-[#f5f3ee] dark:bg-[#141414]">
      <div className="w-full max-w-md">
        <div className="text-center mb-6">
          <HomeButton />
        </div>

        <div className="bg-[#f0ede8] dark:bg-[#1a1a1a] border border-black/[0.08] dark:border-white/[0.07] rounded-xl p-8">
          <div className="text-center mb-7">
            <h1 className="text-2xl font-serif font-normal text-gray-900 dark:text-zinc-100 tracking-tight mb-1">
              {greetings[0]}
            </h1>
            <p className="text-sm text-gray-500 dark:text-zinc-500">
              {greetings[1]}
            </p>
          </div>

          <GithubButton />

          <div className="relative my-5">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-black/[0.08] dark:border-white/[0.07]" />
            </div>
            <div className="relative flex justify-center">
              <span className="px-3 text-xs font-mono text-gray-400 dark:text-zinc-600 bg-[#f0ede8] dark:bg-[#1a1a1a] tracking-widest">
                or continue with email
              </span>
            </div>
          </div>

          {login ? <LoginForm /> : <SignInForm />}

          <p className="text-sm text-gray-500 dark:text-zinc-500 text-center mt-5">
            {greetings[2]}{" "}
            {login ? (
              <Link
                href="/auth/signup"
                className="text-[#dc503c] hover:opacity-80 font-medium transition-opacity"
              >
                Sign up for free
              </Link>
            ) : (
              <Link
                href="/auth/login"
                className="text-[#dc503c] hover:opacity-80 font-medium transition-opacity"
              >
                Sign in here
              </Link>
            )}
          </p>
        </div>
      </div>
    </div>
  );
}
