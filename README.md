# Restaurant Operations Management System

## 📖 Overview
This project is a full‑stack **Restaurant Operations Management System** built with:
- **Go (Gin)** – backend API server
- **PostgreSQL** – relational database
- **Vue 3 + Vite** – modern frontend UI (instead of React)
- **Pinia** – state management (Vue’s equivalent of Redux)
- **JWT + Bcrypt** – secure authentication & password handling

The system supports dashboards, order management, kitchen order tickets (KOT), menu administration, billing, user/role management, and analytics reports.

---

## 🛠️ Backend (Go) Architecture
### Core Packages
| Package | Purpose |
|---|---|
| `config` | Initializes the PostgreSQL connection and runs DB migrations. |
| `handlers` | Implements all HTTP handlers (CRUD for users, orders, tables, menu, etc.). |
| `middleware` | Auth middleware (`AuthMiddleware`) validates JWTs; `RoleMiddleware` enforces role‑based access control. |
| `models` | Struct definitions that map to DB tables (User, Order, Table, MenuItem, …). |
| `db/schema.sql` | Database schema – tables, indexes, constraints. |

### Important Concepts & Algorithms
- **JWT Authentication** – Stateless session tokens signed with a secret. Tokens are issued on login (`Login`) and verified by `AuthMiddleware`. This avoids server‑side session storage.
- **Bcrypt Password Hashing** – Passwords are never stored plain‑text. `Register` hashes with `bcrypt.GenerateFromPassword`; `Login` compares using `bcrypt.CompareHashAndPassword`.
- **Role‑Based Access Control** – `RoleMiddleware` receives a list of allowed roles (e.g., `admin`, `manager`). It checks `user.Role` from the JWT claims and aborts with `403` if unauthorized.
- **CORS Configuration** – Allows the Vue dev server (`localhost:5173`) to call the API.
- **Merge Sort for KOT Priority** – Orders in the kitchen are sorted by priority (e.g., VIP, time‑sensitivity) using a custom merge‑sort implementation.
- **Backtracking Table Allocation** – When a group of guests arrives, the system tries to allocate a combination of tables that exactly fits the party size using backtracking.
- **Binary Search for Menu Price Lookup** – Fast price retrieval when validating order totals.

### Key Functions (selected)
```go
// handlers.Login – authenticates a user and returns a JWT
func (h *Handler) Login(c *gin.Context) {
    var creds LoginRequest
    if err := c.ShouldBindJSON(&creds); err != nil { … }
    // fetch user, compare bcrypt hash, generate JWT
}

// middleware.AuthMiddleware – protects routes
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.GetHeader("Authorization")
        // validate token, set user info in context
    }
}

// config.InitDB – opens a PostgreSQL connection and runs migrations
func InitDB() (*sql.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), …)
    db, err := sql.Open("postgres", dsn)
    // ping & migrate
}
```

### Go Module Dependencies (`go.mod`)
```text
github.com/gin-gonic/gin v1.9.1          // HTTP router & middleware framework
github.com/gin-contrib/cors v1.3.1      // CORS handling for cross‑origin requests
github.com/joho/godotenv v1.4.0         // Loads .env files into environment variables
golang.org/x/crypto/bcrypt v0.0.0‑20230712 // Secure password hashing
github.com/dgrijalva/jwt-go v3.2.0+incompatible // JWT creation & verification
```
- **Gin** provides a lightweight, high‑performance router with built‑in middleware support.
- **CORS** is required because the frontend runs on a different port during development.
- **godotenv** simplifies local development by loading configuration from `.env`.
- **bcrypt** ensures passwords are stored securely.
- **jwt-go** implements the token‑based auth flow.

---

## 🎨 Frontend (Vue 3) Overview for React Developers
The UI lives in `frontend/` and is a **Vite**‑powered Vue 3 project.

### Why Vue?
- **Composition API** (similar to React hooks) – `setup()` function with `ref`, `computed`, and lifecycle hooks.
- **Single‑File Components** (`.vue`) combine template, script, and style – analogous to JSX + CSS modules.
- **Pinia** – the official state‑management library, conceptually like Redux Toolkit but with a simpler API.
- **Vue Router** – declarative routing, comparable to React Router.

