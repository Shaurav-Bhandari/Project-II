# Restaurant Operations Management System

Complete setup and deployment guide for the Restaurant Management System.

## Prerequisites

- **Node.js** 18+ (for frontend)
- **Go** 1.21+ (for backend)
- **PostgreSQL** 14+ (database)

## Quick Start

### 1. Database Setup

```bash
# Create database
createdb restaurant_db

# Run schema
psql -d restaurant_db -f backend/db/schema.sql
```

### 2. Backend Setup

```bash
cd backend

# Copy environment file
cp .env.example .env
# Edit .env with your database credentials

# Run server
go run main.go
```

Server runs on `http://localhost:8080`

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev
```

App runs on `http://localhost:5173`

## Default Login

- **Email**: `admin@restaurant.com`
- **Password**: `admin123`

## Project Structure

```
project II/
├── backend/
│   ├── main.go              # Server entry point
│   ├── config/              # Database config
│   ├── handlers/            # API handlers
│   ├── middleware/          # Auth middleware
│   ├── models/              # Data models
│   └── db/schema.sql        # Database schema
└── frontend/
    ├── src/
    │   ├── views/           # Page components
    │   ├── stores/          # Pinia stores
    │   ├── router/          # Vue router
    │   └── api/             # API client
    └── package.json
```

## Key Features

| Feature | Description |
|---------|-------------|
| **Dashboard** | Real-time metrics, quick actions |
| **Orders** | Create, manage, track orders |
| **KOT** | Kitchen order tickets with priority |
| **Menu** | Item and category management |
| **Billing** | Payment processing, receipts |
| **Users** | Role-based access control |
| **Reports** | Sales analytics, popular items |

## Algorithms Implemented

1. **Bcrypt Hashing** - Secure password storage
2. **JWT Tokens** - Session management
3. **Merge Sort** - KOT priority ordering
4. **Backtracking** - Table allocation for groups
5. **Binary Search** - Menu price lookup
