package php_session_decoder

import (
	"bytes"
	//	"errors"
	"fmt"
	"strconv"

//	"strings"
)

type PhpEncoder struct {
	dest bytes.Buffer
	data PhpSessionData
}

func NewPhpEncoder(sessionData PhpSessionData) *PhpEncoder {
	var buffer bytes.Buffer
	e := &PhpEncoder{
		dest: buffer,
		data: sessionData,
	}
	return e
}

//todo: check if root key will be some not ansi string
func (encoder *PhpEncoder) Encode() (string, error) {
	var resultErr error

	if encoder.data == nil {
		return "", nil
	}

	for k, v := range encoder.data {
		encoder.dest.WriteString(k)
		encoder.dest.WriteRune(VALUE_NAME_SEPARATOR)
		if resultErr = encoder.encodeValue(v); resultErr != nil {
			break
		}
		encoder.dest.WriteRune(VALUES_SEPARATOR)
	}
	return encoder.dest.String(), resultErr
}

func (encoder *PhpEncoder) encodeValue(value PhpValue) error {
	var err error
	switch t := value.(type) {
	default:
		err = fmt.Errorf("Unexpected type %T", t)
	case bool:
		encoder.dest.WriteString("b")
		encoder.dest.WriteRune(TYPE_VALUE_SEPARATOR)
		if bValue, _ := value.(bool); bValue {
			encoder.dest.WriteString("1")
		} else {
			encoder.dest.WriteString("0")
		}
	case int:
		encoder.dest.WriteString("i")
		encoder.dest.WriteRune(TYPE_VALUE_SEPARATOR)
		iValue, _ := value.(int)

		strValue := strconv.Itoa(iValue)
		encoder.dest.WriteString(strValue)
	case float32:
	case float64:
		encoder.dest.WriteString("d")
		encoder.dest.WriteRune(TYPE_VALUE_SEPARATOR)

		fValue, _ := value.(float64)
		strValue := strconv.FormatFloat(fValue, 'f', -1, 64)

		encoder.dest.WriteString(strValue)
	case string:
		encoder.dest.WriteString("s")
		encoder.dest.WriteRune(TYPE_VALUE_SEPARATOR)

		strValue, _ := value.(string)
        encoder.encodeString(strValue)
    case PhpSessionData:
		encoder.dest.WriteString("a")
		encoder.dest.WriteRune(TYPE_VALUE_SEPARATOR)

		arrValue, _ := value.(PhpSessionData)
        encoder.encodeArrayCore(arrValue)
    case *PhpObject:
		encoder.dest.WriteString("O")
		encoder.dest.WriteRune(TYPE_VALUE_SEPARATOR)

		objValue, _ := value.(*PhpObject)
        encoder.encodeString(objValue.GetClassName())

        encoder.encodeArrayCore(objValue.GetMembers())
	}
	return err
}

func (encoder *PhpEncoder) encodeString(strValue string) {
    valLen := strconv.Itoa(len(strValue))
    encoder.dest.WriteString(valLen)
    encoder.dest.WriteRune(TYPE_VALUE_SEPARATOR)

    encoder.dest.WriteRune('"')
    encoder.dest.WriteString(strValue)
    encoder.dest.WriteRune('"')
}

func (encoder *PhpEncoder) encodeArrayCore(arrValue PhpSessionData) error {
	var err error
    
    valLen := strconv.Itoa(len(arrValue))
    encoder.dest.WriteString(valLen)
    encoder.dest.WriteRune(TYPE_VALUE_SEPARATOR)
    
    encoder.dest.WriteRune('{')

    for k, v := range arrValue {
        if intKey, _err := strconv.Atoi(k); _err == nil {
            if err = encoder.encodeValue(intKey); err != nil {
                break
            }
        } else {
            if err = encoder.encodeValue(k); err != nil {
                break
            }
        }
        encoder.dest.WriteRune(VALUES_SEPARATOR)
        if err = encoder.encodeValue(v); err != nil {
            break
        }
        encoder.dest.WriteRune(VALUES_SEPARATOR)
    }

    encoder.dest.WriteRune('}')
	return err
}

