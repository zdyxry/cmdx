package prompt

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
)

func Test_getOptionsFromCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    []string
		wantErr bool
	}{
		{
			name:    "echo command with multiple lines",
			command: "echo -e 'option1\\noption2\\noption3'",
			want:    []string{"option1", "option2", "option3"},
			wantErr: false,
		},
		{
			name:    "echo command with single line",
			command: "echo 'single-option'",
			want:    []string{"single-option"},
			wantErr: false,
		},
		{
			name:    "command with empty lines filtered out",
			command: "echo -e 'option1\\n\\noption2\\n\\n\\noption3'",
			want:    []string{"option1", "option2", "option3"},
			wantErr: false,
		},
		{
			name:    "command with whitespace trimmed",
			command: "echo -e '  option1  \\n  option2  \\n  option3  '",
			want:    []string{"option1", "option2", "option3"},
			wantErr: false,
		},
		{
			name:    "invalid command",
			command: "nonexistentcommand123",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "command with exit code 1",
			command: "sh -c 'exit 1'",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty command output",
			command: "echo -n ''",
			want:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getOptionsFromCommand(tt.command)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_createPrompt(t *testing.T) {
	data := []struct {
		title  string
		prompt Prompt
		exp    survey.Prompt
	}{
		{
			title: "input",
			prompt: Prompt{
				Type:    "input",
				Message: "message",
				Help:    "help",
			},
			exp: &survey.Input{
				Message: "message",
				Help:    "help",
			},
		},
		{
			title: "multiline",
			prompt: Prompt{
				Type:    "multiline",
				Message: "message",
				Help:    "help",
			},
			exp: &survey.Multiline{
				Message: "message",
				Help:    "help",
			},
		},
		{
			title: "password",
			prompt: Prompt{
				Type:    "password",
				Message: "message",
				Help:    "help",
			},
			exp: &survey.Password{
				Message: "message",
				Help:    "help",
			},
		},
		{
			title: "confirm",
			prompt: Prompt{
				Type:    "confirm",
				Message: "message",
				Help:    "help",
			},
			exp: &survey.Confirm{
				Message: "message",
				Help:    "help",
			},
		},
		{
			title: "editor",
			prompt: Prompt{
				Type:    "editor",
				Message: "message",
				Help:    "help",
			},
			exp: &survey.Editor{
				Message:       "message",
				Help:          "help",
				HideDefault:   true,
				AppendDefault: true,
			},
		},
		{
			title: "select",
			prompt: Prompt{
				Type:    "select",
				Message: "message",
				Help:    "help",
				Options: []string{"blue", "green"},
			},
			exp: &survey.Select{
				Message: "message",
				Help:    "help",
				Options: []string{"blue", "green"},
			},
		},
		{
			title: "multi_select",
			prompt: Prompt{
				Type:    "multi_select",
				Message: "message",
				Help:    "help",
				Options: []string{"blue", "green"},
			},
			exp: &survey.MultiSelect{
				Message: "message",
				Help:    "help",
				Options: []string{"blue", "green"},
			},
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			assert.Equal(t, d.exp, Create(d.prompt))
		})
	}
}

func Test_createPromptWithCommand(t *testing.T) {
	tests := []struct {
		name   string
		prompt Prompt
		want   survey.Prompt
	}{
		{
			name: "select with command",
			prompt: Prompt{
				Type:    "select",
				Message: "Choose option",
				Help:    "help text",
				Command: "echo -e 'option1\\noption2\\noption3'",
			},
			want: &survey.Select{
				Message: "Choose option",
				Help:    "help text",
				Options: []string{"option1", "option2", "option3"},
			},
		},
		{
			name: "multi_select with command",
			prompt: Prompt{
				Type:    "multi_select",
				Message: "Choose multiple",
				Help:    "help text",
				Command: "echo -e 'option1\\noption2\\noption3'",
			},
			want: &survey.MultiSelect{
				Message: "Choose multiple",
				Help:    "help text",
				Options: []string{"option1", "option2", "option3"},
			},
		},
		{
			name: "select with command fallback to static options",
			prompt: Prompt{
				Type:    "select",
				Message: "Choose option",
				Help:    "help text",
				Command: "nonexistentcommand123",
				Options: []string{"fallback1", "fallback2"},
			},
			want: &survey.Select{
				Message: "Choose option",
				Help:    "help text",
				Options: []string{"fallback1", "fallback2"},
			},
		},
		{
			name: "multi_select with command fallback to static options",
			prompt: Prompt{
				Type:    "multi_select",
				Message: "Choose multiple",
				Help:    "help text",
				Command: "sh -c 'exit 1'",
				Options: []string{"fallback1", "fallback2"},
			},
			want: &survey.MultiSelect{
				Message: "Choose multiple",
				Help:    "help text",
				Options: []string{"fallback1", "fallback2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Create(tt.prompt)
			assert.Equal(t, tt.want, got)
		})
	}
}
