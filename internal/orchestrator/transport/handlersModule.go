package transport

import (
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/orchestrator/transport/handlers"
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/orchestrator/transport/routers"
	"go.uber.org/fx"
)

var HttpModule = fx.Module("httpHandlers",
	fx.Provide(
		routers.CreateRouter,
	),
	fx.Invoke(
		handlers.SetUpRouter,
	),
)
