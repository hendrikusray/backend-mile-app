package http

import (
	"context"
	"errors"
	"net/http"

	"mile-app-test/domain"
	"mile-app-test/utils"

	"github.com/gin-gonic/gin"
)

type taskHandler struct {
	useCase domain.TaskUsecase
}

func RegisterTaskRoutes(rg *gin.RouterGroup, uc domain.TaskUsecase, auth gin.HandlerFunc) {
	h := &taskHandler{useCase: uc}

	g := rg.Group("/tasks")
	if auth != nil {
		g.Use(auth)
	}

	g.GET("", h.List)
	g.POST("", h.Create)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *taskHandler) List(c *gin.Context) {
	var q domain.TaskListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "invalid query params", nil)
		return
	}

	ownerID := c.GetString("user_id")
	res, err := h.useCase.List(ctx(c), ownerID, q)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "internal error", nil)
		return
	}
	utils.JSONResponse(c, http.StatusOK, "Success", res)
}

func (h *taskHandler) Create(c *gin.Context) {
	var in domain.CreateTaskReq
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "invalid JSON body", nil)
		return
	}

	ownerID := c.GetString("user_id")
	out, err := h.useCase.Create(ctx(c), ownerID, in)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrValidation):
			utils.JSONResponse(c, http.StatusBadRequest, err.Error(), nil)
		default:
			utils.JSONResponse(c, http.StatusInternalServerError, "internal error", nil)
		}
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "Success", out)
}

func (h *taskHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	ownerID := c.GetString("user_id")

	out, err := h.useCase.GetByID(ctx(c), id, ownerID)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrValidation) || err.Error() == "invalid id":
			utils.JSONResponse(c, http.StatusBadRequest, "invalid id", nil)
		case errors.Is(err, utils.ErrForbidden) || err.Error() == "forbidden":
			utils.JSONResponse(c, http.StatusForbidden, "forbidden", nil)
		default:
			utils.JSONResponse(c, http.StatusInternalServerError, "internal error", nil)
		}
		return
	}
	if out == nil {
		utils.JSONResponse(c, http.StatusNotFound, "task not found", nil)
		return
	}
	utils.JSONResponse(c, http.StatusOK, "Success", out)
}

func (h *taskHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var in domain.UpdateTaskReq
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "invalid JSON body", nil)
		return
	}

	ownerID := c.GetString("user_id")
	out, err := h.useCase.Update(ctx(c), id, ownerID, in)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrValidation) || err.Error() == "invalid id":
			utils.JSONResponse(c, http.StatusBadRequest, err.Error(), nil) // "invalid id" / "no fields to update"
		case errors.Is(err, utils.ErrForbidden) || err.Error() == "forbidden":
			utils.JSONResponse(c, http.StatusForbidden, "forbidden", nil)
		default:
			utils.JSONResponse(c, http.StatusInternalServerError, "internal error", nil)
		}
		return
	}
	if out == nil {
		utils.JSONResponse(c, http.StatusNotFound, "task not found", nil)
		return
	}
	utils.JSONResponse(c, http.StatusOK, "Success", out)
}

func (h *taskHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	ownerID := c.GetString("user_id")

	if err := h.useCase.Delete(ctx(c), id, ownerID); err != nil {
		switch {
		case errors.Is(err, utils.ErrValidation) || err.Error() == "invalid id":
			utils.JSONResponse(c, http.StatusBadRequest, "invalid id", nil)
		case errors.Is(err, utils.ErrForbidden) || err.Error() == "forbidden":
			utils.JSONResponse(c, http.StatusForbidden, "forbidden", nil)
		case errors.Is(err, utils.ErrNotFound) || err.Error() == "mongo: no documents in result":
			utils.JSONResponse(c, http.StatusNotFound, "task not found", nil)
		default:
			utils.JSONResponse(c, http.StatusInternalServerError, "internal error", nil)
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func ctx(c *gin.Context) context.Context {
	return c.Request.Context()
}
