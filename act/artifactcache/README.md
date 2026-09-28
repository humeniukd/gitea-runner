# Artifact cache OpenAPI contract

The cache server serves its hand-maintained OpenAPI 3.1.0 document at
`/openapi.json` with the `application/json` content type. It is embedded directly
from `openapi.json`; no code generation or bundled documentation UI is required.

The document describes the legacy dependency-cache API under
`/_apis/artifactcache`. It does not cover the v2 cache service, artifact proxy,
or internal control-plane API. Keep it in sync when changing the documented
routes, request and response formats, or authentication requirements.

Cache management requires a registered per-job `ACTIONS_RUNTIME_TOKEN` bearer
token. Downloads use the signed `archiveLocation` returned by cache lookup.
The specification itself is public and contains no job credentials.

To view the contract interactively, import the JSON into an OpenAPI viewer.
The server does not provide HTML, JavaScript, CSS, or favicon assets.
