import { Routes, Route, Navigate } from "react-router";
import Login from "./pages/Login";
import Projects from "./pages/Projects";
import Signup from "./pages/Signup";
import CreateProject from "./pages/CreateProject.tsx";
import SubmitProposal from "./pages/SubmitProposal.tsx";

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
            <ProtectedRoute>
              <CreateProject />
            </ProtectedRoute>
          }
        />
        <Route
          path="/projects/:projectId/propose"
          element={
            // <ProtectedRoute>
            <SubmitProposal />
            // </ProtectedRoute>
          }
        />
      </Routes>
    </>
  );
}

function Header(props: { title: string }) {
  return <h1>{props.title}</h1>;
}

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const token = localStorage.getItem("token");

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return children;
}

export default App;
