# study-space-api

REST API for [StudySpace](https://github.com/BYUI-Study-Space/study-space), the WDD 231 team project. It serves study locations and stores planned study sessions in Postgres.

Built with Go 1.27, [Fiber v3](https://gofiber.io), [pgx](https://github.com/jackc/pgx) and [goose](https://github.com/pressly/goose) migrations, which run automatically on startup.

## Run locally

Requirements: Go 1.27+ and Docker (OrbStack or Docker Desktop).

```sh
cp .env.example .env   # first time only
make db-up             # start Postgres in Docker
make run               # API on http://localhost:3000
make test              # go vet + unit tests
make db-down           # stop Postgres
```

### Run with Docker

Only Docker is needed; this builds the API image and starts it with Postgres.

```sh
make up                # API on http://localhost:3000
make down              # stop everything
```

To push the image to a registry (Docker Hub shown; build for the server's platform):

```sh
docker buildx build --platform linux/amd64 -t <user>/study-space-api:latest --push .
```

## Endpoints

All responses are JSON. Errors look like `{"error": "..."}`; validation errors also include a `fields` object keyed by input name.

| Method | Path | Description |
|---|---|---|
| GET | `/api/health` | `{"status":"ok"}` when the API and database are up |
| GET | `/api/locations` | All locations. Optional filters: `wifi=true\|false`, `group=true\|false`, `noise=quiet\|moderate\|collaborative` |
| GET | `/api/locations/:id` | One location |
| POST | `/api/sessions` | Plan a session (body below). Returns `201` with the saved session |
| GET | `/api/sessions/:id` | A saved session, by the UUID returned on create |

### POST /api/sessions

```json
{
  "location_id": 2,
  "student_name": "Ada Lovelace",
  "email": "ada@byui.edu",
  "student_id": "123456789",
  "date": "2026-10-20",
  "start_time": "15:00",
  "duration_minutes": 120,
  "session_type": "group",
  "group_size": 4,
  "notes": "CSE 210 review"
}
```

Rules: all fields except `group_size` and `notes` are required. `duration_minutes` must be 30, 60, 90, 120 or 180. The session can't start in the past (campus time) and must fit within the location's hours. Group sessions need a group-friendly location and a `group_size` between 2 and the location's capacity. `notes` can be up to 500 characters.

A validation failure returns `400`:

```json
{ "error": "validation failed", "fields": { "email": "must be a valid email" } }
```

## Configuration

| Variable | Default | Purpose |
|---|---|---|
| `DATABASE_URL` | — (required) | Postgres connection string |
| `PORT` | `3000` | HTTP port (Railway sets it) |
| `ALLOWED_ORIGINS` | `http://localhost:5500,http://127.0.0.1:5500` | Comma-separated CORS origins |
| `TIMEZONE` | `America/Denver` | Campus timezone for "not in the past" checks |

## Deploy to Railway

1. **New Project → Deploy from GitHub repo** and pick `BYUI-Study-Space/study-space-api`. Railway builds the `Dockerfile`.
2. In the project, **+ New → Database → PostgreSQL**.
3. In the API service's **Variables**, add:
   - `DATABASE_URL` = `${{Postgres.DATABASE_URL}}`
   - `ALLOWED_ORIGINS` = `https://byui-study-space.github.io`
4. In **Settings**, set the healthcheck path to `/api/health` and **Generate Domain** under Networking.

Every push to `main` redeploys automatically.
