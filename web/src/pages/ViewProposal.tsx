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
  const [refresh, setRefresh] = useState(0);

  async function handleAccept(proposalId: number) {
    setError("");

    try {
      await api(`/api/proposals/${proposalId}/accept`, {
        method: "PUT",
      });
      setRefresh((r) => r + 1);
    } catch (err) {
      setError((err as Error).message);
    }
  }

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
  }, [projectId, refresh]);

  return (
    <>
      {loading && <p className="muted">Loading...</p>}
      {error && <p className="error">{error}</p>}
      {!loading && !error && proposals.length === 0 && <p>No proposals</p>}
      <ul>
        {proposals.map((p) => (
          <li key={p.proposalId}>
            <div>Name: {p.freelancerName}</div>
            <div>Proposed price: {p.proposedPrice}</div>
            <div>Estimated Duration: {p.estimatedDuration} days</div>
            <div>Status: {p.status}</div>
            {p.status === "pending" && (
              <button onClick={() => handleAccept(p.proposalId)}>Accept</button>
            )}
          </li>
        ))}
      </ul>
    </>
  );
}

export default ViewProposal;
