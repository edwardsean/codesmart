import Link from "next/link";
import { HomeButton } from "@/components/ui/Button";
import LoginForm from "@/components/auth/LoginForm";
import SignInForm from "@/components/auth/SignInForm";

export default function AuthPage({
  greetings,
  login,
}: {
  greetings: string[];
  login: boolean;
}) {
  return (
    <div className="page-container">
      <div className="w-full max-w-md">
        {/* back button */}
        <div className="text-center mb-8">
          <HomeButton />
        </div>
        <div className="card animate-fade-in">
          <div className="text-center mb-8">
            <h1 className="text-2xl font-bold text-gray-900 mb-2">
              {greetings[0]}
            </h1>
            <p className="text-gray-600">{greetings[1]}</p>
          </div>

          {login ? <LoginForm /> : <SignInForm />}

          <div className="text-center mt-6">
            <p className="text-sm text-gray-600">
              {greetings[2]}{" "}
              {login ? (
                <Link
                  href="/signup"
                  className="text-red-600 hover:text-red-700 font-medium transition-colors"
                >
                  Sign up for free
                </Link>
              ) : (
                <Link
                  href="/login"
                  className="text-red-600 hover:text-red-700 font-medium transition-colors"
                >
                  Sign in here
                </Link>
              )}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
