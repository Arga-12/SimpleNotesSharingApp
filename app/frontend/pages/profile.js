import { useEffect, useState } from "react";
import { useRouter } from "next/router";

export default function Profile() {
  const [user, setUser] = useState(null);
  const [loading,setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const fetchMe = async () => {
      const res = await fetch("/api/me", { credentials: "include" });
      if (res.ok) {
        const data = await res.json();
        setUser(data);
      } else {
        router.push("/");
      }
      setLoading(false);
    };
    fetchMe();
  }, []);

  const logout = async () => {
    await fetch("/api/logout", { method: "POST", credentials: "include" });
    router.push("/");
  };

  if (loading) return <div>Loading...</div>;
  if (!user) return null;

  return (
    <div style={{ maxWidth: 640, margin: "40px auto", fontFamily: "system-ui" }}>
      <h1>Profile</h1>
      <p><strong>ID:</strong> {user.id}</p>
      <p><strong>Username:</strong> {user.username}</p>
      <p><strong>Email:</strong> {user.email || "-"}</p>
      <p><strong>Created:</strong> {new Date(user.created_at).toLocaleString()}</p>
      <button onClick={logout}>Logout</button>
    </div>
  );
}