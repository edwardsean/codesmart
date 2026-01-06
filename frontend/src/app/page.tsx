import Link from "next/link";
import AuthGate from "@/components/auth/AuthGate";

export default function Home() {
  return (
    <AuthGate>
      <nav className="navbar">
        <div className="nav-content">
          <div className="logo">CodeSmart</div>
          <div className="flex items-center gap-4">
            <Link
              href="/login"
              className="text-gray-600 hover:text-red-600 font-medium transition-colors"
            >
              Sign In
            </Link>
            <Link href="/signup">
              <button className="bg-red-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-red-700 transition-colors">
                Get Started
              </button>
            </Link>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <div className="page-container pt-20">
        <div className="hero animate-fade-in">
          {/* Main Heading */}
          <h1 className="hero-title">
            AI-Powered Code Reviews
            <br />
            <span className="text-red-600">Made Simple</span>
          </h1>

          {/* Subtitle */}
          <p className="hero-subtitle">
            Catch bugs instantly, collaborate in real-time, and ship code with
            confidence. Experience the future of collaborative development.
          </p>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center max-w-md mx-auto">
            <Link href="/signup" className="w-full sm:w-auto">
              <button className="btn-primary">Start Free Trial</button>
            </Link>
            <Link href="/login" className="w-full sm:w-auto">
              <button className="btn-secondary">Sign In</button>
            </Link>
          </div>
        </div>

        {/* Features Section */}
        <div className="feature-grid mt-16 max-w-6xl mx-auto px-4">
          <div className="feature-card">
            <div className="text-4xl mb-4 text-red-600">🤖</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-3">
              AI Analysis
            </h3>
            <p className="text-gray-600">
              Detect bugs, security issues, and performance problems instantly
              with advanced AI
            </p>
          </div>

          <div className="feature-card">
            <div className="text-4xl mb-4 text-red-600">👥</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-3">
              Real-time Collaboration
            </h3>
            <p className="text-gray-600">
              Video calls, live editing, and instant feedback all in one
              platform
            </p>
          </div>

          <div className="feature-card">
            <div className="text-4xl mb-4 text-red-600">📈</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-3">
              60% Faster Reviews
            </h3>
            <p className="text-gray-600">
              Reduce review time and catch 3x more issues than manual reviews
            </p>
          </div>
        </div>
      </div>
    </AuthGate>
  );
}
