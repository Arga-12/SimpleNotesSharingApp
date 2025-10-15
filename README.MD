# SimpleNotesSharingApp

> A full-stack notes sharing application with authentication, built with **Golang**, **PostgreSQL**, **Next.js**, and **Docker**.

## 📋 Features

### Backend (Golang + PostgreSQL + JWT)
- ✅ User authentication (Register & Login)
- ✅ JWT-based authorization with HttpOnly cookies
- ✅ CRUD operations for notes
- ✅ Owner-only note deletion
- ✅ Comprehensive request/response logging to PostgreSQL
- ✅ Auto-migration system

### Frontend (Next.js)
- ✅ Modern UI with TailwindCSS
- ✅ User registration and login
- ✅ Notes management (Create, Read, Update, Delete)
- ✅ Real-time autosave
- ✅ Responsive design

### Logging System
Every API request is logged with:
- Timestamp (datetime)
- HTTP method + endpoint
- Request headers (Authorization masked)
- Request payload
- Response body
- Response status code
- Request duration (ms)
- User ID (if authenticated)

---

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose installed
- Git

### Setup with Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/Arga-12/SimpleNotesSharingApp.git
   cd SimpleNotesSharingApp
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up --build
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - PgAdmin: http://localhost:5050

4. **Default credentials for PgAdmin**
   - Email: `admin@admin.com`
   - Password: `admin`

---

## 🛠️ Local Development Setup

### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd app/backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set environment variables** (create `.env` file or export)
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=authdb
   export PORT=8080
   export JWT_SECRET=your-secret-key
   ```

4. **Ensure PostgreSQL is running**
   ```bash
   # Using Docker
   docker run -d \
     --name postgres-dev \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=authdb \
     -p 5432:5432 \
     postgres:latest
   ```

5. **Run the backend**
   ```bash
   go run .
   ```

   Backend will be available at `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd app/frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Run development server**
   ```bash
   npm run dev
   ```

   Frontend will be available at `http://localhost:3000`

---

## 📝 Environment Variables

### Backend (.env)

```env
# Database Configuration
DB_HOST=localhost          # Database host (use 'db' for Docker)
DB_PORT=5432              # PostgreSQL port
DB_USER=postgres          # Database username
DB_PASSWORD=postgres      # Database password
DB_NAME=authdb           # Database name

# Server Configuration
PORT=8080                # Backend server port

# JWT Configuration
JWT_SECRET=your-secret-key-here  # Secret key for JWT signing
```

### Frontend (.env.local)

```env
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080  # Backend API URL

# Next.js Configuration
NEXT_TELEMETRY_DISABLED=1
HOST=0.0.0.0
```

### Docker Compose (already configured)

All environment variables are pre-configured in `docker-compose.yml`:
- Backend connects to PostgreSQL service named `db`
- Frontend connects to backend service named `auth-app`
- Persistent volumes for database data

---

## 📚 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/register` | Register new user | ❌ |
| POST | `/api/login` | Login user (returns JWT) | ❌ |
| GET | `/api/me` | Get current user info | ✅ |
| POST | `/api/logout` | Logout user | ✅ |

### Notes

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/notes` | Get all notes | ✅ |
| POST | `/api/notes` | Create new note | ✅ |
| GET | `/api/notes/{id}` | Get note by ID | ✅ |
| PUT | `/api/notes/{id}` | Update note | ✅ |
| DELETE | `/api/notes/{id}` | Delete note (owner only) | ✅ |

---

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password TEXT NOT NULL,
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Notes Table
```sql
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT,
    shared BOOLEAN DEFAULT false,
    favorite BOOLEAN DEFAULT false,
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Logs Table
```sql
CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    datetime TIMESTAMP DEFAULT NOW(),
    method VARCHAR(10) NOT NULL,
    endpoint TEXT NOT NULL,
    request_headers JSONB,
    request_payload TEXT,
    response_body TEXT,
    response_status INTEGER NOT NULL,
    duration_ms INTEGER,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 📊 Example Log Entry

Berikut contoh log yang dihasilkan sistem (dari tabel `logs`):

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
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "Content-Type": "application/json",
    "Content-Length": "51",
    "Sec-Fetch-Mode": "cors",
    "Accept-Encoding": "gzip, deflate, br, zstd"
  },
  "request_payload": "{\"username\":\"abyanpaboss\",\"password\":\"abyanpaboss\"}",
  "response_body": "{\"message\":\"logged in\"}",
  "response_status": 200,
  "duration_ms": 48,
  "user_id": 1,
  "created_at": "2025-10-15 04:19:30.164934"
}
```

