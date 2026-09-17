import { useEffect, useState } from "react";
import { api } from "../api.ts";
import { Link } from "react-router";

function MyProjects() {
  type Project = {
    id: number;
    title: string;
    status: string;
    proposalCount: number;
    category: string;
    budgetMin: number;
    budgetMax: number;
    deadline: string;
  };
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function load() {
      setError("");

      try {
        const data = await api("/api/projects/mine");
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
      <ul>
        {projects.map((p) => (
          <li key={p.id}>
            <div>Title: {p.title}</div>
            <div>Status: {p.status}</div>
            <div>
              <span>{p.proposalCount} proposals</span>{" "}
              <Link to={`/projects/${p.id}/proposals`}>View proposals</Link>
            </div>
            <div>Category: {p.category}</div>
            <div>Budget Min: {p.budgetMin}</div>
            <div>Budget Max: {p.budgetMax}</div>
            <div>Deadline: {p.deadline}</div>
            <br />
          </li>
        ))}
      </ul>
    </>
  );
}

export default MyProjects;
