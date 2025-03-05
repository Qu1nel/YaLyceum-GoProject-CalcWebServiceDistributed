package handlers

import (
	"fmt"
	"net/http"

	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/customError"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (r *Router) GetExpressions(c *gin.Context) {
	var req GetExpressionsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		r.Log.Error("GetExpressions: Failed bind query", zap.Error(err))
		customError.New(http.StatusBadRequest, fmt.Errorf("GetExpressions: error: %w", err)).SendError(c)
		c.Abort()
		return
	}
	if req.Size == 0 {
		req.Size = 10
	}
	exprs, total, cErr := r.service.GetExpressions(req.Size, req.Page)
	if cErr != nil {
		r.Log.Error("GetExpressions: Failed get expressions", zap.Error(cErr))
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	resp := make([]*GetExpressionResp, len(exprs))
	for i, expr := range exprs {
		resp[i] = mapExpression(expr)
	}
	c.JSON(http.StatusOK, resp)
}
