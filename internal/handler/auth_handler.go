package handler

import (
	"auth-service/internal/errors"
	"auth-service/internal/model"
	"auth-service/internal/service"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка валидации запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат данных",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.Register(c.Request.Context(), &req); err != nil {
		switch err {
		case errors.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email уже существует"})
		default:
			log.Printf("Ошибка регистрации: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка валидации запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат данных",
			"details": err.Error(),
		})
		return
	}

	token, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			log.Printf("Пользователь не найден: %s", req.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		case errors.ErrInvalidPassword:
			log.Printf("Неверный пароль для пользователя: %s", req.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		default:
			log.Printf("Ошибка входа для пользователя %s: %v", req.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
