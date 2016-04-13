package migrate

import "github.com/go-ozzo/ozzo-dbx"

// MigrationDescriptor interface
type MigrationDescriptor interface {
	Name() string
	Up() *dbx.Query
	Down() *dbx.Query
}

type Migration struct {
	db       *dbx.DB
	name     string
	up, down *dbx.Query
}

func (d *Migration) Down() *dbx.Query {
	return d.down
}

func (d *Migration) DownQuery(query *dbx.Query) *Migration {
	d.down = query
	return d
}

func (d *Migration) DownSql(sql string) *Migration {
	return d.DownQuery(d.db.NewQuery(sql))
}

func (d *Migration) Name() string {
	return d.name
}

func (d *Migration) Up() *dbx.Query {
	return d.up
}

func (d *Migration) UpQuery(query *dbx.Query) *Migration {
	d.up = query
	return d
}

func (d *Migration) UpSql(sql string) *Migration {
	return d.UpQuery(d.db.NewQuery(sql))
}
