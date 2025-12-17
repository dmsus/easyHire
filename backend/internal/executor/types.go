package executor

import "time"

type ExecuteRequest struct {
	Language       string            `json:"language" binding:"required,oneof=go python javascript"`
	Mode           string            `json:"mode" binding:"required,oneof=go_test python_unittest node_test run"`
	Files          map[string]string `json:"files" binding:"required,min=1"`
	TimeoutSeconds int               `json:"timeout_seconds" binding:"min=1,max=120"`
	CPUs           float64           `json:"cpus" binding:"omitempty,min=0.1,max=4"`
	MemoryMB       int               `json:"memory_mb" binding:"omitempty,min=64,max=2048"`
}

type ExecuteResponse struct {
	OK         bool          `json:"ok"`
	Passed     bool          `json:"passed"`
	ExitCode   int           `json:"exit_code"`
	Stdout     string        `json:"stdout"`
	Stderr     string        `json:"stderr"`
	Duration   time.Duration `json:"duration"`
	Error      string        `json:"error,omitempty"`
	Container  string        `json:"container,omitempty"`
	Image      string        `json:"image,omitempty"`
	Workdir    string        `json:"workdir,omitempty"`
	Truncated  bool          `json:"truncated"`
	OutputSize int           `json:"output_size"`
}
