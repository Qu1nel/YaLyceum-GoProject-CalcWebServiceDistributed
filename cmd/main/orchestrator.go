package main

import (
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/app/orchestrator"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		orchestrator.Orchestrator,
	).Run()
}
