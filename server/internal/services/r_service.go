package services

import (
	"bytes"
	"fmt"
	"text/template"

	"server/internal/models"
)

type RService struct {
}

func (s *RService) generateRobynScript(params *models.RobynParams) (string, error) {
	tmpl, err := template.New("robyn").Parse(rScriptTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
