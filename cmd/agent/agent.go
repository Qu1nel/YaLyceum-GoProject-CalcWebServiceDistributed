package main

import (
	"YaLyceum/app/agent"
	agent2 "YaLyceum/internal/agent"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		agent.Agent,
		fx.Invoke(func(*agent2.Agent) {}),
	).Run()
}
