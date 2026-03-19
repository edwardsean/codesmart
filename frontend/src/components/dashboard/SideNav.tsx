import NavLinks from "@/components/dashboard/Nav-links";

export default function SideNav() {
  return (
    <aside className="w-52 flex-shrink-0 border-r border-black/[0.06] dark:border-white/[0.06] bg-gray-50/50 dark:bg-zinc-900/50 min-h-[calc(100vh-65px)]">
      <nav className="flex flex-col gap-1 p-3">
        <NavLinks />
      </nav>
    </aside>
  );
}
