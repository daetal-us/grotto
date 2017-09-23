package main

import (
  "fmt"
  "encoding/json"
)

func (g *Grotto) getResourcesListFromDb() (data []interface{}, err error) {
  query := `select relname, n_live_tup from pg_stat_user_tables;`
  rows, err := g.DB.Query(query)
  if err == nil {
    defer rows.Close()
    for rows.Next() {
      var resource string
      var count int
      err = rows.Scan(&resource, &count)
      if err != nil {
        return
      }
      data = append(data, map[string]interface{}{"resource": resource, "count": count})
    }
  }
  return
}

func (g *Grotto) getResourcesFromDB(resources string) (data []interface{}, err error) {
  var result string
  template := `
  select jsonb_agg(a) as data from (
    select b.id, '%[1]s' as type, to_jsonb(b.*) as attributes from (
        select * from %[1]s
    ) as b
  ) as a;`
  query := fmt.Sprintf(template, resources)
  err = g.DB.QueryRow(query).Scan(&result)
  if err == nil {
    err = json.Unmarshal([]byte(result), &data)
  }

  return
}

func (g *Grotto) getResourceFromDB(resources string, id string) (data interface{}, err error) {
  var result string
  template := `
  select to_jsonb(a) as data from (
    select b.id, '%[1]s' as type, to_jsonb(b.*) as attributes from (
        select * from %[1]s where id = $1
    ) as b
  ) as a;`
  query := fmt.Sprintf(template, resources)
  err = g.DB.QueryRow(query, id).Scan(&result)
  if err == nil {
    err = json.Unmarshal([]byte(result), &data)
  }
  return
}
