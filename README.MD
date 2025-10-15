# SimpleNotesSharingApp

> Aplikasi berbagi catatan dengan autentikasi: pengguna dapat registrasi dan login (menghasilkan JWT), melihat dan mengedit catatan milik siapa pun, memerlukan JWT untuk membuat/mengubah catatan, serta hanya pemilik yang dapat menghapus catatan.
---

## Setup

### Dengan Docker (direkomendasikan)
- Tujuan: menjalankan seluruh stack (PostgreSQL, Backend, Frontend, PgAdmin) dengan satu perintah.

1. Clone repo dan jalankan:
   ```bash
   git clone https://github.com/Arga-12/SimpleNotesSharingApp.git
   cd SimpleNotesSharingApp
   docker-compose up --build
   ```
2. Tunggu hingga layanan siap (db healthcheck OK), lalu akses:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - PgAdmin: http://localhost:5050 (email: admin@admin.com, password: admin)
3. Catatan penting:
   - `DB_HOST` pada backend sudah disetel menjadi : `db` via `docker-compose.yml`.
   - `NEXT_PUBLIC_API_URL` frontend sudah diarahkan ke `http://auth-app:8080` saat di Docker.
   - Volume `db_data` menyimpan data Postgres agar tidak hilang saat container restart.
4. Perintah umum:
   ```bash
   docker-compose up -d        # start background
   docker-compose logs -f app  # lihat log backend
   docker-compose down         # stop semua
   docker-compose down -v      # stop + hapus volume (reset DB)
   ```

### Setup Lokal (tanpa Docker)
- Tujuan: menjalankan backend dan frontend langsung di host Anda untuk pengembangan cepat.
- Prasyarat: Go 1.20+, Node.js 18+, PostgreSQL 14+ terpasang dan berjalan di 5432.

1. Backend
   ```bash
   cd app/backend
   # siapkan .env sesuai contoh di bawah (DB_HOST=localhost, dst.)
   go mod download
   go run .
   ```
   - Akses Backend: http://localhost:8080
   - DB akan dimigrasi otomatis pada start (folder `app/backend/migrations`).

2. Frontend
   ```bash
   cd app/frontend
   # siapkan .env.local (NEXT_PUBLIC_API_URL=http://localhost:8080)
   npm install
   npm run dev
   ```
   - Akses Frontend: http://localhost:3000

3. Uji cepat (opsional)
   ```bash
   # cek backend up
   curl http://localhost:8080/api/notes -i
   ```

---

## Contoh .env

### Backend (`app/backend/.env`)

- Lokal (tanpa Docker):
```env
# Database (lokal)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=authdb

# Server
PORT=8080

# JWT (ganti dengan secret yang kuat untuk produksi)
JWT_SECRET=your-secret-key-here
```

- Docker Compose (nilai ini sudah diinject via `docker-compose.yml`, tulis di sini hanya jika jalan manual):
```env
# Database (di dalam Docker network)
DB_HOST=db                  # nama service Postgres di docker-compose
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=authdb

# Server
PORT=8080

# JWT
JWT_SECRET=your-docker-secret
```

### Frontend (`app/frontend/.env.local`)

- Lokal (tanpa Docker):
```env
# URL Backend lokal
NEXT_PUBLIC_API_URL=http://localhost:8080

# Opsional (non-telemetry & binding host)
NEXT_TELEMETRY_DISABLED=1
HOST=0.0.0.0
```

- Docker Compose (biasanya sudah diatur di `docker-compose.yml`):
```env
# URL Backend di dalam Docker network (hanya perlu jika override)
NEXT_PUBLIC_API_URL=http://auth-app:8080

# Opsional
NEXT_TELEMETRY_DISABLED=1
HOST=0.0.0.0
```

---

## Screenshot Aplikasi

- Login: ![Halaman Login](https://imgur.com/oSUmjuh.png)
- Register: ![Halaman Registrasi](https://imgur.com/ULlUZCR.png)
- Dashboard Catatan dan Editor Catatan: ![Dashboard Catatan](https://imgur.com/qbv1F0F.png)

---

## Contoh Log
```json
{
  "id": 4,
  "datetime": "2025-10-15 04:19:30.164571",
  "method": "POST",
  "endpoint": "/api/login",
  "request_headers": {
    "Accept": "*/*",
    "Cookie": "***MASKED***",
    "Origin": "http://localhost:3000",
    "Referer": "http://localhost:3000/login",
    "Sec-Ch-Ua": "\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Opera GX\";v=\"122\"",
    "Connection": "close",
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 OPR/122.0.0.0",
    "Content-Type": "application/json",
    "Content-Length": "51",
    "Sec-Fetch-Dest": "empty",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Site": "same-origin",
    "Accept-Encoding": "gzip, deflate, br, zstd",
    "Accept-Language": "en-US,en;q=0.9,id;q=0.8,fr;q=0.7,de;q=0.6,ja;q=0.5,th;q=0.4,ru;q=0.3",
    "Sec-Ch-Ua-Mobile": "?0",
    "X-Forwarded-Host": "localhost:3000",
    "Sec-Ch-Ua-Platform": "\"Windows\""
  },
  "request_payload": "{\"username\":\"redacted\",\"password\":\"redacted\"}",
  "response_body": "{\"message\":\"logged in\"}",
  "response_status": 200,
  "duration_ms": 48,
  "user_id": 1,
  "created_at": "2025-10-15 04:19:30.164934"
}
```

---

**Made by Argandull using Golang, Next.js, and Docker**