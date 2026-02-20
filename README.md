# 🍳 Recipe Sharing API with Image Upload

> A hackathon-ready REST API for sharing recipes with image upload, ingredient-based search, and ratings. Built with Go (Gin + GORM + SQLite).

---

## 🚀 Quick Start

### Prerequisites
- **Go 1.21+** installed ([download](https://go.dev/dl/))
- No database setup needed — SQLite is embedded!

### Setup & Run

```bash
# 1. Clone/navigate to the project
cd D:\Go

# 2. Download dependencies
go mod tidy

# 3. Run the server
go run main.go
```

The server starts at **http://localhost:8080**. Hit `/api/health` to verify.

---

## 📁 Project Structure

```
/
├── .env                          # Environment configuration
├── main.go                       # Entry point (= index.js)
├── go.mod / go.sum               # Dependency management (= package.json)
├── public/temp/                  # Processed image storage
└── src/
    ├── controllers/              # Route handlers
    │   ├── recipe.controller.go  # Recipe CRUD + search
    │   ├── rating.controller.go  # Add & view ratings
    │   └── user.controller.go    # User registration
    ├── db/
    │   └── db.go                 # GORM + SQLite connection
    ├── middlewares/
    │   ├── error.middleware.go   # Global panic recovery
    │   └── upload.middleware.go  # File upload (like multer)
    ├── models/
    │   ├── recipe.model.go       # Recipe schema
    │   ├── rating.model.go       # Rating schema
    │   └── user.model.go         # User schema
    ├── routes/
    │   ├── index.routes.go       # Central route hub
    │   ├── recipe.routes.go      # Recipe endpoints
    │   ├── rating.routes.go      # Rating endpoints
    │   └── user.routes.go        # User endpoints
    └── utils/
        ├── image.util.go         # Resize & compress images
        ├── response.util.go      # Standardized JSON responses
        └── async.util.go         # Safe goroutine wrapper
```

---

## 📡 API Endpoints

### Health Check
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET`  | `/api/health` | Server health check |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/users` | Register a new user |
| `GET`  | `/api/users/:id` | Get user profile + recipes |

### Recipes
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST`   | `/api/recipes` | Create recipe (multipart form + image) |
| `GET`    | `/api/recipes` | List all recipes (paginated) |
| `GET`    | `/api/recipes/:id` | Get single recipe with ratings |
| `GET`    | `/api/recipes/search?ingredients=tomato,onion` | Search by ingredients |
| `PUT`    | `/api/recipes/:id` | Update recipe |
| `DELETE` | `/api/recipes/:id` | Delete recipe + ratings |

### Ratings
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/recipes/:id/ratings` | Rate a recipe (1-5) |
| `GET`  | `/api/recipes/:id/ratings` | Get all ratings for a recipe |

---

## 📝 Example Usage (cURL)

### Register a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"username": "chef_john", "email": "john@example.com"}'
```

### Create a Recipe with Image
```bash
curl -X POST http://localhost:8080/api/recipes \
  -F "title=Spaghetti Bolognese" \
  -F "description=Classic Italian pasta" \
  -F 'ingredients=["spaghetti","tomato","ground beef","onion","garlic"]' \
  -F "prep_time=15" \
  -F "cook_time=30" \
  -F "servings=4" \
  -F "user_id=YOUR_USER_ID" \
  -F "image=@/path/to/photo.jpg"
```

### Search Recipes by Ingredients
```bash
curl "http://localhost:8080/api/recipes/search?ingredients=tomato,garlic"
```

### Rate a Recipe
```bash
curl -X POST http://localhost:8080/api/recipes/RECIPE_ID/ratings \
  -H "Content-Type: application/json" \
  -d '{"user_name": "foodie_jane", "score": 5, "comment": "Best recipe ever!"}'
```

---

## 🛠 Design Decisions

| Decision | Rationale |
|----------|-----------|
| **Gin framework** | Most Express-like Go framework — familiar routing patterns |
| **GORM ORM** | Equivalent to Sequelize — model-driven, auto-migrations |
| **SQLite** | Zero setup, portable, perfect for hackathon evaluation |
| **UUID primary keys** | Better than auto-increment for API resources |
| **JSON ingredients** | Flexible schema for ingredient arrays within SQLite |
| **LIKE-based search** | Pragmatic for SQLite; production would use full-text search |
| **imaging library** | Pure Go, no CGO deps required for image processing |

---

## 📦 Tech Stack

| Component | Technology | Node.js Equivalent |
|-----------|------------|-------------------|
| Language | Go 1.21+ | Node.js |
| Framework | Gin | Express.js |
| ORM | GORM | Sequelize/Mongoose |
| Database | SQLite | MongoDB/PostgreSQL |
| Image Processing | `disintegration/imaging` | `sharp`/`jimp` |
| File Upload | Custom middleware | `multer` |
| Env Config | `godotenv` | `dotenv` |

---

## 📄 License

MIT — Built for hackathon purposes.
