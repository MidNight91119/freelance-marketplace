import { Routes, Route, Navigate } from "react-router";
import Login from "./pages/Login";
import Projects from "./pages/Projects";
import Signup from "./pages/Signup";
import CreateProject from "./pages/CreateProject.tsx";
import SubmitProposal from "./pages/SubmitProposal.tsx";
import ViewProposal from "./pages/ViewProposal.tsx";
import Contracts from "./pages/Contracts.tsx";

function App() {
  return (
    <>
      <Header title="Freelance Marketplace" />
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
      </Routes>
    </>
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
  const token = localStorage.getItem("token");

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  const userRole = localStorage.getItem("role");
  if (role && userRole !== role) {
    return <Navigate to="/projects" replace />;
  }

  return children;
}

export default App;
