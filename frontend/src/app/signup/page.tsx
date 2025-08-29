import AuthPage from "@/components/auth/AuthPage";

export default function SignupPage() {
  const greetings = [
    "Create Your Account",
    "Get started with CodeReview AI today",
    "Already have an account?",
  ];
  return <AuthPage greetings={greetings} login={false} />;
}
