import { Suspense } from "react";
import RepositoryWrapper from "@/components/repositories/RepositoryWrapper";
import Loading from "@/components/ui/Loading";

export default function Repositories() {
  return (
    <>
      {/* stat cards (how many repository) */}

      {/* this is for the repository wrapper */}
      <Suspense fallback={<Loading />}>
        <RepositoryWrapper />
      </Suspense>
    </>
  );
}
