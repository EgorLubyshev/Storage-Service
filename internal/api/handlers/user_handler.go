package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-username/storage-service/internal/models"
)

type UserRepo interface {
	Create(ctx context.Context, name, passwordHash string) (string, error)
	GetByName(ctx context.Context, name string) (*models.User, error)
}

type UserHandler struct {
	users UserRepo
}

func NewUserHandler(users UserRepo) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) LoginForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"html": strings.TrimSpace(loginFormHTML)})
}

func (h *UserHandler) RegisterForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"html": strings.TrimSpace(registerFormHTML)})
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Confirm  string `json:"confirm"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || req.Password == "" || req.Confirm == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name and password are required"})
		return
	}
	if req.Password != req.Confirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}
	// if len(req.Password) < 6 {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
	// 	return
	// }

	if _, err := h.users.GetByName(ctx.Request.Context(), req.Name); err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "name already exists"})
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check user"})
		return
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	if _, err := h.users.Create(ctx.Request.Context(), req.Name, hash); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	user, err := h.users.GetByName(ctx.Request.Context(), req.Name)
	if err == nil {
		ctx.Header("x-user-token", user.ID)
	}
	ctx.JSON(http.StatusCreated, gin.H{"ok": true})
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name and password are required"})
		return
	}

	user, err := h.users.GetByName(ctx.Request.Context(), req.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	if !checkPassword(user.PasswordHash, req.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ctx.Header("x-user-token", user.ID)
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

const loginFormHTML = `
<div class="auth-card">
  <h2>Вход</h2>
  <label>Логин</label>
  <input id="loginName" type="text" autocomplete="username" />
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
  <input id="regName" type="text" autocomplete="username" />
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
