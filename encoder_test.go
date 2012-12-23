package php_session_decoder;

import (
    "testing"
)

//"login_ok|b:1;"
//"login_ok|b:1;inteiro|i:34;"
func TestDecoderFabrica(t *testing.T) {
    decoder := NewPhpDecoder("")
    if decoder == nil {
        t.Error("Can not create decoder object\n")
    }
}

func TestDecodeBooleanValue(t *testing.T) {
    decoder := NewPhpDecoder("login_ok|b:1;")
    if result, err := decoder.Decode(); err != nil {
        t.Errorf("Can not decode boolens value %#v \n", err)
    } else {
        if _, ok := (*result)["login_ok"]; ok {
            t.Errorf("Boolean value was not decoded \n")
        }
    }
}

