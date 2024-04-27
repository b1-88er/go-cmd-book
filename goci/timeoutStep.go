package main

import (
	"context"
	"os/exec"
	"time"
)

type Command func(ctx context.Context, name string, arg ...string) *exec.Cmd
type timeoutStep struct {
	step    // type embedding, so much for lack inheritance lol
	timeout time.Duration
	command Command
}

func newTimeoutStep(name, exe, message, proj string, args []string, timeout time.Duration, command Command) timeoutStep {
	s := timeoutStep{}
	s.step = newStep(name, exe, message, proj, args)
	s.timeout = timeout
	if s.timeout == 0 {
		s.timeout = time.Second * 30
	}
	s.command = command

	return s
}

func (s timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := s.command(ctx, s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &stepErr{
				step:  s.name,
				msg:   "failed time out",
				cause: context.DeadlineExceeded,
			}
		}
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}
	return s.message, nil
}
