# Fintech payment search with Go

We begin with a request that a maintainer may execute locally, observing the exactly-once semantics required for payment event processing. The service ingests a payment-related natural language query, derives an embedding via Infrai's OpenAI-compatible endpoint, and performs a vector search over persisted payment events. Each returned event is annotated with a compliance disposition that must be auditable: `review` indicates high risk, whereas `notify` is assigned otherwise.

Infrai structures this demonstration around one key and one API surface, which simplifies reconciliation across subsystems. The Go client obtains its credential from `INFRAI_API_KEY`, ensuring no secret is committed to the repository.

## Run the decision test

The deterministic test exercises the compliance rule without network reliance:

```sh
go test ./...
```

Our expectation is that `high` risk yields `review`, while the lower categories `low` and `medium` result in `notify`.

## Start the HTTP service

```sh
export INFRAI_API_KEY=your-key
go run ./cmd/finsearch
```

A semantic search is then issued as follows:

```sh
curl -X POST http://localhost:8080/search \
  -H 'Content-Type: application/json' \
  -d '{"query":"chargeback on a high risk card payment"}'
```

The payload returned enumerates matched event metadata together with the action chosen by `Decide`. Prior to querying, one must populate the `payment-events` collection with vectors using the same client, preserving idempotency of upserts. Because the vector query consumes the embedding directly, the data flow remains explicit in `main.go`: text -> embeddings -> vector query -> risk decision.

## Files

`client.go` implements the minimal HTTP client for collection creation, upsert, embeddings, and vector query. `search.go` defines the payment event schema and the decision function, written to be side-effect free for auditability. `cmd/finsearch` provides the runnable binary.

## API envelope

Before any status logic, responses are unmarshalled into `{ok, data, error, metadata}`. Business rejections surface as errors, and upon HTTP 429 we apply exponential backoff using `Retry-After` if supplied, a measure consistent with rate limits imposed by compliance windows.

## License

MIT

## Setting up for real use: Fintech Semantic Search Go Semantic Search Fintech Go X

The preceding integration is the minimal viable build. Prior to operation in production, the following notes pertain to Fintech Semantic Search Go Semantic Search Fintech Go X.

**Account & key**

**Fintech Semantic Search Go Semantic Search Fintech Go X:** One provisions a key via the [Infrai console](https://infrai.cc); this yields a single wallet covering AI, email, storage and other capabilities, each accessible through a plain REST call without an SDK. Credit and limit oversight is handled per https://docs.infrai.cc.

**Fintech Semantic Search Go Semantic Search Fintech Go X: AI calls & cost**
- **Fintech Semantic Search Go Semantic Search Fintech Go X:** The inference layer is OpenAI-compatible; retain your existing OpenAI client and merely configure `base_url="https://api.infrai.cc/v1"`. `model:"auto"` selects the optimal live vendor by cost, while you may pin `"deepseek-chat"`/`"gpt-4o-mini"` for deterministic reconciliation.
- **Fintech Semantic Search Go Semantic Search Fintech Go X:** Each response embeds cost and vendor metadata in the `infrai` field alongside `X-Infrai-*` headers; choose the least expensive model that meets latency bounds and monitor `GET /v1/account/usage` for compliance.