package handlers

import (
	"fmt"
	"net/http"

	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/customError"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (r *Router) CreateExpression(c *gin.Context) {
	var req CreateExpressionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		r.Log.Error("CreateExpression: Failed bind body", zap.Error(err))
		customError.New(http.StatusBadRequest, fmt.Errorf("CreateExpression err: %w", err)).Error()
		c.Abort()
		return
	}

	id, cErr := r.service.CreateExpression(*req.Expression)
	if cErr != nil {
		r.Log.Error("CreateExpression: Failed create expression", zap.Error(cErr.Err))
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated,
		CreateExpressionResp{ID: id},
	)
}
