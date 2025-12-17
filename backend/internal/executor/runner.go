package executor

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultTimeoutSeconds = 15
	defaultMemoryMB       = 256
	defaultCPUs           = 1.0

	maxOutputBytes = 256 * 1024 // 256KB total (stdout+stderr)
)

type Runner struct {
	DockerBin string

	GoImage   string
	PyImage   string
	NodeImage string
}

func NewRunner() *Runner {
	r := &Runner{
		DockerBin: "docker",
		GoImage:   getenv("EXECUTOR_GO_IMAGE", "golang:1.22-alpine"),
		PyImage:   getenv("EXECUTOR_PY_IMAGE", "python:3.12-alpine"),
		NodeImage: getenv("EXECUTOR_NODE_IMAGE", "node:20-alpine"),
	}
	if d := os.Getenv("DOCKER_BIN"); d != "" {
		r.DockerBin = d
	}
	return r
}

func (r *Runner) Execute(ctx context.Context, req ExecuteRequest) ExecuteResponse {
	start := time.Now()

	timeout := req.TimeoutSeconds
	if timeout == 0 {
		timeout = defaultTimeoutSeconds
	}
	mem := req.MemoryMB
	if mem == 0 {
		mem = defaultMemoryMB
	}
	cpus := req.CPUs
	if cpus == 0 {
		cpus = defaultCPUs
	}

	// IMPORTANT:
	// executor runs inside a container and talks to host docker via docker.sock.
	// So we must create the workdir inside a NAMED VOLUME mounted into executor,
	// otherwise "docker run -v <path>" will point to host FS and files won't exist.
	base := "/workspaces"
	_ = os.MkdirAll(base, 0o755)

	workdir, err := os.MkdirTemp(base, "easyhire-exec-*")
	if err != nil {
		return fail("mktemp failed", err, start)
	}
	defer os.RemoveAll(workdir)

	if err := writeFiles(workdir, req.Files); err != nil {
		return fail("write files failed", err, start)
	}

	containerName := "easyhire-exec-" + randHex(8)

	image, cmdLine, err := buildCommand(req)
	if err != nil {
		return fail("invalid request", err, start)
	}

	wd := filepath.Base(workdir) // folder name inside the shared volume

	args := []string{
		"run", "--pull=never", "--rm",
		"--name", containerName,

		// sandbox (MVP)
		"--network", "none",
		"--read-only",
		"--pids-limit", "128",
		"--cap-drop", "ALL",
		"--security-opt", "no-new-privileges",

		// limits
		"--cpus", fmt.Sprintf("%.2f", cpus),
		"--memory", fmt.Sprintf("%dm", mem),
		"--memory-swap", fmt.Sprintf("%dm", mem),

		// give /tmp enough space for runtimes that still use it
		"--tmpfs", "/tmp:rw,noexec,nosuid,size=512m",

		// shared volume with files
		"-v", "easyhire_executor_work:/work:rw",
		"-w", filepath.ToSlash(filepath.Join("/work", wd)),

		// runtime defaults
		"-e", "HOME=/tmp",

		image,
		"sh", "-c", cmdLine,
	}

	execCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	var stdoutBuf, stderrBuf bytes.Buffer
	limStdout := &limitedWriter{W: &stdoutBuf, N: maxOutputBytes}
	limStderr := &limitedWriter{W: &stderrBuf, N: maxOutputBytes}

	cmd := exec.CommandContext(execCtx, r.DockerBin, args...)
	cmd.Stdout = limStdout
	cmd.Stderr = limStderr

	runErr := cmd.Run()

	// If timed out, force-remove container (docker client may die before cleanup)
	if errors.Is(execCtx.Err(), context.DeadlineExceeded) {
		_ = exec.Command(r.DockerBin, "rm", "-f", containerName).Run()
		return ExecuteResponse{
			OK:         false,
			Passed:     false,
			ExitCode:   124,
			Stdout:     stdoutBuf.String(),
			Stderr:     stderrBuf.String(),
			Duration:   time.Since(start),
			Error:      "timeout",
			Container:  containerName,
			Image:      image,
			Workdir:    workdir,
			Truncated:  limStdout.Truncated || limStderr.Truncated,
			OutputSize: stdoutBuf.Len() + stderrBuf.Len(),
		}
	}

	exitCode := exitCodeFromErr(runErr)
	passed := (exitCode == 0)

	ok := true
	errMsg := ""
	if runErr != nil {
		ok = false
		errMsg = runErr.Error()
	}

	return ExecuteResponse{
		OK:         ok,
		Passed:     passed,
		ExitCode:   exitCode,
		Stdout:     stdoutBuf.String(),
		Stderr:     stderrBuf.String(),
		Duration:   time.Since(start),
		Error:      errMsg,
		Container:  containerName,
		Image:      image,
		Workdir:    workdir,
		Truncated:  limStdout.Truncated || limStderr.Truncated,
		OutputSize: stdoutBuf.Len() + stderrBuf.Len(),
	}
}

