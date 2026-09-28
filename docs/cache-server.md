# Cache server URLs

On runners, `cache.external_server` identifies the shared cache server to forward requests to. On `gitea-runner cache-server`, the same key sets the fallback URL advertised in cache download responses:

```yaml
cache:
  dir: /data/actcache
  port: 8088
  external_secret_file: /run/secrets/cache-secret
  external_server: "https://cache.example.com/cache/"
```

`gitea-runner -c cache-server.yaml cache-server --external-url https://cache.example.com/cache/` overrides the configured URL. The URL must be absolute, use HTTP or HTTPS, and contain a hostname without a query or fragment. Path prefixes are allowed; trailing slashes are removed. A reverse proxy must route that prefix to the cache server.

A public URL supplied by a runner when registering a job takes precedence over this fallback, so downloads return through that job's runner proxy. If neither is set, the server advertises its detected outbound address and listening port. Setting the advertised URL does not change the listening address or enable TLS; terminate HTTPS at your reverse proxy.
