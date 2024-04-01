// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	TemplateFilePath string `envconfig:"PLUGIN_TEMPLATE_FILE_PATH"`
}

func Exec(ctx context.Context, args Args) error {
	log := logrus.New()

	templateFilePath := os.Getenv("PLUGIN_TEMPLATE_FILE_PATH")
	if templateFilePath == "" {
		log.Fatal("PLUGIN_TEMPLATE_FILE_PATH environment variable is not set.")
	}

	content, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		log.Fatalf("Error reading the template file: %v", err)
	}

	contentStr := string(content)
	contentStr = replacePlaceholders(contentStr, log)

	err = ioutil.WriteFile(templateFilePath, []byte(contentStr), 0644)
	if err != nil {
		log.Fatalf("Error writing the updated template file: %v", err)
	}

	return nil
}

func replacePlaceholders(content string, log *logrus.Logger) string {
	envVars := extractEnvVars()

	placeholders := findPlaceholders(content)
	for _, placeholder := range placeholders {
		lowerPlaceholder := strings.ToLower(placeholder)
		if val, ok := envVars[lowerPlaceholder]; ok {
			// Using regexp to match both cases with and without `.Values`
			var re = regexp.MustCompile(`\{\{\s*(\.Values\.)?` + regexp.QuoteMeta(placeholder) + `\s*\}\}`)
			content = re.ReplaceAllString(content, val)
			log.Infof("Replaced placeholder '%s' with environment variable value", placeholder)
		}
	}

	return content
}

func extractEnvVars() map[string]string {
	envVars := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "PLUGIN_") {
			// Remove the "PLUGIN_" prefix and convert to lower case for case-insensitive matching
			normalizedKey := strings.ToLower(strings.TrimPrefix(pair[0], "PLUGIN_"))
			envVars[normalizedKey] = pair[1]
		}
	}
	return envVars
}

func findPlaceholders(content string) []string {
	var result []string
	// Adjusted to match both `{{ .Values.VARIABLENAME }}` and `{{ VARIABLENAME }}`
	placeholderRegex := regexp.MustCompile(`\{\{\s*(?:\.Values\.)?(\w+)\s*\}\}`)
	matches := placeholderRegex.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	for _, match := range matches {
		m := match[1]
		if m != "" && !seen[m] {
			result = append(result, m)
			seen[m] = true
		}
	}
	return result
}
