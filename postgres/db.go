package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	db     *sqlx.DB
	ctx    context.Context
	cancel func()
	DSN    string
	Now    func() time.Time
}

func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
		Now: time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open() (err error) {
	if db.DSN == "" {
		return fmt.Errorf("dsn required.")
	}

	db.db, err = sqlx.Open("postgres", db.DSN)
	if err != nil {
		return err
	}

	if err = db.db.Ping(); err != nil {
		return err
	}

	log.Println("🚀 Connected to database successfully")

	return nil
}

func (db *DB) Close() error {
	db.cancel()
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) *Tx {
	tx := db.db.MustBeginTx(ctx, opts)
	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}
}

type Tx struct {
	*sqlx.Tx
	db  *DB
	now time.Time
}

type NullTime time.Time

func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		*(*time.Time)(n) = time.Time{}
		return nil
	} else if value, ok := value.(string); ok {
		*(*time.Time)(n), _ = time.Parse(time.RFC3339, value)
		return nil
	}
	return fmt.Errorf("NullTime: cannot scan time to time.Time: %T", value)
}

func (n *NullTime) Value() (driver.Value, error) {
	if n == nil || (*time.Time)(n).IsZero() {
		return nil, nil
	}
	return (*time.Time)(n).UTC().Format(time.RFC3339), nil
}
