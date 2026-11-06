package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SaloEater/WhatNot-Webhook-Holder/discord"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
)

// DiscordInteractions is Discord's Interactions Endpoint. Registered as "POST /discord/interactions";
// deliberately not wrapped by WrapRoute, because the signature covers the raw body.
func DiscordInteractions(svc *service.Service, publicKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if publicKey == "" || !discord.VerifyInteraction(publicKey, r.Header.Get("X-Signature-Ed25519"), r.Header.Get("X-Signature-Timestamp"), body) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var interaction discord.Interaction
		if err := json.Unmarshal(body, &interaction); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responseType := discord.ResponseTypeDeferredUpdateMsg
		switch interaction.Type {
		case discord.InteractionTypePing:
			responseType = discord.ResponseTypePong
		case discord.InteractionTypeMessageComponent:
			svc.HandleDiscordViewButton(interaction.Data.CustomId)
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprintf(w, `{"type":%d}`, responseType); err != nil {
			fmt.Println("discord interaction: write response: " + err.Error())
		}
	}
}
