package daemon

import (
	"context"
	"sync"
	"time"

	"github.com/Koyo-os/vote-crud-service/internal/entity"
	"github.com/Koyo-os/vote-crud-service/pkg/logger"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

const ChannelSize = 1

type (
	Operation struct {
		Timestamp     time.Time
		OperationType string
		Payload       []byte
	}

	Repository interface {
		Create(*entity.Vote) error
		Update(string, string, interface{}) error
		Delete(string) error
		Get(string) (*entity.Vote, error)
		GetMore(string, interface{}) ([]entity.Vote, error)
	}

	UpdateRequest struct {
		ID    string      `json:"id"`
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}

	DeleteRequest struct {
		ID string `json:"id"`
	}

	Daemon struct {
		mux    *sync.Mutex
		wg     *sync.WaitGroup
		logger *logger.Logger
		repo   Repository
		input  chan Operation
		output chan Operation
	}
)

func NewDaemon(repo Repository, output chan Operation, input chan Operation) *Daemon {
	var (
		mux sync.Mutex
		wg  sync.WaitGroup
	)

	return &Daemon{
		mux:    &mux,
		wg:     &wg,
		repo:   repo,
		logger: logger.Get(),
		output: output,
		input:  input,
	}
}

func (d *Daemon) Create(vote *entity.Vote) error {
	payload, err := sonic.Marshal(vote)
	if err != nil {
		return err
	}

	if err = d.repo.Create(vote); err != nil {
		return err
	}

	d.input <- Operation{
		Timestamp:     time.Now(),
		Payload:       payload,
		OperationType: "created",
	}

	return nil
}

func (d *Daemon) Update(id, key string, value interface{}) error {
	payload, err := sonic.Marshal(&UpdateRequest{
		ID:    id,
		Key:   key,
		Value: value,
	})
	if err != nil {
		return err
	}

	if err = d.repo.Update(id, key, value); err != nil {
		return err
	}

	d.input <- Operation{
		Timestamp:     time.Now(),
		Payload:       payload,
		OperationType: "update",
	}

	return nil
}

func (d *Daemon) Delete(id string) error {
	payload, err := sonic.Marshal(&DeleteRequest{
		ID: id,
	})
	if err != nil {
		return err
	}

	if err = d.repo.Delete(id); err != nil {
		return err
	}

	d.input <- Operation{
		Timestamp:     time.Now(),
		Payload:       payload,
		OperationType: "delete",
	}

	return nil
}

func (d *Daemon) GetByPollID(ctx context.Context, pollID uuid.UUID) chan entity.Vote {
	voteChan := make(chan entity.Vote, ChannelSize)

	go func() {
		for operation := range d.output {
			if operation.OperationType == "created" {
				var vote entity.Vote

				if err := sonic.Unmarshal(operation.Payload, &vote); err != nil {
					continue
				}

				if vote.PollID == pollID {
					voteChan <- vote
				}
			}
		}
	}()

	return voteChan
}

func (d *Daemon) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			d.logger.Info("daemon stopped!")
			return
		case req := <-d.input:
			d.output <- req

			d.logger.Info("vote created")
		}
	}
}
