package embed_test

import (
	"strings"
	"testing"

	"github.com/UT-CTF/landschaft/embed"
)

func TestExecuteTestScript(t *testing.T) {
	output, err := embed.ExecuteScript("test/test.ps1", false)
	if err != nil {
		t.Fatalf("Failed to execute test script: %v", err)
	}

	expected := "Hello Windows!"
	if strings.TrimSpace(output) != expected {
		t.Errorf("Expected output %q, got %q", expected, strings.TrimSpace(output))
	}
}

func TestExecuteArgsTestScript(t *testing.T) {
	output, err := embed.ExecuteScript("test/args_test.ps1", false, "-TestArg", "Hello!")
	if err != nil {
		t.Fatalf("Failed to execute test script: %v", err)
	}

	expected := "Received argument: Hello!"
	if strings.TrimSpace(output) != expected {
		t.Errorf("Expected output %q, got %q", expected, strings.TrimSpace(output))
	}
}
