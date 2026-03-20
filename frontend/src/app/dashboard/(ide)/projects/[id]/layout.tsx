// app/dashboard/(ide)/projects/[id]/layout.tsx
import IDELayout from "@/components/ide/IDELayout";

export default async function ProjectIDELayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  return <IDELayout projectId={id}>{children}</IDELayout>;
}
