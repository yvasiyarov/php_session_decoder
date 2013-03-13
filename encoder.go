package php_session_decoder

import (
	"bytes"
//	"errors"
	"fmt"
//	"strconv"
//	"strings"
)


type PhpEncoder struct {
	dest     bytes.Buffer
	data     PhpSessionData
}

func NewPhpEncoder(sessionData PhpSessionData) *PhpEncoder {
        var buffer bytes.Buffer
	e := &PhpEncoder{
		dest:   buffer,
		data:   sessionData,
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

func (encoder *PhpEncoder) encodeValue(value PhpValue) (error) {
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
    //case string:
//        stringKey, _ := k.(string)
//    #        value[stringKey] = v
//    #case int:
//    #        intKey, _ := k.(int)
//    #        strKey := strconv.Itoa(intKey)
//    #        value[strKey] = v
    }
    return err 
}
