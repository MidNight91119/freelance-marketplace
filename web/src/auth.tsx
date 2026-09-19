import { createContext, useContext, useState } from "react";

type Auth = {
  role: string | null;
  login: (token: string, role: string) => void;
  logout: () => void;
};

const AuthContext = createContext<Auth | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [role, setRole] = useState(localStorage.getItem("role"));

  function login(token: string, role: string) {
    localStorage.setItem("token", token);
    localStorage.setItem("role", role);
    setRole(role);
  }

  function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
    setRole(null);
  }

  return (
    <AuthContext.Provider value={{ role, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside AuthProvider");
  return ctx;
}
