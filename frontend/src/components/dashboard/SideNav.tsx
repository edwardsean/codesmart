import Link from "next/link";
import NavLinks from "@/components/dashboard/Nav-links";

export default function SideNav() {
  return (
    <div className="w-64 bg-white shadow-sm border-r border-gray-200 min-h-screen">
      <nav className="mt-8">
        <div className="px-4 space-y-2">
          <NavLinks />
        </div>
      </nav>
    </div>
  );
}