### Project Structure
```
frontend/
├─ src/
│  ├─ api/          # thin wrapper around fetch/axios for backend calls
│  ├─ assets/       # static images, icons
│  ├─ components/   # reusable UI pieces (buttons, tables, etc.)
│  ├─ router/       # route definitions (src/router/index.js)
│  ├─ stores/       # Pinia stores (auth, orders, kots, etc.)
│  ├─ views/        # page‑level components (Dashboard.vue, Orders.vue, …)
│  └─ App.vue       # root layout with sidebar, header, router‑view
├─ index.html
├─ package.json
└─ vite.config.js
```

### Core Files Explained
#### `src/main.js`
```js
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
```
- Boots the Vue app, registers Pinia (state) and Vue Router (navigation).

#### `src/router/index.js`
```js
import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Login from '../views/Login.vue'
// …other lazy‑loaded routes

const routes = [
  { path: '/', component: Dashboard, meta: { requiresAuth: true } },
  { path: '/login', component: Login },
  // protected routes use a navigation guard (see below)
]

const router = createRouter({ history: createWebHistory(), routes })
router.beforeEach((to, from, next) => {
  const isAuth = !!localStorage.getItem('token')
  if (to.meta.requiresAuth && !isAuth) return next('/login')
  next()
})
export default router
```
- Mirrors React Router’s `<Route>` definitions and a global guard for auth.

#### Pinia Stores (`src/stores/*.js`)
Example: `auth.js`
```js
import { defineStore } from 'pinia'
import api from '../api'

export const useAuthStore = defineStore('auth', {
  state: () => ({ token: null, user: null }),
  actions: {
    async login(email, password) {
      const { data } = await api.post('/auth/login', { email, password })
      this.token = data.token
      this.user = data.user
      localStorage.setItem('token', data.token)
    },
    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
    },
  },
  getters: {
    isAuthenticated: (state) => !!state.token,
    userName: (state) => state.user?.name || ''
  }
})
```
- Similar to a Redux slice: state, actions (async thunks), and selectors (getters).

#### `src/api/index.js`
```js
import axios from 'axios'
const api = axios.create({ baseURL: 'http://localhost:8080/api' })
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})
export default api
```
- Centralised HTTP client; comparable to a custom `fetch` wrapper in React projects.

#### `src/App.vue`
- Provides the **layout**: a collapsible sidebar with navigation links, a top header with notifications, and a `<router-view/>` where page components render.
- Uses Vue’s `computed` to derive values from Pinia stores (`activeOrderCount`, `pendingKotCount`).
- The logout button calls `authStore.logout()` – analogous to a React context logout function.

### Styling & UI
- The project uses **vanilla CSS** with CSS variables for theming (light/dark mode). You can replace it with a UI library (e.g., Vuetify) if desired.
- Micro‑animations are added via CSS transitions on hover/focus to give a premium feel.

### How to Get Started (React‑style checklist)
1. **Install Node** (>=18) – already listed in prerequisites.
2. ```bash
   cd frontend
   npm install   # installs Vue, Pinia, Vue Router, axios, etc.
   npm run dev   # starts Vite dev server on http://localhost:5173
   ```
3. Ensure the **backend** (`go run main.go`) is running on port 8080.
4. Open the dev URL; you’ll be redirected to `/login`. Use the default credentials from the README.

### Key Frontend Dependencies (`package.json`)
| Dependency | Reason |
|---|---|
| `vue@3` | Core framework – composition API, reactivity system. |
| `pinia` | State management; lighter and more ergonomic than Vuex. |
| `vue-router` | Declarative routing, similar to React Router. |
| `axios` | Promise‑based HTTP client for API calls. |
| `vite` | Fast dev server & bundler; comparable to Create‑React‑App but with native ES modules. |

---

## 🚀 Running the Full Stack
```bash
# 1️⃣ Backend
cd backend
go run main.go   # server on http://localhost:8080

# 2️⃣ Frontend (in a separate terminal)
cd ../frontend
npm run dev      # Vite dev server on http://localhost:5173
```
The frontend automatically proxies API requests to the backend because of the CORS settings.

## 📚 Further Reading
- **Gin Documentation** – https://gin-gonic.com/docs/
- **Vue 3 Composition API** – https://vuejs.org/guide/composition-api/introduction.html
- **Pinia Guide** – https://pinia.vuejs.org/
- **JWT Best Practices** – https://jwt.io/introduction/

---

*Happy coding! 🎉*
