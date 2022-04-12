package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/avast/retry-go"
)

type CH struct {
	ok   bool
	conn driver.Conn
	ctx  context.Context
}

type Row struct {
	Dt  time.Time `ch:"dt"`
	Msg string    `ch:"msg"`
}

func (ch *CH) Write(msg string) error {
	batch, err := ch.conn.PrepareBatch(ch.ctx, "INSERT INTO user_creation (dt, msg)")
	if err != nil {
		return err
	}
	for i := 0; i < 1; i++ {
		if err := batch.Append(time.Now(), msg); err != nil {
			return err
		}
	}
	if err := batch.Send(); err != nil {
		return err
	}
	return nil
}

func (ch *CH) Read() (*[]Row, error) {
	var result []Row

	if err := ch.conn.Select(
		ch.ctx,
		&result,
		"SELECT dt, msg FROM user_creation ORDER BY dt DESC LIMIT 10",
	); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ch *CH) Connect(addr, db, user, pass string) {
	var err error
	defer func() {
		ch.ok = err == nil
	}()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: db,
			Username: user,
			Password: pass,
		},
		//Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		return
	}
	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 10,
	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
		log.Println("progress: ", p)
	}))

	err = retry.Do(
		func() error {
			log.Println("Trying to connect to the Clickhouse...")
			if err := conn.Ping(ctx); err != nil {
				if exception, ok := err.(*clickhouse.Exception); ok {
					return errors.New(fmt.Sprintf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace))
				}
				return err
			}
			return nil
		},
		retry.Attempts(5),
		retry.Delay(time.Second),
		retry.DelayType(retry.BackOffDelay),
	)

	ch.conn = conn
	ch.ctx = ctx
}

func (ch *CH) Close() {
	ch.conn.Close()
}
