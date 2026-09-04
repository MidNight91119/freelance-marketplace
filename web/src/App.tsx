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

  async function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    const res = await fetch("http://localhost:8080/api/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    });

    const data = await res.json();
    console.log(data, res.status);
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
      <p>email is: {email}</p>
      <p>pwd is: {password}</p>
    </form>
  );
}

export default App;
