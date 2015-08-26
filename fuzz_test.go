package php_session_decoder

import (
	"testing"
)

func TestFuzzCrashers(t *testing.T) {

	var crashers = []string{
		"|C2984619140625:",
		"|C9478759765625:",
		"|C :590791705756156:",
		"|C298461940625:",
	}

	for _, f := range crashers {
		decoder := NewPhpDecoder(f)
		decoder.Decode()
	}
}
