package store

import (
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Storage interface {
	OrdersGateway
	QueryRow(string, ...interface{}) *sql.Row
	Query(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

type DB struct {
	*goqu.Database
	cache map[string]Order
}

func NewStore(conn, migration string) (Storage, error) {
	var err error

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("open db error: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connect db error: %v", err)
	}

	mg, err := migrate.New(migration, conn)

	if err != nil {
		return nil, fmt.Errorf("open migration error: %v", err)
	}

	err = mg.Up()

	switch err {
	case nil:
		break
	case migrate.ErrNoChange:
		err = nil
	default:
		return nil, fmt.Errorf("apply migration error: %v", err)
	}

	_, _ = mg.Close()
	d := goqu.New("postgres", db)

	rdb := DB{d, map[string]Order{}}
	orderslist, err := rdb.GetAllOrders()
	for _, order := range orderslist {
		rdb.cache[order.OrderUid] = order
	}
	return &rdb, err
}
