import { Repository } from "@/types/entity";
import { getLanguageColor } from "@/lib/utils";
import Link from "next/link";
import { AlertCircle } from "lucide-react";
import { api } from "@/lib/api";

interface RepositoryCardProps {
  repository: Repository;
}

export default async function RepositoryWrapper() {
  try {
    // const cookieStore = cookies();
    // const access_token = (await cookieStore).get("access_token");

    // if (access_token) {
    //   const data = await api.get_repositories(access_token.value);

    //   return (
    //     <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
    //       {data.repositories.map((repo) => (
    //         <RepositoryCard key={repo.id} repository={repo} />
    //       ))}
    //     </div>
    //   );
    // } else {
    //   throw new Error("access token not available");
    // }
    const data = await api.get_repositories();
    console.error("data:", data);

    return (
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {data.repositories.map((repo) => (
          <RepositoryCard key={repo.id} repository={repo} />
        ))}
      </div>
    );
  } catch (err: any) {
    return <p>Error: {err.message}</p>;
  }
}

export function RepositoryCard({ repository }: RepositoryCardProps) {
  return (
    <Link
      href={`/dashboard/repositories/${repository.name}`}
      className="bg-white rounded-lg shadow border border-gray-200 p-6 hover:shadow-lg transition-shadow cursor-pointer"
    >
      {/* Header */}
      <div className="flex justify-between items-start mb-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-900 mb-1">
            {repository.name}
          </h3>
          <p className="text-sm text-gray-500">
            {repository.owner.login}/{repository.name}
          </p>
        </div>
        {/* <button
          onClick={handleExternalClick}
          className="text-gray-400 hover:text-gray-600 transition-colors"
          title="Open repository in GitHub"
        >
          <ExternalLink className="w-5 h-5" />
        </button> */}
      </div>

      {/* Description */}
      <p className="text-gray-600 text-sm mb-4 line-clamp-2">
        {repository.description}
      </p>

      {/* Language & Status */}
      <div className="flex items-center gap-3 mb-4">
        <div className="flex items-center gap-2">
          <div
            className={`w-3 h-3 rounded-full ${getLanguageColor(
              repository.language
            )}`}
          ></div>
          <span className="text-sm text-gray-600">{repository.language}</span>
        </div>
        <span
        //   className={`px-2 py-1 rounded-full text-xs font-medium ${
        //     repository.status === "active"
        //       ? "bg-green-100 text-green-800"
        //       : "bg-gray-100 text-gray-800"
        //   }`}
        >
          repo status
        </span>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 gap-4 mb-4">
        <div className="text-center">
          <p className="text-lg font-semibold text-gray-900">total reviews</p>
          <p className="text-xs text-gray-500">Reviews</p>
        </div>
        <div className="text-center">
          <p className="text-lg font-semibold text-gray-900">team members</p>
          <p className="text-xs text-gray-500">Members</p>
        </div>
      </div>

      {/* Footer */}
      <div className="flex justify-between items-center pt-4 border-t border-gray-100">
        <div className="flex items-center gap-2">
          <AlertCircle className="w-4 h-4 text-yellow-500" />
          <span className="text-sm text-gray-600">active issues</span>
        </div>
        <span className="text-xs text-gray-500">last activity</span>
      </div>
    </Link>
  );
}
