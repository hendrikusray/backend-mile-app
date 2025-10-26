# 🚀 Mile App Backend

Backend service for the Mile App Fullstack Test — built using **Go (Gin)** and **MongoDB**, providing RESTful APIs for authentication and Task CRUD operations.

---
### Authentication
- `POST /v1/login`: Mock login that returns a token if username & password match.
- Token is then used by the frontend to simulate authenticated requests.

### Task CRUD Module
- Endpoints:
  - `GET /v1/tasks`: Get list of tasks (supports filter, sort, pagination).
  - `POST /v1/tasks`: Create new task.
  - `PUT /v1/tasks/:id`: Update task.
  - `DELETE /v1/tasks/:id`: Delete task.
- All responses follow a consistent JSON structure with `meta` info (page, total, limit).

---

## Design Decisions
| Area | Decision | Reason |
|-------|-----------|--------|
| **Language** | Go (Gin Framework) | Fast, lightweight, clean routing and middleware. |
| **Database** | MongoDB | Flexible schema and quick iteration for document data (tasks). |
| **Architecture** | Layered (domain, repository, usecase, delivery) | Clear separation for maintainability and testability. |
| **Response Format** | Standardized JSON | Easy integration for frontend and consistent structure. |
| **Deployment** | Railway | Simple continuous deployment for Go-based apps. |

---
## Strengths of This Module
- Clean and modular code (easy to add new features).
- Supports filtering, sorting, and pagination natively.
- Fast and efficient using MongoDB indexes.
- Consistent REST response with meta information.
- Deployed and accessible online for integration with frontend.

---
## Database Indexes (db/indexes.js)

```js
db.tasks.createIndex({ title: "text", description: "text" });
db.tasks.createIndex({ status: 1 });
db.tasks.createIndex({ created_at: -1 });
db.tasks.createIndex({ owner_id: 1 });

## Struktur
mile-app/
├── app/main.go
├── config/ (Tempat semua konfigurasi global (database, environment, setup awal))
│   ├── config.go
│   └── mongo_db.go
├── domain/ (Berisi definisi struktur data (model) dan interface (kontrak) antar layer)
│   ├── user.go
│   ├── task.go
│   └── interfaces.go
├── internal/ (implementasi login logic)
│   ├── user/
│   └── task/
└── db/ (Tempat semua file yang berhubungan dengan database setup, seeding, dan Docker environment)
    ├── data/
    ├── indexes.js
    └── Dockerfile

