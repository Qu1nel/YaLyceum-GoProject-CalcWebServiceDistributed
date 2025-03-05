package agent

import (
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/agent"
	config2 "github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/config"
	logger2 "github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/pkg/logger"
	"go.uber.org/fx"
)

var Agent = fx.Option(
	fx.Provide(
		logger2.New,
		config2.New,
		agent.New,
	),
)
