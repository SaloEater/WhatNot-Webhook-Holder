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

// ListShowsSince walks /shows newest-first in pages of pageSize, keeping shows until one starts before
// since. A zero since keeps everything. maxPages only guards against an endless cursor walk.
func (c *Client) ListShowsSince(status string, since time.Time, pageSize, maxPages int) ([]Show, error) {
	var shows []Show
	cursor := ""
	for page := 0; page < maxPages; page++ {
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
			return nil, err
		}
		for _, sh := range resp.Data {
			if !since.IsZero() && sh.ScheduledStartAt != nil && sh.ScheduledStartAt.Before(since) {
				return shows, nil
			}
			shows = append(shows, sh)
		}
		if resp.NextCursor == nil || *resp.NextCursor == "" {
			break
		}
		cursor = *resp.NextCursor
	}
	return shows, nil
}

func (c *Client) GetShow(id string) (*Show, error) {
	var show Show
	if err := c.get("/shows/"+url.PathEscape(id), nil, &show); err != nil {
		return nil, err
	}
	return &show, nil
}

func (c *Client) ListOrdersForShow(showId string, notBefore time.Time, maxPages int) ([]Order, error) {
	var orders []Order
	cursor := ""
	for page := 0; page < maxPages; page++ {
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
			break
		}
		if resp.NextCursor == nil || *resp.NextCursor == "" {
			break
		}
		cursor = *resp.NextCursor
	}
	return orders, nil
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
