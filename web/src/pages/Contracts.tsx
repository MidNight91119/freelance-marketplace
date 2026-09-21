import { useEffect, useState } from "react";
import { api } from "../api.ts";

function Contracts() {
  type Contract = {
    id: number;
    projectTitle: string;
    clientName: string;
    freelancerName: string;
    amount: number;
    status: string;
    createdAt: string;
  };
  const [contracts, setContracts] = useState<Contract[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      setError("");

      try {
        const data = await api("/api/contracts");
        setContracts(data);
      } catch (err) {
        setError((err as Error).message);
      }
      setLoading(false);
    }
    load();
  }, []);

  return (
    <>
      <h1>Contracts</h1>
      {loading && <p className="muted">Loading...</p>}
      {error && <p className="error">{error}</p>}
      {!loading && !error && contracts.length === 0 && <p>No contracts</p>}

      <ul>
        {contracts.map((c) => (
          <li key={c.id}>
            <div>{c.projectTitle}</div>
            <div>{c.clientName}</div>
            <div>{c.freelancerName}</div>
            <div>{c.amount}</div>
            <div>{c.status}</div>
            <div>{c.createdAt}</div>
          </li>
        ))}
      </ul>
    </>
  );
}

export default Contracts;
