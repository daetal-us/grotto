package main

import (
  "errors"
  "flag"
  "fmt"
  "database/sql"
  "net/http"

  "github.com/labstack/echo"
  _ "github.com/lib/pq"
)

type (
  Grotto struct {
    HTTP *echo.Echo
    DB   *sql.DB
  }
  Payload struct {
    Data []interface{} `json:"data"`
    Meta *Meta          `json:"meta"`
  }
  Meta struct {
    Resource string `json:"resource"`
  }
)

// Configure Grotto with specified flags, connect to database and
// start the HTTP server.
func main() {
  conn :=  flag.String("db", "", "Postgres URI")
  port := flag.String("port", ":8008", "Port for HTTP Listener")
  flag.Parse()
  g := NewGrotto(conn)
  fmt.Println(fmt.Sprintf("Grotto now available @ %s", *port))
  g.HTTP.Logger.Fatal(g.HTTP.Start(*port))
  defer g.DB.Close()
}

// Generate a new Grotto instance, instantiate a database reference with the
// specified database connection string, instantiate a new HTTP server and
// setup the routes
func NewGrotto(conn *string) *Grotto {
  uri := fmt.Sprintf("postgres://%s", *conn)
  db, err := sql.Open("postgres", uri)
  if err != nil {
    panic(err)
  }
  e := echo.New()
  e.HideBanner = true
  g := &Grotto{e, db}
  g.Route()
  return g
}

// Configure the HTTP methods, routes and handlers
func (g *Grotto) Route() {
  g.HTTP.GET("/", g.welcome)
  g.HTTP.GET("/:resources", g.getResources)
  g.HTTP.GET("/:resources/:id", g.getResource)
  g.HTTP.POST("/:resources/:id", g.saveResource)
  g.HTTP.PUT("/:resources/:id", g.updateResource)
  g.HTTP.DELETE("/:resources/:id", g.deleteResource)
}

// A generic welcome message for the index
func (g *Grotto) welcome(c echo.Context) error {
  message := []byte(`{"meta":{"message":"Aloha!","works":true}}`)
  return c.JSONBlob(http.StatusOK, message)
}

// Get resources HTTP handler
func (g *Grotto) getResources(c echo.Context) error {
  resources, err := getResourcesFromContext(c)
  if err != nil {
    return err
  }
  return g.getResourcesFromDB(c, resources)
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
  return g.getResourceFromDB(c, resources, id)
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

func (g *Grotto) getResourcesFromDB(c echo.Context, resources string) error {
  return c.JSON(http.StatusOK, &Payload{
    make([]interface{}, 0),
    &Meta{resources},
  })
}

func (g *Grotto) getResourceFromDB(c echo.Context, resources string, id string) error {
  record := map[string]interface{}{}
  record["id"] = id
  data := []interface{}{record}
  return c.JSON(http.StatusOK, &Payload{
    data,
    &Meta{resources},
  })
}

// Extract resources parameter from HTTP request context
func getResourcesFromContext(c echo.Context) (resources string, err error) {
  resources = c.Param("resources")
  if resources == "" {
    err = errors.New("No resource specified.")
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
