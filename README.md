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
