package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
	"time"

	"YaLyceum/internal/models"

	config2 "YaLyceum/internal/pkg/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Agent struct {
	ComputingPower int
	Jobs           chan models.Task
	Results        chan models.Task
	Wg             *sync.WaitGroup
	Log            *zap.Logger
	Shutdown       chan struct{}
	URL            string
}
type PostResult struct {
	ID           int64   `json:"id"`
	ExpressionID int64   `json:"expression_id"`
	Result       float64 `json:"result,omitempty"`
	Error        *string `json:"error,omitempty"`
}

func (a *Agent) Recieve() {

	client := http.DefaultClient
	for {
		select {
		case <-a.Shutdown:
			return
		default:
			req, _ := http.NewRequest(http.MethodGet, a.URL, nil)
			resp, err := client.Do(req)
			if err != nil {
				a.Log.Error("Agent Request failed", zap.Error(err))
				time.Sleep(1 * time.Second)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				time.Sleep(1 * time.Second)
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				a.Log.Error("Agent Reading body failed", zap.Error(err))
			}
			var task models.Task
			err = json.Unmarshal(body, &task)
			if err != nil {
				a.Log.Error("Agent Unmarshaling body failed", zap.Error(err))
			}
			a.Log.Info("Agent task received", zap.Any("task", task))
			a.Jobs <- task
		}
	}
}
func (a *Agent) Send() {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	for {
		select {
		case <-a.Shutdown:
			for {
				select {
				case task, ok := <-a.Results:
					if !ok {
						return
					}
					s := PostResult{
						ID:           task.ID,
						ExpressionID: task.ExpressionID,
						Result:       task.Result,
						Error:        task.Error,
					}
					jsonData, err := json.Marshal(s)
					if err != nil {
						a.Log.Error("Agent Marshaling task failed", zap.Error(err), zap.Any("task", task))
						continue
					}
					sended := false
					for i := 0; i < 3; i++ {
						req, _ := http.NewRequest(http.MethodPost, a.URL, bytes.NewBuffer(jsonData))
						resp, err := client.Do(req)
						if err != nil {
							a.Log.Error("Agent Request failed", zap.Error(err))
							time.Sleep(1 * time.Second)
							continue
						}
						if resp.StatusCode != http.StatusOK {
							a.Log.Error("Agent failed send result", zap.Any("task", task))
							break
						}
						sended = true
						break
					}
					if sended {
						a.Log.Info("Agent result send", zap.Any("task", task))
					} else {
						a.Log.Info("Agent failed send result", zap.Any("task", task))
					}
				default:
					return
				}
			}
		case task, ok := <-a.Results:
			if !ok {
				return
			}
			s := PostResult{
				ID:           task.ID,
				ExpressionID: task.ExpressionID,
				Result:       task.Result,
				Error:        task.Error,
			}
			jsonData, err := json.Marshal(s)
			if err != nil {
				a.Log.Error("Agent Marshaling task failed", zap.Error(err), zap.Any("task", task))
				continue
			}
			sended := false
			for i := 0; i < 3; i++ {
				req, _ := http.NewRequest(http.MethodPost, a.URL, bytes.NewBuffer(jsonData))
				resp, err := client.Do(req)
				if err != nil {
					a.Log.Error("Agent Request failed", zap.Error(err))
					time.Sleep(1 * time.Second)
					continue
				}
				if resp.StatusCode != http.StatusOK {
					a.Log.Error("Agent failed send result", zap.Any("task", task))
					break
				}
				sended = true
				break
			}
			if sended {
				a.Log.Info("Agent result send", zap.Any("task", task))
			} else {
				a.Log.Info("Agent failed send result", zap.Any("task", task))
			}
		}
	}

}
func (a *Agent) Start(ctx context.Context) error {
	a.Log.Info("Starting agent")
	a.Wg.Add(a.ComputingPower)
	for i := 0; i < a.ComputingPower; i++ {
		go a.Worker()
	}
	go a.Recieve()
	go a.Send()
	return nil
}
func (a *Agent) Stop(ctx context.Context) error {
	close(a.Shutdown)
	close(a.Jobs)
	defer close(a.Results)
	a.Wg.Wait()
	return nil
}
func (a *Agent) Worker() {
	defer a.Wg.Done()
	for j := range a.Jobs {
		dur, err := time.ParseDuration(fmt.Sprintf("%dms", j.OperationTime))
		if err != nil {
			continue
		}
		time.Sleep(dur)
		switch j.Operation {
		case models.Addition:
			res := j.Arg1 + j.Arg2
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
				Result:       res,
				Error:        nil,
			}
			a.Results <- t
		case models.Subtraction:
			res := j.Arg1 - j.Arg2
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
				Result:       res,
				Error:        nil,
			}
			a.Results <- t
		case models.Multiplication:
			res := j.Arg1 * j.Arg2
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
				Result:       res,
				Error:        nil,
			}
			a.Results <- t
		case models.Division:
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
			}
			if j.Arg2 == 0 {
				errMsg := ErrDivisionByZero.Error()
				t.Error = &errMsg
				t.Result = 0.0
			} else {
				res := j.Arg1 / j.Arg2
				t.Result = res
				t.Error = nil
			}
			a.Results <- t
		case models.Exponentiation:
			res := math.Pow(j.Arg1, j.Arg2)
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
				Result:       res,
				Error:        nil,
			}
			a.Results <- t
		case models.UnaryMinus:
			res := -j.Arg1
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
				Result:       res,
				Error:        nil,
			}
			a.Results <- t
		case models.Logarithm:
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
			}
			if j.Arg1 <= 0 || j.Arg1 == 1 {
				errMsg := ErrLogNotDefinedFor.Error()
				t.Error = &errMsg
				t.Result = .0
			} else if j.Arg2 <= 0.0 {
				errMsg := ErrLogOutOfFuncDomain.Error()
				t.Error = &errMsg
				t.Result = 0.0
			} else {
				res := math.Log(j.Arg2) / math.Log(j.Arg1)
				t.Result = res
				t.Error = nil
			}
			a.Results <- t
		case models.SquareRoot:
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
			}
			if j.Arg1 < 0 {
				errMsg := ErrSqrtOutOfDomain.Error()
				t.Error = &errMsg
				t.Result = 0.0
			} else {
				res := math.Sqrt(j.Arg1)
				t.Result = res
				t.Error = nil
			}
			a.Results <- t
		}
	}
}
func New(config config2.Config, lc fx.Lifecycle, log *zap.Logger) *Agent {
	agent := &Agent{
		Jobs:           make(chan models.Task, 100),
		Results:        make(chan models.Task, 100),
		ComputingPower: config.ComputingPower,
		Wg:             &sync.WaitGroup{},
		Log:            log,
		Shutdown:       make(chan struct{}),
		URL:            fmt.Sprintf("http://%s:%d/api/v1/internal/task", config.ServerHost, config.ServerPort),
	}
	log.Info("Agent created")
	lc.Append(fx.Hook{
		OnStart: agent.Start,
		OnStop:  agent.Stop,
	})
	log.Info("Agent started")
	return agent
}