func buildCommand(req ExecuteRequest) (image string, cmdLine string, err error) {
	switch req.Mode {
	case "go_test":
		// run fully inside /work/<run_id> (rw volume), do NOT use tiny /tmp for caches
		return "golang:1.22-alpine", `
set -e
mkdir -p ./tmp ./.cache/go-build ./.cache/gomod
export TMPDIR="$PWD/tmp"
export GOTMPDIR="$PWD/tmp"
export GOCACHE="$PWD/.cache/go-build"
export GOMODCACHE="$PWD/.cache/gomod"
go test ./... -count=1
`, nil

	case "python_unittest":
		return "python:3.12-alpine", `python -m unittest discover -s . -p "test_*.py" -q`, nil

	case "node_test":
		return "node:20-alpine", `node --test`, nil

	case "run":
		switch req.Language {
		case "go":
			return "golang:1.22-alpine", `
set -e
mkdir -p ./tmp ./.cache/go-build ./.cache/gomod
export TMPDIR="$PWD/tmp"
export GOTMPDIR="$PWD/tmp"
export GOCACHE="$PWD/.cache/go-build"
export GOMODCACHE="$PWD/.cache/gomod"
go run .
`, nil
		case "python":
			return "python:3.12-alpine", `python main.py`, nil
		case "javascript":
			return "node:20-alpine", `node main.js`, nil
		default:
			return "", "", fmt.Errorf("unsupported language: %s", req.Language)
		}

	default:
		return "", "", fmt.Errorf("unsupported mode: %s", req.Mode)
	}
}

func writeFiles(root string, files map[string]string) error {
	for p, content := range files {
		if strings.TrimSpace(p) == "" {
			return fmt.Errorf("empty filename")
		}
		clean := filepath.Clean(p)
		if strings.HasPrefix(clean, "..") || filepath.IsAbs(clean) {
			return fmt.Errorf("invalid path: %q", p)
		}

		full := filepath.Join(root, clean)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func exitCodeFromErr(err error) int {
	if err == nil {
		return 0
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return 1
}

type limitedWriter struct {
	W         io.Writer
	N         int
	Truncated bool
}

func (lw *limitedWriter) Write(p []byte) (int, error) {
	if lw.N <= 0 {
		lw.Truncated = true
		return len(p), nil
	}
	if len(p) > lw.N {
		lw.Truncated = true
		p = p[:lw.N]
	}
	n, err := lw.W.Write(p)
	lw.N -= n
	return len(p), err
}

func randHex(nbytes int) string {
	b := make([]byte, nbytes)
	if _, err := rand.Read(b); err != nil {
		// fallback, but should never happen
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func fail(msg string, err error, start time.Time) ExecuteResponse {
	return ExecuteResponse{
		OK:         false,
		Passed:     false,
		ExitCode:   1,
		Stdout:     "",
		Stderr:     "",
		Duration:   time.Since(start),
		Error:      fmt.Sprintf("%s: %v", msg, err),
		Container:  "",
		Image:      "",
		Workdir:    "",
		Truncated:  false,
		OutputSize: 0,
	}
}
