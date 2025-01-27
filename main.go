package main

import (
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение к PostgreSQL
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer db.Close()

	// Проверка подключения к БД
	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки подключения к PostgreSQL: %v", err)
	}

	// Подключение к Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
	})

	// Проверка подключения к Redis
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	defer rdb.Close()

	// Проверка наличия JWT-секрета
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET не установлен")
	}

	// Инициализация зависимостей
	userRepo := repository.NewUserRepository(db, rdb)
	authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	authHandler := handler.NewAuthHandler(authService)

	// Настройка роутера
	r := gin.Default()

	// Роуты
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// Добавляем восстановление после паники
	r.Use(gin.Recovery())

	// Добавляем логирование
	r.Use(gin.Logger())

	log.Println("Сервер запущен на порту 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
