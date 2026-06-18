# Backend Implementation Plan — Series Prices (PLAN_SERIES_PRICE)

Spec source: `WhatNot-Webhook-Holder/PLAN_SERIES_PRICE.md` (original one-liner spec above this plan).

---

## Summary

Replace the manual `series_team_price` table with a `price` field on each `Photo` record.
A new computed endpoint sums unsold photo prices per team for a given series and caches the result.
The cache is invalidated whenever a photo is marked sold or unsold.

---

## Step 1 — Delete the `series_team_price` stack

Delete these 9 source files (migration 20260504000011 stays on disk; Migration G in Step 2 writes the revert):

```
entity/series_team_price.go
repository/series_team_price_repository.go
repository/repository_sqlx/series_team_price_repository.go
service/series_team_price_get.go
service/series_team_price_set.go
service/series_team_price_get_last_prices.go
api/series_team_price_get.go
api/series_team_price_set.go
api/series_team_price_get_last_prices.go
```

---

## Step 2 — Database Migrations

### Migration F — Add `price` to `photo`

File: `20260504000012_add_price_to_photo.up.sql` / `.down.sql`

```sql
-- up
ALTER TABLE photo ADD COLUMN price INTEGER NOT NULL DEFAULT 0;

-- down
ALTER TABLE photo DROP COLUMN price;
```

### Migration G — Drop `series_team_price`

File: `20260504000013_drop_series_team_price.up.sql` / `.down.sql`

```sql
-- up
DROP TABLE IF EXISTS series_team_price;

-- down
CREATE TABLE series_team_price (
    id        BIGSERIAL PRIMARY KEY,
    series_id BIGINT    NOT NULL REFERENCES series(id),
    team      VARCHAR   NOT NULL,
    price     NUMERIC   NOT NULL DEFAULT 0,
    UNIQUE (series_id, team)
);
```

---

## Step 3 — Entities

### Update `entity/photo.go`

Add field:

```go
Price int64 `json:"price" db:"price"`
```

### New `entity/series_team_total.go`

```go
package entity

type SeriesTeamTotal struct {
    Team  string `json:"team"  db:"team"`
    Price int64  `json:"price" db:"price"`
}
```

---

## Step 4 — Repository

### `repository/photo_repository.go` — add method to interface

```go
GetPricesBySeriesId(seriesId int64) ([]*entity.SeriesTeamTotal, error)
```

### `repository/repository_sqlx/photo_repository.go` — implement it

```go
func (r *PhotoRepository) GetPricesBySeriesId(seriesId int64) ([]*entity.SeriesTeamTotal, error) {
    totals := []*entity.SeriesTeamTotal{}
    err := r.DB.Select(&totals, `
        SELECT team, SUM(price) AS price
        FROM photo
        WHERE series_id = $1
          AND is_sold = false
          AND is_deleted = false
        GROUP BY team
        ORDER BY team
    `, seriesId)
    return totals, err
}
```

---

## Step 5 — Service

### Update `service/service.go`

Remove:
```go
repository.SeriesTeamPriceRepositorier
```

Add cache field:
```go
SeriesPricesCache cacheInterface.Cache[[]*entity.SeriesTeamTotal]
```

### New `service/series_get_prices.go`

```go
package service

import (
    "github.com/SaloEater/WhatNot-Webhook-Holder/cache"
    "github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type SeriesGetPricesRequest struct {
    SeriesId int64 `json:"series_id"`
}

func (s *Service) SeriesGetPrices(r *SeriesGetPricesRequest) ([]*entity.SeriesTeamTotal, error) {
    key := cache.IdToKey(r.SeriesId)
    if s.SeriesPricesCache.Has(key) {
        cached, _ := s.SeriesPricesCache.Get(key)
        return cached, nil
    }
    totals, err := s.PhotoRepositorier.GetPricesBySeriesId(r.SeriesId)
    if err == nil {
        s.SeriesPricesCache.Set(key, totals)
    }
    return totals, err
}
```

### Update `service/photo_mark_sold.go` — add SeriesId, evict cache

Add `SeriesId int64 \`json:"series_id"\`` to `PhotoMarkSoldRequest`.

After a successful `MarkSold` call, add:
```go
s.SeriesPricesCache.Delete(cache.IdToKey(r.SeriesId))
```

---

## Step 6 — API Handler

### New `api/series_get_prices.go`

Uses Go 1.22 path values to read `series_id` from the URL:

```go
package api

import (
    "net/http"
    "strconv"
    "github.com/SaloEater/WhatNot-Webhook-Holder/service"
)

func (a *API) SeriesGetPrices(w http.ResponseWriter, r *http.Request) (any, error) {
    id, err := strconv.ParseInt(r.PathValue("series_id"), 10, 64)
    if err != nil {
        return nil, err
    }
    return a.Service.SeriesGetPrices(&service.SeriesGetPricesRequest{SeriesId: id})
}
```

---

## Step 7 — `main.go` Wiring

Remove:
```go
SeriesTeamPriceRepositorier: &repository_sqlx.SeriesTeamPriceRepository{DB: db},
```

Remove routes:
```go
http.HandleFunc("/api/series/team_prices", ...)
http.HandleFunc("/api/series/team_price/set", ...)
http.HandleFunc("/api/series/team_price/last_prices", ...)
```

Add cache init alongside breakCache/streamCache/channelCache:
```go
seriesPricesCache := go_cache.CreateCache[[]*entity.SeriesTeamTotal](10 * time.Hour)
```

Add to `svc` struct:
```go
SeriesPricesCache: &seriesPricesCache,
```

Add route:
```go
http.HandleFunc("/api/series/{series_id}/prices", routeBuilder.WrapRoute(apiO.SeriesGetPrices, api.HttpGet, true))
```

---

## Verification

```bash
go build ./...
```

1. Upload photo to series 1 with `price=100`, `team="NYY"` → confirm `Photo.Price` persists.
2. `GET /api/series/1/prices` → `{"data":[{"team":"NYY","price":100}],"error":""}`.
3. Mark that photo sold (`POST /api/photo/mark_sold` with `series_id=1, id=<id>, sold=true`).
4. `GET /api/series/1/prices` → `{"data":[{"team":"NYY","price":0}],"error":""}` (cache evicted, unsold sum is 0).
