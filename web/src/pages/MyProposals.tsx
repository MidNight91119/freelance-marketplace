import { useEffect, useState } from "react";
import { api } from "../api.ts";

function MyProposals() {
  type Proposal = {
    proposalId: number;
    projectTitle: string;
    coverLetter: string;
    proposedPrice: number;
    estimatedDuration: number;
    status: string;
    createdAt: string;
  };

  const [proposals, setProposals] = useState<Proposal[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function load() {
      setError("");

      try {
        const data = await api("/api/proposals/mine");
        setProposals(data);
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
      {!loading && !error && proposals.length === 0 && <p>No proposals</p>}
      <ul>
        {proposals.map((pr) => (
          <li key={pr.proposalId}>
            <div>Project Title: {pr.projectTitle}</div>
            <div>Proposed price: {pr.proposedPrice}</div>
            <div>Estimated Duration: {pr.estimatedDuration} days</div>
            <div>Status: {pr.status}</div>
            <div>Created at: {pr.createdAt}</div>
            <div>Cover Letter: {pr.coverLetter}</div>
            <br />
          </li>
        ))}
      </ul>
    </>
  );
}

export default MyProposals;
