package resources

import (
  "fmt"
  "encoding/json"
  "database/sql"
  "strings"
  "errors"
)

var ErrNotFound = errors.New("Resource not found.")

func List(db *sql.DB) (tables []interface{}, err error) {
  query := `select relname, n_live_tup from pg_stat_user_tables;`
  rows, err := db.Query(query)
  if err == nil {
    defer rows.Close()
    for rows.Next() {
      var resource string
      var count int
      err = rows.Scan(&resource, &count)
      if err != nil {
        return
      }
      tables = append(tables, map[string]interface{}{"resource": resource, "count": count})
    }
  }
  return
}

func All(db *sql.DB, table string) (resources []interface{}, err error) {
  template := `
  select jsonb_agg(a) as data from (
    select * from %[1]s
  ) as a;`
  query := fmt.Sprintf(template, table)
  var r sql.NullString
  err = db.QueryRow(query).Scan(&r)
  if err == nil {
    if r.Valid {
      err = json.Unmarshal([]byte(r.String), &resources)
    }
  } else {
    if err == sql.ErrNoRows {
      err = nil
    }
  }
  return
}

func Get(db *sql.DB, table string, id string) (resource map[string]interface{}, err error) {
  template := `
  select to_jsonb(a) as data from (
    select * from %[1]s where id = $1
  ) as a;`
  query := fmt.Sprintf(template, table)
  var r sql.NullString
  err = db.QueryRow(query, id).Scan(&r)
  if err == nil {
    if r.Valid {
      err = json.Unmarshal([]byte(r.String), &resource)
    }
  } else {
    if err == sql.ErrNoRows {
      err = ErrNotFound
    }
  }
  return
}

func Create(db *sql.DB, table string, resource string) (err error) {
  template := `insert into %[1]s (select * from json_populate_record(null::%[1]s, $1))`
  query := fmt.Sprintf(template, table)
  transaction, err := db.Begin()
  if err != nil {
    return
  }
  statement, err := transaction.Prepare(query)
  if err != nil {
    transaction.Rollback()
    return
  }
  defer statement.Close()
  _, err = statement.Exec(resource)
  if err != nil {
    transaction.Rollback()
    return
  }
  transaction.Commit()
  return
}

func Update(db *sql.DB, table string, id string, resource string) (err error) {
  var d map[string]interface{}
  err = json.Unmarshal([]byte(resource), &d)
  if err != nil {
    return
  }
  sets, values, err := setsAndValues(d)
  if err != nil {
    return
  }
  values = append([]interface{}{id}, values...)
  transaction, err := db.Begin()
  if err != nil {
    return
  }
  template := "update %s set %s where id = $1"
  query := fmt.Sprintf(template, table, sets)
  statement, err := transaction.Prepare(query)
  if err != nil {
    transaction.Rollback()
    return
  }
  defer statement.Close()
  _, err = statement.Exec(values...)
  if err != nil {
    transaction.Rollback()
    return
  }
  transaction.Commit()
  return
}

func Delete(db *sql.DB, table string, id string) (count int64, err error) {
  query := fmt.Sprintf("delete from %v where id = $1", table)
  result, err := db.Exec(query, id)
  if err == nil {
    count, err = result.RowsAffected()
  }
  return
}

func setsAndValues(data map[string]interface{}) (sets string, values []interface{}, err error) {
  fragments := []string{}
  values = []interface{}{}
  for key, value := range data {
    switch value.(type) {
    default:
      values = append(values, value)
    case []interface{}, map[string]interface{}:
      v, err := json.Marshal(value)
      if err != nil {
        return ``, []interface{}{}, err
      }
      values = append(values, string(v))
    }
    count := len(values) + 1
    fragment := fmt.Sprintf("%s = $%d", key, count)
    fragments = append(fragments, fragment)
  }
  sets = strings.Join(fragments, ", ")
  return
}
