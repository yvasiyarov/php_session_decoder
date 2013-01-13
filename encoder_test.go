package php_session_decoder;

import (
    "testing"
)

const BOOLEAN_VALUE_ENCODED = "login_ok|b:1;"
const BOOLEAN_VALUE_ENCODED_WITHOUT_NAME = "b:1;"
const INT_VALUE_ENCODED = "inteiro|i:34;"
const BOOLEAN_AND_INT_ENCODED = "login_ok|b:1;inteiro|i:34;"
const FLOAT_VALUE_ENCODED = "float_test|d:34.467999999900002;"
const STRING_VALUE_ENCODED = "name|s:9:\"some text\";"
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

func TestDecodeBooleanValue(t *testing.T) {
    decoder := NewPhpDecoder(BOOLEAN_VALUE_ENCODED)
    if result, err := decoder.Decode(); err != nil {
        t.Errorf("Can not decode boolens value %#v \n", err)
    } else {
        if v, ok := (*result)["login_ok"]; !ok {
            t.Errorf("Boolean value was not decoded \n")
        } else if v != true {
            t.Errorf("Boolean value was incorrectly decoded \n")
        }
    }
}

func TestDecodeIntValue(t *testing.T) {
    decoder := NewPhpDecoder(INT_VALUE_ENCODED)
    if result, err := decoder.Decode(); err != nil {
        t.Errorf("Can not decode int value %#v \n", err)
    } else {
        if v, ok := (*result)["inteiro"]; !ok {
            t.Errorf("Int value was not decoded \n")
        } else if v != 34 {
            t.Errorf("Int value was decoded incorrectly: %v\n", v)
        }
    }
}

func TestDecodeBooleanAndIntValue(t *testing.T) {
    decoder := NewPhpDecoder(BOOLEAN_AND_INT_ENCODED)
    if result, err := decoder.Decode(); err != nil {
        t.Errorf("Can not decode int value %#v \n", err)
    } else {
        if v, ok := (*result)["inteiro"]; !ok {
            t.Errorf("Int value was not decoded \n")
        } else if v != 34 {
            t.Errorf("Int value was decoded incorrectly: %v\n", v)
        }
    }
}

func TestDecodeFloatValue(t *testing.T) {
    decoder := NewPhpDecoder(FLOAT_VALUE_ENCODED)
    if result, err := decoder.Decode(); err != nil {
        t.Errorf("Can not decode float value %#v \n", err)
    } else {
        if v, ok := (*result)["float_test"]; !ok {
            t.Errorf("Float value was not decoded \n")
        } else if v != 34.467999999900002 {
            t.Errorf("Float value was decoded incorrectly: %v\n", v)
        }
    }
}

