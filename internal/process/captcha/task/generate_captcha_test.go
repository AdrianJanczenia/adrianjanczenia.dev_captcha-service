package task

import (
	"fmt"
	"testing"
)

func TestGenerateCaptchaTask_Execute(t *testing.T) {
	task := NewGenerateCaptchaTask()

	id, b64, answer, err := task.Execute()

	fmt.Printf("%s", b64)
	fmt.Printf("\n%s\n", answer)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" || b64 == "" || answer == "" {
		t.Errorf("missing data: id=%s, b64=%s, answer=%s", id, b64, answer)
	}
}
