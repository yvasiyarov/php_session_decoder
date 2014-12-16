package php_serialize

import (
	"fmt"
	"strconv"
)

func Serialize(v PhpValue) (string, error) {
	encoder := NewSerializer()
	encoder.SetSerializedEncodeFunc(SerializedEncodeFunc(Serialize))
	return encoder.Encode(v)
}

type Serializer struct {
	lastErr    error
	encodeFunc SerializedEncodeFunc
}

func NewSerializer() *Serializer {
	return &Serializer{}
}

func (self *Serializer) SetSerializedEncodeFunc(f SerializedEncodeFunc) {
	self.encodeFunc = f
}

func (self *Serializer) Encode(v PhpValue) (string, error) {
	var value string

	switch t := v.(type) {
	default:
		self.saveError(fmt.Errorf("php_serialize: Unknown type %T with value %#v", t, v))
	case nil:
		value = self.encodeNull()
	case bool:
		value = self.encodeBool(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		value = self.encodeNumber(v)
	case string:
		value = self.encodeString(v, rune(DELIMITER_STRING_LEFT), rune(DELIMITER_STRING_RIGHT), true)
	case PhpArray, map[PhpValue]PhpValue:
		value = self.encodeArray(v, true)
	case *PhpObject:
		value = self.encodeObject(v)
	case *PhpObjectSerialized:
		value = self.encodeSerialized(v)
	}

	return value, self.lastErr
}

func (self *Serializer) encodeNull() string {
	return string(TOKEN_NULL) + string(SEPARATOR_VALUES)
}

func (self *Serializer) encodeBool(v PhpValue) string {
	val := "0"
	if bVal, ok := v.(bool); ok && bVal == true {
		val = "1"
	}
	return string(TOKEN_BOOL) + string(SEPARATOR_VALUE_TYPE) + val + string(SEPARATOR_VALUES)
}

func (self *Serializer) encodeNumber(v PhpValue) (res string) {
	var val string

	isFloat := false

	switch v.(type) {
	default:
		val = "0"
	case int:
		intVal, _ := v.(int)
		val = strconv.FormatInt(int64(intVal), 10)
	case int8:
		intVal, _ := v.(int8)
		val = strconv.FormatInt(int64(intVal), 10)
	case int16:
		intVal, _ := v.(int16)
		val = strconv.FormatInt(int64(intVal), 10)
	case int32:
		intVal, _ := v.(int32)
		val = strconv.FormatInt(int64(intVal), 10)
	case int64:
		intVal, _ := v.(int64)
		val = strconv.FormatInt(int64(intVal), 10)
	case uint:
		intVal, _ := v.(uint)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint8:
		intVal, _ := v.(uint8)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint16:
		intVal, _ := v.(uint16)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint32:
		intVal, _ := v.(uint32)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint64:
		intVal, _ := v.(uint64)
		val = strconv.FormatUint(uint64(intVal), 10)
	// PHP has precision = 17 by default
	case float32:
		floatVal, _ := v.(float32)
		val = strconv.FormatFloat(float64(floatVal), byte(FORMATTER_FLOAT), FORMATTER_PRECISION, 32)
		isFloat = true
	case float64:
		floatVal, _ := v.(float64)
		val = strconv.FormatFloat(float64(floatVal), byte(FORMATTER_FLOAT), FORMATTER_PRECISION, 64)
		isFloat = true
	}

	if isFloat {
		res = string(TOKEN_FLOAT)
	} else {
		res = string(TOKEN_INT)
	}

	res += string(SEPARATOR_VALUE_TYPE) + val + string(SEPARATOR_VALUES)
	return
}

func (self *Serializer) encodeString(v PhpValue, left, right rune, isFinal bool) (res string) {
	val, _ := v.(string)

	if isFinal {
		res = string(TOKEN_STRING)
	}
	res += self.prepareLen(len(val)) + string(left) + val + string(right)
	if isFinal {
		res += string(SEPARATOR_VALUES)
	}
	return
}

func (self *Serializer) encodeArray(v PhpValue, isFinal bool) (res string) {
	var (
		arrLen  int
		data, s string
	)

	if isFinal {
		res = string(TOKEN_ARRAY)
	}

	switch v.(type) {
	case PhpArray:
		arrVal, _ := v.(PhpArray)
		arrLen = len(arrVal)
		for k, v := range arrVal {
			s, _ = self.Encode(k)
			data += s
			s, _ = self.Encode(v)
			data += s
		}

	case map[PhpValue]PhpValue:
		arrVal, _ := v.(map[PhpValue]PhpValue)
		arrLen = len(arrVal)
		for k, v := range arrVal {
			s, _ = self.Encode(k)
			data += s
			s, _ = self.Encode(v)
			data += s
		}
	}

	res += self.prepareLen(arrLen) + string(DELIMITER_OBJECT_LEFT) + data + string(DELIMITER_OBJECT_RIGHT)
	return
}

func (self *Serializer) encodeObject(v PhpValue) string {
	obj, _ := v.(*PhpObject)
	return string(TOKEN_OBJECT) + self.prepareClassName(obj.className) + self.encodeArray(obj.members, false)
}

func (self *Serializer) encodeSerialized(v PhpValue) (res string) {
	var serialized string

	obj, _ := v.(*PhpObjectSerialized)
	res = string(TOKEN_OBJECT_SERIALIZED) + self.prepareClassName(obj.className)

	if self.encodeFunc == nil {
		serialized = obj.GetData()
	} else {
		var err error
		if serialized, err = self.encodeFunc(obj.GetValue()); err != nil {
			self.saveError(err)
		}
	}

	res += self.encodeString(serialized, rune(DELIMITER_OBJECT_LEFT), rune(DELIMITER_OBJECT_RIGHT), false)
	return
}

func (self *Serializer) prepareLen(l int) string {
	return string(SEPARATOR_VALUE_TYPE) + strconv.Itoa(l) + string(SEPARATOR_VALUE_TYPE)
}

func (self *Serializer) prepareClassName(name string) string {
	return self.encodeString(name, rune(DELIMITER_STRING_LEFT), rune(DELIMITER_STRING_RIGHT), false)
}

func (self *Serializer) saveError(err error) {
	if self.lastErr == nil {
		self.lastErr = err
	}
}

func wrapWithRune(s string, left, right rune) string {
	return string(left) + s + string(right)
}
