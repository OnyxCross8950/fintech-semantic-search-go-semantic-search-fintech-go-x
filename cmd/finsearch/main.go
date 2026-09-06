package main

import (
	"encoding/json"
	"fintech-semantic-search"
	"log"
	"net/http"
)

func main() {
	client, err := fintechsearch.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		var in struct {
			Query string `json:"query"`
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.Query == "" {
			http.Error(w, "query required", http.StatusBadRequest)
			return
		}
		emb, err := client.Embed(in.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		items, err := client.Query("payment-events", emb, 5)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		type result struct {
			Event    fintechsearch.PaymentEvent `json:"event"`
			Decision fintechsearch.Decision     `json:"decision"`
		}
		resp := make([]result, 0, len(items))
		for _, item := range items {
			resp = append(resp, result{item, fintechsearch.Decide(item)})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	log.Println("fintech search listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
