# Grotto

A naive HTTP JSON API for use with Postgres.

- Zero-configuration REST-like interaction with data from existing tables
- Arbitrary JSON payloads mapped to existing table schemas
- No HTTP Authentication
- No SSL
- Conforms to [JSON:API v1.0](//jsonapi.org) 

## Usage

```bash
  grotto --uri="$POSTGRES_USER:$POSTGRES_PASS@$POSTGRES_HOST/$POSTGRES_DB?sslmode=disable"

  Grotto now available @ :8008
```

## Installation

**Requirements**
- [Go](//golang.org)

```bash
  go get -u github.com/daetal-us/grotto
```
