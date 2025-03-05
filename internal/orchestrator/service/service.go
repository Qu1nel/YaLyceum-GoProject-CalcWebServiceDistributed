package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/models"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/orchestrator/repository"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/calculator"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/customError"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/token"
	"go.uber.org/zap"
)

type Service struct {
	Log        *zap.Logger
	Repository repository.Repo
	Calculator *calculator.Calculator
}

func New(
	repo repository.Repo,
	log *zap.Logger,
	calculator *calculator.Calculator,
) *Service {
	return &Service{
		Log:        log,
		Repository: repo,
		Calculator: calculator,
	}
}

func (s *Service) GetExpressions(size, page int) ([]*models.Expression, int64, *customError.CustomError) {
	exprs, total, err := s.Repository.GetExpressions(size, page)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, 0, customError.New(http.StatusNotFound, fmt.Errorf("Service.GetExpressions: error: %w", err))
		}
		return nil, 0, customError.New(http.StatusInternalServerError, fmt.Errorf("Service.GetExpressions: unknown errorerror: %w", err))
	}
	return exprs, total, nil
}
func (s *Service) GetExpression(id int64) (*models.Expression, *customError.CustomError) {
	expr, err := s.Repository.GetExpression(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, customError.New(http.StatusNotFound, fmt.Errorf("Service.GetExpression: error: %w", err))
		}
		return nil, customError.New(http.StatusInternalServerError, fmt.Errorf("Service.GetExpression: unknown errorerror: %w", err))
	}
	return expr, nil
}
func (s *Service) CreateExpression(expression string) (int64, *customError.CustomError) {
	expressionNT, err := token.TokenizeExpression(expression)
	if err != nil {
		return 0, customError.New(http.StatusUnprocessableEntity, fmt.Errorf("Service.CreateExpression: error: %w", err))
	}

	id, err := s.Repository.CreateExpression(expression)
	if err != nil {
		return 0, customError.New(http.StatusInternalServerError, fmt.Errorf("Service.CreateExpression: error: %w", err))
	}
	go s.Calculator.Calc(*expressionNT, id)
	return id, nil
}
func (s *Service) SetResult(id, expressionID int64, result float64, error *string) *customError.CustomError {
	task := &models.Task{
		ID:           id,
		ExpressionID: expressionID,
		Result:       result,
		Error:        error,
	}
	if !s.Calculator.Exists(id) {
		return customError.New(http.StatusNotFound, fmt.Errorf("Service.SerResult: task not found"))
	}
	s.Calculator.ReceiveResult(task)
	return nil
}
func (s *Service) GetTask() (*models.Task, *customError.CustomError) {
	task, cErr := s.Calculator.GetTask()
	if cErr != nil {
		return nil, cErr
	}
	return task, nil
}
