import { Routes, Route, Navigate, Link, useNavigate } from "react-router";
import Login from "./pages/Login";
import Projects from "./pages/Projects";
import Signup from "./pages/Signup";
import CreateProject from "./pages/CreateProject.tsx";
import SubmitProposal from "./pages/SubmitProposal.tsx";
import ViewProposal from "./pages/ViewProposal.tsx";
import Contracts from "./pages/Contracts.tsx";
import MyProjects from "./pages/MyProjects.tsx";
import MyProposals from "./pages/MyProposals.tsx";
import { useAuth } from "./auth.tsx";

function App() {
  return (
    <>
      <Header title="Freelance Marketplace" />
      <Nav />
      <main>
        <Routes>
          <Route path="/signup" element={<Signup />} />
          <Route path="/login" element={<Login />} />
          <Route
            path="/projects"
            element={
              <ProtectedRoute>
                <Projects />
              </ProtectedRoute>
            }
          />
          <Route
            path="/projects/new"
            element={
              <ProtectedRoute role="client">
                <CreateProject />
              </ProtectedRoute>
            }
          />
          <Route
            path="/projects/:projectId/propose"
            element={
              <ProtectedRoute role="freelancer">
                <SubmitProposal />
              </ProtectedRoute>
            }
          />
          <Route
            path="/projects/:projectId/proposals"
            element={
              <ProtectedRoute role="client">
                <ViewProposal />
              </ProtectedRoute>
            }
          />
          <Route
            path="/contracts"
            element={
              <ProtectedRoute>
                <Contracts />
              </ProtectedRoute>
            }
          />
          <Route
            path="/projects/mine"
            element={
              <ProtectedRoute role="client">
                <MyProjects />
              </ProtectedRoute>
            }
          />
          <Route
            path="/proposals/mine"
            element={
              <ProtectedRoute role="freelancer">
                <MyProposals />
              </ProtectedRoute>
            }
          />
        </Routes>
      </main>
    </>
  );
}

function Nav() {
  const { role, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/login");
  }

  return (
    <nav>
      {role && <Link to="/projects">Browse</Link>}{" "}
      {role && <Link to="/contracts">Contracts</Link>}{" "}
      {role === "client" && <Link to="/projects/mine">My projects</Link>}{" "}
      {role === "freelancer" && <Link to="/proposals/mine">My proposals</Link>}{" "}
      {!role && (
        <>
          <Link to="/login">Login</Link> <Link to="/signup">Signup</Link>
        </>
      )}{" "}
      {role && <button onClick={handleLogout}>Logout</button>}
    </nav>
  );
}

function Header(props: { title: string }) {
  return <h1>{props.title}</h1>;
}

function ProtectedRoute({
  children,
  role,
}: {
  children: React.ReactNode;
  role?: "client" | "freelancer";
}) {
  const { role: userRole } = useAuth();
  const token = localStorage.getItem("token");

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  if (role && userRole !== role) {
    return <Navigate to="/projects" replace />;
  }

  return children;
}

export default App;
