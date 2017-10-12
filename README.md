# Grotto

_A naive HTTP JSON API for Postgres_

Zero-configuration, REST-like interaction with existing data in a postgres database.

## Installation
```bash
go get github.com/daetal-us/grotto
```
## Requirements

- Requires [Go](//golang.org).
- Requires an accessible instance of [Postgres](//postgresql.org).
- Only table names containing alphanumeric, dashes and underscore characters will be accessible.
- Primary key columns for all tables assumed to be `id`.

## Usage

```bash
grotto -addr :8080 -dsn postgres://$POSTGRES_USER:$POSTGRES_PASS@$POSTGRES_HOST/$POSTGRES_DB?sslmode=disable
```
## Conventions

### Paths

| path | description | `GET` (read) | `POST` (create) | `PUT` (update) | `DELETE` |
| --- | --- | :-: | :-: |  :-: | :-: |
| `/` | all available tables | √ | | | |
| `/:table` | all resources in table | √ | | | |
| `/:table/:id` | specific resource in table | √ | √ | √ | √ |

### Responses

Successful responses:

```
GET /users

{
  "data": [
    {
      ...
    },
    ...
  ],
  "meta": {
    "table": "users"
  }
}
```

```
GET /users/1234

{
  "data": {
    "id": 1234,
    ...
  },
  "meta": {
    "table": "users"
  }
}
```

Unsuccessful responses:

```
GET /nonexistant-table

{
  "error": {
    "status": 500,
    "message": "Some error message."
  }
}
```
