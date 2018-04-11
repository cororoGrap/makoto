makoto
========
Simple migration tool for PostgreSQL

Install
-------
Install makoto CLI
```bash
go get -u github.com/cororoGrap/makoto/cmd/makoto
```
Install makoto migrator
```bash
go get github.com/cororoGrap/makoto
```

Structure
----------
makoto will create a directory named migration under your project. All sql migration sql will placed under this directory.
Migration scripts should be named as
```bash
v[numeric version number]_[script name].sql
e.g.
v1_basic.sql
```

CLI
------
Init migration directory
```bash
makoto init
```

Create new migration sql script
```bash
makoto new [script_name]
```

Generate golang migration collection, a golang file 'collection.go' will be created under the migration directory
```bash
makoto collect
```

Integrate with Golang
-----
First generate the collection file with CLI.

Initialize the makoto migrator
```go
migrator := makoto.New(db, "postgres") // pass the DB pointer and DB driver name
```
Pass the migration collection to migrator
```go
migrator.SetCollection(migration.GetCollection())
```
Perform migration
```go
migrator.Up() // migrate to latest version
// or
migrator.EnsureSchema("v10") // migrate to a given version
```

#### Example
```go
func startMigration(db *sql.DB) {
    migrator := makoto.New(db, "postgres")
    migrator.SetCollection(migration.GetCollection())
    migrator.EnsureSchema("V2") 
}
```
