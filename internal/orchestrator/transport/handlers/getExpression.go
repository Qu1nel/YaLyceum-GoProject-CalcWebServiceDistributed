package handlers

import (
	"fmt"
	"net/http"

	"YaLyceum/internal/models"
	"YaLyceum/internal/pkg/customError"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (r *Router) GetExpression(c *gin.Context) {
	var req GetExpressionReq
	if err := c.ShouldBindUri(&req); err != nil {
		r.Log.Error("GetExpression: Failed bind uri", zap.Error(err))
		customError.New(http.StatusBadRequest, fmt.Errorf("GetExpression: error: %w", err)).SendError(c)
		c.Abort()
		return
	}
	expr, cErr := r.service.GetExpression(*req.ID)
	if cErr != nil {
		r.Log.Error("GetExpression: Failed get expression", zap.Error(cErr))
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, mapExpression(expr))
}
func mapExpression(expression *models.Expression) *GetExpressionResp {
	return &GetExpressionResp{
		ID:         expression.ID,
		Expression: expression.Expression,
		Status:     expression.Status,
		Result:     expression.Result,
	}
}
