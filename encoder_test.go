package php_session_decoder;

import (
    "testing"
)

const BOOLEAN_VALUE_ENCODED = "login_ok|b:1;"
const BOOLEAN_VALUE_ENCODED_WITHOUT_NAME = "b:1;"
const BOOLEAN_AND_INT_ENCODED = "login_ok|b:1;inteiro|i:34;"

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
/*
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
    decoder := NewPhpDecoder(BOOLEAN_AND_INT_ENCODED)
    if result, err := decoder.Decode(); err != nil {
        t.Errorf("Can not decode boolens value %#v \n", err)
    } else {
        if _, ok := (*result)["login_ok"]; ok {
            t.Errorf("Boolean value was not decoded \n")
        }
        if _, ok := (*result)["inteiro"]; ok {
            t.Errorf("Int value was not decoded \n")
        }
    }
}
*/


