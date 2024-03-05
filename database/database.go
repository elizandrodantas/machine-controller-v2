package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const (
	defaultMaxConns          = int32(5)         // Limite máximo de conexões no pool
	defaultMinConns          = int32(0)         // Número mínimo de conexões no pool
	defaultMaxConnLifetime   = time.Minute * 30 // Tempo máximo de vida de uma conexão no pool
	defaultMaxConnIdleTime   = time.Minute * 5  // Tempo máximo que uma conexão pode ficar inativa no pool
	defaultHealthCheckPeriod = time.Minute * 10 // Período de verificação de saúde do pool
	defaultConnectTimeout    = time.Second * 10 // Tempo máximo para conectar ao banco de dados
)

// NewPostgresConnect creates a new connection pool to a PostgreSQL database.
// It takes an environment configuration object as input and returns a *pgxpool.Pool and an error.
func NewPostgresConnect(e *config.Env) (*pgxpool.Pool, error) {
	var sql string
	if e.DB_URL != "" {
		sql = e.DB_URL
	} else {
		sql = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", e.DB_USER, e.DB_PASS, e.DB_HOST, e.DB_PORT, e.DB_NAME)
	}

	config, err := pgxpool.ParseConfig(sql)
	if err != nil {
		log.Fatalf("error while parsing configuration: %s\n", err.Error())
	}

	config.MaxConns = defaultMaxConns
	config.MinConns = defaultMinConns
	config.MaxConnLifetime = defaultMaxConnLifetime
	config.MaxConnIdleTime = defaultMaxConnIdleTime
	config.HealthCheckPeriod = defaultHealthCheckPeriod
	config.ConnConfig.ConnectTimeout = defaultConnectTimeout

	client, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("error when creating pool connection to database: %s\n", err.Error())
	}

	for i := 0; i < 5; i++ {
		err = client.Ping(context.Background())
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
	}

	return client, err
}
