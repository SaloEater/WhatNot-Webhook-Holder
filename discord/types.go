package discord

import "encoding/json"

const (
	InteractionTypePing             = 1
	InteractionTypeMessageComponent = 3

	ResponseTypePong              = 1
	ResponseTypeDeferredUpdateMsg = 6
	ComponentTypeActionRow        = 1
	ComponentTypeButton           = 2
	ButtonStylePrimary            = 1
	ButtonStyleSecondary          = 2
	ButtonStyleLink               = 5
)

// Message is the payload of create/edit message calls.
type Message struct {
	Embeds          []Embed         `json:"embeds"`
	Components      []ActionRow     `json:"components"`
	AllowedMentions AllowedMentions `json:"allowed_mentions"`
	// Content and MessageReference are used only by the alert ping.
	Content          string            `json:"content,omitempty"`
	MessageReference *MessageReference `json:"message_reference,omitempty"`
}

type MessageReference struct {
	MessageId       string `json:"message_id"`
	FailIfNotExists bool   `json:"fail_if_not_exists"`
}

// AllowedMentions marshals as {"users": [...]} when Users is set, else as {"parse": []} (nobody is notified).
type AllowedMentions struct {
	Parse []string
	Users []string
}

func (a AllowedMentions) MarshalJSON() ([]byte, error) {
	if len(a.Users) > 0 {
		return json.Marshal(struct {
			Users []string `json:"users"`
		}{a.Users})
	}
	parse := a.Parse
	if parse == nil {
		parse = []string{}
	}
	return json.Marshal(struct {
		Parse []string `json:"parse"`
	}{parse})
}

type Embed struct {
	Title       string       `json:"title,omitempty"`
	URL         string       `json:"url,omitempty"`
	Description string       `json:"description,omitempty"`
	Color       int          `json:"color,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Image       *EmbedImage  `json:"image,omitempty"`
	Footer      *EmbedFooter `json:"footer,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type EmbedImage struct {
	URL string `json:"url"`
}

type EmbedFooter struct {
	Text string `json:"text"`
}

type ActionRow struct {
	Type       int         `json:"type"`
	Components []Component `json:"components"`
}

type Component struct {
	Type     int    `json:"type"`
	Style    int    `json:"style"`
	Label    string `json:"label"`
	CustomId string `json:"custom_id,omitempty"`
	URL      string `json:"url,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

// Interaction is the part of an incoming interaction this backend reads.
type Interaction struct {
	Type int `json:"type"`
	Data struct {
		CustomId string `json:"custom_id"`
	} `json:"data"`
}
