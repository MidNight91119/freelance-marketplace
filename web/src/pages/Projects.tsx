import { useEffect, useState } from "react";
import { Link } from "react-router";
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

  useEffect(() => {
    async function load() {
      setError("");

      try {
        const data = await api("/api/projects");
        setProjects(data);
      } catch (err) {
        setError((err as Error).message);
      }
      setLoading(false);
    }
    load();
  }, []);

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
