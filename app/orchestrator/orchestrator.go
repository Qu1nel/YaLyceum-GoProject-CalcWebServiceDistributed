package orchestrator

import (
	"YaLyceum/internal/orchestrator/repository"
	"YaLyceum/internal/orchestrator/repository/db"
	"YaLyceum/internal/orchestrator/service"
	"YaLyceum/internal/orchestrator/transport"
	"YaLyceum/internal/orchestrator/transport/routers"
	"YaLyceum/internal/pkg/cache"
	"YaLyceum/internal/pkg/calculator"
	"YaLyceum/internal/pkg/config"
	"YaLyceum/internal/pkg/counter"
	"YaLyceum/internal/pkg/http"
	"YaLyceum/internal/pkg/logger"
	"YaLyceum/internal/pkg/postgres"

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