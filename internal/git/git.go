// Package git provides an integration with the git command
package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/apex/log"
)

// IsRepo returns true if current folder is a git repository.
func IsRepo() bool {
	out, err := Run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// RunEnv runs a git command with the specified env vars and returns its output or errors.
func RunEnv(env map[string]string, args ...string) (string, error) {
	// TODO: use exex.CommandContext here and refactor.
	var extraArgs = []string{
		"-c", "log.showSignature=false",
	}
	args = append(extraArgs, args...)
	/* #nosec */
	var cmd = exec.Command("git", args...)

	if env != nil {
		cmd.Env = []string{}
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(stderr.String())
	}

	log.WithField("args", args).Debug("running git")
	log.WithField("stdout", stdout.String()).
		WithField("stderr", stderr.String()).
		Debug("git result")
	return string(stdout.String()), nil
}

// Run runs a git command and returns its output or errors.
func Run(args ...string) (string, error) {
	return RunEnv(nil, args...)
}

// Clean the output.
func Clean(output string, err error) (string, error) {
	output = strings.Replace(strings.Split(output, "\n")[0], "'", "", -1)
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}
	return output, err
}
