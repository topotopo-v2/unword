## UNWORD
**Words without a perfect translation, one word a day.** 

[Unword](https://unword-6i0m.onrender.com/) is a minimalist word of the day focused on difficult to translate words from different languages. 
<br>

### Features
* Display one word per day
* Show native script or word, pronunciation, language, country of origin, definition, and source.
* Save and unsave words locally (in browser's localStorage)

### Tech Stack
**Frontend:** React, TypeScript, Vite<br>
**Backend:** Go, pgx, REST<br>
**Database:** PostgresSQL<br>
**Infra:** Docker, Render for Static Site, Web Service, PostgresSQL

### Local Setup
#### Backend
1. Go to backend directory
````
cd backend
````
2. Create a local `.env.` Refer to `.env.example`<br>
3. Start PostgresSQL
````
docker compose up -d postgres
````
4. Run migrations
````
migrate \ -path ./migrations \ -database "postgres://unword:unword@localhost:5432/unword?sslmode=disable" \ up
````
5. Seed the database
````
go run ./cmd/seed
````
6. Start
````
go run ./cmd/server
````

### Frontend
1. Go to backend directory
````
cd frontend
````
2. Create local `.env`. Refer to `.env.example`
3. Install dependencies
````
npm install
````
4. Start Vite
````
npm run dev
````

### How I made this

I built this website as a way to ease back into coding after almost a year-long break. My goal wasn’t to build everything completely by hand, but to learn and relearn the fundamentals of taking a website from an idea to production.

I started by brainstorming the idea with ChatGPT, then created a rough UI in Figma. From there, I worked through the backend setup using Go, PostgreSQL, database migrations, APIs, and Docker before moving on to the frontend.

Throughout the project, I used ChatGPT more like a guide and pair-programming partner. I asked for explanations, broke larger tasks into smaller steps, discussed implementation options, and got help debugging when I got stuck. I still wrote and tested the code myself, made the implementation decisions, and tried to understand each part before moving on.

The goal of this project is less about building something complex and more about rebuilding my confidence and understanding the full process of developing and shipping a small web application.

