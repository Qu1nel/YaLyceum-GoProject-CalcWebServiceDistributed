package routers

import (
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/models"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/customError"
	"github.com/gin-gonic/gin"
)

type Service interface {
	GetExpressions(size, page int) ([]*models.Expression, int64, *customError.CustomError)
	GetExpression(id int64) (*models.Expression, *customError.CustomError)
	CreateExpression(expression string) (int64, *customError.CustomError)
	SetResult(id, expressionID int64, result float64, error *string) *customError.CustomError
	GetTask() (*models.Task, *customError.CustomError)
}
type Routers struct {
	Public *gin.RouterGroup
}

func CreateRouter(g *gin.Engine) *Routers {
	public := g.Group("/api/v1")
	return &Routers{
		Public: public,
	}
}
