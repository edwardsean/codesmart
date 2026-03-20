import Header from "@/components/ui/Header";
import SideNav from "@/components/dashboard/SideNav";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-white dark:bg-[#141414]">
      <Header showAuth />
      <div className="flex pt-[65px]">
        <SideNav />
        <main className="flex-1 min-h-[calc(100vh-65px)] py-8 px-12">
          <div className="max-w-3xl mx-auto">{children}</div>
        </main>
      </div>
    </div>
  );
}
