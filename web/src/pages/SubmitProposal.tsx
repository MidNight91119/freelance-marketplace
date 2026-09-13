import { useState } from "react";
import { api } from "../api";
import { useNavigate, useParams } from "react-router";

function SubmitProposal() {
  const [coverLetter, setCoverLetter] = useState("");
  const [proposedPrice, setProposedPrice] = useState(0);
  const [estimatedDuration, setEstimatedDuration] = useState(0);
  const [error, setErrror] = useState("");
  const { projectId } = useParams();
  const navigate = useNavigate();

  async function handleSubmitProposal(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();
    setErrror("");

    try {
      await api(`/api/projects/${projectId}/proposals`, {
        method: "POST",
        body: JSON.stringify({
          coverLetter: coverLetter,
          proposedPrice: proposedPrice,
          estimatedDuration: estimatedDuration,
        }),
      });
      navigate("/projects");
    } catch (err) {
      setErrror((err as Error).message);
    }
  }

  return (
    <>
      <h1>Create Proposal</h1>
      <form onSubmit={handleSubmitProposal}>
        <div>
          <textarea
            placeholder="cover letter"
            value={coverLetter}
            onChange={(e) => setCoverLetter(e.target.value)}
          />
        </div>
        <div>
          <input
            placeholder="proposed price"
            type="number"
            value={proposedPrice}
            onChange={(e) => setProposedPrice(e.target.valueAsNumber)}
          />
          <input
            placeholder="estimated duration"
            type="number"
            value={estimatedDuration}
            onChange={(e) => setEstimatedDuration(e.target.valueAsNumber)}
          />
        </div>
        <br />
        <div>
          <button>Submit</button>
        </div>
      </form>
      {error && <p>{error}</p>}
    </>
  );
}

export default SubmitProposal;
