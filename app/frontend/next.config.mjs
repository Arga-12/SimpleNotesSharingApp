/** @type {import('next').NextConfig} */
const nextConfig = {
  // Hanya masukkan konfigurasi yang diizinkan oleh next.js
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        // Ini adalah konfigurasi proxy yang sudah benar ke Golang di port 8080
        destination: 'http://auth-app:8080/api/:path*', 
      },
    ]
  },
};

// Gunakan 'export default' untuk ES Module
export default nextConfig;