package main

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "testing"
  "net/http"
  "net/http/httptest"

  "github.com/labstack/echo"
  "github.com/DATA-DOG/go-sqlmock"
)

func mockGrotto() (g *Grotto, m sqlmock.Sqlmock) {
  db, m, _ := sqlmock.New()
  g := &Grotto{
    HTTP: echo.New(),
    DB: db,
  }
  g.Route()
  return
}

func mockServe(g *Grotto, request *http.Request) *httptest.ResponseRecorder {
  recorder := httptest.NewRecorder()
  g.HTTP.ServeHTTP(recorder, request)
  return recorder
}

func TestReadResource(t *testing.T) {
  g, m := mockGrotto()
  defer g.DB.Close()

  json := `[{"id":1,"key":"value"}]`
  m.ExpectQuery("SELECT json_agg\\(resources\\) as data FROM resources where id = (.+)").
    WithArgs("1").
    WillReturnRows(m.NewRows([]string{"data"}).AddRow(json))

  r, _ := http.NewRequest("GET", "/resources/1", nil)
  response := mockServe(g, r)

  expectResponseCode(t, 200, response.Code)
  expectResponseBody(t, json, response.Body)
  if err := m.ExpectationsWereMet(); err != nil {
    t.Errorf("there were unfulfilled expections: %s", err)
  }
}

func expectResponseBody(t *testing.T, expected string, data *bytes.Buffer) {
  bs, _ := ioutil.ReadAll(data)
  result := string(bs)
  if expected != result {
    t.Errorf("Expected response body:`%s` Received: `%s`\n", expected, result)
  }
}

func expectResponseCode(t *testing.T, expected int, result int) {
  if expected != result {
    t.Errorf("Expected response code %d. Received %d\n", expected, result)
  }
}
