package migrate

import (
	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"os"
	"log"
	"testing"
)

func TestAll(t *testing.T) {
	db := getDB()
	e := NewExecutor(db)
	e.LogFunc = log.Printf
	e.TableName = "migrate_test"
	e.NewMigration("item_test").
		UpSql("CREATE TABLE item_test (id INT, name VARCHAR(10))").
		DownSql("DROP TABLE item_test")
	e.NewMigration("category_test").
		UpSql("CREATE TABLE category_test (id INT, name VARCHAR(10))").
		DownSql("DROP TABLE category_test")

	err := e.Up()
	assert.Nil(t, err)

	var c int
	db.NewQuery("SELECT COUNT(*) FROM migrate_test").Row(&c)
	assert.Equal(t, 2, c)

}

func getDB() *dbx.DB {
	db, err := dbx.Open("postgres", os.Getenv("MIGRATION_TEST_DSN"))
	if err != nil {
		panic(err)
	}
	_, err = db.NewQuery(`
		DROP TABLE IF EXISTS migrate_test;
		DROP TABLE IF EXISTS item_test;
		DROP TABLE IF EXISTS category_test;
	`).Execute()
	if err != nil {
		panic(err)
	}
	return db
}
