# 📖 Complete Explanation — Recipe Sharing API

---

## Table of Contents

1. [What Does This Project Do?](#1-what-does-this-project-do)
2. [How The Database Works](#2-how-the-database-works)
3. [All API Endpoints (Full List)](#3-all-api-endpoints)
4. [How To Test Every Endpoint Using Postman](#4-how-to-test-every-endpoint-using-postman)
5. [How Image Upload & Compression Works](#5-how-image-upload--compression-works)
6. [How The Project Is Organized](#6-how-the-project-is-organized)
7. [How To Run This Project](#7-how-to-run-this-project)
8. [Presentation Tips (7-Minute Code Review)](#8-presentation-tips)

---

## 1. What Does This Project Do?

This is a REST API (backend only, no frontend) that lets users:

- **Register** as a user with a username and email
- **Post recipes** with a title, description, list of ingredients, cooking times, and an optional photo
- **Search recipes** by typing ingredient names (e.g., "give me all recipes that use tomato and garlic")
- **Rate recipes** that other people posted (score from 1 to 5 stars, with an optional comment)
- **View recipes** with their average rating calculated automatically

The API stores everything in a **SQLite database** (a single file called `recipe.db` — no need to install MySQL, PostgreSQL, or MongoDB).

---

## 2. How The Database Works

### What is SQLite?

SQLite is a database that stores everything in **one single file** (`recipe.db`). You don't need to install any database server. When you run the app, it automatically creates this file and all the tables inside it.

### What is GORM?

GORM is the library we use to talk to the database. Instead of writing raw SQL queries like `SELECT * FROM recipes WHERE id = '123'`, we write Go code like `db.DB.First(&recipe, "id = ?", id)`. GORM translates this into SQL for us.

### What is AutoMigrate?

When the app starts, this line runs:

```go
DB.AutoMigrate(&models.User{}, &models.Recipe{}, &models.Rating{})
```

This looks at our Go structs (User, Recipe, Rating) and **automatically creates the database tables** to match them. If you add a new field to a struct later, it will add that column to the table. You never need to manually create tables.

### The 3 Tables

#### Table 1: `users`

This table stores people who use the app.

| Column | Type | Description |
|--------|------|-------------|
| `id` | TEXT (UUID) | Unique ID like `"a1b2c3d4-e5f6-..."`. Auto-generated. |
| `username` | TEXT | Must be unique. Example: `"chef_john"` |
| `email` | TEXT | Must be unique. Example: `"john@example.com"` |
| `created_at` | DATETIME | When the user registered. Auto-set. |
| `updated_at` | DATETIME | When the user was last modified. Auto-set. |

**Relationships:** One user can have many recipes (one-to-many).

#### Table 2: `recipes`

This is the main table. It stores all the recipes.

| Column | Type | Description |
|--------|------|-------------|
| `id` | TEXT (UUID) | Unique ID. Auto-generated. |
| `title` | TEXT | Recipe name. Example: `"Spaghetti Bolognese"` |
| `description` | TEXT | A short description of the recipe. |
| `image_url` | TEXT | Path to the uploaded photo. Example: `"/uploads/recipe_a1b2c3d4.jpg"` |
| `ingredients` | TEXT | A JSON array stored as a string. Example: `'["tomato","onion","garlic"]'` |
| `prep_time` | INTEGER | Preparation time in minutes. Example: `15` |
| `cook_time` | INTEGER | Cooking time in minutes. Example: `30` |
| `servings` | INTEGER | How many people it serves. Example: `4` |
| `average_rating` | REAL (float) | Calculated automatically from all ratings. Example: `4.5` |
| `user_id` | TEXT | The ID of the user who created this recipe (foreign key to `users` table). |
| `created_at` | DATETIME | Auto-set. |
| `updated_at` | DATETIME | Auto-set. |

**Relationships:**
- Each recipe belongs to one user (`user_id` → `users.id`)
- Each recipe can have many ratings (one-to-many)

**Why are ingredients stored as a JSON string?**  
SQLite doesn't have a native JSON/Array column type. So we store ingredients as a plain text string in JSON format like `'["tomato","onion","garlic"]'`. When searching, we use `LIKE '%tomato%'` to find recipes that contain "tomato" anywhere in that string.

#### Table 3: `ratings`

This table stores star ratings and comments for recipes.

| Column | Type | Description |
|--------|------|-------------|
| `id` | TEXT (UUID) | Unique ID. Auto-generated. |
| `recipe_id` | TEXT | Which recipe is being rated (foreign key to `recipes.id`). |
| `user_name` | TEXT | Who rated it. Example: `"foodie_jane"` |
| `score` | INTEGER | Star rating from 1 to 5. |
| `comment` | TEXT | Optional comment. Example: `"Loved it!"` |
| `created_at` | DATETIME | Auto-set. |
| `updated_at` | DATETIME | Auto-set. |

**What happens when a rating is added?**  
Every time someone rates a recipe, the app automatically recalculates the `average_rating` on the recipe. It runs `AVG(score)` across all ratings for that recipe and updates the recipes table. So the `average_rating` is always up-to-date.

### How The Tables Are Connected (Visual)

```
┌──────────┐         ┌──────────────┐         ┌──────────────┐
│  USERS   │         │   RECIPES    │         │   RATINGS    │
├──────────┤         ├──────────────┤         ├──────────────┤
│ id (PK)  │◄────────│ user_id (FK) │         │ id (PK)      │
│ username │   1:N   │ id (PK)      │◄────────│ recipe_id(FK)│
│ email    │         │ title        │   1:N   │ user_name    │
│          │         │ ingredients  │         │ score (1-5)  │
│          │         │ image_url    │         │ comment      │
│          │         │ avg_rating   │         │              │
└──────────┘         └──────────────┘         └──────────────┘

PK = Primary Key (unique identifier)
FK = Foreign Key (links to another table)
1:N = One-to-Many relationship
```

---

## 3. All API Endpoints

The API has **11 endpoints** organized into 4 groups:

### Health Check (1 endpoint)

| # | Method | URL | What It Does |
|---|--------|-----|-------------|
| 1 | `GET` | `/api/health` | Check if the server is running |

### User Endpoints (2 endpoints)

| # | Method | URL | What It Does |
|---|--------|-----|-------------|
| 2 | `POST` | `/api/users` | Register a new user |
| 3 | `GET` | `/api/users/:id` | Get a user's profile and their recipes |

### Recipe Endpoints (6 endpoints)

| # | Method | URL | What It Does |
|---|--------|-----|-------------|
| 4 | `POST` | `/api/recipes` | Create a new recipe (with optional image) |
| 5 | `GET` | `/api/recipes` | Get all recipes (with pagination) |
| 6 | `GET` | `/api/recipes/:id` | Get one recipe by its ID (includes ratings) |
| 7 | `GET` | `/api/recipes/search?ingredients=tomato,garlic` | Search recipes by ingredients |
| 8 | `PUT` | `/api/recipes/:id` | Update a recipe |
| 9 | `DELETE` | `/api/recipes/:id` | Delete a recipe and all its ratings |

### Rating Endpoints (2 endpoints)

| # | Method | URL | What It Does |
|---|--------|-----|-------------|
| 10 | `POST` | `/api/recipes/:id/ratings` | Add a rating to a recipe |
| 11 | `GET` | `/api/recipes/:id/ratings` | Get all ratings for a recipe |

> **Note:** `:id` means you replace it with the actual ID. For example, `/api/recipes/a1b2c3d4` where `a1b2c3d4` is the recipe's ID.

---

## 4. How To Test Every Endpoint Using Postman

### Initial Setup

1. **Start the server first:**
   ```bash
   cd D:\Go
   go mod tidy
   go run main.go
   ```
   You should see: `🍳 Recipe Sharing API` and `📍 Running on: http://localhost:8080`

2. **Open Postman** (download from [postman.com](https://www.postman.com/downloads/) if you don't have it)

3. **Import the collection:**
   - Click "Import" button in Postman
   - Select the file `D:\Go\postman_collection.json`
   - All requests will appear organized in folders

4. **Or create requests manually** (follow the steps below for each endpoint)

---

### Endpoint 1: Health Check

Check if the server is running.

- **Method:** `GET`
- **URL:** `http://localhost:8080/api/health`
- **Body:** None
- **Headers:** None

**In Postman:**
1. Click "New" → "HTTP Request"
2. Set method to `GET`
3. Type `http://localhost:8080/api/health` in the URL bar
4. Click "Send"

**Expected Response (200 OK):**
```json
{
    "success": true,
    "message": "🚀 Recipe API is running!",
    "data": {
        "status": "healthy",
        "version": "1.0.0"
    }
}
```

---

### Endpoint 2: Register a User

Create a user account. You need this before creating recipes.

- **Method:** `POST`
- **URL:** `http://localhost:8080/api/users`
- **Headers:** `Content-Type: application/json`
- **Body:** JSON

**In Postman:**
1. Set method to `POST`
2. URL: `http://localhost:8080/api/users`
3. Go to the **"Body"** tab
4. Select **"raw"**
5. In the dropdown next to "raw", select **"JSON"**
6. Paste this in the body:

```json
{
    "username": "chef_john",
    "email": "john@example.com"
}
```

7. Click "Send"

**Expected Response (201 Created):**
```json
{
    "success": true,
    "message": "User registered successfully! 👤",
    "data": {
        "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
        "username": "chef_john",
        "email": "john@example.com",
        "created_at": "2026-02-20T22:30:00Z",
        "updated_at": "2026-02-20T22:30:00Z"
    }
}
```

> ⚠️ **IMPORTANT:** Copy the `"id"` value from the response! You'll need it for creating recipes. In this example it's `"f47ac10b-58cc-4372-a567-0e02b2c3d479"` but yours will be different.

---

### Endpoint 3: Get User by ID

View a user's profile along with all their recipes.

- **Method:** `GET`
- **URL:** `http://localhost:8080/api/users/{paste-the-user-id-here}`

**In Postman:**
1. Set method to `GET`
2. URL: `http://localhost:8080/api/users/f47ac10b-58cc-4372-a567-0e02b2c3d479` (use YOUR user's ID)
3. Click "Send"

**Expected Response (200 OK):**
```json
{
    "success": true,
    "message": "User fetched successfully",
    "data": {
        "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
        "username": "chef_john",
        "email": "john@example.com",
        "created_at": "2026-02-20T22:30:00Z",
        "updated_at": "2026-02-20T22:30:00Z",
        "recipes": []
    }
}
```

---

### Endpoint 4: Create a Recipe (WITH Image)

This is the most complex endpoint. It uses `multipart/form-data` (not JSON) because we're uploading a file.

- **Method:** `POST`
- **URL:** `http://localhost:8080/api/recipes`
- **Body:** form-data (NOT JSON!)

**In Postman (Step by Step):**
1. Set method to `POST`
2. URL: `http://localhost:8080/api/recipes`
3. Go to the **"Body"** tab
4. Select **"form-data"** (NOT "raw", NOT "JSON")
5. Add these key-value pairs one by one:

| Key | Value | Type |
|-----|-------|------|
| `title` | `Spaghetti Bolognese` | Text |
| `description` | `Classic Italian pasta with rich meat sauce` | Text |
| `ingredients` | `["spaghetti","tomato","ground beef","onion","garlic","olive oil","basil"]` | Text |
| `prep_time` | `15` | Text |
| `cook_time` | `30` | Text |
| `servings` | `4` | Text |
| `user_id` | `f47ac10b-58cc-4372-a567-0e02b2c3d479` | Text |
| `image` | *(select a .jpg or .png file from your computer)* | **File** |

> **For the `image` field:** In the "Key" column, after typing `image`, hover over the right side of the Key cell — you'll see a dropdown that says "Text". Change it to **"File"**. Then in the "Value" column, click "Select Files" and choose any JPEG or PNG image from your computer.

> **For `ingredients`:** This MUST be a valid JSON array string. Always use double quotes inside the brackets: `["tomato","onion"]`. Single quotes will cause an error.

6. Click "Send"

**Expected Response (201 Created):**
```json
{
    "success": true,
    "message": "Recipe created successfully! 🎉",
    "data": {
        "id": "b2c3d4e5-f6a7-8901-bcde-f01234567890",
        "title": "Spaghetti Bolognese",
        "description": "Classic Italian pasta with rich meat sauce",
        "image_url": "/uploads/recipe_a1b2c3d4.jpg",
        "ingredients": "[\"spaghetti\",\"tomato\",\"ground beef\",\"onion\",\"garlic\",\"olive oil\",\"basil\"]",
        "prep_time": 15,
        "cook_time": 30,
        "servings": 4,
        "average_rating": 0,
        "user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
        "created_at": "2026-02-20T22:31:00Z",
        "updated_at": "2026-02-20T22:31:00Z"
    }
}
```

> ⚠️ **Copy the recipe `"id"` from this response!** You'll need it for rating, updating, and deleting.

**To view the uploaded image:** Open your browser and go to `http://localhost:8080/uploads/recipe_a1b2c3d4.jpg` (use the filename from `image_url`).

---

### Endpoint 4b: Create a Recipe (WITHOUT Image)

Same as above but skip the `image` field. Everything else stays the same.

**In Postman:**
1. Same steps as above, but DO NOT add the `image` key
2. The recipe will be created with `"image_url": ""`

Example body (form-data):

| Key | Value | Type |
|-----|-------|------|
| `title` | `Tomato Soup` | Text |
| `description` | `A warm and comforting tomato soup` | Text |
| `ingredients` | `["tomato","onion","garlic","cream","salt","pepper"]` | Text |
| `prep_time` | `10` | Text |
| `cook_time` | `25` | Text |
| `servings` | `2` | Text |
| `user_id` | `f47ac10b-58cc-4372-a567-0e02b2c3d479` | Text |

---

### Endpoint 5: Get All Recipes (Paginated)

Lists all recipes, with newest first.

- **Method:** `GET`
- **URL:** `http://localhost:8080/api/recipes?page=1&per_page=10`

**Query Parameters:**
| Parameter | Default | Description |
|-----------|---------|-------------|
| `page` | `1` | Which page of results |
| `per_page` | `10` | How many recipes per page (max 100) |

**In Postman:**
1. Set method to `GET`
2. URL: `http://localhost:8080/api/recipes?page=1&per_page=10`
3. Click "Send"

**Expected Response (200 OK):**
```json
{
    "success": true,
    "message": "Recipes fetched successfully",
    "data": [
        {
            "id": "b2c3d4e5-...",
            "title": "Spaghetti Bolognese",
            "description": "Classic Italian pasta with rich meat sauce",
            "image_url": "/uploads/recipe_a1b2c3d4.jpg",
            "ingredients": "[\"spaghetti\",\"tomato\",\"ground beef\"]",
            "prep_time": 15,
            "cook_time": 30,
            "servings": 4,
            "average_rating": 0,
            "user_id": "f47ac10b-...",
            "created_at": "2026-02-20T22:31:00Z",
            "updated_at": "2026-02-20T22:31:00Z"
        },
        {
            "id": "c3d4e5f6-...",
            "title": "Tomato Soup",
            "..."
        }
    ],
    "page": 1,
    "per_page": 10,
    "total_count": 2
}
```

---

### Endpoint 6: Get One Recipe by ID

Gets a single recipe with all its ratings included.

- **Method:** `GET`
- **URL:** `http://localhost:8080/api/recipes/{recipe-id}`

**In Postman:**
1. Set method to `GET`
2. URL: `http://localhost:8080/api/recipes/b2c3d4e5-f6a7-8901-bcde-f01234567890` (use YOUR recipe ID)
3. Click "Send"

**Expected Response (200 OK):**
```json
{
    "success": true,
    "message": "Recipe fetched successfully",
    "data": {
        "id": "b2c3d4e5-...",
        "title": "Spaghetti Bolognese",
        "description": "Classic Italian pasta with rich meat sauce",
        "image_url": "/uploads/recipe_a1b2c3d4.jpg",
        "ingredients": "[\"spaghetti\",\"tomato\",\"ground beef\",\"onion\",\"garlic\"]",
        "prep_time": 15,
        "cook_time": 30,
        "servings": 4,
        "average_rating": 4.5,
        "user_id": "f47ac10b-...",
        "created_at": "2026-02-20T22:31:00Z",
        "updated_at": "2026-02-20T22:31:00Z",
        "ratings": [
            {
                "id": "d4e5f6a7-...",
                "recipe_id": "b2c3d4e5-...",
                "user_name": "foodie_jane",
                "score": 5,
                "comment": "Best recipe ever!",
                "created_at": "2026-02-20T22:35:00Z"
            }
        ]
    }
}
```

---

### Endpoint 7: Search Recipes by Ingredients

Search for recipes that contain specific ingredients. You can search for multiple ingredients separated by commas.

- **Method:** `GET`
- **URL:** `http://localhost:8080/api/recipes/search?ingredients=tomato,garlic`

**How the search works:**
- It looks for recipes where the `ingredients` field contains ANY of the search terms
- The search is case-insensitive (`Tomato` and `tomato` both work)
- It uses SQL `LIKE` operator: `WHERE ingredients LIKE '%tomato%' OR ingredients LIKE '%garlic%'`

**In Postman:**
1. Set method to `GET`
2. URL: `http://localhost:8080/api/recipes/search?ingredients=tomato,garlic`
3. Click "Send"

**More search examples:**
- Search for one ingredient: `?ingredients=chicken`
- Search for multiple: `?ingredients=tomato,onion,garlic`
- Search for something specific: `?ingredients=ground beef`

**Expected Response (200 OK):**
```json
{
    "success": true,
    "message": "Search results fetched successfully",
    "data": {
        "query": "tomato,garlic",
        "count": 2,
        "recipes": [
            {
                "id": "b2c3d4e5-...",
                "title": "Spaghetti Bolognese",
                "ingredients": "[\"spaghetti\",\"tomato\",\"ground beef\",\"onion\",\"garlic\"]",
                "..."
            },
            {
                "id": "c3d4e5f6-...",
                "title": "Tomato Soup",
                "ingredients": "[\"tomato\",\"onion\",\"garlic\",\"cream\"]",
                "..."
            }
        ]
    }
}
```

---

### Endpoint 8: Update a Recipe

Update specific fields of a recipe. You only need to send the fields you want to change.

- **Method:** `PUT`
- **URL:** `http://localhost:8080/api/recipes/{recipe-id}`
- **Headers:** `Content-Type: application/json`
- **Body:** JSON

**In Postman:**
1. Set method to `PUT`
2. URL: `http://localhost:8080/api/recipes/b2c3d4e5-...` (use YOUR recipe ID)
3. Go to **"Body"** → **"raw"** → select **"JSON"**
4. Paste:

```json
{
    "title": "Updated Spaghetti Bolognese",
    "description": "My improved version with secret ingredients",
    "servings": 6
}
```

5. Click "Send"

> You can update any combination of fields: `title`, `description`, `ingredients`, `prep_time`, `cook_time`, `servings`, `image_url`. You CANNOT update `id`, `created_at`, or `average_rating`.

**Expected Response (200 OK):**
```json
{
    "success": true,
    "message": "Recipe updated successfully",
    "data": {
        "id": "b2c3d4e5-...",
        "title": "Updated Spaghetti Bolognese",
        "description": "My improved version with secret ingredients",
        "servings": 6,
        "..."
    }
}
```

---

### Endpoint 9: Delete a Recipe

Permanently deletes a recipe AND all ratings associated with it.

- **Method:** `DELETE`
- **URL:** `http://localhost:8080/api/recipes/{recipe-id}`

**In Postman:**
1. Set method to `DELETE`
2. URL: `http://localhost:8080/api/recipes/b2c3d4e5-...` (use YOUR recipe ID)
3. Click "Send"

**Expected Response (200 OK):**
```json
{
    "success": true,
    "message": "Recipe deleted successfully"
}
```

> ⚠️ This is permanent! The recipe and all ratings for it are removed from the database.

---

### Endpoint 10: Rate a Recipe

Add a star rating (1-5) with an optional comment.

- **Method:** `POST`
- **URL:** `http://localhost:8080/api/recipes/{recipe-id}/ratings`
- **Headers:** `Content-Type: application/json`
- **Body:** JSON

**In Postman:**
1. Set method to `POST`
2. URL: `http://localhost:8080/api/recipes/b2c3d4e5-.../ratings` (replace with YOUR recipe ID)
3. Go to **"Body"** → **"raw"** → select **"JSON"**
4. Paste:

```json
{
    "user_name": "foodie_jane",
    "score": 5,
    "comment": "Absolutely delicious! Best recipe I have tried."
}
```

5. Click "Send"

**Rules:**
- `user_name` is required
- `score` is required and must be between 1 and 5
- `comment` is optional

**Expected Response (201 Created):**
```json
{
    "success": true,
    "message": "Rating added successfully! ⭐",
    "data": {
        "rating": {
            "id": "d4e5f6a7-...",
            "recipe_id": "b2c3d4e5-...",
            "user_name": "foodie_jane",
            "score": 5,
            "comment": "Absolutely delicious! Best recipe I have tried.",
            "created_at": "2026-02-20T22:35:00Z"
        },
        "recipe": {
            "id": "b2c3d4e5-...",
            "title": "Spaghetti Bolognese",
            "average_rating": 5,
            "..."
        }
    }
}
```

> The response includes the updated recipe with the new `average_rating`. Try adding more ratings with different scores to see the average change!

---

### Endpoint 11: Get All Ratings for a Recipe

View all ratings and the average for a specific recipe.

- **Method:** `GET`
- **URL:** `http://localhost:8080/api/recipes/{recipe-id}/ratings`

**In Postman:**
1. Set method to `GET`
2. URL: `http://localhost:8080/api/recipes/b2c3d4e5-.../ratings` (use YOUR recipe ID)
3. Click "Send"

**Expected Response (200 OK):**
```json
{
    "success": true,
    "message": "Ratings fetched successfully",
    "data": {
        "recipe_id": "b2c3d4e5-...",
        "average_rating": 4.33,
        "count": 3,
        "ratings": [
            {
                "id": "d4e5f6a7-...",
                "user_name": "foodie_jane",
                "score": 5,
                "comment": "Best recipe ever!",
                "created_at": "2026-02-20T22:35:00Z"
            },
            {
                "id": "e5f6a7b8-...",
                "user_name": "home_cook_mike",
                "score": 4,
                "comment": "Pretty good, needs more salt",
                "created_at": "2026-02-20T22:36:00Z"
            },
            {
                "id": "f6a7b8c9-...",
                "user_name": "pasta_lover",
                "score": 4,
                "comment": null,
                "created_at": "2026-02-20T22:37:00Z"
            }
        ]
    }
}
```

---

## 5. How Image Upload & Compression Works

### The Full Flow (Step by Step)

```
You upload a photo (2MB, 3000x2000px)
        │
        ▼
Upload Middleware checks:
  ✓ Is it JPEG or PNG? (rejects GIF, WebP, etc.)
  ✓ Is it under 10MB? (configurable in .env)
        │
        ▼
Saves raw file to: public/temp/raw_a1b2c3d4.jpg
        │
        ▼
Controller calls ProcessImage():
  1. Opens the raw image
  2. Resizes to max 800px wide (height auto-calculated to keep proportions)
  3. Re-encodes as JPEG at 80% quality
  4. Saves to: public/temp/recipe_e5f6a7b8.jpg
  5. Deletes the raw file
        │
        ▼
Result: 200KB, 800x533px image (was 2MB, 3000x2000px)
        │
        ▼
Filename stored in database: "recipe_e5f6a7b8.jpg"
Image accessible at: http://localhost:8080/uploads/recipe_e5f6a7b8.jpg
```

### What gets reduced?

| Property | Before | After |
|----------|--------|-------|
| File size | ~2MB | ~200KB (90% smaller) |
| Width | 3000px | 800px |
| Height | 2000px | 533px (auto-calculated) |
| Format | JPEG/PNG | Always JPEG |
| Quality | 100% | 80% |

### Configuration (in .env file)

```env
UPLOAD_DIR=./public/temp     # Where images are saved
MAX_UPLOAD_SIZE=10           # Maximum upload size in MB
IMG_MAX_WIDTH=800            # Max width after resize (in pixels)
IMG_QUALITY=80               # JPEG quality (1-100, higher = better quality but bigger file)
```

---

## 6. How The Project Is Organized

```
D:\Go\
│
├── main.go                        ← App starts here. Loads config, connects DB, starts server.
├── .env                           ← Settings (port, DB path, image settings)
├── go.mod                         ← Lists all dependencies (like package.json)
├── recipe.db                      ← Database file (created automatically when you run the app)
│
├── public/temp/                   ← Processed images are saved here
│
└── src/
    ├── db/
    │   └── db.go                  ← Connects to SQLite, creates tables
    │
    ├── models/                    ← Define the shape of our data
    │   ├── user.model.go          ← User struct (id, username, email)
    │   ├── recipe.model.go        ← Recipe struct (title, ingredients, image, etc.)
    │   └── rating.model.go        ← Rating struct (score 1-5, comment)
    │
    ├── controllers/               ← The actual logic for each endpoint
    │   ├── recipe.controller.go   ← Create, Read, Update, Delete, Search recipes
    │   ├── rating.controller.go   ← Add ratings, get ratings, calculate average
    │   └── user.controller.go     ← Register user, get user profile
    │
    ├── middlewares/                ← Code that runs BEFORE the controller
    │   ├── error.middleware.go    ← Catches crashes and returns clean error JSON
    │   └── upload.middleware.go   ← Handles file uploads (like multer in Node.js)
    │
    ├── routes/                    ← Maps URLs to controller functions
    │   ├── index.routes.go        ← Combines all routes + health check
    │   ├── recipe.routes.go       ← /api/recipes/* routes
    │   ├── rating.routes.go       ← /api/recipes/:id/ratings routes
    │   └── user.routes.go         ← /api/users/* routes
    │
    └── utils/                     ← Helper functions used across the app
        ├── image.util.go          ← Resize and compress images
        ├── response.util.go       ← Standard JSON response format
        └── async.util.go          ← Safe background task runner
```

### How a request flows through the code:

```
Client sends request
    → main.go (Gin router receives it)
    → routes/ (finds the matching URL pattern)
    → middlewares/ (runs upload middleware if needed, error handler wraps everything)
    → controllers/ (executes the business logic)
    → models/ + db/ (reads/writes to the database)
    → utils/ (uses helpers for image processing, response formatting)
    → Response sent back to client
```

---

## 7. How To Run This Project

### Prerequisites

1. Install Go from [https://go.dev/dl/](https://go.dev/dl/)
2. Verify: open terminal and type `go version` — you should see something like `go version go1.21.0 windows/amd64`

### Steps

```bash
# Navigate to the project
cd D:\Go

# Download all dependencies (this reads go.mod and downloads the packages)
go mod tidy

# Run the server
go run main.go
```

### What happens when you run it:

1. Reads `.env` file for configuration
2. Creates/connects to `recipe.db` (SQLite database file)
3. Auto-creates the `users`, `recipes`, and `ratings` tables
4. Creates the `public/temp/` directory for image uploads
5. Starts the HTTP server on port 8080
6. Prints a startup banner with the URL

### You should see this output:

```
✅ Database connected successfully (SQLite)
✅ Database tables migrated successfully
==============================================
  🍳 Recipe Sharing API
  📍 Running on: http://localhost:8080
  📦 Database:   SQLite (./recipe.db)
  📁 Uploads:    ./public/temp
==============================================
```

---

## 8. Presentation Tips

### Recommended Testing Order in Postman:

1. ✅ Hit Health Check → shows server is alive
2. ✅ Register a User → get the user ID
3. ✅ Create Recipe with Image → show image compression working
4. ✅ Create Recipe without Image → show it's optional
5. ✅ Get All Recipes → show pagination
6. ✅ Search by Ingredients → "find me recipes with tomato"
7. ✅ Rate a Recipe → add 2-3 ratings with different scores
8. ✅ Get Ratings → show average calculation
9. ✅ Get Recipe by ID → show ratings embedded in recipe
10. ✅ Update Recipe → partial update
11. ✅ Delete Recipe → clean deletion with cascade

### Key Points to Mention:

- **"Zero database setup"** — SQLite is embedded, no installation needed
- **"Automatic image optimization"** — Uploaded images are resized and compressed automatically (saves 60-70% storage)
- **"Smart search"** — Search recipes by ingredient names
- **"Auto-calculated ratings"** — Average rating updates automatically when new ratings are added
- **"Consistent API"** — Every response follows the same `{ success, message, data }` format
