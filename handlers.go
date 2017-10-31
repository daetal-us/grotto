package main

import (
  "net/http"
  "regexp"
  "io/ioutil"

  "github.com/daetal-us/grotto/resources"
  "github.com/labstack/echo"
)

type Response struct {
  Data interface{} `json:"data,omitempty"`
  Error string     `json:"error,omitempty"`
}

// Error handler
func (g *Grotto) errorHandler(err error, c echo.Context) {
  c.Logger().Error(err)
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
  JSONError(c, code, err.Error())
}

// A generic welcome message for the index
func (g *Grotto) index(c echo.Context) error {
  data, err := resources.List(g.DB)
  if err != nil {
    return err
  }
  return c.JSON(http.StatusOK, &Response{
    Data: data,
  })
}

// Get tables HTTP handler
func (g *Grotto) getRows(c echo.Context) error {
  table, err := tableFromContext(c)
  if err != nil {
    return err
  }
  data, err := resources.All(g.DB, table)
  if err != nil {
    if err == resources.ErrNotFound {
      return JSONError(c, http.StatusNotFound, "Table not found.")
    }
    return err
  }
  return JSONData(c, data)
}

// Get row HTTP handler
func (g *Grotto) getRow(c echo.Context) error {
  table, err := tableFromContext(c)
  if err != nil {
    return err
  }
  id, err := idFromContext(c)
  if err != nil {
    return err
  }
  data, err := resources.Get(g.DB, table, id)
  if err != nil {
    if err == resources.ErrNotFound {
      return JSONError(c, http.StatusNotFound, "Row not found.")
    }
    return err
  }
  return JSONData(c, data)
}

// Save row HTTP handler
func (g *Grotto) createRow(c echo.Context) error {
  table, err := tableFromContext(c)
  if err != nil {
    return err
  }

  r := c.Request()
  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    return err
  }
  err = resources.Create(g.DB, table, string(data))
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err)
  }
  return c.NoContent(http.StatusOK)
}
// Update row HTTP handler
func (g *Grotto) updateRow(c echo.Context) error {
  table, err := tableFromContext(c)
  if err != nil {
    return err
  }
  id, err := idFromContext(c)
  if err != nil {
    return err
  }
  body, err := ioutil.ReadAll(c.Request().Body)
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err)
  }
  err = resources.Update(g.DB, table, id, string(body))
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err)
  }
  return c.NoContent(http.StatusOK)
}
// Delete row HTTP handler
func (g *Grotto) deleteRow(c echo.Context) error {
  table, err := tableFromContext(c)
  if err != nil {
    return err
  }
  id, err := idFromContext(c)
  if err != nil {
    return err
  }
  count, err := resources.Delete(g.DB, table, id)
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err)
  }
  if count < 1 {
    return resources.ErrNotFound
  }
  return c.NoContent(http.StatusOK)
}

// Helper for JSON error response
func JSONError(c echo.Context, code int, message string) error {
  return c.JSON(code, Response{Error: message})
}

// Helper for JSON data response
func JSONData(c echo.Context, data interface{}) error {
  return c.JSON(http.StatusOK, Response{Data: data})
}

// Extract tables parameter from HTTP request context
func tableFromContext(c echo.Context) (table string, err error) {
  table = c.Param("table")
  if table == "" {
    err = echo.NewHTTPError(http.StatusBadRequest, "No table specified.")
  }
  re := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
  invalid := !re.MatchString(table)
  if invalid {
    err = echo.NewHTTPError(http.StatusBadRequest, "Invalid table specified.")
  }
  return
}

// Extract id parameter from HTTP request context
func idFromContext(c echo.Context) (id string, err error) {
  id = c.Param("id")
  if id == "" {
    err = echo.NewHTTPError(http.StatusBadRequest, "No id specified.")
  }
  return
}
