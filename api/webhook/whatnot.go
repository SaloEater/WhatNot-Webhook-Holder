package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"github.com/SaloEater/WhatNot-Webhook-Holder/whatnot"
)

// Whatnot receives signed Whatnot webhook deliveries: verify -> persist -> 204 -> process async.
// Registered as "POST /webhook/whatnot/{channel_id}"; deliberately not wrapped by WrapRoute.
func Whatnot(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channelId, err := strconv.ParseInt(r.PathValue("channel_id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		channel, err := svc.ChannelRepositorier.Get(channelId)
		if err != nil || channel == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if channel.WhatnotWebhookSecret == "" {
			fmt.Println(fmt.Sprintf("whatnot webhook: channel %d has no whatnot_webhook_secret", channelId))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if err != nil {
			fmt.Println("whatnot webhook: read body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := whatnot.VerifyWebhook(channel.WhatnotWebhookSecret, r.Header, body, time.Now()); err != nil {
			fmt.Println(fmt.Sprintf("whatnot webhook: channel %d: verify: %s", channelId, err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var envelope struct {
			Type      string          `json:"type"`
			Timestamp string          `json:"timestamp"`
			SellerId  json.RawMessage `json:"seller_id"`
		}
		webhookId := r.Header.Get("webhook-id")
		if err := json.Unmarshal(body, &envelope); err != nil || envelope.Type == "" || len(webhookId) > 64 {
			fmt.Println(fmt.Sprintf("whatnot webhook: channel %d: bad envelope", channelId))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ev := &entity.WhatnotWebhookEvent{
			WebhookId: webhookId,
			ChannelId: channelId,
			Type:      envelope.Type,
			SellerId:  strings.Trim(string(envelope.SellerId), `"`),
			Payload:   body,
		}
		if envelope.SellerId == nil || string(envelope.SellerId) == "null" {
			ev.SellerId = ""
		}
		if t, err := time.Parse(time.RFC3339Nano, envelope.Timestamp); err == nil {
			ev.EventTime = &t
		}

		inserted, err := svc.InsertWebhookEvent(ev)
		if err != nil {
			fmt.Println(fmt.Sprintf("whatnot webhook: channel %d: insert: %s", channelId, err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		if !inserted {
			return
		}
		go svc.ProcessWhatnotWebhook(ev.Id)
	}
}
