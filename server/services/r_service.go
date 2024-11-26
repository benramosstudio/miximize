package services

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

type RService struct {
	scriptPath string
}

type RRequest struct {
	Variables []float64 `json:"variables"`
}

func NewRService(scriptPath string) *RService {
	return &RService{
		scriptPath: scriptPath,
	}
}

func (s *RService) ProcessRScript(req RRequest) ([]byte, error) {
	// Process R template
	tmpl, err := template.ParseFiles(s.scriptPath)
	if err != nil {
		return nil, fmt.Errorf("template parsing error: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "script-*.r")
	if err != nil {
		return nil, fmt.Errorf("temp file creation error: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := tmpl.Execute(tmpFile, req); err != nil {
		return nil, fmt.Errorf("template execution error: %w", err)
	}

	// Run R script
	cmd := exec.Command("Rscript", tmpFile.Name())
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("r script failed: %v, output: %s", err, output)
	}

	results, err := os.ReadFile("output.json")
	if err != nil {
		return nil, fmt.Errorf("reading results error: %w", err)
	}

	return results, nil
}
