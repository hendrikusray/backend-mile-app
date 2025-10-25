package http

import (
	"errors"
	"mile-app-test/domain"
	"mile-app-test/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	useCase domain.UserUseCase
}

func NewUserHandler(rg *gin.RouterGroup, useCase domain.UserUseCase) {
	handler := &userHandler{
		useCase: useCase,
	}
	rg.POST("/login", handler.Login)
}

func (h *userHandler) Login(c *gin.Context) {
	var in domain.Login
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "invalid JSON body", nil)
		return
	}

	tok, err := h.useCase.Login(c, &in)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrValidation):
			utils.JSONResponse(c, http.StatusBadRequest, err.Error(), nil)
		case errors.Is(err, utils.ErrInvalidCredentials):
			utils.JSONResponse(c, http.StatusUnauthorized, "invalid credentials", nil)
		default:
			utils.JSONResponse(c, http.StatusInternalServerError, "internal error", nil)
		}
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Success", tok)
}
