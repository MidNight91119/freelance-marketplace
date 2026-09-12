import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router";

function Projects() {
  type Project = {
    id: number;
    title: string;
  };
  const [projects, setProjects] = useState<Project[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

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

        if (res.status == 401) {
          localStorage.removeItem("token");
          navigate("/login");
        }

        return;
      }
      setProjects(data);
      console.log(data);
    }
    load();
  }, [navigate]);

  return (
    <>
      {loading && <p>Loading...</p>}
      {error && <p>{error}</p>}
      {!loading && !error && projects.length === 0 && <p>No projects</p>}
      <div>
        <Link to="/projects/new">Create new project</Link>
      </div>
      <ul>
        {projects.map((p) => (
          <li key={p.id}>{p.title}</li>
        ))}
      </ul>
    </>
  );
}

export default Projects;
