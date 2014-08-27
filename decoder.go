package php_session_decoder

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"log"
)

type PhpDecoder struct {
	DecodeFunc	SerializableDecodeFunc
	RawData		string
	source		*strings.Reader
	data		PhpSessionData
}

func NewPhpDecoder(phpSession string) *PhpDecoder {
	sessionData := make(PhpSessionData)
	d := &PhpDecoder{
		RawData:	phpSession,
		source:		strings.NewReader(phpSession),
		data:		sessionData,
	}
	return d
}

func (decoder *PhpDecoder) Decode() (PhpSessionData, error) {
	var resultErr error
	for {
		if valueName, err := decoder.readUntil(VALUE_NAME_SEPARATOR); err == nil {
			if value, err := decoder.DecodeValue(); err == nil {
				decoder.data[valueName] = value
			} else {
				resultErr = fmt.Errorf("Can not read variable(%v) value:%v", valueName, err)
				break
			}
		} else {
			break
		}
	}
	return decoder.data, resultErr
}

func (decoder *PhpDecoder) DecodeValue() (PhpValue, error) {
	var (
		value PhpValue
		err   error
	)

	if token, _, err := decoder.source.ReadRune(); err == nil {
		decoder.expect(TYPE_VALUE_SEPARATOR)
		switch token {
		case 'N':
			value = nil
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
					err = fmt.Errorf("Can not convert %v to Int:%v", rawValue, _err)
				}
			} else {
				err = errors.New("Can not read int value")
			}
		case 'd':
			if rawValue, _err := decoder.readUntil(VALUES_SEPARATOR); _err == nil {
				if value, _err = strconv.ParseFloat(rawValue, 64); _err != nil {
					err = fmt.Errorf("Can not convert %v to Float:%v", rawValue, _err)
				}
			} else {
				err = errors.New("Can not read float value")
			}
		case 's':
			value, err = decoder.decodeString()
			decoder.expect(VALUES_SEPARATOR)
		case 'a':
			value, err = decoder.decodeArray()
			decoder.allow(VALUES_SEPARATOR)
		case 'O':
			value, err = decoder.decodeObject()
		case 'C':
			value, err = decoder.decodeSerializableObject()
		default:
			log.Panicf("Undefined token: %v [%#U]", token, token)
		}
	}
	return value, err
}

func (decoder *PhpDecoder) decodeObject() (*PhpObject, error) {
	value := &PhpObject{}
	var err error

	if value.className, err = decoder.decodeString(); err == nil {
		decoder.expect(TYPE_VALUE_SEPARATOR)
		value.members, err = decoder.decodeArray()
	}
	return value, err
}

func (decoder *PhpDecoder) decodeSerializableObject() (*PhpObject, error) {
	value := &PhpObject{}
	var err error

	if value.className, err = decoder.decodeString(); err == nil {
		decoder.expect(TYPE_VALUE_SEPARATOR)
		value.RawData, err = decoder.decodeStringWithDelimiters('{', '}')
	}

	if decoder.DecodeFunc != nil {
		value.members, err = decoder.DecodeFunc(value.RawData)
	}

	return value, err
}

func (decoder *PhpDecoder) decodeArray() (PhpSessionData, error) {
	value := make(PhpSessionData)
	var err error
	if rawArrlen, _err := decoder.readUntil(TYPE_VALUE_SEPARATOR); _err == nil {
		if arrLen, _err := strconv.Atoi(rawArrlen); _err != nil {
			err = fmt.Errorf("Can not convert array length %v to int:%v", rawArrlen, _err)
		} else {
			decoder.expect('{')
			for i := 0; i < arrLen; i++ {
				if k, _err := decoder.DecodeValue(); err != nil {
					err = fmt.Errorf("Can not read array key %v", _err)
				} else if v, _err := decoder.DecodeValue(); err != nil {
					err = fmt.Errorf("Can not read array value %v", _err)
				} else {
					switch t := k.(type) {
					default:
						err = fmt.Errorf("Unexpected key type %T", t)
					case string:
						stringKey, _ := k.(string)
						value[stringKey] = v
					case int:
						intKey, _ := k.(int)
						strKey := strconv.Itoa(intKey)
						value[strKey] = v
					}
				}
			}
			decoder.expect('}')
		}
	} else {
		err = errors.New("Can not read array length")
	}
	return value, err
}

func (decoder *PhpDecoder) decodeString() (string, error) {
	return decoder.decodeStringWithDelimiters('"', '"')
}

func (decoder *PhpDecoder) decodeStringWithDelimiters(left, right rune) (string, error) {
	var (
		value string
		err   error
	)
	if rawStrlen, _err := decoder.readUntil(TYPE_VALUE_SEPARATOR); _err == nil {
		if strLen, _err := strconv.Atoi(rawStrlen); _err != nil {
			err = fmt.Errorf("Can not convert string length %v to int:%v", rawStrlen, _err)
		} else {
			decoder.expect(left)
			tmpValue := make([]byte, strLen, strLen)
			if nRead, _err := decoder.source.Read(tmpValue); _err != nil || nRead != strLen {
				err = fmt.Errorf("Can not read string content %v. Read only: %v from %v", _err, nRead, strLen)
			} else {
				value = string(tmpValue)
				decoder.expect(right)
			}
		}
	} else {
		err = fmt.Errorf("Can not read string length with delimiters L:%v [%#U], R:%v [%#U]", left, left, right, right)
	}
	return value, err
}

func (decoder *PhpDecoder) readUntil(stopByte byte) (string, error) {
	result := new(bytes.Buffer)
	var (
		token byte
		err   error
	)
	for {
		if token, err = decoder.source.ReadByte(); err != nil || token == stopByte {
			break
		} else {
			result.WriteByte(token)
		}
	}
	return result.String(), err
}

func (decoder *PhpDecoder) expect(expectRune rune) error {
	token, _, err := decoder.source.ReadRune()
	if err != nil {
		err = fmt.Errorf("Can not read expected: %v", expectRune)
	} else if token == expectRune {
		err = fmt.Errorf("Read %v, but expected: %v", token, expectRune)
	}
	return err
}

func (decoder *PhpDecoder) allow(expectRune rune) error {
	token, _, err := decoder.source.ReadRune()
	if err != nil {
		err = errors.New("Can not read next rune")
	} else if token != expectRune {
		err = decoder.source.UnreadRune()
	}
	return err
}

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
