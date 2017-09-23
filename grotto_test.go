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
  g = &Grotto{
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

  resource := "resources"
  id := 1
  data := fmt.Sprintf(`[{"id":%d,"key":"value"}]`, id)

  m.ExpectQuery(".*").
    WithArgs(string(id)).
    WillReturnRows(sqlmock.NewRows([]string{"data"}).AddRow(data))

  // Generate test request
  r, _ := http.NewRequest("GET", fmt.Sprintf("/resources/%d", id), nil)
  response := mockServe(g, r)

  expectResponseCode(t, 200, response.Code)
  body := fmt.Sprintf(`{"data":%s,"meta":{"resource":"%s"}}`, data, resource)
  expectResponseBody(t, body, response.Body)
  if err := m.ExpectationsWereMet(); err != nil {
    t.Errorf("there were unfulfilled expections: %s", err)
  }
}

func expectResponseBody(t *testing.T, expected string, data *bytes.Buffer) {
  bs, _ := ioutil.ReadAll(data)
  received := string(bs)
  if expected != received {
    t.Errorf("\nUnexpected response body:\n%s\n\nExpected:\n%s\n", received, expected)
  }
}

func expectResponseCode(t *testing.T, expected int, received int) {
  if expected != received {
    t.Errorf("\nUnexpected response code: %d\nExpected %d\n", received, expected)
  }
}
