package errors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("пользователь уже существует")
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrInvalidPassword   = errors.New("неверный пароль")
	ErrInvalidToken      = errors.New("неверный токен")
	ErrDatabaseError     = errors.New("ошибка базы данных")
	ErrRedisError        = errors.New("ошибка Redis")
)
