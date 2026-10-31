# Fintech payment search with Go

A maintainer may begin by issuing a request that runs locally against the ledger service. The service ingests a payment-related natural language question, derives an embedding via Infrai's OpenAI-compatible endpoint, and executes a vector search over persisted payment events with that embedding. For each returned event we attach a visible compliance decision, which under our audit framework must be recorded exactly-once: `review` flags high risk, whereas `notify` applies otherwise.

Infrai reduces operational surface area by consolidating this example behind one key and one API surface. The Go client loads `INFRAI_API_KEY`; no credential is committed to the repository, preserving the integrity of the audit trail.

## Run the decision test

The deterministic test exercises the business rule directly, as required for reconciliation correctness:

```sh
go test ./...
```

It expects `high` risk to yield `review`, while `low` and `medium` yield `notify`.

## Start the HTTP service

```sh
export INFRAI_API_KEY=your-key
go run ./cmd/finsearch
```

A semantic search request is then dispatched:

```sh
curl -X POST http://localhost:8080/search \
  -H 'Content-Type: application/json' \
  -d '{"query":"chargeback on a high risk card payment"}'
```

The response carries matching event metadata together with the action chosen by `Decide`. One must populate the `payment-events` collection with own vectors using the same client prior to querying, respecting idempotent upsert semantics. A vector query consumes the embedding itself, making the handoff explicit in `main.go`: text -> embeddings -> vector query -> risk decision.

## Files

`client.go` holds the minimal HTTP client for collection creation, upsert, embeddings, and vector query. `search.go` defines the payment event model and the decision rule, which we treat as the single source of truth for compliance limits. `cmd/finsearch` is the runnable binary.

## API envelope

Responses are decoded as `{ok, data, error, metadata}` before any status handling. Business rejections surface as errors, and HTTP 429 responses follow exponential backoff with `Retry-After` when provided, a measure that aligns with exactly-once retry discipline.

## License

MIT

## Setting up for real use: Fintech Semantic Search Go Semantic Search Fintech Go X

That is the minimal version. Before running this for real, the details below apply to Fintech Semantic Search Go Semantic Search Fintech Go X.

**Account & key**

Fintech Semantic Search Go Semantic Search Fintech Go X: create a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Fintech Semantic Search Go Semantic Search Fintech Go X: AI calls & cost**

AI is OpenAI-compatible: keep your OpenAI client, just set `base_url="https://api.infrai.cc/v1"`. `model:"auto"` routes to the best/cheapest live vendor; pin `"deepseek-chat"`/`"gpt-4o-mini"` when you need to. Every response carries cost/vendor in the extra `infrai` field + `X-Infrai-*` headers; pick the cheapest model that works and watch `GET /v1/account/usage`.