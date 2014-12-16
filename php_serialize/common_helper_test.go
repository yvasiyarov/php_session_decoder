package php_serialize

import "testing"

func TestPhpValueString(t *testing.T) {
	var (
		val      PhpValue = "string"
		expected string   = "string"
	)
	if newVal := PhpValueString(val); newVal != expected {
		t.Errorf("Expected %q but got %q", expected, newVal)
	}
}

func TestPhpValueBool(t *testing.T) {
	var (
		val      PhpValue = true
		expected bool     = true
	)
	if newVal := PhpValueBool(val); newVal != expected {
		t.Errorf("Expected %t but got %t", expected, newVal)
	}
}

func TestPhpValueInt(t *testing.T) {
	var (
		val      PhpValue = 10
		expected int      = 10
	)
	if newVal := PhpValueInt(val); newVal != expected {
		t.Errorf("Expected %d but got %d", expected, newVal)
	}
}

func TestPhpValueInt64(t *testing.T) {
	var (
		val      PhpValue = int64(10)
		expected int64    = 10
	)
	if newVal := PhpValueInt64(val); newVal != expected {
		t.Errorf("Expected %d but got %d", expected, newVal)
	}
}

func TestPhpValueUInt(t *testing.T) {
	var (
		val      PhpValue = uint(10)
		expected uint     = 10
	)
	if newVal := PhpValueUInt(val); newVal != expected {
		t.Errorf("Expected %d but got %d", expected, newVal)
	}
}

func TestPhpValueUInt64(t *testing.T) {
	var (
		val      PhpValue = uint64(10)
		expected uint64   = 10
	)
	if newVal := PhpValueUInt64(val); newVal != expected {
		t.Errorf("Expected %d but got %d", expected, newVal)
	}
}

func TestPhpValueFloat64(t *testing.T) {
	var (
		val      PhpValue = float64(10.0)
		expected float64  = 10.0
	)
	if newVal := PhpValueFloat64(val); newVal != expected {
		t.Errorf("Expected %v but got %v", expected, newVal)
	}
}
