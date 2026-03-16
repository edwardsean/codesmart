import AuthPage from "@/components/auth/AuthPage";

export default function LoginPage() {
  const greetings = [
    "Welcome Back",
    "Sign in to your account to continue",
    "Don't have an account?",
  ];
  return <AuthPage greetings={greetings} login={true} />;
}
