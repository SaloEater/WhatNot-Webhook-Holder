package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/SaloEater/WhatNot-Webhook-Holder/discord"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/SaloEater/WhatNot-Webhook-Holder/whatnot"
)

// /backfill-show is temporary and is meant to be removed together with this file.
const backfillShowCommand = "backfill-show"

// RegisterDiscordCommands registers the slash commands as server commands in every server the bot is in.
// Silent no-op without the bot token or application id; failures are only logged. Meant to run in a goroutine.
func (s *Service) RegisterDiscordCommands(applicationId string) {
	defer func() {
		if r := recover(); r != nil {
			discordLogf("register commands: panic: %v", r)
		}
	}()
	if !s.discordEnabled() || applicationId == "" {
		return
	}
	guilds, err := s.Discord.ListGuildIds()
	if err != nil {
		discordLogf("register commands: %s", err.Error())
		return
	}
	commands := []map[string]any{{
		"name":        backfillShowCommand,
		"description": "Add Whatnot shows to this feed (all scheduled, or one by id)",
		"type":        1,
		"options": []map[string]any{{
			"type":        3,
			"name":        "id",
			"description": "Whatnot show id",
			"required":    false,
		}},
	}}
	for _, g := range guilds {
		if err := s.Discord.SetGuildCommands(applicationId, g, commands); err != nil {
			discordLogf("register commands in guild %s: %s", g, err.Error())
		}
	}
}

// HandleDiscordCommand runs an application command in the background and always ends by replacing the deferred
// reply with a final message. The caller answers the interaction itself.
func (s *Service) HandleDiscordCommand(in discord.Interaction) {
	go func() {
		result := "Something went wrong"
		defer func() {
			if r := recover(); r != nil {
				discordLogf("command %q: panic: %v", in.Data.Name, r)
				result = "Something went wrong"
			}
			if err := s.Discord.EditInteractionResponse(in.ApplicationId, in.Token, truncateRunes(result, 2000)); err != nil {
				discordLogf("command %q: reply: %s", in.Data.Name, err.Error())
			}
		}()
		if in.Data.Name != backfillShowCommand {
			result = "Unknown command"
			return
		}
		showId := ""
		for _, o := range in.Data.Options {
			if o.Name == "id" {
				if v, ok := o.Value.(string); ok {
					showId = strings.TrimSpace(v)
				}
			}
		}
		result = s.backfillShows(in.ChannelId, showId)
	}()
}

type backfillFailure struct{ title, reason string }

// backfillShows posts the show (or all scheduled shows) of the feed's seller account to the feed and returns the summary.
func (s *Service) backfillShows(discordChannelId, showId string) string {
	feed, err := s.DiscordFeedChannelRepositorier.GetByDiscordChannelId(discordChannelId)
	if err != nil {
		discordLogf("backfill: feed lookup: %s", err.Error())
		return "Could not look up the feed: " + err.Error()
	}
	if feed == nil {
		return "This command must be run in a show feed channel."
	}
	client, err := s.whatnotClientForChannel(feed.ChannelId)
	if err != nil {
		return "Whatnot is not available for this feed: " + err.Error()
	}

	var shows []whatnot.Show
	if showId != "" {
		wn, err := client.GetShow(showId)
		if err != nil {
			return "Could not fetch the show: " + err.Error()
		}
		if wn.Id == "" {
			wn.Id = showId
		}
		shows = []whatnot.Show{*wn}
	} else {
		cursor := ""
		for page := 0; page < whatnotShowsMaxPages; page++ {
			data, next, err := client.ListShowsPage("scheduled", cursor, whatnotShowsPageSize)
			if err != nil {
				return "Could not list scheduled shows: " + err.Error()
			}
			shows = append(shows, data...)
			if next == "" {
				break
			}
			cursor = next
		}
		sort.SliceStable(shows, func(i, j int) bool {
			a, b := shows[i].ScheduledStartAt, shows[j].ScheduledStartAt
			if a == nil || b == nil {
				return a != nil && b == nil
			}
			return a.Before(*b)
		})
	}
	if len(shows) == 0 {
		return "No scheduled shows found."
	}

	var added []string
	var failed []backfillFailure
	skippedExisting, skippedCancelled := 0, 0
	for i := range shows {
		wn := &shows[i]
		title := wn.Title
		if title == "" {
			title = wn.Id
		}
		outcome, reason := s.backfillShow(feed.ChannelId, wn)
		switch outcome {
		case "added":
			added = append(added, title)
		case "existing":
			skippedExisting++
		case "cancelled":
			skippedCancelled++
		default:
			failed = append(failed, backfillFailure{title, reason})
		}
	}
	return backfillSummary(added, skippedExisting, skippedCancelled, failed)
}

