package bot

import (
	"chatgpt/pkg/chatgpt"
	"testing"
	"time"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "test list character",
			args: []string{"character", "list"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatSession := chatgpt.NewChat(&chatgpt.RepeatedBot{}, chatgpt.NewMemoryLimitHistory(32, time.Minute*10))
			command := BuildCommand(chatSession)
			got, err := RunCmd(command, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
