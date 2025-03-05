package main

import (
	"YaLyceum/app/orchestrator"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		orchestrator.Orchestrator,
	).Run()
}
