package whatnot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const baseURL = "https://api.whatnot.com/seller-api/rest/v2"
const minInterval = 110 * time.Millisecond

type Client struct {
	token string
	http  *http.Client

	mu   sync.Mutex
	last time.Time
}

func New(token string) *Client {
	return &Client{token: token, http: &http.Client{Timeout: 30 * time.Second}}
}

type Show struct {
	Id               string     `json:"id"`
	Title            string     `json:"title"`
	Status           string     `json:"status"`
	ScheduledStartAt *time.Time `json:"scheduled_start_at"`
	EndedAt          *time.Time `json:"ended_at"`
	Description      string     `json:"description"`
	Visibility       string     `json:"visibility"`
	// CreatedAt and UpdatedAt stay raw strings so an odd value never fails decoding the whole show.
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Thumbnail *struct {
		Url string `json:"url"`
	} `json:"thumbnail"`
}

type Order struct {
	Id              string
	Status          string
	ShowId          string
	TransactionType string
	Username        string
	Total           float64
	CreatedAt       time.Time
}

type OrderItem struct {
	Title     string
	Quantity  int
	UnitPrice float64
}

type money struct {
	Amount string `json:"amount"`
}

func (m money) float() float64 {
	f, _ := strconv.ParseFloat(m.Amount, 64)
	return f
}

func (c *Client) throttle() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if wait := minInterval - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func (c *Client) get(path string, q url.Values, out any) error {
	if c.token == "" {
		return errors.New("whatnot_api_key is not set")
	}
	u := baseURL + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	for attempt := 0; ; attempt++ {
		c.throttle()
		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt == 0 {
			secs, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
			if secs <= 0 {
				secs = 1
			}
			time.Sleep(time.Duration(secs) * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("whatnot api %s: status %d: %s", path, resp.StatusCode, string(body))
		}
		return json.Unmarshal(body, out)
	}
}

// ListShowsPage fetches one page of /shows and returns the cursor of the next page ("" when there is none).
func (c *Client) ListShowsPage(status, cursor string, pageSize int) ([]Show, string, error) {
	q := url.Values{}
	q.Set("status", status)
	q.Set("limit", strconv.Itoa(pageSize))
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var resp struct {
		Data       []Show  `json:"data"`
		NextCursor *string `json:"next_cursor"`
	}
	if err := c.get("/shows", q, &resp); err != nil {
		return nil, "", err
	}
	if resp.NextCursor == nil {
		return resp.Data, "", nil
	}
	return resp.Data, *resp.NextCursor, nil
}

// ListShowsSince walks /shows newest-first in pages of pageSize, keeping shows until one starts before
// since. A zero since keeps everything. maxPages only guards against an endless cursor walk.
//
// When hasMore is true, resumeCursor is where older shows continue. If the walk stopped at since, that is
// the cursor of the page it stopped on, because the rest of that page hasn't been returned yet. Resuming
// therefore repeats the newer shows of that page, and the cursor is "" if it stopped on the first page.
func (c *Client) ListShowsSince(status string, since time.Time, pageSize, maxPages int) (shows []Show, resumeCursor string, hasMore bool, err error) {
	cursor := ""
	for page := 0; page < maxPages; page++ {
		data, next, err := c.ListShowsPage(status, cursor, pageSize)
		if err != nil {
			return nil, "", false, err
		}
		for _, sh := range data {
			if !since.IsZero() && sh.ScheduledStartAt != nil && sh.ScheduledStartAt.Before(since) {
				return shows, cursor, true, nil
			}
			shows = append(shows, sh)
		}
		if next == "" {
			return shows, "", false, nil
		}
		cursor = next
	}
	return shows, cursor, true, nil
}

func (c *Client) GetShow(id string) (*Show, error) {
	var show Show
	if err := c.get("/shows/"+url.PathEscape(id), nil, &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// ListOrdersForShow walks the seller's orders newest-first and keeps the ones of the show. There is no page
// limit: the walk ends at notBefore or at the last page, so an old show takes longer but is never cut short.
func (c *Client) ListOrdersForShow(showId string, notBefore time.Time) ([]Order, error) {
	var orders []Order
	cursor := ""
	for {
		q := url.Values{}
		q.Set("limit", "50")
		if cursor != "" {
			q.Set("cursor", cursor)
		}
		var resp struct {
			Data []struct {
				Id              string `json:"id"`
				Status          string `json:"status"`
				ShowId          string `json:"show_id"`
				TransactionType string `json:"transaction_type"`
				Buyer           struct {
					Username string `json:"username"`
				} `json:"buyer"`
				Amounts struct {
					Total money `json:"total"`
				} `json:"amounts"`
				CreatedAt time.Time `json:"created_at"`
			} `json:"data"`
			NextCursor *string `json:"next_cursor"`
		}
		if err := c.get("/orders", q, &resp); err != nil {
			return nil, err
		}
		oldest := time.Time{}
		for _, o := range resp.Data {
			if oldest.IsZero() || o.CreatedAt.Before(oldest) {
				oldest = o.CreatedAt
			}
			if o.ShowId != showId {
				continue
			}
			orders = append(orders, Order{
				Id:              o.Id,
				Status:          o.Status,
				ShowId:          o.ShowId,
				TransactionType: o.TransactionType,
				Username:        o.Buyer.Username,
				Total:           o.Amounts.Total.float(),
				CreatedAt:       o.CreatedAt,
			})
		}
		if !oldest.IsZero() && oldest.Before(notBefore) {
			return orders, nil
		}
		if len(resp.Data) == 0 || resp.NextCursor == nil || *resp.NextCursor == "" {
			return orders, nil
		}
		cursor = *resp.NextCursor
	}
}

func (c *Client) ListOrderItems(orderId string) ([]OrderItem, error) {
	q := url.Values{}
	q.Set("limit", "50")
	var resp struct {
		Data []struct {
			PurchasedItem struct {
				Title string `json:"title"`
			} `json:"purchased_item"`
			Quantity  int   `json:"quantity"`
			UnitPrice money `json:"unit_price"`
		} `json:"data"`
	}
	if err := c.get("/orders/"+url.PathEscape(orderId)+"/items", q, &resp); err != nil {
		return nil, err
	}
	items := make([]OrderItem, len(resp.Data))
	for i, d := range resp.Data {
		items[i] = OrderItem{Title: d.PurchasedItem.Title, Quantity: d.Quantity, UnitPrice: d.UnitPrice.float()}
	}
	return items, nil
}
