package fintechsearch

import "testing"

func TestDecideRiskAction(t *testing.T) {
	cases := []struct{ risk, action string }{{"high", "review"}, {"low", "notify"}, {"medium", "notify"}}
	for _, tc := range cases {
		got := Decide(PaymentEvent{Risk: tc.risk})
		if got.Action != tc.action {
			t.Fatalf("risk %q: got %q", tc.risk, got.Action)
		}
	}
}
