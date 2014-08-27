package php_session_decoder

import (
	"io/ioutil"
	"testing"
)

func TestDecoderFabrica(t *testing.T) {
	decoder := NewPhpDecoder("")
	if decoder == nil {
		t.Error("Can not create decoder object\n")
	}
}

const BOOLEAN_VALUE_ENCODED_WITHOUT_NAME = "b:1;"

func TestDecodeBooleanValueWithoutName(t *testing.T) {
	decoder := NewPhpDecoder(BOOLEAN_VALUE_ENCODED_WITHOUT_NAME)
	if result, err := decoder.DecodeValue(); err != nil {
		t.Errorf("Can not decode boolens value %#v \n", err)
	} else {
		if v, ok := (result).(bool); !ok {
			t.Errorf("Boolean value was not decoded \n")
		} else if v != true {
			t.Errorf("Boolean value was incorrectly decoded \n")
		}
	}
}

const BOOLEAN_VALUE_ENCODED = "login_ok|b:1;"

func TestDecodeBooleanValue(t *testing.T) {
	decoder := NewPhpDecoder(BOOLEAN_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode boolens value %#v \n", err)
	} else {
		if v, ok := (result)["login_ok"]; !ok {
			t.Errorf("Boolean value was not decoded \n")
		} else if v != true {
			t.Errorf("Boolean value was incorrectly decoded \n")
		}
	}
}

const INT_VALUE_ENCODED = "inteiro|i:34;"

func TestDecodeIntValue(t *testing.T) {
	decoder := NewPhpDecoder(INT_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode int value %#v \n", err)
	} else {
		if v, ok := (result)["inteiro"]; !ok {
			t.Errorf("Int value was not decoded \n")
		} else if v != 34 {
			t.Errorf("Int value was decoded incorrectly: %v\n", v)
		}
	}
}

const BOOLEAN_AND_INT_ENCODED = "login_ok|b:1;inteiro|i:34;"

func TestDecodeBooleanAndIntValue(t *testing.T) {
	decoder := NewPhpDecoder(BOOLEAN_AND_INT_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode int value %#v \n", err)
	} else {
		if v, ok := (result)["inteiro"]; !ok {
			t.Errorf("Int value was not decoded \n")
		} else if v != 34 {
			t.Errorf("Int value was decoded incorrectly: %v\n", v)
		}
	}
}

const FLOAT_VALUE_ENCODED = "float_test|d:34.4679999999;"

func TestDecodeFloatValue(t *testing.T) {
	decoder := NewPhpDecoder(FLOAT_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode float value %#v \n", err)
	} else {
		if v, ok := (result)["float_test"]; !ok {
			t.Errorf("Float value was not decoded \n")
		} else if v != 34.4679999999 {
			t.Errorf("Float value was decoded incorrectly: %v\n", v)
		}
	}
}

const STRING_VALUE_ENCODED = "name|s:9:\"some text\";"

func TestDecodeStringValue(t *testing.T) {
	decoder := NewPhpDecoder(STRING_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode string value %#v \n", err)
	} else {
		if v, ok := (result)["name"]; !ok {
			t.Errorf("String value was not decoded \n")
		} else if v != "some text" {
			t.Errorf("String value was decoded incorrectly: %v\n", v)
		}
	}
}

const ARRAY_VALUE_ENCODED = "arr|a:3:{s:4:\"test\";b:1;i:0;i:5;s:5:\"test2\";N;};"

func TestDecodeArrayValue(t *testing.T) {
	decoder := NewPhpDecoder(ARRAY_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		if v, ok := (result)["arr"]; !ok {
			t.Errorf("Array value was not decoded \n")
		} else if arrValue, ok := v.(PhpSessionData); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if value1, ok := arrValue["test"]; !ok || value1 != true {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		} else if value2, ok := arrValue["0"]; !ok || value2 != 5 {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		} else if value3, ok := arrValue["test2"]; !ok || value3 != nil {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		}
	}
}

const OBJECT_VALUE_ENCODED = "obj|O:10:\"TestObject\":3:{s:1:\"a\";i:5;s:13:\"\x00TestObject\x00b\";s:4:\"priv\";s:4:\"\x00*\x00c\";i:8;}"

func TestDecodeObjectValue(t *testing.T) {
	decoder := NewPhpDecoder(OBJECT_VALUE_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode object value %#v \n", err)
	} else {
		if v, ok := (result)["obj"]; !ok {
			t.Errorf("Object value was not decoded \n")
		} else if objValue, ok := v.(*PhpObject); ok != true {
			t.Errorf("Object value was decoded incorrectly: %#v \n", v)
		} else if objValue.GetClassName() != "TestObject" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", objValue.GetClassName())
		} else if value1, ok := objValue.GetPublicMemberValue("a"); !ok || value1 != 5 {
			t.Errorf("Public member of object was decoded incorrectly: %#v\n", objValue.GetMembers())
		} else if value2, ok := objValue.GetPrivateMemberValue("b"); !ok || value2 != "priv" {
			t.Errorf("Private member of object was decoded incorrectly: %#v\n", objValue.GetMembers())
		} else if value3, ok := objValue.GetProtectedMemberValue("c"); !ok || value3 != 8 {
			t.Errorf("Protected member of object was decoded incorrectly: %#v\n", objValue.GetMembers())
		}
	}
}

