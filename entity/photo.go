package entity

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Photo struct {
	Id        int64     `json:"id"         db:"id"`
	SeriesId  int64     `json:"series_id"  db:"series_id"`
	Name      string    `json:"name"       db:"name"`
	Team      string    `json:"team"       db:"team"`
	Url       string    `json:"url"        db:"url"`
	Thumbnail string    `json:"thumbnail"  db:"thumbnail"`
	IsSold    bool      `json:"is_sold"    db:"is_sold"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
	Price     int64     `json:"price"      db:"price"`
	Rotation  int64     `json:"rotation"   db:"rotation"`
	Version   int64     `json:"-"          db:"version"`
	ArtX      int64     `json:"-"          db:"art_x"`
	ArtY      int64     `json:"-"          db:"art_y"`
	ArtW      int64     `json:"-"          db:"art_w"`
	ArtH      int64     `json:"-"          db:"art_h"`
	LabelX    int64     `json:"-"          db:"label_x"`
	LabelY    int64     `json:"-"          db:"label_y"`
	LabelW    int64     `json:"-"          db:"label_w"`
	LabelH    int64     `json:"-"          db:"label_h"`
	LabelText string    `json:"-"          db:"label_text"`
}

type Bounds struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
	W int64 `json:"w"`
	H int64 `json:"h"`
}

func (p Photo) ArtBounds() *Bounds {
	if p.ArtW == 0 || p.ArtH == 0 {
		return nil
	}
	return &Bounds{X: p.ArtX, Y: p.ArtY, W: p.ArtW, H: p.ArtH}
}

func (p Photo) LabelBounds() *Bounds {
	if p.LabelW == 0 || p.LabelH == 0 {
		return nil
	}
	return &Bounds{X: p.LabelX, Y: p.LabelY, W: p.LabelW, H: p.LabelH}
}

type LabelText struct {
	Kind string       `json:"kind"`
	Rows [4][2]string `json:"rows"`
}

var labelKindRe = regexp.MustCompile(`^[a-z0-9_-]{1,16}$`)

// Normalize trims cells and validates the shape in place.
func (l *LabelText) Normalize() error {
	l.Kind = strings.TrimSpace(l.Kind)
	if !labelKindRe.MatchString(l.Kind) {
		return fmt.Errorf("kind must be a short lowercase token")
	}
	for i := range l.Rows {
		for j := range l.Rows[i] {
			c := strings.TrimSpace(l.Rows[i][j])
			if utf8.RuneCountInString(c) > 64 {
				return fmt.Errorf("cell [%d][%d] longer than 64 chars", i, j)
			}
			l.Rows[i][j] = c
		}
	}
	return nil
}

// ParseLabelText validates a label_text JSON document (exactly 4 rows x 2 cells).
func ParseLabelText(s string) (*LabelText, error) {
	var raw struct {
		Kind string     `json:"kind"`
		Rows [][]string `json:"rows"`
	}
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		return nil, err
	}
	if len(raw.Rows) != 4 {
		return nil, fmt.Errorf("rows must have exactly 4 entries")
	}
	l := &LabelText{Kind: raw.Kind}
	for i, r := range raw.Rows {
		if len(r) != 2 {
			return nil, fmt.Errorf("row %d must have exactly 2 cells", i)
		}
		copy(l.Rows[i][:], r)
	}
	if err := l.Normalize(); err != nil {
		return nil, err
	}
	return l, nil
}

func (p Photo) MarshalJSON() ([]byte, error) {
	type plain Photo
	var labelText *LabelText
	if p.LabelText != "" {
		if lt, err := ParseLabelText(p.LabelText); err == nil {
			labelText = lt
		}
	}
	return json.Marshal(struct {
		plain
		Version     int64      `json:"version"`
		ArtBounds   *Bounds    `json:"art_bounds"`
		LabelBounds *Bounds    `json:"label_bounds"`
		LabelText   *LabelText `json:"label_text"`
	}{plain(p), p.Version, p.ArtBounds(), p.LabelBounds(), labelText})
}
