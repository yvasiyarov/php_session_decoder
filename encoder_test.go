package php_session_decoder

import (
	"strings"
	"testing"
)

func TestEncoderFabrica(t *testing.T) {
	var data PhpSessionData
	encoder := NewPhpEncoder(data)
	if encoder == nil {
		t.Error("Can not create encoder object\n")
	}
}

func TestEncodeBooleanValue(t *testing.T) {
	data := make(PhpSessionData)
	data["login_ok"] = true
	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode boolens value %#v \n", err)
	} else {
		if result != BOOLEAN_VALUE_ENCODED {
			t.Errorf("Boolean value was encoded incorrectly %v \n", result)
		}
	}
}

func TestEncodeIntValue(t *testing.T) {
	data := make(PhpSessionData)
	data["inteiro"] = 34
	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode int value %#v \n", err)
	} else {
		if result != INT_VALUE_ENCODED {
			t.Errorf("Int value was encoded incorrectly %v \n", result)
		}
	}
}

func TestEncodeFloatValue(t *testing.T) {
	data := make(PhpSessionData)
	data["float_test"] = 34.4679999999
	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode float value %#v \n", err)
	} else {
		if result != FLOAT_VALUE_ENCODED {
			t.Errorf("Float value was encoded incorrectly %v\n", result)
		}
	}
}

func TestEncodeStringValue(t *testing.T) {
	data := make(PhpSessionData)
	data["name"] = "some text"
	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode string value %#v \n", err)
	} else {
		if result != STRING_VALUE_ENCODED {
			t.Errorf("String value was encoded incorrectly %v\n", result)
		}
	}
}

func TestEncodeArrayValue(t *testing.T) {
	data := make(PhpSessionData)
	data2 := make(PhpSessionData)
	data2["test"] = true
	data2["0"] = 5
	data2["test2"] = nil
	data["arr"] = data2

	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode array value %#v \n", err)
	} else {
		if !strings.Contains(result, "i:0;i:5;") || !strings.Contains(result, "s:4:\"test\";b:1") || !strings.Contains(result, "s:5:\"test2\";N") {
			t.Errorf("Array value was encoded incorrectly %v, %v\n", result, ARRAY_VALUE_ENCODED)
		}
    }
}

func TestEncodeObjectValue(t *testing.T) {
	data := make(PhpSessionData)

	obj := NewPhpObject()
	obj.SetClassName("TestObject")
	obj.SetPublicMemberValue("a", 5)
	obj.SetProtectedMemberValue("c", 8)
	obj.SetPrivateMemberValue("b", "priv")

	data["obj"] = obj

	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode object value %#v \n", err)
	} else {
		if !strings.Contains(result, "s:1:\"a\";i:5") || !strings.Contains(result, "10:\"TestObject\"") || !strings.Contains(result, "s:13:\"\x00TestObject\x00b\";s:4:\"priv\"") || !strings.Contains(result, "s:4:\"\x00*\x00c\";i:8") {
			t.Errorf("Object value was encoded incorrectly %v\n", result)
		}
	}
}
