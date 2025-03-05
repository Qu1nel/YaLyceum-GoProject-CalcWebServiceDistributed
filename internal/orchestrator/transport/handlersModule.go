package transport

import (
	"YaLyceum/internal/orchestrator/transport/handlers"
	"YaLyceum/internal/orchestrator/transport/routers"

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
