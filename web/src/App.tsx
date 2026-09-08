import { Routes, Route } from "react-router";
import Login from "./pages/Login";
import Projects from "./pages/Projects";
import Signup from "./pages/Signup";

function App() {
  return (
    <>
      <Header title="Freelance Marketplace" />
      <Routes>
        <Route path="/signup" element={<Signup />} />
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
