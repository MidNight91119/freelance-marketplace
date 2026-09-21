import React, { useState } from "react";
import { useNavigate } from "react-router";
import { api } from "../api";

function Signup() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("client");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  async function handleSignup(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");

    try {
      await api("/api/auth/signup", {
        method: "POST",
        body: JSON.stringify({
          name: name,
          email: email,
          password: password,
          role: role,
        }),
      });
      //   console.log(data);
      navigate("/login");
    } catch (err) {
      setError((err as Error).message);
    }
  }

  return (
    <>
      <form onSubmit={handleSignup}>
        <input
          placeholder="name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <input
          placeholder="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          placeholder="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <select value={role} onChange={(e) => setRole(e.target.value)}>
          <option value="client">Client</option>
          <option value="freelancer">Freelancer</option>
        </select>
        <button>Sign up</button>
      </form>
      {error && <p className="error">{error}</p>}
    </>
  );
}

export default Signup;
