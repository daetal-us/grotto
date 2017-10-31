<img src="http://daetal.us/static/media/grotto.png" align="right">

# Grotto

_A naive HTTP JSON API for Postgres_

Minimally-configurable, REST-like interaction with a postgres database.

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
grotto -addr :8080 -dsn postgres://user:password@host/db
```

## Interface

The interface consists of standard HTTP [request methods](//en.wikipedia.org/wiki/Hypertext_Transfer_Protocol#Request_methods) paired with the following path conventions:

| path | description | `GET` (read) | `POST` (create) | `PUT` (update) | `DELETE` |
| --- | --- | :-: | :-: |  :-: | :-: |
| `/` | all available tables | √ | | | |
| `/:table` | all rows in table | √ | | | |
| `/:table/:id` | specific row in table | √ | √ | √ | √ |

## Responses

All rows in hypothetical `users` table:

```
GET /users HTTP/1.1
...
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
...
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

Single row in hypothetical `users` table:

```
GET /users/1234 HTTP/1.1
...
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
...
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

Table not found:

```
GET /nonexistant-table HTTP/1.1
...
HTTP/1.1 404 Not Found
Content-Type: application/json; charset=UTF-8
...
{
  "error": "Table not found."
}
```
