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
  id := "1"
  data := fmt.Sprintf(`{"id":%s,"key":"value"}`, id)

  m.ExpectQuery(".*").
    WithArgs(id).
    WillReturnRows(sqlmock.NewRows([]string{"data"}).AddRow(data))

  // Generate test request
  r, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s", resource, id), nil)
  response := mockServe(g, r)

  expectResponseCode(t, 200, response.Code)
  body := fmt.Sprintf(`{"data":%s}`, data)
  expectResponseBody(t, body, response.Body)
  if err := m.ExpectationsWereMet(); err != nil {
    t.Errorf("there were unfulfilled expections: %s", err)
  }
}

func TestReadNotFound(t *testing.T) {
  g, m := mockGrotto()
  defer g.DB.Close()

  resource := "resources"
  id := "2"

  m.ExpectQuery(".*").
    WithArgs(id).
    WillReturnRows(sqlmock.NewRows([]string{"data"}))

  // Generate test request
  r, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s", resource, id), nil)
  response := mockServe(g, r)

  expectResponseRowNotFound(t, response)
  if err := m.ExpectationsWereMet(); err != nil {
    t.Errorf("there were unfulfilled expections: %s", err)
  }
}

func testAllResources(t *testing.T) {
  g, m := mockGrotto()
  defer g.DB.Close()

  resource := "resources"
  data := fmt.Sprintf(`[{"id":1,"key":"value"},{"id":2,"key":"value2"}]`)

  m.ExpectQuery(".*").
    WillReturnRows(sqlmock.NewRows([]string{"data"}).AddRow(data))

  // Generate test request
  r, _ := http.NewRequest("GET", fmt.Sprintf("/%s", resource), nil)
  response := mockServe(g, r)

  expectResponseCode(t, 200, response.Code)
  body := fmt.Sprintf(`{"data":%s}`, data)
  expectResponseBody(t, body, response.Body)
  if err := m.ExpectationsWereMet(); err != nil {
    t.Errorf("there were unfulfilled expections: %s", err)
  }
}

func testAllNotFound(t *testing.T) {
  g, m := mockGrotto()
  defer g.DB.Close()

  resource := "resources"

  m.ExpectQuery(".*").
    WillReturnError(fmt.Errorf("sql: table not found"))

  // Generate test request
  r, _ := http.NewRequest("GET", fmt.Sprintf("/%s", resource), nil)
  response := mockServe(g, r)
  expectResponseTableNotFound(t, response)
  if err := m.ExpectationsWereMet(); err != nil {
    t.Errorf("there were unfulfilled expections: %s", err)
  }
}

func testAllEmpty(t *testing.T) {
  g, m := mockGrotto()
  defer g.DB.Close()

  resource := "resources"

  m.ExpectQuery(".*").
    WillReturnRows(sqlmock.NewRows([]string{"data"}))

  // Generate test request
  r, _ := http.NewRequest("GET", fmt.Sprintf("/%s", resource), nil)
  response := mockServe(g, r)

  expectResponseCode(t, 200, response.Code)
  body := `{"data":[]}`
  expectResponseBody(t, body, response.Body)
  if err := m.ExpectationsWereMet(); err != nil {
    t.Errorf("there were unfulfilled expections: %s", err)
  }
}

func expectResponseTableNotFound(t *testing.T, r *httptest.ResponseRecorder) {
  expectResponseCode(t, 404, r.Code)
  body := `{"error":"Table not found."}`
  expectResponseBody(t, body, r.Body)
}

func expectResponseRowNotFound(t *testing.T, r *httptest.ResponseRecorder) {
  expectResponseCode(t, 404, r.Code)
  body := `{"error":"Row not found."}`
  expectResponseBody(t, body, r.Body)
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
