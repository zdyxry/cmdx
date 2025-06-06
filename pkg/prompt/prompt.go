package prompt

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

const (
	confirmPromptType = "confirm"
)

type Prompt struct {
	Type    string   `json:"type"`
	Message string   `json:"message,omitempty"`
	Help    string   `json:"help,omitempty"`
	Options []string `json:"options,omitempty"`
	Command string   `json:"command,omitempty"`
}

func Create(prompt Prompt) survey.Prompt {
	switch prompt.Type {
	case "":
		return nil
	case "input":
		return &survey.Input{
			Message: prompt.Message,
			Help:    prompt.Help,
		}
	case "multiline":
		return &survey.Multiline{
			Message: prompt.Message,
			Help:    prompt.Help,
		}
	case "password":
		return &survey.Password{
			Message: prompt.Message,
			Help:    prompt.Help,
		}
	case "confirm":
		return &survey.Confirm{
			Message: prompt.Message,
			Help:    prompt.Help,
		}
	case "select":
		options := prompt.Options
		if prompt.Command != "" {
			cmdOptions, err := getOptionsFromCommand(prompt.Command)
			if err != nil {
				// Fall back to static options if command fails
				if len(options) == 0 {
					// If no fallback options, return nil to indicate prompt creation failed
					return nil
				}
			} else {
				options = cmdOptions
			}
		}
		return &survey.Select{
			Message: prompt.Message,
			Help:    prompt.Help,
			Options: options,
		}
	case "multi_select":
		options := prompt.Options
		if prompt.Command != "" {
			cmdOptions, err := getOptionsFromCommand(prompt.Command)
			if err != nil {
				// Fall back to static options if command fails
				if len(options) == 0 {
					// If no fallback options, return nil to indicate prompt creation failed
					return nil
				}
			} else {
				options = cmdOptions
			}
		}
		return &survey.MultiSelect{
			Message: prompt.Message,
			Help:    prompt.Help,
			Options: options,
		}
	case "editor":
		return &survey.Editor{
			Message:       prompt.Message,
			Help:          prompt.Help,
			HideDefault:   true,
			AppendDefault: true,
		}
	}
	return nil
}

func GetValue(prompt survey.Prompt, typ string) (any, error) {
	switch typ {
	case confirmPromptType:
		ans := false
		err := survey.AskOne(prompt, &ans)
		return ans, err
	case "select":
		ans := ""
		if err := survey.AskOne(prompt, &ans); err != nil {
			return nil, err
		}
		return ans, nil
	case "multi_select":
		ans := []string{}
		if err := survey.AskOne(prompt, &ans); err != nil {
			return nil, err
		}
		return ans, nil
	default:
		ans := ""
		return ans, survey.AskOne(prompt, &ans)
	}
}

func getOptionsFromCommand(command string) ([]string, error) {
	// Use bash to execute the command (similar to how cmdx executes scripts)
	var cmd *exec.Cmd
	if _, err := exec.LookPath("bash"); err != nil {
		cmd = exec.Command("sh", "-c", command)
	} else {
		cmd = exec.Command("bash", "-euo", "pipefail", "-c", command)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split output by lines and filter empty lines
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var options []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			options = append(options, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for consistency
	if options == nil {
		options = []string{}
	}

	return options, nil
}
