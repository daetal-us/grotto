package main

import (
  "fmt"
  "errors"
  "net/http"
  "regexp"

  "github.com/labstack/echo"
)

type Response struct {
  Data interface{}     `json:"data,omitempty"`
  Errors []interface{} `json:"errors,omitempty"`
  Meta *Meta           `json:"meta,omitempty"`
}

type Meta struct {
  Resource string `json:"resource"`
}

func (g *Grotto) errorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
  c.JSON(code, &Response{
    Errors: []interface{}{map[string]interface{}{
      "status": fmt.Sprintf(`%d`, code),
      "title": err.Error(),
    }},
  })
  c.Logger().Error(err)
}


// A generic welcome message for the index
func (g *Grotto) index(c echo.Context) error {
  data, err := g.getResourcesListFromDb()
  if err != nil {
    return err
  }
  return c.JSON(http.StatusOK, &Response{
    Data: data,
  })
}

// Get resources HTTP handler
func (g *Grotto) getResources(c echo.Context) error {
  resources, err := getResourcesFromContext(c)
  if err != nil {
    return err
  }
  data, err := g.getResourcesFromDB(resources)
  if err != nil {
    return err
  }
  return c.JSON(http.StatusOK, &Response{
    Data: data,
    Meta: &Meta{resources},
  })
}

// Get resource HTTP handler
func (g *Grotto) getResource(c echo.Context) error {
  resources, err := getResourcesFromContext(c)
  if err != nil {
    return err
  }
  id, err := getIdFromContext(c)
  if err != nil {
    return err
  }
  data, err := g.getResourceFromDB(resources, id)
  if err != nil {
    return err
  }
  return c.JSON(http.StatusOK, &Response{
    Data: data,
    Meta: &Meta{resources},
  })
}

// Save resource HTTP handler
func (g *Grotto) saveResource(c echo.Context) error {
  return nil
}
// Update resources HTTP handler
func (g *Grotto) updateResource(c echo.Context) error {
  return nil
}
// Delete resources HTTP handler
func (g *Grotto) deleteResource(c echo.Context) error {
  return nil
}

// Extract resources parameter from HTTP request context
func getResourcesFromContext(c echo.Context) (resources string, err error) {
  resources = c.Param("resources")
  if resources == "" {
    err = errors.New("No resource specified.")
  }
  re := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
  invalid := !re.MatchString(resources)
  if invalid {
    resources = ""
    err = errors.New("Invalid resource specified.")
  }
  return
}

// Extract id parameter from HTTP request context
func getIdFromContext(c echo.Context) (id string, err error) {
  id = c.Param("id")
  if id == "" {
    err = errors.New("No id specified.")
  }
  return
}
