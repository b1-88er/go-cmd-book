package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mockCmdContext(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess"}

	cs = append(cs, exe)
	cs = append(cs, args...)

	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func mockCmdTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cmd := mockCmdContext(ctx, exe, args...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("GO_HELPER_TIMEOUT") == "1" {
		time.Sleep(10 * time.Second)
	}

	if os.Args[2] == "git" {
		fmt.Fprintln(os.Stdout, "everything up-to-date")
		os.Exit(0)
	}

	os.Exit(1)
}

func setupGit(t *testing.T, projPath string) (string, func()) {
	t.Helper()

	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	gitRemote := t.TempDir()
	proj := t.TempDir()

	// copy entire testdata to keep the paths as temp/testdata/<projPath>
	copyCmd := exec.Command("cp", "-r", "testdata", proj)
	// now the project is under temp path
	proj = filepath.Join(proj, projPath)

	if err := copyCmd.Run(); err != nil {
		t.Fatal(err)
	}

	remoteUri := fmt.Sprintf("file://%s", gitRemote)
	gitCmdList := []struct {
		args []string
		dir  string
		env  []string
	}{
		{[]string{"init", "--bare"}, gitRemote, nil},
		{[]string{"init"}, proj, nil},
		{[]string{"remote", "add", "origin", remoteUri}, proj, nil},
		{[]string{"add", "."}, proj, nil},
		{[]string{"commit", "-m", "test"}, proj, []string{
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@example.com",
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@example.com"},
		},
	}

	for _, g := range gitCmdList {
		gitCmd := exec.Command(gitExec, g.args...)
		gitCmd.Dir = g.dir

		if g.env != nil {
			gitCmd.Env = append(os.Environ(), g.env...)
		}

		if err := gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}

	return proj, func() {
		os.RemoveAll(proj)
	}
}
func TestRun(t *testing.T) {
	testCases := []struct {
		name      string
		proj      string
		out       string
		stderr    string
		expErr    error
		setpupGit bool
		mockCmd   func(ctx context.Context, name string, args ...string) *exec.Cmd
	}{
		{
			name:      "success",
			proj:      "testdata/tool",
			out:       "Go build: SUCCESS\nGo test: SUCCESS\nGo fmt: SUCCESS\nGit push: SUCCESS\n",
			stderr:    "",
			expErr:    nil,
			setpupGit: true,
			mockCmd:   exec.CommandContext,
		},
		{
			name:      "MockSuccess",
			proj:      "testdata/tool",
			out:       "Go build: SUCCESS\nGo test: SUCCESS\nGo fmt: SUCCESS\nGit push: SUCCESS\n",
			stderr:    "",
			expErr:    nil,
			setpupGit: true,
			mockCmd:   mockCmdContext,
		},
		{
			name:      "validation error",
			proj:      "testdata/toolErr",
			out:       "",
			expErr:    &stepErr{step: "go build"},
			setpupGit: false,
			mockCmd:   exec.CommandContext,
		},
		{
			name:      "format error",
			proj:      "testdata/toolFmtErr",
			out:       "Go build: SUCCESS\nGo test: SUCCESS\n",
			expErr:    &stepErr{step: "go fmt"},
			setpupGit: false,
			mockCmd:   exec.CommandContext,
		},
		{
			name:      "failTimeout",
			proj:      "./testdata/tool",
			out:       "Go build: SUCCESS\nGo test: SUCCESS\nGo fmt: SUCCESS\n",
			expErr:    context.DeadlineExceeded,
			setpupGit: false,
			mockCmd:   mockCmdTimeout,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setpupGit {
				if _, err := exec.LookPath("git"); err != nil {
					t.Skip("Git not installed, skipping")
				}
				copiedProject, cleanup := setupGit(t, testCase.proj)
				testCase.proj = copiedProject
				defer cleanup()
			}

			out := bytes.Buffer{}
			err := run(testCase.proj, &out, testCase.mockCmd)

			if err != nil {
				assert.ErrorIs(t, err, testCase.expErr)

			}
			assert.Equal(t, testCase.out, out.String())
		})
	}
}

func TestSignal(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		sig    syscall.Signal
		expErr error
	}{
		{"SIGINT", "./testdata/tool", syscall.SIGINT, ErrSignal},
		{"SIGTERM", "./testdata/tool", syscall.SIGTERM, ErrSignal},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			errCh := make(chan error)

			go func() {
				errCh <- run(testCase.proj, io.Discard, mockCmdContext)
			}()

			go func() {
				// Kill actually sends a signal, odd
				syscall.Kill(os.Getpid(), testCase.sig)
			}()
			assert.ErrorIs(t, <-errCh, testCase.expErr)
		})
	}

}
