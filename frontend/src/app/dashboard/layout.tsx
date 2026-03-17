import Header from "@/components/ui/Header";
import SideNav from "@/components/dashboard/SideNav";
import styles from "./dashboard.module.css";

export default async function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className={styles.root}>
      <Header showAuth />
      <div className={styles.body}>
        <SideNav />
        <main className={styles.main}>{children}</main>
      </div>
    </div>
  );
}
