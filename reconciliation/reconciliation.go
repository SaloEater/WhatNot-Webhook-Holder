// Package reconciliation holds the pure (no I/O) matching rules used to compare Whatnot orders
// against the break/event rows of a stream. It mirrors the frontend whatnot-verify page.
package reconciliation

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	MatchedByNumber = "number"
	MatchedByNone   = "none"
	MatchedByManual = "manual"
	MatchedBySkip   = "skipped"

	StatusMismatch      = "mismatch"
	StatusPriceMismatch = "price_mismatch"

	// GiveawayGap is the silence after the last team order that makes a giveaway belong to the next break.
	GiveawayGap = 10 * time.Minute
)

func Norm(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

// ---- break matching ----

type DBBreak struct {
	Id   int64
	Name string
}

type BreakMatch struct {
	BreakId   int64 // 0 when not matched
	BreakName string
	By        string // number | none | manual | skipped
}

var breakPhrase = regexp.MustCompile(`(?i)\bbreak\s+0*(\d+)\b`)

// breakNumber returns N from the first "Break N" in the name ("Break 03 Titans 3" -> "3"), ok=false when
// the phrase is absent. Other numbers in the name are ignored.
func breakNumber(name string) (string, bool) {
	m := breakPhrase.FindStringSubmatch(name)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// SkipBreak is the manual-map value for a Whatnot break the human chose to skip.
const SkipBreak int64 = -1

// MatchBreaksWithOverrides applies the human's decisions first, then auto-matches the remaining names
// against the remaining DB breaks. An override pointing to a break that no longer exists is ignored.
func MatchBreaksWithOverrides(wnNames []string, dbBreaks []DBBreak, overrides map[string]int64) []BreakMatch {
	out := make([]BreakMatch, len(wnNames))
	byId := map[int64]DBBreak{}
	for _, b := range dbBreaks {
		byId[b.Id] = b
	}
	used := map[int64]bool{}
	restIdx := []int{}
	restNames := []string{}
	for i, name := range wnNames {
		id, has := overrides[name]
		b, known := byId[id]
		switch {
		case has && id == SkipBreak:
			out[i] = BreakMatch{BreakId: SkipBreak, By: MatchedBySkip}
		case has && known && !used[id]:
			used[id] = true
			out[i] = BreakMatch{BreakId: b.Id, BreakName: b.Name, By: MatchedByManual}
		default:
			restIdx = append(restIdx, i)
			restNames = append(restNames, name)
		}
	}
	restBreaks := []DBBreak{}
	for _, b := range dbBreaks {
		if !used[b.Id] {
			restBreaks = append(restBreaks, b)
		}
	}
	for j, m := range MatchBreaks(restNames, restBreaks) {
		out[restIdx[j]] = m
	}
	return out
}

// MatchBreaks maps every Whatnot break name to a DB break by the "Break N" phrase in both names. A number
// carried by more than one Whatnot name or more than one DB break stays unmatched for a human to resolve,
// because a wrong guess would write fixes to the wrong break. The result is aligned with wnNames.
func MatchBreaks(wnNames []string, dbBreaks []DBBreak) []BreakMatch {
	wnCount := map[string]int{}
	for _, name := range wnNames {
		if n, ok := breakNumber(name); ok {
			wnCount[n]++
		}
	}
	dbByNumber := map[string][]DBBreak{}
	for _, b := range dbBreaks {
		if n, ok := breakNumber(b.Name); ok {
			dbByNumber[n] = append(dbByNumber[n], b)
		}
	}
	out := make([]BreakMatch, len(wnNames))
	for i, name := range wnNames {
		out[i].By = MatchedByNone
		n, ok := breakNumber(name)
		if !ok || wnCount[n] != 1 || len(dbByNumber[n]) != 1 {
			continue
		}
		b := dbByNumber[n][0]
		out[i] = BreakMatch{BreakId: b.Id, BreakName: b.Name, By: MatchedByNumber}
	}
	return out
}

// ---- teams ----

type Order struct {
	OrderId  string
	Team     string
	Username string
	Price    float64
}

type Event struct {
	Id         int64
	Team       string
	Customer   string
	Price      float64
	IsGiveaway bool
}

type Fix struct {
	Order  Order
	Event  Event
	Status string // mismatch | price_mismatch
}

type TeamResult struct {
	Matched int
	Fixes   []Fix
	ToAdd   []Order // no free event with this team: added to the break as a new event
	Extra   []Event
}

func cents(p float64) int64 { return int64(math.Round(p * 100)) }

func isRealCustomer(c string) bool {
	c = strings.TrimSpace(c)
	return c != "" && c != "?"
}

// CompareTeams: each order claims the first unclaimed non-giveaway event with the same team.
func CompareTeams(orders []Order, events []Event) TeamResult {
	res := TeamResult{Fixes: []Fix{}, ToAdd: []Order{}, Extra: []Event{}}
	team := make([]Event, 0, len(events))
	for _, e := range events {
		if !e.IsGiveaway {
			team = append(team, e)
		}
	}
	claimed := make([]bool, len(team))
	for _, o := range orders {
		idx := -1
		for i, e := range team {
			if !claimed[i] && e.Team == o.Team {
				idx = i
				break
			}
		}
		if idx < 0 {
			res.ToAdd = append(res.ToAdd, o)
			continue
		}
		claimed[idx] = true
		ev := team[idx]
		switch {
		case Norm(ev.Customer) != Norm(o.Username):
			res.Fixes = append(res.Fixes, Fix{Order: o, Event: ev, Status: StatusMismatch})
		case cents(ev.Price) != cents(o.Price):
			res.Fixes = append(res.Fixes, Fix{Order: o, Event: ev, Status: StatusPriceMismatch})
		default:
			res.Matched++
		}
	}
	for i, e := range team {
		if !claimed[i] && isRealCustomer(e.Customer) {
			res.Extra = append(res.Extra, e)
		}
	}
	return res
}

// ---- giveaways ----

type TeamOrder struct {
	WnBreak   string
	CreatedAt time.Time
}

type Giveaway struct {
	OrderId   string
	Title     string
	Username  string
	CreatedAt time.Time
}

// BucketGiveaways returns, aligned with giveaways, the Whatnot break name each giveaway belongs to
// ("" when there is no team order at all), going by the nearest team orders around it. Between two orders
// of the same break, it's that break. Between orders of two breaks, the timestamp decides: more than
// GiveawayGap after the last order means the next break, otherwise the previous one. With an order on only
// one side, it's that order's break.
// Timeline = team orders + giveaways + other timestamps
// (unparsed orders, which are neither before nor after anchors); sorted by time, stable.
// gap(last team order before) > GiveawayGap and a team order exists after -> the one after, else
// the one before (else the one after).
func BucketGiveaways(teamOrders []TeamOrder, giveaways []Giveaway, otherTimes []time.Time) []string {
	type entry struct {
		at   time.Time
		wn   string // team order break name
		give int    // giveaway index, -1 otherwise
	}
	tl := make([]entry, 0, len(teamOrders)+len(giveaways)+len(otherTimes))
	for _, o := range teamOrders {
		tl = append(tl, entry{at: o.CreatedAt, wn: o.WnBreak, give: -1})
	}
	for i, g := range giveaways {
		tl = append(tl, entry{at: g.CreatedAt, give: i})
	}
	for _, t := range otherTimes {
		tl = append(tl, entry{at: t, give: -1})
	}
	sort.SliceStable(tl, func(a, b int) bool { return tl[a].at.Before(tl[b].at) })

	out := make([]string, len(giveaways))
	for i, e := range tl {
		if e.give < 0 {
			continue
		}
		var before, after *entry
		for j := i - 1; j >= 0 && before == nil; j-- {
			if tl[j].wn != "" {
				before = &tl[j]
			}
		}
		for j := i + 1; j < len(tl) && after == nil; j++ {
			if tl[j].wn != "" {
				after = &tl[j]
			}
		}
		pick := before
		if before == nil || (e.at.Sub(before.at) > GiveawayGap && after != nil) {
			pick = after
		}
		if pick != nil {
			out[e.give] = pick.wn
		}
	}
	return out
}

type GiveawayDiff struct {
	Missing []Giveaway
	Extra   []Event
}

// DiffGiveaways compares the buyers of the giveaways placed in a break with its is_giveaway events
// (multiset by exact buyer/customer).
func DiffGiveaways(placed []Giveaway, events []Event) GiveawayDiff {
	d := GiveawayDiff{Missing: []Giveaway{}, Extra: []Event{}}
	byCustomer := map[string][]Event{}
	for _, e := range events {
		if e.IsGiveaway {
			byCustomer[e.Customer] = append(byCustomer[e.Customer], e)
		}
	}
	for _, g := range placed {
		if list := byCustomer[g.Username]; len(list) > 0 {
			byCustomer[g.Username] = list[1:]
			continue
		}
		d.Missing = append(d.Missing, g)
	}
	// keep deterministic order: by event id
	for _, list := range byCustomer {
		d.Extra = append(d.Extra, list...)
	}
	sort.Slice(d.Extra, func(a, b int) bool { return d.Extra[a].Id < d.Extra[b].Id })
	return d
}

// IsSlabTitle: giveaway_type Slab when the title contains "slb" or "mag" (case-insensitive), else Pack.
func IsSlabTitle(title string) bool {
	t := strings.ToLower(title)
	return strings.Contains(t, "slb") || strings.Contains(t, "mag")
}
