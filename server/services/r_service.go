package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/benramosstudio/miximize/internal/models"
)

type RService struct {
	scriptPath string
}

func NewRService(scriptPath string) *RService {
	return &RService{
		scriptPath: scriptPath,
	}
}

func (s *RService) ProcessRScript(req models.RobynParams) ([]byte, error) {
	// Get absolute path for script template
	absScriptPath, err := filepath.Abs(s.scriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Process R template
	tmpl, err := template.ParseFiles(absScriptPath)
	if err != nil {
		return nil, fmt.Errorf("template parsing error: %w", err)
	}

	// Create temp file in system's temp directory
	tmpFile, err := os.CreateTemp("/tmp", "robyn-*.r")
	if err != nil {
		return nil, fmt.Errorf("temp file creation error: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := tmpl.Execute(tmpFile, req); err != nil {
		return nil, fmt.Errorf("template execution error: %w", err)
	}

	// Ensure file is written and closed
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("error closing temp file: %w", err)
	}

	// Run R script with full path to Rscript
	cmd := exec.Command("/usr/bin/Rscript", tmpFile.Name())
	output, err := cmd.CombinedOutput()

	// Always print R script output for debugging
	fmt.Printf("R Script Output:\n%s\n", string(output))

	if err != nil {
		return nil, fmt.Errorf("r script failed: %v, output: %s", err, string(output))
	}

	// Read results using absolute path
	results, err := os.ReadFile(filepath.Join(req.OutputDirectory, "output.json"))
	if err != nil {
		return nil, fmt.Errorf("reading results error: %w, R output was: %s", err, string(output))
	}

	return results, nil
}
