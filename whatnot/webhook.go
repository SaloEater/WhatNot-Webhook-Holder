package whatnot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	webhookSecretPrefix = "wn_wh_"
	webhookTolerance    = 5 * time.Minute
)

// VerifyWebhook checks a Whatnot (Standard Webhooks) delivery: returns nil when one of the v1
// signatures matches HMAC-SHA256(base64decode(secret), "{id}.{timestamp}.{body}") and the timestamp
// is within +-5 minutes of now.
func VerifyWebhook(secret string, headers http.Header, body []byte, now time.Time) error {
	id := headers.Get("webhook-id")
	ts := headers.Get("webhook-timestamp")
	sigHeader := headers.Get("webhook-signature")
	if id == "" || ts == "" || sigHeader == "" {
		return errors.New("missing headers")
	}

	key, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(secret, webhookSecretPrefix))
	if err != nil || len(key) == 0 {
		return errors.New("bad secret")
	}

	sec, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return errors.New("bad timestamp")
	}
	diff := now.Sub(time.Unix(sec, 0))
	if diff > webhookTolerance || diff < -webhookTolerance {
		return errors.New("timestamp out of tolerance")
	}

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(id + "." + ts + "."))
	mac.Write(body)
	expected := []byte(base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	for _, entry := range strings.Fields(sigHeader) {
		version, sig, ok := strings.Cut(entry, ",")
		if !ok || version != "v1" {
			continue
		}
		if hmac.Equal([]byte(sig), expected) {
			return nil
		}
	}
	return errors.New("signature mismatch")
}
