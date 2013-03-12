package php_session_decoder

import (
	"testing"
)

func TestEncoderFabrica(t *testing.T) {
        var data PhpSessionData
	encoder := NewPhpEncoder(&data)
	if encoder == nil {
		t.Error("Can not create encoder object\n")
	}
}

func TestEncodeBooleanValue(t *testing.T) {
        var data PhpSessionData
        data["login_ok"] = true 
	encoder := NewPhpDecoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode boolens value %#v \n", err)
	} else {
		if result != BOOLEAN_VALUE_ENCODED {
			t.Errorf("Boolean value was encoded incorrectly \n")
		}
	}
}


