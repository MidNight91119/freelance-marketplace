import { useState } from "react";

function App() {
  return (
    <>
      <Header title="Freelance Marketplace" />
      <Login />
    </>
  );
}

function Header(props: { title: string }) {
  return <h1>{props.title}</h1>;
}

function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");

    const res = await fetch("http://localhost:8080/api/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    });

    const data = await res.json();
    if (!res.ok) {
      setError(data.message);
      return;
    }

    console.log("logged in", data.accessToken);
  }

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <button>Submit</button>
      {error && <p>{error}</p>}
      <p>email is: {email}</p>
      <p>pwd is: {password}</p>
    </form>
  );
}

export default App;