### Query untuk melihat logs:
```sql
-- Lihat 10 log terakhir
SELECT * FROM logs ORDER BY datetime DESC LIMIT 10;

-- Lihat log berdasarkan user
SELECT * FROM logs WHERE user_id = 1;

-- Lihat log error (status >= 400)
SELECT * FROM logs WHERE response_status >= 400;

-- Lihat log endpoint tertentu
SELECT * FROM logs WHERE endpoint LIKE '/api/notes%';
```

---

## 🖼️ Screenshots

### Login Page
![Login Page](docs/screenshots/login.png)

### Register Page
![Register Page](docs/screenshots/register.png)

### Notes Dashboard
![Notes Dashboard](docs/screenshots/notes.png)

### Note Editor
![Note Editor](docs/screenshots/editor.png)

---

## 🐳 Docker Services

| Service | Container Name | Port | Description |
|---------|---------------|------|-------------|
| Backend | `auth-app` | 8080 | Golang API server |
| Frontend | `auth-frontend` | 3000 | Next.js application |
| Database | `auth-db` | 5432 | PostgreSQL database |
| PgAdmin | `auth-pgadmin` | 5050 | Database management UI |

### Docker Commands

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f

# Rebuild services
docker-compose up --build

# Remove volumes (reset database)
docker-compose down -v
```

---

## 🔧 Troubleshooting

### Backend tidak bisa connect ke database
- Pastikan PostgreSQL sudah running
- Cek environment variables `DB_HOST`, `DB_PORT`, dll.
- Untuk Docker: gunakan `DB_HOST=db` (nama service)
- Untuk local: gunakan `DB_HOST=localhost`

### Frontend tidak bisa fetch API
- Pastikan backend sudah running di port 8080
- Cek CORS settings di `middlewares/cors.go`
- Pastikan cookie credentials diaktifkan

### Migration error
- Pastikan folder `migrations/` ada di container
- Cek Dockerfile: `COPY --from=builder /app/migrations ./migrations`
- Lihat log: `docker-compose logs auth-app`

---

## 📦 Project Structure

```
.
├── app/
│   ├── backend/              # Golang backend
│   │   ├── handlers/         # HTTP handlers
│   │   ├── middlewares/      # Auth, CORS, Logging
│   │   ├── models/           # Data models
│   │   ├── migrations/       # SQL migrations
│   │   ├── main.go          # Entry point
│   │   ├── db.go            # Database connection
│   │   ├── Dockerfile       # Backend container
│   │   └── go.mod           # Go dependencies
│   └── frontend/            # Next.js frontend
│       ├── pages/           # Next.js pages
│       ├── lib/             # API utilities
│       ├── styles/          # CSS styles
│       └── package.json     # Node dependencies
├── docker-compose.yml       # Docker orchestration
└── README.md               # This file
```

---

## 🧪 Testing

### Test Authentication
```bash
# Register
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  -c cookies.txt

# Get current user
curl http://localhost:8080/api/me -b cookies.txt
```

### Test Notes
```bash
# Create note
curl -X POST http://localhost:8080/api/notes \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"title":"My Note","content":"Hello World","shared":false,"favorite":false}'

# Get all notes
curl http://localhost:8080/api/notes -b cookies.txt

# Delete note
curl -X DELETE http://localhost:8080/api/notes/1 -b cookies.txt
```

---

## 👨‍💻 Development

### Tech Stack
- **Backend**: Go 1.25, PostgreSQL, JWT
- **Frontend**: Next.js 15, React 19, TailwindCSS 4
- **Deployment**: Docker, Docker Compose

### Dependencies

**Backend:**
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/lib/pq` - PostgreSQL driver
- `golang.org/x/crypto` - Password hashing

**Frontend:**
- `next` - React framework
- `react` & `react-dom` - UI library
- `tailwindcss` - CSS framework

---

## 📄 License

MIT License - feel free to use this project for learning purposes.

---

## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first.

---

## 📧 Contact

- GitHub: [@Arga-12](https://github.com/Arga-12)
- Repository: [SimpleNotesSharingApp](https://github.com/Arga-12/SimpleNotesSharingApp)

---

**Made with ❤️ using Golang, Next.js, and Docker**