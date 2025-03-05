package orchestrator

import (
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/orchestrator/repository"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/orchestrator/repository/db"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/orchestrator/service"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/orchestrator/transport"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/orchestrator/transport/routers"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/cache"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/calculator"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/config"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/counter"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/http"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/logger"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/postgres"
	"go.uber.org/fx"
)

var Orchestrator = fx.Options(
	fx.Provide(
		config.New,
		logger.New,
		http.New,
		counter.New,
		cache.New,
		postgres.New,
		fx.Annotate(db.New,
			fx.As(new(repository.Repo)),
		),
		calculator.New,
		fx.Annotate(service.New,
			fx.As(new(routers.Service)),
		),
	),
	fx.Invoke(
		postgres.MigrateDB,
	),
	transport.HttpModule,
)