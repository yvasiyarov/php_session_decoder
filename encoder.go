package php_session_decoder;

import (
    "errors"
    "bytes"
    "strings"
    "fmt"
)
const VALUES_SEPARATOR = '|'

/*
func encode() {
}
*/


type PhpValue interface{}

type PhpSessionData map[string]PhpValue

type PhpDecoder struct {
    source *strings.Reader
    position int
    data *PhpSessionData
}

func NewPhpDecoder(phpSession string) (*PhpDecoder) {
    sessionData := make(PhpSessionData)
    d := &PhpDecoder{
        source: strings.NewReader(phpSession), 
        position: 0,
        data: &sessionData,
    }
    return d
}

/*
func (decoder *PhpDecoder)Decode() (*PhpSessionData, error) {
    var (
        token byte
        err error
    )
    
    if valueName, err := decoder.readUntil(VALUES_SEPARATOR, false); err == nil {
        if value, err:= decoder.DecodeValue(); err == nil {
            decoder.data[string(valueName)] = value
        } else {
            errors.New("Can not read variable name:" + string(decoder.readUntil(VALUES_SEPARATOR, true)))
        }
    } else {
       err = errors.New("Can not read variable name")
    }
    return decoder.data, err
}
*/

func (decoder *PhpDecoder)DecodeValue() (PhpValue, error) {
    var (
        value PhpValue
        err error
    )
    
    if token, err := decoder.source.ReadByte(); err == nil {
        decoder.expect(':')
        switch token {
            case 'b': 
                if rawValue, _, _err := decoder.source.ReadRune(); _err == nil {
                    value = rawValue == '1'
                    err = errors.New("Can not read boolean value:")
                } else {
                    err = errors.New("Can not read boolean value")
                }
 
                decoder.expect(';')
        }
    }
    return value, err
}

func (decoder *PhpDecoder) readUntil(stopByte byte) ([]byte, error) {
    result := new(bytes.Buffer)
    var (
        token byte
        err error
    )
    for {
        if token, err = decoder.source.ReadByte(); err != nil || token == stopByte {
            break;
        } else {
            result.WriteByte(token)
        }
    }
    return result.Bytes(), err
}

func (decoder *PhpDecoder) expect(expectRune rune) (error) {
    token, _, err := decoder.source.ReadRune(); 
    if err != nil {
        err = errors.New(fmt.Sprintf("Can not read expected: %v", expectRune))
    } else if token == expectRune {
        err = errors.New(fmt.Sprintf("Read %v, but expected: %v", token, expectRune))
    }
    return err
}

