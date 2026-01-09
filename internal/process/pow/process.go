package pow

import (
	"context"
)

type CreateSignedSeedTask interface {
	Execute() (string, string, error)
}

type Response struct {
	Seed      string `json:"seed"`
	Signature string `json:"signature"`
}

type Process struct {
	createSignedSeedTask CreateSignedSeedTask
}

func NewProcess(task CreateSignedSeedTask) *Process {
	return &Process{
		createSignedSeedTask: task,
	}
}

func (p *Process) Process(ctx context.Context) (*Response, error) {
	seed, signature, err := p.createSignedSeedTask.Execute()
	if err != nil {
		return nil, err
	}

	return &Response{
		Seed:      seed,
		Signature: signature,
	}, nil
}
