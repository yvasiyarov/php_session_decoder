package php_serialize

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"log"
)

func UnSerialize(s string) (PhpValue, error) {
	decoder := NewUnSerializer(s)
	decoder.SetSerializedDecodeFunc(SerializedDecodeFunc(UnSerialize))
	return decoder.Decode()
}

type UnSerializer struct {
	source     string
	r          *strings.Reader
	lastErr    error
	decodeFunc SerializedDecodeFunc
}

func NewUnSerializer(data string) *UnSerializer {
	return &UnSerializer{
		source: data,
	}
}

func (self *UnSerializer) SetReader(r *strings.Reader) {
	self.r = r
}

func (self *UnSerializer) SetSerializedDecodeFunc(f SerializedDecodeFunc) {
	self.decodeFunc = f
}

func (self *UnSerializer) Decode() (PhpValue, error) {
	if self.r == nil {
		self.r = strings.NewReader(self.source)
	}

	var value PhpValue

	if token, _, err := self.r.ReadRune(); err == nil {
		switch token {
		default:
			self.saveError(fmt.Errorf("php_serialize: Unknown token %#U", token))
		case TOKEN_NULL:
			value = self.decodeNull()
		case TOKEN_BOOL:
			value = self.decodeBool()
		case TOKEN_INT:
			value = self.decodeNumber(false)
		case TOKEN_FLOAT:
			value = self.decodeNumber(true)
		case TOKEN_STRING:
			value = self.decodeString(rune(DELIMITER_STRING_LEFT), rune(DELIMITER_STRING_RIGHT), true)
		case TOKEN_ARRAY:
			value = self.decodeArray()
		case TOKEN_OBJECT:
			value = self.decodeObject()
		case TOKEN_OBJECT_SERIALIZED:
			value = self.decodeSerialized()
		}
	}

	return value, self.lastErr
}

func (self *UnSerializer) decodeNull() PhpValue {
	self.expect(rune(SEPARATOR_VALUES))
	return nil
}

func (self *UnSerializer) decodeBool() PhpValue {
	var (
		raw rune
		err error
	)
	self.expect(rune(SEPARATOR_VALUE_TYPE))

	if raw, _, err = self.r.ReadRune(); err != nil {
		self.saveError(fmt.Errorf("php_serialize: Error while reading bool value: %v", err))
	}

	self.expect(rune(SEPARATOR_VALUES))
	return raw == '1'
}

func (self *UnSerializer) decodeNumber(isFloat bool) PhpValue {
	var (
		raw string
		err error
		val PhpValue
	)
	self.expect(rune(SEPARATOR_VALUE_TYPE))

	if raw, err = self.readUntil(byte(SEPARATOR_VALUES)); err != nil {
		self.saveError(fmt.Errorf("php_serialize: Error while reading number value: %v", err))
	} else {
		if isFloat {
			if val, err = strconv.ParseFloat(raw, 64); err != nil {
				self.saveError(fmt.Errorf("php_serialize: Unable to convert %s to float: %v", raw, err))
			}
		} else {
			if val, err = strconv.Atoi(raw); err != nil {
				self.saveError(fmt.Errorf("php_serialize: Unable to convert %s to int: %v", raw, err))
			}
		}
	}

	return val
}

func (self *UnSerializer) decodeString(left, right rune, isFinal bool) PhpValue {
	var (
		err     error
		val     PhpValue
		strLen  int
		readLen int
	)

	strLen = self.readLen()
	self.expect(left)

	if strLen > 0 {
		buf := make([]byte, strLen, strLen)
		if readLen, err = self.r.Read(buf); err != nil {
			self.saveError(fmt.Errorf("php_serialize: Error while reading string value: %v", err))
		} else {
			if readLen != strLen {
				self.saveError(fmt.Errorf("php_serialize: Unable to read string. Expected %d but have got %d bytes", strLen, readLen))
			} else {
				val = string(buf)
			}
		}
	}

	self.expect(right)
	if isFinal {
		self.expect(rune(SEPARATOR_VALUES))
	}
	return val
}

func (self *UnSerializer) decodeArray() PhpValue {
	var arrLen int
	val := make(PhpArray)

	arrLen = self.readLen()
	self.expect(rune(DELIMITER_OBJECT_LEFT))

	for i := 0; i < arrLen; i++ {
		k, errKey := self.Decode()
		v, errVal := self.Decode()

		if errKey == nil && errVal == nil {
			val[k] = v
			/*switch t := k.(type) {
			default:
				self.saveError(fmt.Errorf("php_serialize: Unexpected key type %T", t))
			case string:
				stringKey, _ := k.(string)
				val[stringKey] = v
			case int:
				intKey, _ := k.(int)
				val[strconv.Itoa(intKey)] = v
			}*/
		} else {
			self.saveError(fmt.Errorf("php_serialize: Error while reading key or(and) value of array"))
		}
	}

	self.expect(rune(DELIMITER_OBJECT_RIGHT))
	return val
}

func (self *UnSerializer) decodeObject() PhpValue {
	val := &PhpObject{
		className: self.readClassName(),
	}

	rawMembers := self.decodeArray()
	val.members, _ = rawMembers.(PhpArray)

	return val
}

func (self *UnSerializer) decodeSerialized() PhpValue {
	val := &PhpObjectSerialized{
		className: self.readClassName(),
	}

	rawData := self.decodeString(rune(DELIMITER_OBJECT_LEFT), rune(DELIMITER_OBJECT_RIGHT), false)
	val.data, _ = rawData.(string)

	if self.decodeFunc != nil && val.data != "" {
		var err error
		if val.value, err = self.decodeFunc(val.data); err != nil {
			self.saveError(err)
		}
	}

	return val
}

func (self *UnSerializer) expect(expected rune) {
	if token, _, err := self.r.ReadRune(); err != nil {
		self.saveError(fmt.Errorf("php_serialize: Error while reading expected rune %#U: %v", expected, err))
	} else if token != expected {
		if debugMode {
			log.Printf("php_serialize: source\n%s\n", self.source)
			log.Printf("php_serialize: reader info\n%#v\n", self.r)
		}
		self.saveError(fmt.Errorf("php_serialize: Expected %#U but have got %#U", expected, token))
	}
}

func (self *UnSerializer) readUntil(stop byte) (string, error) {
	var (
		token byte
		err   error
	)
	buf := bytes.NewBuffer([]byte{})

	for {
		if token, err = self.r.ReadByte(); err != nil || token == stop {
			break
		} else {
			buf.WriteByte(token)
		}
	}

	return buf.String(), err
}

func (self *UnSerializer) readLen() int {
	var (
		raw string
		err error
		val int
	)
	self.expect(rune(SEPARATOR_VALUE_TYPE))

	if raw, err = self.readUntil(byte(SEPARATOR_VALUE_TYPE)); err != nil {
		self.saveError(fmt.Errorf("php_serialize: Error while reading lenght of value: %v", err))
	} else {
		if val, err = strconv.Atoi(raw); err != nil {
			self.saveError(fmt.Errorf("php_serialize: Unable to convert %s to int: %v", raw, err))
		}
	}
	return val
}

func (self *UnSerializer) readClassName() (res string) {
	rawClass := self.decodeString(rune(DELIMITER_STRING_LEFT), rune(DELIMITER_STRING_RIGHT), false)
	res, _ = rawClass.(string)
	return
}

func (self *UnSerializer) saveError(err error) {
	if self.lastErr == nil {
		self.lastErr = err
	}
}
