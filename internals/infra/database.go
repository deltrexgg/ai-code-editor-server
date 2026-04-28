package infra
import (
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
/*
	Postgress databse connection setup
	with pool configuration

	use pgBouncer in production to manage connection pool
	set sqlDB.SetConnMaxLifetime(0) and sqlDB.SetMaxOpenConns(0)

*/

var (
	DataBaseClient *gorm.DB
)

func InitDB(dsn string) {
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	// Load pool configs from env (with defaults)
	maxOpen := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	maxIdle := getEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	maxLife := getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)
	maxIdleTime := getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 1*time.Minute)

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLife)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	log.Println("PostgreSQL connected with connection pooling enabled")

	DataBaseClient = db
}

func getEnvAsInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}