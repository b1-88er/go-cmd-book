package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		t.Logf("running git %v", g.args)
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
		expErr    *stepErr
		setpupGit bool
	}{
		{
			name:      "success",
			proj:      "testdata/tool",
			out:       "Go build: SUCCESS\nGo test: SUCCESS\nGo fmt: SUCCESS\nGit push: SUCCESS\n",
			stderr:    "",
			expErr:    nil,
			setpupGit: true,
		},
		{
			name:      "validation error",
			proj:      "testdata/toolErr",
			out:       "",
			expErr:    &stepErr{step: "go build"},
			setpupGit: false,
		},
		{
			name:      "format error",
			proj:      "testdata/toolFmtErr",
			out:       "Go build: SUCCESS\nGo test: SUCCESS\n",
			expErr:    &stepErr{step: "go fmt"},
			setpupGit: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setpupGit {
				copiedProject, cleanup := setupGit(t, testCase.proj)
				testCase.proj = copiedProject
				defer cleanup()
			}
			out := bytes.Buffer{}
			err := run(testCase.proj, &out)
			assert.Equal(t, testCase.out, out.String())

			// both assertions do the same thing
			if err != nil {
				assert.ErrorIs(t, testCase.expErr, err)

			}
			if expErr, ok := (err).(*stepErr); ok {
				assert.Equal(t, expErr.step, testCase.expErr.step)
			}

		})
	}
}