// backfillShow returns "added", "existing", "cancelled" or "failed" (with the reason).
func (s *Service) backfillShow(channelId int64, wn *whatnot.Show) (string, string) {
	row, err := s.WhatnotShowRepositorier.GetByShowId(channelId, wn.Id)
	if err != nil {
		return "failed", err.Error()
	}
	if row != nil {
		if row.Status == entity.WhatnotShowStatusCancelled {
			return "cancelled", ""
		}
		if s.hasDiscordMessage(row.Id) {
			return "existing", ""
		}
	}
	if wn.Status == "cancelled" {
		return "cancelled", ""
	}
	show := whatnotShowFromAPI(channelId, wn)
	if err := s.WhatnotShowRepositorier.Upsert(show); err != nil {
		return "failed", err.Error()
	}
	s.invalidateShowStreams(show.Id)
	s.SyncDiscordShow(show.Id)
	if !s.hasDiscordMessage(show.Id) {
		return "failed", "message was not posted, see the backend log"
	}
	return "added", ""
}

func (s *Service) hasDiscordMessage(whatnotShowId int64) bool {
	msg, err := s.DiscordShowMessageRepositorier.Get(whatnotShowId)
	return err == nil && msg != nil && msg.DiscordMessageId != nil && *msg.DiscordMessageId != ""
}

// backfillSummary keeps the counts and shrinks the title lists until the text fits a Discord message.
func backfillSummary(added []string, skippedExisting, skippedCancelled int, failed []backfillFailure) string {
	build := func(limit int) string {
		var parts []string
		line := fmt.Sprintf("Added %d", len(added))
		if len(added) > 0 && limit > 0 {
			line += ": " + joinLimited(added, limit)
		}
		parts = append(parts, line+".")
		parts = append(parts, fmt.Sprintf("Skipped %d already in the feed.", skippedExisting))
		if skippedCancelled > 0 {
			parts = append(parts, fmt.Sprintf("Skipped %d cancelled.", skippedCancelled))
		}
		if len(failed) > 0 {
			items := make([]string, len(failed))
			for i, f := range failed {
				items[i] = truncateRunes(f.title, 60) + " (" + truncateRunes(f.reason, 120) + ")"
			}
			line := fmt.Sprintf("Failed %d", len(failed))
			if limit > 0 {
				line += ": " + joinLimited(items, limit)
			}
			parts = append(parts, line+".")
		}
		return strings.Join(parts, " ")
	}
	for limit := 50; limit > 0; limit /= 2 {
		if out := build(limit); len([]rune(out)) <= 1900 {
			return out
		}
	}
	return truncateRunes(build(0), 2000)
}

func joinLimited(items []string, limit int) string {
	shown := items
	if len(shown) > limit {
		shown = shown[:limit]
	}
	out := make([]string, len(shown))
	for i, t := range shown {
		out[i] = truncateRunes(t, 150)
	}
	text := strings.Join(out, ", ")
	if len(items) > len(shown) {
		text += fmt.Sprintf(" and %d more", len(items)-len(shown))
	}
	return text
}
