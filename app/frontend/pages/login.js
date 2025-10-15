import { useState } from "react";
import { useRouter } from "next/router";
import Link from "next/link";

export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [msg, setMsg] = useState("");
  const router = useRouter();
  const noHTML = /^[^<>]*$/; // No HTML tags allowed

  const handleLogin = async (e) => {
    e.preventDefault();
    setMsg("");

    if (!username || !password) {
      setMsg("Username dan password wajib diisi");
      return;
    }

    // Validate no HTML tags
    if (!noHTML.test(username) || !noHTML.test(password)) {
      setMsg("Dont try to sanitize us you naughty");
      return;
    }

    try {
      const res = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
        credentials: "include",
      });

      const data = await res.json().catch(() => ({ error: "Invalid response" }));

      if (res.ok) {
        router.push("/notes");
      } else {
        setMsg(data.error || "Login gagal");
      }
    } catch (err) {
      setMsg("Terjadi kesalahan koneksi");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-6" style={{ background: '#0E2148' }}>
      <div className="max-w-md w-full">
        <div className="text-center mb-6">
          <h1 className="text-3xl font-bold mb-2" style={{ color: '#E3D095' }}>NotesShare</h1>
          <p className="text-gray-300">Login to your account</p>
        </div>

        <div className="rounded-lg p-8 shadow-xl" style={{ background: 'rgba(255, 255, 255, 0.1)', borderColor: 'rgba(227, 208, 149, 0.2)', border: '1px solid' }}>
          <h2 className="text-2xl font-bold text-white mb-6">Login</h2>
          
          <form onSubmit={handleLogin} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Username</label>
              <input 
                value={username} 
                onChange={(e) => setUsername(e.target.value)}
                className="w-full px-4 py-2 rounded-lg focus:outline-none text-white" 
                style={{ background: 'rgba(255, 255, 255, 0.1)', border: '1px solid rgba(227, 208, 149, 0.3)' }}
                placeholder="Enter username"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Password</label>
              <input 
                type="password" 
                value={password} 
                onChange={(e) => setPassword(e.target.value)}
                className="w-full px-4 py-2 rounded-lg focus:outline-none text-white"
                style={{ background: 'rgba(255, 255, 255, 0.1)', border: '1px solid rgba(227, 208, 149, 0.3)' }}
                placeholder="Enter password"
              />
            </div>
            <button 
              type="submit"
              className="w-full py-3 text-white rounded-lg font-semibold transition"
              style={{ background: '#7965C1' }}
            >
              Sign In
            </button>
          </form>

          {msg && (
            <div className="mt-4 p-3 rounded-lg text-sm" style={{ background: 'rgba(239, 68, 68, 0.2)', color: '#ffcccc' }}>
              {msg}
            </div>
          )}

          <p className="mt-6 text-center text-sm text-gray-400">
            Don't have an account?{' '}
            <Link href="/register" className="font-semibold hover:underline" style={{ color: '#E3D095' }}>
              Register
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}