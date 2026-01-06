import SideNav from "@/components/dashboard/SideNav";
import Header from "@/components/dashboard/Header";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <div className="flex">
        <SideNav />
        <div className="flex-1">{children}</div>
      </div>
    </div>
  );
}
