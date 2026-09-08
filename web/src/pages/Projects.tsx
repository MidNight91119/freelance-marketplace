import { useEffect, useState } from "react";

function Projects() {
  type Project = {
    id: number;
    title: string;
  };
  const [projects, setProjects] = useState<Project[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      setError("");

      const res = await fetch("http://localhost:8080/api/projects", {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      setLoading(false);

      const data = await res.json();
      if (!res.ok) {
        setError(data.message);
        console.log(data.message);
        return;
      }
      setProjects(data);
      console.log(data);
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
          <li key={p.id}>{p.title}</li>
        ))}
      </ul>
    </>
  );
}

export default Projects;
