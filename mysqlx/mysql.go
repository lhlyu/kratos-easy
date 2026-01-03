package mysqlx

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
)

// NewClient 新建 mysql 客户端
//
// source 格式示例：username:password@tcp(domain.com:3306)/dbname?parseTime=True&loc=Local
func NewClient(logger log.Logger, source string, opts ...Option) (*sql.DB, func(), error) {
	l := log.NewHelper(logger)
	db, err := sql.Open("mysql", source)
	if err != nil {
		l.Errorw(
			"open database failed",
			"error", err,
		)
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		l.Errorw(
			"ping database failed",
			"error", err,
		)
		return nil, nil, err
	}

	o := &options{
		maxOpen:      20,
		maxIdleCount: 10,
		maxLifetime:  30 * time.Minute,
		maxIdleTime:  10 * time.Minute,
	}

	for _, opt := range opts {
		opt(o)
	}

	db.SetConnMaxLifetime(o.maxLifetime)
	db.SetConnMaxIdleTime(o.maxIdleTime)
	db.SetMaxOpenConns(o.maxOpen)
	db.SetMaxIdleConns(o.maxIdleCount)

	l.Infow("database connected")

	cleanup := func() {
		if err := db.Close(); err != nil {
			l.Errorw("close database failed", "error", err)
		}
	}

	return db, cleanup, nil
}
