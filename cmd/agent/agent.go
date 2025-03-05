package main

import (
	"github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/app/agent"
	agent2 "github.com/Qu1nel/YaLyceum-GoProject-CalcWebServiceDistributed/internal/agent"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		agent.Agent,
		fx.Invoke(func(*agent2.Agent) {}),
	).Run()
}
