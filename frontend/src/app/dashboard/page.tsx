export default function Dashboard() {
  return (
    <>
      {true && (
        <div className="p-8">
          {/* Welcome Section */}
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-gray-900">
              Welcome back username!
            </h1>
            <p className="text-gray-600 mt-2">
              Ready to review some code? Lets get started.
            </p>
          </div>

          {/* <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            {projects.map((project) => (
              <ProjectCard key={project.id} project={project}></ProjectCard>
            ))}
          </div> */}
        </div>
      )}
    </>
  );
}
