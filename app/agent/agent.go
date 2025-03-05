package agent

import (
	"YaLyceum/internal/agent"
	config2 "YaLyceum/internal/pkg/config"
	logger2 "YaLyceum/internal/pkg/logger"

	"go.uber.org/fx"
)

var Agent = fx.Option(
	fx.Provide(
		logger2.New,
		config2.New,
		agent.New,
	),
)
