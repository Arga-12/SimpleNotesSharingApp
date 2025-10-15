import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen flex items-center justify-center p-6" style={{ background: '#0E2148' }}>
      <div className="max-w-md w-full text-center">
        <h1 className="text-4xl font-bold mb-3" style={{ color: '#E3D095' }}>NotesShare</h1>
        <p className="text-gray-300 mb-8">Share your thoughts, collaborate seamlessly</p>
        
        <div className="rounded-lg p-8 shadow-xl" style={{ background: 'rgba(255, 255, 255, 0.1)', borderColor: 'rgba(227, 208, 149, 0.2)', border: '1px solid' }}>
          <h2 className="text-2xl font-bold text-white mb-6">Get Started</h2>
          <div className="space-y-3">
            <Link href="/login">
              <button className="w-full px-6 mb-2 py-3 text-white rounded-lg transition" style={{ background: '#7965C1' }}>
                Login
              </button>
            </Link>
            <Link href="/register">
              <button className="w-full px-6 py-3 text-white rounded-lg transition" style={{ background: '#E3D095', color: '#0E2148' }}>
                Register
              </button>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}