package dal

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
	struct: DataSource
	description: stores several data sources
*/
type DataSource struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

/*
	func: InitDS
	description: Initialize database: Connection & Migration
*/
func InitDS() (*DataSource, error) {
	log.Info("Initializing data sources\n")

	// initialize postgres
	log.Info("Connecting to PostgreSQL\n")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DB")
	pgSSL := os.Getenv("PG_SSL")
	pgTimeZone := os.Getenv("PG_TIMEZONE")
	posDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		pgHost, pgUser, pgPassword, pgDB, pgPort, pgSSL, pgTimeZone)
	db, err := gorm.Open(postgres.Open(posDsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error while connecting to postgreSQL: %w", err)
	}

	// migrate to postgres
	log.Info("Migrating to PostgreSQL\n")
	err = dBMigrator(db)
	if err != nil {
		return nil, fmt.Errorf("error while migrating to postgreSQL: %w", err)
	}

	// initializing redis
	log.Info("Connecting to Redis\n")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	// test connection
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("error connectinf to redis: %w", err)
	}

	return &DataSource{
		DB:          db,
		RedisClient: rdb,
	}, nil
}

/*
	func: CloseDS
	description: close connections to databases
*/
func CloseDS(ds *DataSource) error {
	// close connection to postgres
	// nothing

	// close connection to redis
	if err := ds.RedisClient.Close(); err != nil {
		return fmt.Errorf("error closing Redis Client: %w", err)
	}
	return nil
}
