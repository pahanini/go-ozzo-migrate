package migrate

import (
	"errors"
	"fmt"
	"github.com/go-ozzo/ozzo-dbx"
)

type Executor struct {
	db         *dbx.DB
	TableName  string
	LogFunc    dbx.LogFunc
	migrations []MigrationDescriptor
}

func NewExecutor(db *dbx.DB) *Executor {
	m := make([]MigrationDescriptor, 0)
	return &Executor{
		db,
		"migrate",
		nil,
		m,
	}
}

func (m *Executor) appliedMap() (r map[string]int, err error) {
	var (
		rows *dbx.Rows
		name string
		id   int
	)
	r = make(map[string]int)
	rows, err = m.db.Select("id", "name").From(m.TableName).Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			return
		}
		r[name] = id
	}
	return
}

func (m *Executor) Add(d MigrationDescriptor) *Executor {
	m.migrations = append(m.migrations, d)
	return m
}

// Creates new Description and adds it to migrations
func (m *Executor) NewMigration(name string) *Migration {
	d := &Migration{m.db, name, nil, nil}
	m.Add(d)
	return d
}

func (m *Executor) Up() error {

	m.db.CreateTable(m.TableName, map[string]string{
		"id":   "int primary key",
		"name": "varchar(100)",
	}).Execute()

	applied, err := m.appliedMap()
	if err != nil {
		return err
	}

	for i, d := range m.migrations {
		if _, ok := applied[d.Name()]; ok {
			continue
		}
		q := d.Up()
		if q == nil {
			return errors.New(fmt.Sprintf("Invalid up query %v", d.Name()))
		}
		_, err := d.Up().Execute()
		if err != nil {
			return err
		}
		_, err = m.db.Insert(m.TableName, dbx.Params{
			"id":   i,
			"name": d.Name(),
		}).Execute()
		if err != nil {
			return err
		}
		if m.LogFunc != nil {
			m.LogFunc("Successfully migrate up %v", d.Name())
		}
	}
	return nil
}
