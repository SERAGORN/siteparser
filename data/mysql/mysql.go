package mysql

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

// InitMysql ...
func InitMysql(ctx context.Context, host, port, user, password, dbName string) (*sqlx.DB, error) {
	cfg := mysql.NewConfig()

	cfg.Net = "tcp"
	cfg.Addr = host
	cfg.User = user
	cfg.Passwd = password
	cfg.DBName = dbName
	cfg.ParseTime = true
	cfg.Timeout = time.Second * 2

	dsn := cfg.FormatDSN()

	fmt.Println(dsn, cfg.Addr)
	return establishConnection(ctx, "mysql", dsn)
}


func establishConnection(ctx context.Context, driverName string, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, driverName, dsn) // sqlx.Connect performs ping under the hood

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(200)
	return db, nil
}