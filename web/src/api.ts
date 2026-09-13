const BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export async function api(path: string, options: RequestInit = {}) {
  const token = localStorage.getItem("token");

  const res = await fetch(BASE + path, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  });

  const data = await res.json();

  if (!res.ok) {
    if (res.status === 401 && token) {
      localStorage.removeItem("token");
      window.location.href = "/login";
    }
    throw new Error(data.message);
  }

  return data;
}
