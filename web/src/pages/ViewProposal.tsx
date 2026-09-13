import { useEffect, useState } from "react";
import { useParams } from "react-router";
import { api } from "../api.ts";

function ViewProposal() {
  type Proposal = {
    proposalId: number;
    freelancerName: string;
    proposedPrice: number;
    estimatedDuration: number;
    status: string;
  };
  const [proposals, setProposals] = useState<Proposal[]>([]);
  const [error, setError] = useState("");
  const { projectId } = useParams();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      setError("");

      try {
        const data = await api(`/api/projects/${projectId}/proposals`);
        setProposals(data);
      } catch (err) {
        setError((err as Error).message);
      }
      setLoading(false);
    }
    load();
  }, [projectId]);

  return (
    <>
      {loading && <p>Loading...</p>}
      {error && <p>{error}</p>}
      {!loading && !error && proposals.length === 0 && <p>No proposals</p>}
      <ul>
        {proposals.map((p) => (
          <li key={p.proposalId}>
            <div>Name: {p.freelancerName}</div>
            <div>Proposed price: {p.proposedPrice}</div>
            <div>Estimated Duration: {p.estimatedDuration} days</div>
            <div>Status: {p.status}</div>
            <br />
          </li>
        ))}
      </ul>
    </>
  );
}

export default ViewProposal;
