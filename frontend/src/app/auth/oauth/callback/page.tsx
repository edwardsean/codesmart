import { Suspense } from "react";
import LoadingScreen from "@/components/ui/Loading";
import OAuthCallbackContent from "@/app/auth/oauth/callback/OAuthCallbackContent";

export default async function OAuthCallbackPage({
  searchParams,
}: {
  searchParams: Promise<{ code?: string; redirect?: string }>;
}) {
  const params = await searchParams;

  return (
    //we use suspense, because we are using useSearchParams in the content, and without suspense,
    //nextjs cant prerender that part of the apge bcs it doesnt know the search params until the client loads
    //so suspense tells prerender everything outside of this suspense first, this code inside will be filled
    //when ready
    <Suspense fallback={<LoadingScreen message="Signing you in..." />}>
      <OAuthCallbackContent
        code={params.code}
        redirect={
          params.redirect ? decodeURIComponent(params.redirect) : "/dashboard"
        }
      />
    </Suspense>
  );
}
