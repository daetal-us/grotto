package main

import (
  "flag"
  "database/sql"
  "log"

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
  dsn :=  flag.String("dsn", "", "Postgres DSN (e.g. postgres://user:pass@host/db)")
  addr := flag.String("addr", ":8008", "TCP network address for HTTP Listener")
  flag.Parse()

  if *dsn == "" {
    log.Fatal("No DSN provided.")
  }
  // initialize a new Grotto instance
  g := NewGrotto(dsn)
  // Start serving
  g.Serve(addr)
}

// Generate a new Grotto instance
// - instantiate a database connection reference
// - instantiate a new HTTP server reference
// - configure the HTTP routes
func NewGrotto(dsn *string) *Grotto {
  // Database Connection
  db, err := sql.Open("postgres", *dsn)
  if err != nil {
    log.Fatal(err)
  }
  // Verify connection
  if err = db.Ping(); err != nil {
    log.Fatal("Database connection refused.")
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
  g.HTTP.HTTPErrorHandler =       g.errorHandler
  g.HTTP.GET(   "/",              g.index)
  g.HTTP.GET(   "/:resource",     g.getResources)
  g.HTTP.GET(   "/:resource/:id", g.getResource)
  g.HTTP.POST(  "/:resource/:id", g.createResource)
  g.HTTP.PUT(   "/:resource/:id", g.updateResource)
  g.HTTP.DELETE("/:resource/:id", g.deleteResource)
}

// Serve HTTP requests
func (g *Grotto) Serve(addr *string) {
  // Display the Grotto banner
  log.Printf("Starting Grotto @ %s", *addr)
  log.Print("kūkulu pono, kamaʻāina!")
  // Hide the Echo banner
  g.HTTP.HideBanner = true
  // Start serving...
  log.Fatal(g.HTTP.Start(*addr))
  // Eventually close the connection to the DB
  defer g.DB.Close()
}
