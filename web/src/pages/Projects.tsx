import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router";
import { api } from "../api.ts";

function Projects() {
  type Project = {
    id: number;
    title: string;
  };
  const [projects, setProjects] = useState<Project[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const role = localStorage.getItem("role");
  const [searchParams, setSearchParams] = useSearchParams();
  const [category, setCategory] = useState(searchParams.get("category") ?? "");
  const [minBudget, setMinBudget] = useState(
    searchParams.get("minBudget") ?? "",
  );
  const [maxBudget, setMaxBudget] = useState(
    searchParams.get("maxBudget") ?? "",
  );

  function handleFilter(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();
    setSearchParams({ category, minBudget, maxBudget });
  }

  useEffect(() => {
    async function load() {
      setError("");

      try {
        const data = await api(`/api/projects?${searchParams.toString()}`);
        setProjects(data);
      } catch (err) {
        setError((err as Error).message);
      }
      setLoading(false);
    }
    load();
  }, [searchParams]);

  return (
    <>
      {loading && <p>Loading...</p>}
      {error && <p>{error}</p>}
      {!loading && !error && projects.length === 0 && <p>No projects</p>}
      <div>
        {role === "client" && (
          <Link to="/projects/new">Create new project</Link>
        )}
      </div>
      <div>
        <h3>Filter projects</h3>
        <form onSubmit={handleFilter}>
          <input
            type="text"
            placeholder="category"
            value={category}
            onChange={(e) => setCategory(e.target.value)}
          />
          <input
            type="text"
            placeholder="minBudget"
            value={minBudget}
            onChange={(e) => setMinBudget(e.target.value)}
          />
          <input
            type="text"
            placeholder="maxBudget"
            value={maxBudget}
            onChange={(e) => setMaxBudget(e.target.value)}
          />
          <button>Submit</button>
        </form>
      </div>
      <ul>
        {projects.map((p) => (
          <li key={p.id}>
            <div>{p.title}</div>
            {role === "freelancer" && (
              <Link to={`/projects/${p.id}/propose`}>Propose</Link>
            )}
            <div>
              {role === "client" && (
                <Link to={`/projects/${p.id}/proposals`}>View Proposals</Link>
              )}
            </div>
            <br />
          </li>
        ))}
      </ul>
    </>
  );
}

export default Projects;
