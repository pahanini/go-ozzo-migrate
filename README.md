# go-ozzo-migrate

Db migration tool based on ozzo-dbx

Not ready for production


## How to use

``` go
	e := NewExecutor(db)

	e.NewMigration("item_test").
		UpSql("CREATE TABLE item_ (id INT, name VARCHAR(10))").
		DownSql("DROP TABLE item");

	e.NewMigration("category_test").
		UpSql("CREATE TABLE category (id INT, name VARCHAR(10))").
		DownSql("DROP TABLE category");

    // Apply migrations
	err := e.Up()

```
