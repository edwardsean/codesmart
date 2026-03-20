"use client";

import { Network } from "lucide-react";
import EmptyState from "@/components/ui/EmptyState";

export default function SystemDesignPage() {
  return (
    <div className="h-full flex items-center justify-center">
      <EmptyState
        icon={Network}
        title="System design"
        description="Design your system architecture here. Smarty can help you think through components, services, and how they connect."
        action={{ label: "Start with Smarty", href: "#" }}
      />
    </div>
  );
}
