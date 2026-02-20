# 🗣 Prompts Used to Build This Project

This document lists the prompts that were used to generate the entire Recipe Sharing API codebase.

---

## Prompt 1 — Initial Project Generation

```
Act as an expert Go (Golang) backend developer. I need to build a "Recipe Sharing API 
with Image Upload" for a hackathon project within a 24-hour deadline. 

I am transitioning from Node.js/Express.js, and I want this Go project to STRICTLY 
follow my standard Node.js project structure to make it easier for me to maintain 
and present. 

PROJECT REQUIREMENTS:
1. Core Features: Post recipes with images, search for recipes by ingredients, and 
   rate others' recipes.
2. Image Handling: Handle image uploads, process them (resize/compress to optimize 
   storage), and store the recipe data. You can save images locally in a `public/temp` 
   folder or mock a cloud upload, but the resizing logic must be in Go.
3. Tech Stack: 
   - Language: Go (Golang)
   - Framework: Gin (for Express-like routing)
   - ORM: GORM 
   - Database: SQLite (for zero-setup portability during the hackathon evaluation)
   - Image Processing: Use standard Go packages or a lightweight library like 
     `github.com/disintegration/imaging`.

REQUIRED FILE STRUCTURE:
Please generate the complete codebase mapping standard Go idioms to this exact 
directory structure:
/
├── .env
├── main.go               # Entry point (equivalent to index.js)
├── go.mod
├── go.sum
└── src/
    ├── controllers/      # Handlers for recipes, ratings, users
    ├── db/               # Database connection and auto-migration setup
    ├── middlewares/       # Gin middlewares for error handling or multer-like parsing
    ├── models/           # GORM structs (Recipe, Rating, Ingredient)
    ├── routes/           # Gin router setups separated by domain
    └── utils/            # Helpers for image compression, error responses, async wrappers

DELIVERABLES TO GENERATE:
1. All Go source code files mapped to the structure above.
2. A `README.md` containing API documentation, setup instructions, and design decisions.
3. A `design.md` document explaining the database schema, architecture flow, and how 
   the image compression is handled.
4. A raw JSON format of a Postman/cURL collection so I can import and test all 
   endpoints immediately.

Please write clean, modular, and heavily commented code so I can easily explain the 
design decisions in a 7-minute code review presentation. Start by generating `main.go`, 
the `db` setup, and the `models`.

Create a md file for all explanation of what we have done.
And another md file for what all prompts I used to do so.
Start doing.
```

---

## What This Prompt Achieved

From this single prompt, the following was generated:

### Source Code (16 files)
| File | Description |
|------|-------------|
| `main.go` | Application entry point with CORS, static files, route mounting |
| `.env` | Environment configuration |
| `go.mod` | Go module with all dependencies |
| `src/db/db.go` | GORM + SQLite connection and auto-migration |
| `src/models/user.model.go` | User struct with UUID hook |
| `src/models/recipe.model.go` | Recipe struct with JSON ingredients |
| `src/models/rating.model.go` | Rating struct (1-5 scale) |
| `src/controllers/recipe.controller.go` | Recipe CRUD + ingredient search |
| `src/controllers/rating.controller.go` | Rating add + average recalculation |
| `src/controllers/user.controller.go` | User registration + profile |
| `src/middlewares/error.middleware.go` | Global panic recovery |
| `src/middlewares/upload.middleware.go` | Multer-like file upload handling |
| `src/routes/recipe.routes.go` | Recipe route definitions |
| `src/routes/rating.routes.go` | Rating route definitions |
| `src/routes/user.routes.go` | User route definitions |
| `src/routes/index.routes.go` | Central route hub + health check |
| `src/utils/image.util.go` | Image resize + JPEG compression |
| `src/utils/response.util.go` | Standardized API response helpers |
| `src/utils/async.util.go` | Safe goroutine wrapper |

### Documentation (4 files)
| File | Description |
|------|-------------|
| `README.md` | Setup instructions, API docs, design decisions |
| `design.md` | DB schema (Mermaid ER), architecture flow, image pipeline |
| `postman_collection.json` | Postman v2.1 collection for all endpoints |
| `explanation.md` | Detailed technical explanation of every component |

---

## Follow-Up Prompts (If Needed)

Here are some useful follow-up prompts for extending the project:

### Add Authentication
```
Add JWT-based authentication to the Recipe API. Create a src/middlewares/auth.middleware.go 
that validates JWT tokens, and update the user model to include password hashing with bcrypt. 
Protect recipe creation, update, and delete routes — only the recipe owner should be able 
to modify their recipes.
```

### Add Pagination Metadata
```
Update the GetAllRecipes endpoint to include pagination links (next, prev) in the response.
```

### Add Categories/Tags
```
Add a Category model and a many-to-many relationship between recipes and categories. 
Add an endpoint to filter recipes by category.
```

### Deploy to Cloud
```
Create a Dockerfile and docker-compose.yml for the Recipe API. Use a multi-stage build 
to minimize the final image size. Include instructions for deploying to Railway or Render.
```
