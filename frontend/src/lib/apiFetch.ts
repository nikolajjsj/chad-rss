export const apiFetch = async (url: string, options: RequestInit = {}) => {
  const token = localStorage.getItem("token");
  const headers = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };
  const response = await fetch(`http://localhost:8080${url}`, {
    ...options,
    headers,
  });
  if (!response.ok) {
    throw new Error(response.statusText);
  }
  return response.json();
};
