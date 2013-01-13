package php_session_decoder;

import (
    "errors"
    "bytes"
    "strings"
    "fmt"
    "strconv"
)
const VALUE_NAME_SEPARATOR = '|'
const TYPE_VALUE_SEPARATOR = ':'
const VALUES_SEPARATOR     = ';'
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

func (decoder *PhpDecoder)Decode() (*PhpSessionData, error) {
    var resultErr error 
    for {
        if valueName, err := decoder.readUntil(VALUE_NAME_SEPARATOR); err == nil {
            if value, err:= decoder.DecodeValue(); err == nil {
                (*decoder.data)[valueName] = value
            } else {
                resultErr = errors.New(fmt.Sprintf("Can not read variable(%v) value:%v", valueName, err))
                break;
            }
        } else {
            break;
        }
    }
    return decoder.data, resultErr
}

func (decoder *PhpDecoder)DecodeValue() (PhpValue, error) {
    var (
        value PhpValue
        err error
    )
    
    if token, _, err := decoder.source.ReadRune(); err == nil {
        decoder.expect(TYPE_VALUE_SEPARATOR)
        switch token {
            case 'b': 
                if rawValue, _, _err := decoder.source.ReadRune(); _err == nil {
                    value = rawValue == '1'
                    err = errors.New("Can not read boolean value")
                } else {
                    err = errors.New("Can not read boolean value")
                }
 
                decoder.expect(VALUES_SEPARATOR)
            case 'i': 
                if rawValue, _err := decoder.readUntil(VALUES_SEPARATOR); _err == nil {
                    if value, _err = strconv.Atoi(rawValue); _err != nil {
                        err = errors.New(fmt.Sprintf("Can not convert %v to Int:%v", rawValue, _err))
                    }
                } else {
                    err = errors.New("Can not read int value")
                }
            case 'd': 
                if rawValue, _err := decoder.readUntil(VALUES_SEPARATOR); _err == nil {
                    if value, _err = strconv.ParseFloat(rawValue, 64); _err != nil {
                        err = errors.New(fmt.Sprintf("Can not convert %v to Float:%v", rawValue, _err))
                    }
                } else {
                    err = errors.New("Can not read float value")
                }
            case 's': 
                if rawStrlen, _err := decoder.readUntil(TYPE_VALUE_SEPARATOR); _err == nil {
                    if strLen, _err = strconv.Atoi(rawStrlen); _err != nil {
                        err = errors.New(fmt.Sprintf("Can not convert string length %v to int:%v", rawStrlen, _err))
                    } else {
                        decoder.expect('"')
                        tmpValue := new([]byte, 0, strLen)
                        if nRead, _err := decoder.source.Read(tmpValue); _err != nil || nRead != strLen {
                            err = errors.New(fmt.Sprintf("Can not read string content %v. Read only: %v from %v", _err, nRead, strLen))
                        } else {
                            value = string(tmpValue)
                            decoder.expect('"')
                            decoder.expect(VALUES_SEPARATOR)
                        }
                    }
                } else {
                    err = errors.New("Can not read string length")
                }
        }
    }
    return value, err
}

func (decoder *PhpDecoder) readUntil(stopByte byte) (string, error) {
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
    return result.String(), err
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

