# ⭐ Review Service

> The trust layer of **AirBnB-Node** — lets guests review a hotel, but only once they've actually stayed there.

Review Service owns one thing: reviews. But it enforces a rule that matters more than the CRUD around it — you can only leave a review for a **confirmed booking** at a **real hotel** by a **real user**, and it checks all three against the other services before ever writing a row, rather than trusting whatever the client sends.

---

## Where this fits in the system

```
                POST /api/v1/reviews
                          │
                          ▼
        ┌───────────────────────────────┐
        │   Review Service (Go)           │   ← this repo
        └───────────────┬─────────────────┘
                          │  validates against, in order:
          ┌───────────────┼────────────────┐
          ▼                ▼                 ▼
 ┌─────────────────┐ ┌─────────────────┐ ┌────────────────────┐
 │  API Gateway       │ │  HotelService     │ │  BookingService      │
 │  does user exist?  │ │  does hotel exist? │ │  does booking exist,  │
 │                    │ │                   │ │  is it "confirmed"?  │
 └─────────────────┘ └─────────────────┘ └────────────────────┘
                          │
                          ▼
                     ┌───────────┐
                     │  MySQL      │
                     │  reviews    │
                     └───────────┘
```

Like `BookingService`, this service owns no user, hotel, or booking data itself — every fact it needs about those comes from an HTTP call to the service that actually owns it.

---

## What it does

- **Three-way precondition check before writing a review.** `CreateReview` calls the API Gateway to confirm the user exists, `HotelService` to confirm the hotel exists, and `BookingService` to confirm the booking exists — and doesn't stop there.
- **"Confirmed bookings only" enforcement.** After confirming the booking exists, the service parses the booking response and checks `status == "confirmed"`. A `pending` or `cancelled` booking gets a `403 Forbidden` — you can't review a stay you never actually completed.
- **Hotel-scoped review listings**, with the same hotel-existence check applied before querying, so a typo'd hotel ID fails fast with a clear `404` instead of silently returning an empty list.
- **Rating and text constraints enforced at the DTO level.** Ratings are bounded `1–5`; review text must be between 50 and 1000 characters — long enough to be a real review, short enough to stay a review and not an essay.
- **An `is_synced` flag on every review**, hinting at a downstream consumer (e.g. a search index, an aggregate-rating cache, or an analytics pipeline) that reads unsynced reviews and marks them synced once processed — the sync mechanism itself isn't in this repo.

---

## Review creation flow

```
POST /api/v1/reviews
        │
        ▼
  validate payload (rating 1–5, text 50–1000 chars)
        │
        ▼
  GET {API_GATEWAY}/users/{userId}         ──► 404 if user doesn't exist
        │
        ▼
  GET {HOTEL_SERVICE}/hotels/{hotelId}      ──► 404 if hotel doesn't exist
        │
        ▼
  GET {BOOKING_SERVICE}/bookings/{bookingId} ──► 404 if booking doesn't exist
        │
        ▼
  booking.status == "confirmed"?  ──► 403 if not
        │
        ▼
  INSERT review row
```

---

## Tech stack