const COMPLEX_ARRAY_ENCODED = "arr2|a:6:{s:10:\"bool_false\";b:0;s:7:\"neg_int\";i:-5;s:9:\"neg_float\";d:-5;s:6:\"quotes\";s:22:\"test\" and 'v' and `q` \";s:8:\"not_ansi\";s:8:\"тест\";s:5:\"test3\";s:15:\"@@@ test $$$ \\ \";}"

func TestDecodeComplexArrayValue(t *testing.T) {
	decoder := NewPhpDecoder(COMPLEX_ARRAY_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		if v, ok := (result)["arr2"]; !ok {
			t.Errorf("Array value was not decoded \n")
		} else if arrValue, ok := v.(PhpSessionData); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if value1, ok := arrValue["bool_false"]; !ok || value1 != false {
			t.Errorf("Bool false value was decoded incorrectly: %#v\n", v)
		} else if value2, ok := arrValue["neg_int"]; !ok || value2 != -5 {
			t.Errorf("Negative int value was decoded incorrectly: %#v\n", v)
		} else if value3, ok := arrValue["quotes"]; !ok || value3 != "test\" and 'v' and `q` " {
			t.Errorf("String with quotes was decoded incorrectly: %#v\n", v)
		} else if value4, ok := arrValue["not_ansi"]; !ok || value4 != "тест" {
			t.Errorf("String with not ansi symbols was decoded incorrectly: %#v\n", v)
		} else if value5, ok := arrValue["test3"]; !ok || value5 != "@@@ test $$$ \\ " {
			t.Errorf("String with special symbols was decoded incorrectly: %#v\n", v)
		}
	}
}

const MULTIDIMENSIONAL_ARRAY_ENCODED = "arr3|a:1:{s:4:\"dim1\";a:5:{i:0;s:4:\"dim2\";i:1;i:0;i:2;i:3;i:3;i:5;i:4;a:2:{i:0;s:4:\"dim3\";i:1;i:5;}}}"

func TestDecodeMultidimensionalArrayValue(t *testing.T) {
	decoder := NewPhpDecoder(MULTIDIMENSIONAL_ARRAY_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		if v, ok := (result)["arr3"]; !ok {
			t.Errorf("Array value was not decoded \n")
		} else if arrValue, ok := v.(PhpSessionData); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if dim1, ok := arrValue["dim1"]; !ok {
			t.Errorf("Second dimension of array was decoded incorrectly: %#v\n", dim1)
		} else if dim1Value, ok := dim1.(PhpSessionData); ok != true {
			t.Errorf("Second dimension of array was decoded incorrectly: %#v\n", dim1Value)
		} else if dim1Value, ok := dim1.(PhpSessionData); ok != true {
			t.Errorf("Second dimension of array was decoded incorrectly: %#v\n", dim1Value)
		} else if value1, ok := dim1Value["0"]; !ok || value1 != "dim2" {
			t.Errorf("Second dimension of array was decoded incorrectly: %#v\n", value1)
		} else if value2, ok := dim1Value["3"]; !ok || value2 != 5 {
			t.Errorf("Second dimension of array was decoded incorrectly: %#v\n", value2)
		} else if dim2, ok := dim1Value["4"]; !ok {
			t.Errorf("Third dimension of array was decoded incorrectly: %#v\n", dim2)
		} else if dim2Value, ok := dim2.(PhpSessionData); ok != true {
			t.Errorf("Third dimension of array was decoded incorrectly: %#v\n", dim2Value)
		} else if value3, ok := dim2Value["0"]; !ok || value3 != "dim3" {
			t.Errorf("Third dimension of array was decoded incorrectly: %#v\n", value3)
		}
	}
}

const MULTIPLE_ARRAYS_ENCODED_WITH_SEMICOLONS = "array1|a:1:{s:5:\"test1\";b:1;};array2|a:1:{s:5:\"test2\";b:1;};"

func TestDecodeMultipleArraysWithSemicolons(t *testing.T) {
	decoder := NewPhpDecoder(MULTIPLE_ARRAYS_ENCODED_WITH_SEMICOLONS)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		if v, ok := result["array1"]; !ok {
			t.Errorf("First array was not decoded \n")
		} else if arrValue, ok := v.(PhpSessionData); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if value1, ok := arrValue["test1"]; !ok || value1 != true {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		}

		if v, ok := result["array2"]; !ok {
			t.Errorf("Second array was not decoded \n")
		} else if arrValue, ok := v.(PhpSessionData); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if value2, ok := arrValue["test2"]; !ok || value2 != true {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		}
	}
}

