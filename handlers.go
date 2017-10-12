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
  Error *Error     `json:"error,omitempty"`
  Meta *Meta       `json:"meta,omitempty"`
}

type Error struct {
  Status int `json:"status"`
  Detail string `json:"detail"`
}

type Meta struct {
  Table string `json:"table,omitempty"`
}

func (g *Grotto) errorHandler(err error, c echo.Context) {
  c.Logger().Error(err)
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
  errorJSON(c, code, err.Error())
}

func errorJSON(c echo.Context, code int, detail string) error {
  r := Response{
    Error: &Error{
      code,
      detail,
    },
  }
  return c.JSON(code, r)
}

func resourceDataJSON(c echo.Context, resource string, data interface{}) error {
  return c.JSON(http.StatusOK, Response{
    Data: data,
    Meta: &Meta{resource},
  })
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

// Get resources HTTP handler
func (g *Grotto) getResources(c echo.Context) error {
  r, err := resourceFromContext(c)
  if err != nil {
    return err
  }
  data, err := resources.All(g.DB, r)
  if err != nil {
    return err
  }
  return resourceDataJSON(c, r, data)
}

// Get resource HTTP handler
func (g *Grotto) getResource(c echo.Context) error {
  r, err := resourceFromContext(c)
  if err != nil {
    return err
  }
  id, err := idFromContext(c)
  if err != nil {
    return err
  }
  data, err := resources.Get(g.DB, r, id)
  if err != nil {
    if err == resources.ErrNotFound {
      return errorJSON(c, http.StatusNotFound, err.Error())
    }
    return err
  }
  return resourceDataJSON(c, r, data)
}

// Save resource HTTP handler
func (g *Grotto) createResource(c echo.Context) error {
  resource, err := resourceFromContext(c)
  if err != nil {
    return err
  }

  r := c.Request()
  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    return err
  }
  err = resources.Create(g.DB, resource, string(data))
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err)
  }
  return c.NoContent(http.StatusCreated)
}
// Update resources HTTP handler
func (g *Grotto) updateResource(c echo.Context) error {
  resource, err := resourceFromContext(c)
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
  err = resources.Update(g.DB, resource, id, string(body))
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err)
  }
  return g.getResource(c)
}
// Delete resources HTTP handler
func (g *Grotto) deleteResource(c echo.Context) error {
  resource, err := resourceFromContext(c)
  if err != nil {
    return err
  }
  id, err := idFromContext(c)
  if err != nil {
    return err
  }
  count, err := resources.Delete(g.DB, resource, id)
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err)
  }
  if count < 1 {
    return errorJSON(c, http.StatusNotFound, "Resource not found.")
  }
  return c.NoContent(http.StatusOK)
}

// Extract resources parameter from HTTP request context
func resourceFromContext(c echo.Context) (resource string, err error) {
  resource = c.Param("resource")
  if resource == "" {
    err = echo.NewHTTPError(http.StatusBadRequest, "No resource specified.")
  }
  re := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
  invalid := !re.MatchString(resource)
  if invalid {
    err = echo.NewHTTPError(http.StatusBadRequest, "Invalid resource specified.")
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
