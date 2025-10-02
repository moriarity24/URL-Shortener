# URL Shortener

A high-performance URL shortener built with Go, PostgreSQL, and Redis.

## Features

- 🚀 Fast URL shortening
- 📊 Click tracking
- ⚡ Redis caching for lightning-fast redirects
- 🔄 Idempotent (same URL always returns same short code)
- 🐳 Fully Dockerized

## Tech Stack

- **Go** - Backend language
- **Gin** - Web framework
- **PostgreSQL** - Database
- **Redis** - Cache
- **Docker** - Containerization
- **GORM** - ORM

## Project Structure
url-shortener/
├── cmd/api/          # Application entry point
├── internal/
│   ├── config/       # Configuration
│   ├── database/     # Database connections
│   ├── handlers/     # HTTP handlers
│   ├── models/       # Data models
│   ├── repository/   # Database operations
│   └── service/      # Business logic
└── docker-compose.yml

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/url-shortener.git
cd url-shortener

Install dependencies

bashgo mod download

Start databases

bashdocker-compose up -d

Run the application

bashgo run cmd/api/main.go
Server will start on http://localhost:8080
API Endpoints
Create Short URL
bashPOST /api/shorten
Content-Type: application/json

{
  "url": "https://www.google.com"
}
Response:
json{
  "short_code": "aB3xK9",
  "short_url": "http://localhost:8080/aB3xK9",
  "original_url": "https://www.google.com"
}
Redirect
bashGET /:shortCode
Redirects to original URL and increments click counter.
Health Check
bashGET /health
Response:
json{
  "status": "healthy"
}
Architecture
Clean Architecture with 3 layers:
Handler → Service → Repository → Database

Handlers: HTTP request/response
Service: Business logic
Repository: Data access

Environment Variables
See .env.example for configuration options.
License
MIT
EOF

---

## 🎉 WE'RE DONE BUILDING! 

Now let's test it! 🚀

---

## 📋 Step 17: Test the Application!

### 17.1: Make sure Docker is running
```bash
docker-compose ps
Should show both postgres and redis running.
If not running:
bashdocker-compose up -d

17.2: Run the application!
bashgo run cmd/api/main.go
You should see:
✅ Connected to PostgreSQL and Redis
[GIN-debug] GET    /health                   --> main.main.func2 (4 handlers)
[GIN-debug] POST   /api/shorten              --> main.main.func2 (4 handlers)
[GIN-debug] GET    /:shortCode               --> main.main.func2 (4 handlers)
🚀 Server starting on http://localhost:8080