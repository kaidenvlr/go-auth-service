package repository

import (
	"auth-service/internal/errors"
	"auth-service/internal/model"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lib/pq"
)

type UserRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewUserRepository(db *sql.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{db: db, redis: redis}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// Код 23505 - нарушение уникального ограничения
			if pqErr.Code == "23505" {
				return errors.ErrUserAlreadyExists
			}
		}
		log.Printf("Ошибка создания пользователя: %v", err)
		return errors.ErrDatabaseError
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// Проверяем Redis
	userJSON, err := r.redis.Get(ctx, "user:"+email).Result()
	if err == nil {
		log.Printf("Получены данные из Redis для %s: %s", email, userJSON)
		var user model.User
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			// Если ошибка десериализации, логируем и удаляем некорректные данные из кэша
			log.Printf("Ошибка десериализации из Redis: %v", err)
			r.redis.Del(ctx, "user:"+email)
		} else if user.Password != "" {
			// Возвращаем пользователя только если данные корректны
			return &user, nil
		}
	}

	// Ищем в PostgreSQL
	user := &model.User{}
	query := `SELECT id, email, password FROM users WHERE email = $1`
	err = r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		log.Printf("Ошибка получения пользователя из БД: %v", err)
		return nil, errors.ErrDatabaseError
	}

	// Сохраняем в Redis
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		log.Printf("Ошибка сериализации пользователя для Redis: %v", err)
	} else {
		if err := r.redis.Set(ctx, "user:"+email, string(jsonBytes), 30*time.Minute).Err(); err != nil {
			log.Printf("Ошибка сохранения в Redis: %v", err)
		}
	}

	return user, nil
}

func (r *UserRepository) ClearCache(ctx context.Context, email string) error {
	return r.redis.Del(ctx, "user:"+email).Err()
}
