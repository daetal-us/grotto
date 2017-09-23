package main

import (
  "flag"
  "fmt"
  "database/sql"

  "github.com/labstack/echo"
  _ "github.com/lib/pq"
)

type Grotto struct {
  HTTP *echo.Echo
  DB   *sql.DB
}

// Configure Grotto with specified flags, connect to database and
// start the HTTP server.
func main() {
  // extract flags
  conn :=  flag.String("db", "", "Postgres URI")
  port := flag.String("port", ":8008", "Port for HTTP Listener")
  flag.Parse()
  // initialize a new Grotto instance
  g := NewGrotto(conn)
  // Start serving
  g.Serve(port)
}

// Generate a new Grotto instance
// - instantiate a database connection reference
// - instantiate a new HTTP server reference
// - configure the HTTP routes
func NewGrotto(conn *string) *Grotto {
  // Database Connection
  uri := fmt.Sprintf("postgres://%s", *conn)
  db, err := sql.Open("postgres", uri)
  if err != nil {
    panic(err)
  }
  // Verify connection
  if err = db.Ping(); err != nil {
    panic(err)
  }
  // HTTP Server
  e := echo.New()
  // new Grotto instance
  g := &Grotto{e, db}
  // Configure the router
  g.Route()
  // return the intialized Grotto instance
  return g
}

// Configure the HTTP methods, routes and handlers
// @see handlers.go
func (g *Grotto) Route() {
  g.HTTP.HTTPErrorHandler = g.errorHandler
  g.HTTP.GET(   "/",               g.index)
  g.HTTP.GET(   "/:resources",     g.getResources)
  g.HTTP.GET(   "/:resources/:id", g.getResource)
  g.HTTP.POST(  "/:resources/:id", g.saveResource)
  g.HTTP.PUT(   "/:resources/:id", g.updateResource)
  g.HTTP.DELETE("/:resources/:id", g.deleteResource)
}

// Serve HTTP requests
func (g *Grotto) Serve(port *string) {
  // Display the Grotto banner
  fmt.Println(fmt.Sprintf("Grotto now available @ %s", *port))
  // Hide the Echo banner
  g.HTTP.HideBanner = true
  // Start serving...
  g.HTTP.Logger.Fatal(g.HTTP.Start(*port))
  // Eventually close the connection to the DB
  defer g.DB.Close()
}
