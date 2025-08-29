import AuthPage from "@/components/auth/AuthPage";

export default function LoginPage() {
  const greetings = [
    "Welcome Back",
    "Sign in to your account to continue",
    "Don't have an account?",
  ];
  return (
    <AuthPage greetings={greetings} login={true} />
    // <div className="page-container">
    //   <div className="w-full max-w-md">
    //     <div className="text-center mb-8">
    //       <Link
    //         href="/"
    //         className="text-gray-600 hover:text-red-600 text-sm font-medium transition-colors"
    //       >
    //         ← Back to Home
    //       </Link>
    //     </div>

    //     <div className="card animate-fade-in">
    //       <div className="text-center mb-8">
    //         <h1 className="text-2xl font-bold text-gray-900 mb-2">
    //           Welcome Back
    //         </h1>
    //         <p className="text-gray-600">Sign in to your account to continue</p>
    //       </div>

    //       <LoginForm />

    //       <div className="text-center mt-6">
    //         <p className="text-sm text-gray-600">
    //           Dont have an account?{" "}
    //           <Link
    //             href="/signup"
    //             className="text-red-600 hover:text-red-700 font-medium transition-colors"
    //           >
    //             Sign up for free
    //           </Link>
    //         </p>
    //       </div>
    //     </div>
    //   </div>
    // </div>
  );
}
