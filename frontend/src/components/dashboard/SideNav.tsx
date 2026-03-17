import NavLinks from "@/components/dashboard/Nav-links";
import styles from "./SideNav.module.css";

export default function SideNav() {
  return (
    <aside className={styles.aside}>
      <nav className={styles.nav}>
        <NavLinks />
      </nav>
    </aside>
  );
}
