import { useState } from "react";
import { useNavigate } from "react-router";
import { api } from "../api";

function CreateProject() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [category, setCategory] = useState("");
  const [deadline, setDeadline] = useState("");
  const [budgetMin, setBudgetMin] = useState(0);
  const [budgetMax, setBudgetMax] = useState(0);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  async function handleCreateProject(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");

    try {
      await api("/api/projects", {
        method: "POST",
        body: JSON.stringify({
          title: title,
          description: description,
          category: category,
          deadline: deadline,
          budgetMin: budgetMin,
          budgetMax: budgetMax,
        }),
      });
      navigate("/projects");
      //   console.log(data);
    } catch (err) {
      setError((err as Error).message);
    }
  }

  return (
    <>
      <h3>Create Project</h3>
      <form onSubmit={handleCreateProject}>
        <div>
          <input
            placeholder="title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
        </div>
        <div>
          <textarea
            placeholder="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </div>
        <div>
          <input
            placeholder="category"
            type="text"
            value={category}
            onChange={(e) => setCategory(e.target.value)}
          />
        </div>
        <div>
          <input
            placeholder="deadline"
            type="date"
            value={deadline}
            onChange={(e) => setDeadline(e.target.value)}
          />
        </div>
        <div>
          <input
            placeholder="min budget"
            type="number"
            value={budgetMin}
            onChange={(e) => setBudgetMin(e.target.valueAsNumber)}
          />
          <input
            placeholder="max budget"
            type="number"
            value={budgetMax}
            onChange={(e) => setBudgetMax(e.target.valueAsNumber)}
          />
        </div>
        <div>
          <button>Create project</button>
        </div>
      </form>
      {error && <p className="error">{error}</p>}
    </>
  );
}

export default CreateProject;
