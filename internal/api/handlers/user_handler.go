package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-username/storage-service/internal/models"
)

type UserRepo interface {
	Create(ctx context.Context, login, passwordHash string) (int64, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
}

type UserHandler struct {
	users UserRepo
}

func NewUserHandler(users UserRepo) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) LoginForm(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"html": strings.TrimSpace(loginFormHTML)})
}

func (h *UserHandler) RegisterForm(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"html": strings.TrimSpace(registerFormHTML)})
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Confirm  string `json:"confirm"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	req.Login = strings.TrimSpace(req.Login)
	if req.Login == "" || req.Password == "" || req.Confirm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login and password are required"})
		return
	}
	if req.Password != req.Confirm {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}
	// if len(req.Password) < 6 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
	// 	return
	// }

	if _, err := h.users.GetByLogin(c.Request.Context(), req.Login); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "login already exists"})
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check user"})
		return
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	if _, err := h.users.Create(c.Request.Context(), req.Login, hash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	user, err := h.users.GetByLogin(c.Request.Context(), req.Login)
	if err == nil {
		c.Header("x-user-token", strconv.FormatInt(user.ID, 10))
	}
	c.JSON(http.StatusCreated, gin.H{"ok": true})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	req.Login = strings.TrimSpace(req.Login)
	if req.Login == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login and password are required"})
		return
	}

	user, err := h.users.GetByLogin(c.Request.Context(), req.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	if !checkPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.Header("x-user-token", strconv.FormatInt(user.ID, 10))
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

const loginFormHTML = `
<div class="auth-card">
  <h2>Вход</h2>
  <label>Логин</label>
  <input id="loginUsername" type="text" autocomplete="username" />
  <label>Пароль</label>
  <input id="loginPassword" type="password" autocomplete="current-password" />
  <div class="actions">
    <button onclick="LoginClick()">Войти</button>
    <button onclick="GetRegisterForm()">Регистрация</button>
  </div>
</div>
`

const registerFormHTML = `
<div class="auth-card">
  <h2>Регистрация</h2>
  <label>Логин</label>
  <input id="regLogin" type="text" autocomplete="username" />
  <label>Пароль</label>
  <input id="regPassword" type="password" autocomplete="new-password" />
  <label>Подтверждение</label>
  <input id="regConfirmPassword" type="password" autocomplete="new-password" />
  <div class="actions">
    <button onclick="RegClick()">Создать аккаунт</button>
    <button onclick="GetLoginForm()">Назад</button>
  </div>
</div>
`
