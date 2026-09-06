package fintechsearch

import "strings"

type PaymentEvent struct {
	ID       string            `json:"id"`
	Text     string            `json:"text"`
	Risk     string            `json:"risk"`
	Metadata map[string]string `json:"metadata"`
}

type Decision struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

func Decide(event PaymentEvent) Decision {
	risk := strings.ToLower(event.Risk)
	if risk == "high" {
		return Decision{Action: "review", Reason: "high-risk payment requires analyst review"}
	}
	return Decision{Action: "notify", Reason: "payment can receive an audit-friendly notification"}
}
