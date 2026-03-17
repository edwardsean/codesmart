import Link from "next/link";
import { HomeButton } from "@/components/ui/Button";
import LoginForm from "@/components/auth/LoginForm";
import SignInForm from "@/components/auth/SignInForm";
import styles from "./Auth.module.css";

export default function AuthPage({
  greetings,
  login,
}: {
  greetings: string[];
  login: boolean;
}) {
  return (
    <div className={styles.container}>
      <div style={{ width: "100%", maxWidth: "28rem" }}>
        <div style={{ textAlign: "center", marginBottom: "24px" }}>
          <HomeButton />
        </div>

        <div className={styles.card}>
          <div style={{ textAlign: "center", marginBottom: "28px" }}>
            <h1 className={styles.heading}>{greetings[0]}</h1>
            <p className={styles.subheading}>{greetings[1]}</p>
          </div>

          {login ? (
            <LoginForm styles={styles} />
          ) : (
            <SignInForm styles={styles} />
          )}

          <div className={styles.footer}>
            {greetings[2]}{" "}
            {login ? (
              <Link href="/auth/signup" className={styles.footerLink}>
                Sign up for free
              </Link>
            ) : (
              <Link href="/auth/login" className={styles.footerLink}>
                Sign in here
              </Link>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
