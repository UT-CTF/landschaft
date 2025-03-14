package embed_test

import (
	"strings"
	"testing"

	"github.com/UT-CTF/landschaft/embed"
)

func TestExecuteTestScript(t *testing.T) {
	output, err := embed.ExecuteScript("test/test.sh", false)
	if err != nil {
		t.Fatalf("Failed to execute test script: %v", err)
	}

	expected := "Hello Linux!"
	if strings.TrimSpace(output) != expected {
		t.Errorf("Expected output %q, got %q", expected, strings.TrimSpace(output))
	}
}

func TestExecuteArgsTestScript(t *testing.T) {
	output, err := embed.ExecuteScript("test/args_test.sh", false, "arg1test", "arg2test")
	if err != nil {
		t.Fatalf("Failed to execute test script: %v", err)
	}

	expected := "All arguments: arg1test arg2test"
	if strings.TrimSpace(output) != expected {
		t.Errorf("Expected output %q, got %q", expected, strings.TrimSpace(output))
	}
}
