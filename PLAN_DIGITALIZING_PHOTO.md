# Backend Implementation Plan — Card Digitalization (DIGITALIZING_PHOTO_1)

Spec sources: `No-Mod-Livestream/DIGITALIZING_PHOTO_1.md`, `SERIES_AND_IMAGE_FLOW.md`, `CARDS_BOARD_PAGE.md`, `Card-Scanner/PLANNING.md`.

---

## Summary

Add two new domain entities — **Series** and **Photo** — plus the API surface needed by:

1. The Card-Scanner desktop app (create series, upload photos, close series).
2. The Spots Page on the frontend (list series, link a series to a break or stream).
3. The Cards Board Page (fetch unsold photos for a channel's active livestream; mark a photo sold/unsold).
4. The Series Team Prices page (set and retrieve per-team prices for a series).

---

## Step 1 — Database Migrations

Create three migration pairs in `db/migrations/`. Use the next available prefix after `20260504000006`.

### Migration A — `series` table

```sql
-- up
CREATE TABLE series (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT      NOT NULL,
    status     TEXT      NOT NULL DEFAULT 'open',  -- 'open' | 'closed'
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_deleted BOOLEAN   NOT NULL DEFAULT false
);

-- down
DROP TABLE series;
```

### Migration B — `photo` table

```sql
-- up
CREATE TABLE photo (
    id          BIGSERIAL PRIMARY KEY,
    series_id   BIGINT    NOT NULL REFERENCES series(id),
    name        TEXT      NOT NULL DEFAULT '',
    team        VARCHAR   NOT NULL DEFAULT '',
    url         TEXT      NOT NULL,
    is_sold     BOOLEAN   NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_deleted  BOOLEAN   NOT NULL DEFAULT false
);

-- down
DROP TABLE photo;
```

### Migration D — add `team` to `photo`

```sql
-- up
ALTER TABLE photo ADD COLUMN team VARCHAR NOT NULL DEFAULT '';

-- down
ALTER TABLE photo DROP COLUMN team;
```

Files: `20260504000010_add_team_to_photo.up.sql` / `20260504000010_add_team_to_photo.down.sql`

### Migration E — `series_team_price` table

```sql
-- up
CREATE TABLE series_team_price (
    id        BIGSERIAL PRIMARY KEY,
    series_id BIGINT    NOT NULL REFERENCES series(id),
    team      VARCHAR   NOT NULL,
    price     NUMERIC   NOT NULL DEFAULT 0,
    UNIQUE (series_id, team)
);

-- down
DROP TABLE series_team_price;
```

Files: `20260504000011_create_series_team_price.up.sql` / `20260504000011_create_series_team_price.down.sql`

### Migration C — add `series_id` to `break`

```sql
-- up
ALTER TABLE break ADD COLUMN series_id BIGINT REFERENCES series(id);

-- down
ALTER TABLE break DROP COLUMN series_id;
```

---

## Step 2 — Entities

### `entity/series.go`

```go
package entity

import "time"

type Series struct {
    Id        int64     `json:"id"         db:"id"`
    Name      string    `json:"name"       db:"name"`
    Status    string    `json:"status"     db:"status"` // "open" | "closed"
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}
```

### `entity/photo.go`

```go
package entity

import "time"

type Photo struct {
    Id          int64     `json:"id"          db:"id"`
    SeriesId    int64     `json:"series_id"   db:"series_id"`
    Name        string    `json:"name"        db:"name"`
    Team        string    `json:"team"        db:"team"`
    Url         string    `json:"url"         db:"url"`
    IsSold      bool      `json:"is_sold"     db:"is_sold"`
    CreatedAt   time.Time `json:"created_at"  db:"created_at"`
    IsDeleted   bool      `json:"is_deleted"  db:"is_deleted"`
}
```

### `entity/series_team_price.go`

```go
package entity

type SeriesTeamPrice struct {
    Id       int64   `json:"id"        db:"id"`
    SeriesId int64   `json:"series_id" db:"series_id"`
    Team     string  `json:"team"      db:"team"`
    Price    float64 `json:"price"     db:"price"`
}
```

Also add `SeriesId *int64` to `entity/break.go`:

```go
SeriesId *int64 `json:"series_id" db:"series_id"`
```

---

## Step 3 — Repository Interfaces

### `repository/series_repository.go`

```go
package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type SeriesRepositorier interface {
    Create(name string) (int64, error)
    Get(id int64) (*entity.Series, error)
    GetList() ([]*entity.Series, error)
    Update(id int64, name string) error
    Close(id int64) error
    Delete(id int64) error
    CountOpen() (int, error)  // enforces "one open series at a time" rule
}
```

### `repository/series_team_price_repository.go`

```go
package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type SeriesTeamPriceRepositorier interface {
    Set(seriesId int64, team string, price float64) error
    GetBySeriesId(seriesId int64) ([]*entity.SeriesTeamPrice, error)
}
```

### `repository/photo_repository.go`

```go
package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type PhotoRepositorier interface {
    Create(p *entity.Photo) (int64, error)
    GetBySeriesId(seriesId int64) ([]*entity.Photo, error)
    GetUnsoldByChannelActiveStream(channelId int64) ([]*entity.Photo, error)
    MarkSold(id int64, sold bool) error
    Delete(id int64) error
}
```

---

## Step 4 — Repository Implementations (sqlx)

### `repository/repository_sqlx/series_team_price_repository.go`

- `Set`: `INSERT INTO series_team_price (series_id, team, price) VALUES ($1, $2, $3) ON CONFLICT (series_id, team) DO UPDATE SET price = EXCLUDED.price`
- `GetBySeriesId`: `SELECT * FROM series_team_price WHERE series_id = $1 ORDER BY team`

### `repository/repository_sqlx/series_repository.go`

Implement `SeriesRepositorier`. Key queries:

- `Create`: `INSERT INTO series (name) VALUES ($1) RETURNING id`
- `Get`: `SELECT * FROM series WHERE id = $1 AND is_deleted = false`
- `GetList`: `SELECT * FROM series WHERE is_deleted = false ORDER BY created_at DESC`
- `Update`: `UPDATE series SET name = $1 WHERE id = $2`
- `Close`: `UPDATE series SET status = 'closed' WHERE id = $1`
- `Delete`: `UPDATE series SET is_deleted = true WHERE id = $1`
- `CountOpen`: `SELECT COUNT(*) FROM series WHERE status = 'open' AND is_deleted = false`

### `repository/repository_sqlx/photo_repository.go`

Implement `PhotoRepositorier`. Key queries:

- `Create`: `INSERT INTO photo (series_id, name, team, url) VALUES (...) RETURNING id`
- `GetBySeriesId`: `SELECT * FROM photo WHERE series_id = $1 AND is_deleted = false`
- `GetUnsoldByChannelActiveStream`: join chain `channel → active stream → breaks → series (status='closed') → photo (is_sold=false, is_deleted=false)`:

```sql
SELECT p.*
FROM photo p
JOIN series s ON s.id = p.series_id
JOIN break b ON b.series_id = s.id
JOIN stream st ON st.id = b.day_id
JOIN channel c ON c.active_stream_id = st.id
WHERE c.id = $1
  AND s.status = 'closed'
  AND p.is_sold = false
  AND p.is_deleted = false
  AND s.is_deleted = false
  AND b.is_deleted = false
```

- `MarkSold`: `UPDATE photo SET is_sold = $1 WHERE id = $2`
- `Delete`: `UPDATE photo SET is_deleted = true WHERE id = $1`

---

## Step 5 — Digital Ocean Extension

Add a `SaveCardPhoto` method to `DigitalOcean` in `digital_ocean/digital_ocean.go`:

```go
func (d *DigitalOcean) SaveCardPhoto(data []byte, seriesID int64, filename string) (string, error) {
    filePath := fmt.Sprintf("cards/%d/%s", seriesID, filename)
    object := s3.PutObjectInput{
        Bucket: aws.String("mount-olympus-storage"),
        Key:    aws.String(filePath),
        Body:   bytes.NewReader(data),
        ACL:    aws.String("public-read"),
        Metadata: map[string]*string{
            "x-amz-meta-success": aws.String("ok"),
        },
    }
    _, err := d.s3Client.PutObject(&object)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("%s/%s", d.spacesURL, filePath), nil
}
```

Extend the `DigitalOceaner` interface in `service/digital_ocean.go`:

```go
type DigitalOceaner interface {
    SaveLabel(bytes.Buffer, string) (string, error)
    SaveCardPhoto(data []byte, seriesID int64, filename string) (string, error)
}
```

---

## Step 6 — Service Methods

Add `SeriesRepositorier`, `PhotoRepositorier`, and `SeriesTeamPriceRepositorier` to the `Service` struct in `service/service.go`.

One service file per operation (existing convention):

| File | Method | Notes |
|------|--------|-------|
| `service/series_create.go` | `SeriesCreate(name string)` | Reject if `CountOpen() > 0` |
| `service/series_get.go` | `SeriesGet(id int64)` | |
| `service/series_get_list.go` | `SeriesGetList()` | |
| `service/series_update.go` | `SeriesUpdate(id int64, name string)` | |
| `service/series_close.go` | `SeriesClose(id int64)` | |
| `service/series_delete.go` | `SeriesDelete(id int64)` | |
| `service/photo_upload.go` | `PhotoUpload(seriesID int64, data []byte, name, team, filename string)` | Calls `DigitalOceaner.SaveCardPhoto`, then `PhotoRepositorier.Create` |
| `service/photo_get_by_series.go` | `PhotoGetBySeries(seriesID int64)` | |
| `service/photo_delete.go` | `PhotoDelete(id int64)` | |
| `service/photo_mark_sold.go` | `PhotoMarkSold(id int64, sold bool)` | |
| `service/photo_get_for_board.go` | `PhotoGetForBoard(channelID int64)` | Returns unsold photos from closed series linked to channel's active stream |
| `service/break_set_series.go` | `BreakSetSeries(breakID, seriesID int64)` | Updates `break.series_id` |
| `service/series_team_price_get.go` | `SeriesTeamPriceGet(seriesId int64)` | Returns all team prices for a series |
| `service/series_team_price_set.go` | `SeriesTeamPriceSet(seriesId int64, team string, price float64)` | Upserts one team price |

---

## Step 7 — API Handlers

One handler file per operation. Register all routes in `main.go`.

### Series endpoints

| File | Route | Method | Handler |
|------|-------|--------|---------|
| `api/series_create.go` | `POST /api/series/create` | POST | `SeriesCreate` |
| `api/series_get.go` | `POST /api/series/get` | POST | `SeriesGet` |
| `api/series_get_list.go` | `GET /api/series/list` | GET | `SeriesGetList` |
| `api/series_update.go` | `POST /api/series/update` | POST | `SeriesUpdate` |
| `api/series_close.go` | `POST /api/series/close` | POST | `SeriesClose` |
| `api/series_delete.go` | `POST /api/series/delete` | POST | `SeriesDelete` |

### Photo endpoints

| File | Route | Method | Handler | Notes |
|------|-------|--------|---------|-------|
| `api/photo_upload.go` | `POST /api/photo/upload` | POST | `PhotoUpload` | Multipart form: `file`, `series_id`, `name` (optional), `team` (optional) |
| `api/photo_get_by_series.go` | `POST /api/photo/list` | POST | `PhotoGetBySeries` | Body: `{series_id}` |
| `api/photo_delete.go` | `POST /api/photo/delete` | POST | `PhotoDelete` | Body: `{id}` |
| `api/photo_mark_sold.go` | `POST /api/photo/mark_sold` | POST | `PhotoMarkSold` | Body: `{id, sold}` |
| `api/photo_get_for_board.go` | `POST /api/photo/board` | POST | `PhotoGetForBoard` | Body: `{channel_id}` — used by Cards Board Page |

### Break series link endpoint

| File | Route | Method | Handler |
|------|-------|--------|---------|
| `api/break_set_series.go` | `POST /api/break/set_series` | POST | `BreakSetSeries` |

### Series team price endpoints

| File | Route | Method | Handler | Notes |
|------|-------|--------|---------|-------|
| `api/series_team_price_get.go` | `POST /api/series/team_prices` | POST | `SeriesTeamPriceGet` | Body: `{series_id}` |
| `api/series_team_price_set.go` | `POST /api/series/team_price/set` | POST | `SeriesTeamPriceSet` | Body: `{series_id, team, price}` |

### Photo upload handler detail

`photo_upload.go` must parse a multipart form instead of a JSON body (unlike all other handlers). Use `r.ParseMultipartForm(32 << 20)` then `r.FormFile("file")` to read the image bytes. Other fields (`series_id`, `name`, `team`) come from `r.FormValue`.

---

## Step 8 — Wire Up in `main.go`

Add the new repositories to the `Service` struct initialisation:

```go
svc := &service.Service{
    // ... existing fields ...
    SeriesRepositorier: &repository_sqlx.SeriesRepository{DB: db},
    PhotoRepositorier:  &repository_sqlx.PhotoRepository{DB: db},
}
```

Register new routes:

```go
http.HandleFunc("/api/series/create",  routeBuilder.WrapRoute(apiO.SeriesCreate,      api.HttpPost, true))
http.HandleFunc("/api/series/get",     routeBuilder.WrapRoute(apiO.SeriesGet,         api.HttpPost, true))
http.HandleFunc("/api/series/list",    routeBuilder.WrapRoute(apiO.SeriesGetList,     api.HttpGet,  true))
http.HandleFunc("/api/series/update",  routeBuilder.WrapRoute(apiO.SeriesUpdate,      api.HttpPost, true))
http.HandleFunc("/api/series/close",   routeBuilder.WrapRoute(apiO.SeriesClose,       api.HttpPost, true))
http.HandleFunc("/api/series/delete",  routeBuilder.WrapRoute(apiO.SeriesDelete,      api.HttpPost, true))

http.HandleFunc("/api/photo/upload",   routeBuilder.WrapRoute(apiO.PhotoUpload,       api.HttpPost, true))
http.HandleFunc("/api/photo/list",     routeBuilder.WrapRoute(apiO.PhotoGetBySeries,  api.HttpPost, true))
http.HandleFunc("/api/photo/delete",   routeBuilder.WrapRoute(apiO.PhotoDelete,       api.HttpPost, true))
http.HandleFunc("/api/photo/mark_sold",routeBuilder.WrapRoute(apiO.PhotoMarkSold,     api.HttpPost, true))
http.HandleFunc("/api/photo/board",    routeBuilder.WrapRoute(apiO.PhotoGetForBoard,  api.HttpPost, true))

http.HandleFunc("/api/break/set_series", routeBuilder.WrapRoute(apiO.BreakSetSeries, api.HttpPost, true))

http.HandleFunc("/api/series/team_prices",    routeBuilder.WrapRoute(apiO.SeriesTeamPriceGet, api.HttpPost, true))
http.HandleFunc("/api/series/team_price/set", routeBuilder.WrapRoute(apiO.SeriesTeamPriceSet, api.HttpPost, true))
```

---

## Implementation Order

1. Migrations (A, B, C)
2. Entities (`series.go`, `photo.go`, add `SeriesId` to `break.go`)
3. Repository interfaces
4. Repository sqlx implementations
5. Digital Ocean extension (`SaveCardPhoto` + interface update)
6. Service struct update + service files
7. API handler files
8. Wire up in `main.go`

Each step compiles cleanly before moving to the next — the interfaces act as the compile-time contract.

---

## Open Questions

- Should `CountOpen()` enforce one-open-series globally, or per-channel? Spec implies global (single operator app). Implement globally for now.
- What happens to the DO Spaces object when a photo is deleted? Current plan: soft-delete in DB only; DO cleanup is out of scope.
- `channel.active_stream_id` — verify this column name matches the actual schema (check `entity/channel.go` and the migration that added it).
- The multipart photo upload bypasses `WrapRoute`'s JSON response wrapping for the request body parsing, but the response can still use the same `Response{Data, Error}` envelope. No change to `WrapRoute` needed.
