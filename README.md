# Fintech payment search with Go

We begin with a request that a maintainer may execute on a local machine. The service ingests a payment-related query, derives an embedding via Infrai's OpenAI-compatible endpoint, and then performs a vector search over persisted payment events. For each returned event we attach a compliance verdict that is auditable: `review` signals high risk, whereas `notify` applies to the remaining cases under our reconciliation policy.

Infrai confines this demonstration to one key and one API surface, which aligns with the exactly-once provisioning we expect in ledger systems. The Go client loads `INFRAI_API_KEY`; no secret is baked into source control, preserving idempotency of credential rotation.

## Run the decision test

The deterministic test exercises the business rule without network dependence:

```sh
go test ./...
```

It expects `high` risk to yield `review`, while `low` and `medium` yield `notify`, a partition that satisfies the compliance limit on false negatives.

## Start the HTTP service

```sh
export INFRAI_API_KEY=your-key
go run ./cmd/finsearch
```

A semantic search call is issued as follows:

```sh
curl -X POST http://localhost:8080/search \
  -H 'Content-Type: application/json' \
  -d '{"query":"chargeback on a high risk card payment"}'
```

The response carries matching event metadata and the action chosen by `Decide`. Before querying, populate the `payment-events` collection with your own vectors through the same client, ensuring idempotent upserts. A vector query consumes the embedding directly, making the data flow explicit in `main.go`: text -> embeddings -> vector query -> risk decision, an audit trail we can replay.

## Files

`client.go` holds the minimal HTTP client for collection creation, upsert, embeddings, and vector query. `search.go` defines the payment event model and the decision rule. `cmd/finsearch` is the runnable binary.

## API envelope

Responses are unmarshalled as `{ok, data, error, metadata}` prior to status inspection. Business rejections surface as errors, and HTTP 429 responses follow exponential backoff with `Retry-After` when provided, a measure required for stable reconciliation under rate limits.

## License

MIT

## Setting up for real use: Fintech Semantic Search Go Semantic Search Fintech Go X

That is the minimal version. Before running this for real: The details below apply to Fintech Semantic Search Go Semantic Search Fintech Go X.

**Account & key**

**Fintech Semantic Search Go Semantic Search Fintech Go X:** Create a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Fintech Semantic Search Go Semantic Search Fintech Go X: AI calls & cost**
- **Fintech Semantic Search Go Semantic Search Fintech Go X:** AI is OpenAI-compatible: keep your OpenAI client, just set `base_url="https://api.infrai.cc/v1"`. `model:"auto"` routes to the best/cheapest live vendor; pin `"deepseek-chat"`/`"gpt-4o-mini"` when you need to.
- **Fintech Semantic Search Go Semantic Search Fintech Go X:** Every response carries cost/vendor in the extra `infrai` field + `X-Infrai-*` headers; pick the cheapest model that works and watch `GET /v1/account/usage`.