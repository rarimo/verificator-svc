## Global headers (only for OpenAPI 2)

When using OpenAPI 2 you can minimize headers duplications by using `headers` global object (similar to `definitions` or `responses`).
During build process all references to global `headers` will be inlined and `headers` will be removed from the resulting spec so spec will be valid (global `headers` are not allowed by OpenAPI 2 spec):

Example:
```yaml
...
headers:
  Rate-Limit-Limit:
    description: The number of allowed requests in the current period
    type: integer
...
paths:
  /api-keys:
    get:
      summary: Retrieve a list of api keys
      responses:
        200:
          description: A list of api keys was retrieved successfully
          headers:
            Rate-Limit-Limit:
              $ref: "#/headers/Rate-Limit-Limit"
```

## Deprecation Notice

**The light verification flow is deprecated and will be removed in a future version.**
All endpoints under the path `/integrations/verificator-svc/light/` are affected by this deprecation.
Please use the full verification flow instead.
