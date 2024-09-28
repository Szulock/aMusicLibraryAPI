package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	db     *sqlx.DB
	logger = logrus.New()
)

func Init() {
	initLogger()
	loadEnv()
	connectDatabase()
	migrate()
}

func initLogger() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Не удалось загрузить .env файл:", err)
	}
}

func connectDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Fatal("Не удалось подключиться к базе данных:", err)
	}
}

func migrate() {
	schema := `
		CREATE TABLE IF NOT EXISTS songs (
			id SERIAL PRIMARY KEY,
			group_name TEXT NOT NULL,
			song TEXT NOT NULL,
			release_date TEXT NOT NULL,
			text TEXT NOT NULL,
			link TEXT
		);`
	if _, err := db.Exec(schema); err != nil {
		logger.Fatal("Ошибка миграции базы данных:", err)
	}
	logger.Info("База данных успешно мигрирована.")
}

// GetDB возвращает инициализированную базу данных
func GetDB() *sqlx.DB {
	return db
}
