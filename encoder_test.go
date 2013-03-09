package php_session_decoder

import (
	"testing"
)

//TODO: write bool false parse test
// negative int parse test
// negative float parse test
// string parse test
// string with quotes parse test
// string with unicode parse
// string with $/@ and etc 
func TestDecoderFabrica(t *testing.T) {
	decoder := NewPhpDecoder("")
	if decoder == nil {
		t.Error("Can not create decoder object\n")
	}
}

const BOOLEAN_VALUE_ENCODED = "login_ok|b:1;"

func TestDecodeBooleanValueWithoutName(t *testing.T) {
	decoder := NewPhpDecoder(BOOLEAN_VALUE_ENCODED_WITHOUT_NAME)
	if result, err := decoder.DecodeValue(); err != nil {
		t.Errorf("Can not decode boolens value %#v \n", err)
	} else {
		if v, ok := (result).(bool); !ok {
			t.Errorf("Boolean value was not decoded \n")
		} else if v != true {
			t.Errorf("Boolean value was incorrectly decoded \n")
		}
	}
}

const BOOLEAN_VALUE_ENCODED_WITHOUT_NAME = "b:1;"

func TestDecodeBooleanValue(t *testing.T) {
	decoder := NewPhpDecoder(BOOLEAN_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode boolens value %#v \n", err)
	} else {
		if v, ok := (result)["login_ok"]; !ok {
			t.Errorf("Boolean value was not decoded \n")
		} else if v != true {
			t.Errorf("Boolean value was incorrectly decoded \n")
		}
	}
}

const INT_VALUE_ENCODED = "inteiro|i:34;"

func TestDecodeIntValue(t *testing.T) {
	decoder := NewPhpDecoder(INT_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode int value %#v \n", err)
	} else {
		if v, ok := (result)["inteiro"]; !ok {
			t.Errorf("Int value was not decoded \n")
		} else if v != 34 {
			t.Errorf("Int value was decoded incorrectly: %v\n", v)
		}
	}
}

const BOOLEAN_AND_INT_ENCODED = "login_ok|b:1;inteiro|i:34;"

func TestDecodeBooleanAndIntValue(t *testing.T) {
	decoder := NewPhpDecoder(BOOLEAN_AND_INT_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode int value %#v \n", err)
	} else {
		if v, ok := (result)["inteiro"]; !ok {
			t.Errorf("Int value was not decoded \n")
		} else if v != 34 {
			t.Errorf("Int value was decoded incorrectly: %v\n", v)
		}
	}
}

const FLOAT_VALUE_ENCODED = "float_test|d:34.467999999900002;"

func TestDecodeFloatValue(t *testing.T) {
	decoder := NewPhpDecoder(FLOAT_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode float value %#v \n", err)
	} else {
		if v, ok := (result)["float_test"]; !ok {
			t.Errorf("Float value was not decoded \n")
		} else if v != 34.467999999900002 {
			t.Errorf("Float value was decoded incorrectly: %v\n", v)
		}
	}
}

const STRING_VALUE_ENCODED = "name|s:9:\"some text\";"

func TestDecodeStringValue(t *testing.T) {
	decoder := NewPhpDecoder(STRING_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode string value %#v \n", err)
	} else {
		if v, ok := (result)["name"]; !ok {
			t.Errorf("String value was not decoded \n")
		} else if v != "some text" {
			t.Errorf("String value was decoded incorrectly: %v\n", v)
		}
	}
}

const ARRAY_VALUE_ENCODED = "arr|a:2:{s:4:\"test\";b:1;i:0;i:5;}"

func TestDecodeArrayValue(t *testing.T) {
	decoder := NewPhpDecoder(ARRAY_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		if v, ok := (result)["arr"]; !ok {
			t.Errorf("Array value was not decoded \n")
		} else if arrValue, ok := v.(PhpSessionData); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if value1, ok := arrValue["test"]; !ok || value1 != true {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		} else if value2, ok := arrValue["0"]; !ok || value2 != 5 {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		}
	}
}

const OBJECT_VALUE_ENCODED = "obj|O:10:\"TestObject\":3:{s:1:\"a\";i:5;s:13:\"\x00TestObject\x00b\";s:4:\"priv\";s:4:\"\x00*\x00c\";i:8;}"
func TestDecodeObjectValue(t *testing.T) {
	decoder := NewPhpDecoder(OBJECT_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode object value %#v \n", err)
	} else {
		if v, ok := (result)["obj"]; !ok {
			t.Errorf("Object value was not decoded \n")
                } else if objValue, ok := v.(*PhpObject); ok != true {
			t.Errorf("Object value was decoded incorrectly: %#v \n", v)
		} else if objValue.className != "TestObject" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", objValue.className)
		} else if value1, ok := objValue.GetPublicMemberValue("a"); !ok || value1 != 5 {
			t.Errorf("Public member of object was decoded incorrectly: %#v\n", v)
		}
	}
}

