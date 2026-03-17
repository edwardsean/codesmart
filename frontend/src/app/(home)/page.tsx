import Link from "next/link";
import styles from "./home.module.css";
import Header from "@/components/ui/Header";

export default function Home() {
  return (
    <div className={styles.root}>
      <Header showNav />

      <section className={styles.hero}>
        <p className={styles.eyebrow}>AI-powered learning</p>
        <h1 className={styles.h1}>
          Learn to code
          <br />
          <em>the right way.</em>
        </h1>
        <p className={styles.sub}>
          Turn any project or repository into a structured learning path. Real
          code, real feedback, real progress.
        </p>
        <div className={styles.actions}>
          <Link href="/auth/signup" className={styles.btnPrimary}>
            Start for free
          </Link>
          <Link href="/auth/login" className={styles.btnGhost}>
            sign in →
          </Link>
        </div>
      </section>

      <div className={styles.divider} />

      <div className={styles.features}>
        <div className={styles.feature}>
          <p className={styles.featureNum}>01</p>
          <svg className={styles.featureIcon} viewBox="0 0 36 36" fill="none">
            <rect
              x="4"
              y="8"
              width="28"
              height="20"
              rx="3"
              stroke="currentColor"
              strokeWidth="1.5"
            />
            <path
              d="M10 16l4 4-4 4"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
            <path
              d="M18 24h8"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
            />
          </svg>
          <h3>Project analysis</h3>
          <p>
            Paste a GitHub repo or describe an idea — Smarty breaks it into
            levels you can actually learn from.
          </p>
        </div>

        <div className={styles.feature}>
          <p className={styles.featureNum}>02</p>
          <svg className={styles.featureIcon} viewBox="0 0 36 36" fill="none">
            <circle
              cx="18"
              cy="18"
              r="12"
              stroke="currentColor"
              strokeWidth="1.5"
            />
            <path
              d="M18 12v6l4 2"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
            <path
              d="M8 18h2M26 18h2M18 8v2M18 26v2"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
            />
          </svg>
          <h3>Guided progress</h3>
          <p>
            Level-based structure with hints, AI chat, and test validation so
            you always know what to do next.
          </p>
        </div>

        <div className={styles.feature}>
          <p className={styles.featureNum}>03</p>
          <svg className={styles.featureIcon} viewBox="0 0 36 36" fill="none">
            <path
              d="M6 10h24M6 18h16M6 26h20"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
            />
            <circle
              cx="28"
              cy="26"
              r="4"
              stroke="currentColor"
              strokeWidth="1.5"
            />
            <path
              d="M26.5 26l1 1 2-2"
              stroke="currentColor"
              strokeWidth="1.2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
          <h3>Run &amp; validate</h3>
          <p>
            Write and execute code directly in the browser. Tests run instantly
            — no setup, no friction.
          </p>
        </div>
      </div>

      <footer className={styles.footer}>
        <span className={styles.footerCopy}>© 2026 codesmart</span>
        <span className={styles.footerCopy}>
          built for developers who learn by doing
        </span>
      </footer>
    </div>
  );
}
