package php_session_decoder

import "errors"

type CustomDecodeFunc func(string) (PhpSessionData, error)
type CustomEncodeFunc func(PhpValue) (string, error)

func SerializableDecode(s string) (valueData PhpSessionData, err error) {
	var (
		value		PhpValue
		ok			bool
	)

	decoder := NewPhpDecoder(s)
	decoder.DecodeFunc = SerializableDecode

	if value, err = decoder.DecodeValue(); err == nil {
		valueData, ok = value.(PhpSessionData)
		if !ok {
			err = errors.New("Error casting to PhpSessionData")
		}
	}

	return
}