| Concern | Choice |
|---|---|
| Language | Go 1.26 |
| Router | `go-chi/chi` |
| Rate limiting | `go-chi/httprate` |
| Validation | `go-playground/validator` |
| Database | MySQL (`database/sql` + `go-sql-driver/mysql`) |
| Logging | `go.uber.org/zap` |
| Migrations | Goose-style timestamped SQL files |
| Live reload (dev) | [Air](https://github.com/air-verse/air) |

---

## Project structure

```
cmd/app/               # App composition: wires config, DB, router into an http.Server
main.go                # Process entrypoint
internal/
├── config/               # Env loading, server config, DB config, zap logger setup
├── controllers/           # HTTP handlers for reviews
├── services/              # Business logic — cross-service validation lives here
├── repositories/           # SQL queries against MySQL
├── database/
│   ├── models/               # ReviewModel row struct
│   └── migrations/            # Timestamped .sql migrations (Goose)
├── dtos/                  # Request/response payload structs (incl. the minimal
│                            #  FetchBookingDTO used to read a booking's status)
├── middlewares/            # Rate limiting, body/param validation
├── routers/                # Route registration
├── utils/                  # Errors, JSON responses
```

Same layered shape as `AuthService`/`API Gateway` and consistent with the Node services elsewhere in the system: **router → controller → service → repository → model**, with interfaces at each boundary.

---

## API surface (v1)

| Method | Route | Purpose |
|---|---|---|
| `POST` | `/api/v1/reviews` | Create a review (validates user, hotel, and confirmed booking first) |
| `GET` | `/api/v1/reviews/hotel/{hotel_id}` | List all reviews for a hotel |
| `GET` | `/api/v1/reviews/{id}` | Get a single review |
| `PUT` | `/api/v1/reviews/{id}` | Update a review's rating/text |
| `DELETE` | `/api/v1/reviews/{id}` | Delete a review |

---

## Getting started

### Prerequisites

- Go 1.26+
- MySQL running locally or reachable
- [Air](https://github.com/air-verse/air) (optional, for live reload)
- The API Gateway, `HotelService`, and `BookingService` reachable, since creating or listing reviews calls out to all three

### Configure environment

Create a `.env` file in the project root:

```bash
# Server
ADDR=:3030
READ_TIMEOUT=15
WRITE_TIMEOUT=15
IDLE_TIMEOUT=180
APP_ENV=development
REQUESTS_PER_MINUTE=100
JWT_SECRET_KEY=change-me

# Database
DB_USERNAME=admin
DB_PASSWORD=admin
DB_NET=tcp
DB_ADDRESS=127.0.0.1:3306
DB_NAME=dev_db

# Goose (migrations)
GOOSE_DRIVER=mysql
GOOSE_DBSTRING=admin:admin@tcp(127.0.0.1:3306)/dev_db
GOOSE_MIGRATION_DIR=internal/database/migrations

# Inter-service
API_GATEWAY_BASE_URL=http://localhost:3020/api/v1
HOTEL_SERVICE_BASE_URL=http://localhost:3000/api/v1
BOOKING_SERVICE_BASE_URL=http://localhost:3010/api/v1
```

### Run migrations

```bash
goose status         # or: task goose:status
goose up              # or: task goose:up-all
```

### Start the server

```bash
go run main.go
```

Or via the Taskfile:

```bash
task server:run
```

For live reload during development:

```bash
air
# or: task server:dev
```

By default the server listens on `:3030` (override with `ADDR`).

---

## Design notes worth knowing

- **Review creation is read-heavy on purpose.** Three synchronous HTTP calls happen before a single row is written. That's a deliberate trade of latency for correctness — a fake or premature review is worse for the platform than a slightly slower write.
- **Booking status is the actual gatekeeper, not booking existence.** It would be easy to stop at "does this booking ID exist" — the service goes one step further and decodes the booking response to check `status == "confirmed"`, which is what actually prevents someone from reviewing a stay that was cancelled or never happened.
- **No ownership check between user and booking (yet).** The current flow confirms the user, hotel, and booking each exist and that the booking is confirmed, but doesn't cross-check that `reviewPayload.UserID` is the same user who made `reviewPayload.BookingID`. Worth keeping in mind if reviews need to be locked to "the person who actually booked this."
- **`is_synced` is a seam, not a feature.** It's persisted but nothing in this repo flips it — it's set up for a future/external consumer (search indexing, rating aggregation, etc.) to pick up unsynced reviews and mark them processed.
- **The health router is currently commented out** in `routers/main.go` (`//SetupHealthRouter(router)`) — worth re-enabling if you're wiring this into orchestration/monitoring that expects a `/health` endpoint, matching the other Go services in this system.