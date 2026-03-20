"use client";

import { Database } from "lucide-react";
import EmptyState from "@/components/ui/EmptyState";

export default function DatabaseDesignPage() {
  return (
    <div className="h-full flex items-center justify-center">
      <EmptyState
        icon={Database}
        title="Database design"
        description="Design your database schema here. Smarty can help you define tables, relationships, and indexes."
        action={{ label: "Start with Smarty", href: "#" }}
      />
    </div>
  );
}
