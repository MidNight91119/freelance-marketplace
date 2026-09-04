import { Routes, Route } from "react-router";
import Login from "./pages/Login";
import Projects from "./pages/Projects";

function App() {
  return (
    <>
      <Header title="Freelance Marketplace" />
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/projects" element={<Projects />} />
      </Routes>
    </>
  );
}

function Header(props: { title: string }) {
  return <h1>{props.title}</h1>;
}

export default App;