const MULTIPLE_ARRAYS_ENCODED_WITHOUT_SEMICOLONS = "array1|a:1:{s:5:\"test1\";b:1;}array2|a:1:{s:5:\"test2\";b:1;}"

func TestDecodeMultipleArraysWithoutSemicolons(t *testing.T) {
	decoder := NewPhpDecoder(MULTIPLE_ARRAYS_ENCODED_WITHOUT_SEMICOLONS)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		if v, ok := result["array1"]; !ok {
			t.Errorf("First array was not decoded \n")
		} else if arrValue, ok := v.(PhpSessionData); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if value1, ok := arrValue["test1"]; !ok || value1 != true {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		}

		if v, ok := result["array2"]; !ok {
			t.Errorf("Second array was not decoded \n")
		} else if arrValue, ok := v.(PhpSessionData); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if value2, ok := arrValue["test2"]; !ok || value2 != true {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		}
	}
}

const SERIALIZABLE_OBJECT_VALUE_NO_FUNC_ENCODED = "obj|C:10:\"TestObject\":49:{a:3:{s:1:\"a\";i:5;s:1:\"b\";s:4:\"priv\";s:1:\"c\";i:8;}}"

func TestDecodeSerializableObjectValueNoFunc(t *testing.T) {
	decoder := NewPhpDecoder(SERIALIZABLE_OBJECT_VALUE_NO_FUNC_ENCODED)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode object value %#v \n", err)
	} else {
		if v, ok := (result)["obj"]; !ok {
			t.Errorf("Object value was not decoded \n")
		} else if objValue, ok := v.(*PhpObject); ok != true {
			t.Errorf("Object value was decoded incorrectly: %#v \n", v)
		} else if objValue.GetClassName() != "TestObject" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", objValue.GetClassName())
		} else if objValue.RawData != "a:3:{s:1:\"a\";i:5;s:1:\"b\";s:4:\"priv\";s:1:\"c\";i:8;}" {
			t.Errorf("RawData of object was decoded incorrectly: %#v\n", objValue.RawData)
		}
	}
}

const SERIALIZABLE_OBJECT_VALUE_ENCODED = "object|C:10:\"TestObject\":96:{a:1:{s:4:\"item\";O:8:\"AbcClass\":3:{s:1:\"a\";i:5;s:11:\"\x00AbcClass\x00b\";s:7:\"private\";s:4:\"\x00*\x00c\";i:8;}}}"

func TestDecodeSerializableObjectValue(t *testing.T) {
	decoder := NewPhpDecoder(SERIALIZABLE_OBJECT_VALUE_ENCODED)
	decoder.DecodeFunc = SerializableDecode
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode object value %#v \n", err)
	} else {
		if v, ok := (result)["object"]; !ok {
			t.Errorf("Object value was not decoded \n")
		} else if objValue, ok := v.(*PhpObject); ok != true {
			t.Errorf("Object value was decoded incorrectly: %#v \n", v)
		} else if objValue.GetClassName() != "TestObject" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", objValue.GetClassName())
		} else if objValue.RawData != "a:1:{s:4:\"item\";O:8:\"AbcClass\":3:{s:1:\"a\";i:5;s:11:\"\x00AbcClass\x00b\";s:7:\"private\";s:4:\"\x00*\x00c\";i:8;}}" {
			t.Errorf("RawData of object was decoded incorrectly: %#v\n", objValue.RawData)
		} else if itemValue, ok := objValue.GetPublicMemberValue("item"); !ok {
			t.Errorf("Public member of object was decoded incorrectly: %#v\n", objValue.GetMembers())
		} else if itemObjValue, ok := itemValue.(*PhpObject); ok != true {
			t.Errorf("Object value was decoded incorrectly: %#v \n", v)
		} else if itemObjValue.GetClassName() != "AbcClass" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", itemObjValue.GetClassName())
		} else if value1, ok := itemObjValue.GetPublicMemberValue("a"); !ok || value1 != 5 {
			t.Errorf("Public member of object was decoded incorrectly: %#v\n", objValue.GetMembers())
		} else if value2, ok := itemObjValue.GetPrivateMemberValue("b"); !ok || value2 != "private" {
			t.Errorf("Private member of object was decoded incorrectly: %#v\n", objValue.GetMembers())
		} else if value3, ok := itemObjValue.GetProtectedMemberValue("c"); !ok || value3 != 8 {
			t.Errorf("Protected member of object was decoded incorrectly: %#v\n", objValue.GetMembers())
		}
	}
}

func TestDecodeRealData(t *testing.T) {
	testData, _ := ioutil.ReadFile("./data/test.session")
	decoder := NewPhpDecoder(string(testData))
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		rootKeys := []string{"product_last_viewed", "core", "customer", "checkout", "store_default", "catalog", "object"}
		for _, v := range rootKeys {
			if _, ok := result[v]; !ok {
				t.Errorf("Can not find %v key\n", v)
			}
		}
	}
}
