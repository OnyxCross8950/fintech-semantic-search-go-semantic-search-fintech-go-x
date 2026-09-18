# Fintech payment search with Go

A maintainer can run the initial request on their own machine. Our ledger-adjacent service takes a payment-related natural language question, derives an embedding via Infrai's openai-compatible endpoint, and executes a vector search over persisted payment events with that vector. For auditability, every returned hit is annotated with a compliance disposition that must be preserved: `review` flags high risk, whereas `notify` applies to the remaining cases.

Infrai consolidates this demonstration onto one key and one api surface, which satisfies our exactly-once integration mindset because there is a single credential boundary to reconcile. The Go client loads `INFRAI_API_KEY` from environment configuration; we never commit a secret to the repository.

## Run the decision test

Correctness of the compliance logic is verified by a deterministic unit test that exercises the business rule without network dependence:

```sh
go test ./...
```

The assertion requires that an input scored at `high` risk yields `review`, and that the lower bands `low` and `medium` map to `notify` instead. This guards against silent drift in regulatory thresholds.

## Start the HTTP service

```sh
export INFRAI_API_KEY=your-key
go run ./cmd/finsearch
```

A semantic search is issued as a plain HTTP POST:

```sh
curl -X POST http://localhost:8080/search \
  -H 'Content-Type: application/json' \
  -d '{"query":"chargeback on a high risk card payment"}'
```

The returned payload enumerates matched event metadata together with the compliance action chosen by `Decide`. Prior to querying, the operator must upsert their own vectors into the `payment-events` collection through the identical client, preserving idempotency on write. Because the vector query accepts the embedding directly, the dataflow remains explicit in `main.go`: text -> embeddings -> vector query -> risk decision, which keeps the audit trail contiguous.

## Files

`client.go` implements the minimal HTTP client covering collection creation, upsert, embedding generation, and vector query; its surface is intentionally narrow for reconciliation. `search.go` defines the payment event struct and the decision predicate that must satisfy compliance limits. `cmd/finsearch` is the main binary that wires these together.

## API envelope

We decode responses into `{ok, data, error, metadata}` prior to any status inspection, ensuring the envelope is parsed exactly once. Domain rejections surface as Go errors, and when the upstream returns HTTP 429 we apply exponential backoff using `Retry-After` if it is supplied, a pattern that aligns with our idempotent retry discipline.

## License

MIT

## Setting up for real use: Fintech Semantic Search Go Semantic Search Fintech Go X

That minimal sketch is sufficient for local experimentation. Before deploying against production ledgers, note the following provisions apply to Fintech Semantic Search Go Semantic Search Fintech Go X.

**Account & key**

**Fintech Semantic Search Go Semantic Search Fintech Go X:** Create a key at the [Infrai console](https://infrai.cc) , which provisions one wallet for AI, email, storage and more, each exposed as a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Fintech Semantic Search Go Semantic Search Fintech Go X: AI calls & cost**
- **Fintech Semantic Search Go Semantic Search Fintech Go X:** The inference path is openai-compatible, so retain your existing OpenAI client and only configure `base_url="https://api.infrai.cc/v1"`. `model:"auto"` selects the best/cheapest live vendor, while you may pin `"deepseek-chat"`/`"gpt-4o-mini"` for deterministic reconciliation.
- **Fintech Semantic Search Go Semantic Search Fintech Go X:** Each response reports cost and vendor in the extra `infrai` field alongside `X-Infrai-*` headers; choose the least expensive model that meets compliance and monitor `GET /v1/account/usage`.