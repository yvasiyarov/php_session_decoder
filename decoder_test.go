package php_session_decoder

import (
	"testing"
	"io/ioutil"
	"encoding/json"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
)

func TestDecodeBooleanValue(t *testing.T) {
	decoder := NewPhpDecoder("login_ok|b:1;")
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

func TestDecodeIntValue(t *testing.T) {
	decoder := NewPhpDecoder("inteiro|i:34;")
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

func TestDecodeBooleanAndIntValue(t *testing.T) {
	decoder := NewPhpDecoder("login_ok|b:1;inteiro|i:34;")
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

func TestDecodeFloatValue(t *testing.T) {
	decoder := NewPhpDecoder("float_test|d:34.4679999999;")
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

func TestDecodeStringValue(t *testing.T) {
	decoder := NewPhpDecoder("name|s:9:\"some text\";")
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

func TestDecodeArrayValue(t *testing.T) {
	decoder := NewPhpDecoder("arr|a:3:{s:4:\"test\";b:1;i:0;i:5;s:5:\"test2\";N;};")
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		if v, ok := (result)["arr"]; !ok {
			t.Errorf("Array value was not decoded \n")
		} else if arrValue, ok := v.(php_serialize.PhpArray); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if value1, ok := arrValue["test"]; !ok || value1 != true {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		} else if value2, ok := arrValue[php_serialize.PhpValue(0)]; !ok || value2 != 5 {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		} else if value3, ok := arrValue["test2"]; !ok || value3 != nil {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		}
	}
}

func TestDecodeObjectValue(t *testing.T) {
	decoder := NewPhpDecoder("obj|O:10:\"TestObject\":3:{s:1:\"a\";i:5;s:13:\"\x00TestObject\x00b\";s:4:\"priv\";s:4:\"\x00*\x00c\";i:8;}")
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode object value %#v \n", err)
	} else {
		if v, ok := (result)["obj"]; !ok {
			t.Errorf("Object value was not decoded \n")
		} else if objValue, ok := v.(*php_serialize.PhpObject); ok != true {
			t.Errorf("Object value was decoded incorrectly: %#v \n", v)
		} else if objValue.GetClassName() != "TestObject" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", objValue.GetClassName())
		} else if value1, ok := objValue.GetPublic("a"); !ok || value1 != 5 {
			t.Errorf("Public member of object was decoded incorrectly: %#v\n", objValue.GetMembers())
		} else if value2, ok := objValue.GetPrivate("b"); !ok || value2 != "priv" {
			t.Errorf("Private member of object was decoded incorrectly: %#v\n", objValue.GetMembers())
		} else if value3, ok := objValue.GetProtected("c"); !ok || value3 != 8 {
			t.Errorf("Protected member of object was decoded incorrectly: %#v\n", objValue.GetMembers())
		}
	}
}

func TestDecodeComplexArrayValue(t *testing.T) {
	decoder := NewPhpDecoder("arr2|a:6:{s:10:\"bool_false\";b:0;s:7:\"neg_int\";i:-5;s:9:\"neg_float\";d:-5;s:6:\"quotes\";s:22:\"test\" and 'v' and `q` \";s:8:\"not_ansi\";s:8:\"тест\";s:5:\"test3\";s:15:\"@@@ test $$$ \\ \";}")
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		if v, ok := (result)["arr2"]; !ok {
			t.Errorf("Array value was not decoded \n")
		} else if arrValue, ok := v.(php_serialize.PhpArray); ok != true {
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

func TestDecodeMultidimensionalArrayValue(t *testing.T) {
	decoder := NewPhpDecoder("arr3|a:1:{s:4:\"dim1\";a:5:{i:0;s:4:\"dim2\";i:1;i:0;i:2;i:3;i:3;i:5;i:4;a:2:{i:0;s:4:\"dim3\";i:1;i:5;}}}")
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		if v, ok := (result)["arr3"]; !ok {
			t.Errorf("Array value was not decoded \n")
		} else if arrValue, ok := v.(php_serialize.PhpArray); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if dim1, ok := arrValue["dim1"]; !ok {
			t.Errorf("Second dimension of array was decoded incorrectly: %#v\n", dim1)
		} else if dim1Value, ok := dim1.(php_serialize.PhpArray); ok != true {
			t.Errorf("Second dimension of array was decoded incorrectly: %#v\n", dim1Value)
		} else if value1, ok := dim1Value[php_serialize.PhpValue(0)]; !ok || value1 != "dim2" {
			t.Errorf("Second dimension of array was decoded incorrectly: %#v\n", value1)
		} else if value2, ok := dim1Value[php_serialize.PhpValue(3)]; !ok || value2 != 5 {
			t.Errorf("Second dimension of array was decoded incorrectly: %#v\n", value2)
		} else if dim2, ok := dim1Value[php_serialize.PhpValue(4)]; !ok {
			t.Errorf("Third dimension of array was decoded incorrectly: %#v\n", dim2)
		} else if dim2Value, ok := dim2.(php_serialize.PhpArray); ok != true {
			t.Errorf("Third dimension of array was decoded incorrectly: %#v\n", dim2Value)
		} else if value3, ok := dim2Value[php_serialize.PhpValue(0)]; !ok || value3 != "dim3" {
			t.Errorf("Third dimension of array was decoded incorrectly: %#v\n", value3)
		}
	}
}

func TestDecodeMultipleArraysWithoutSemicolons(t *testing.T) {
	decoder := NewPhpDecoder("array1|a:1:{s:5:\"test1\";b:1;}array2|a:1:{s:5:\"test2\";b:1;}")
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode array value %#v \n", err)
	} else {
		if v, ok := result["array1"]; !ok {
			t.Errorf("First array was not decoded \n")
		} else if arrValue, ok := v.(php_serialize.PhpArray); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if value1, ok := arrValue["test1"]; !ok || value1 != true {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		}

		if v, ok := result["array2"]; !ok {
			t.Errorf("Second array was not decoded \n")
		} else if arrValue, ok := v.(php_serialize.PhpArray); ok != true {
			t.Errorf("Array value was decoded incorrectly: %#v \n", v)
		} else if value2, ok := arrValue["test2"]; !ok || value2 != true {
			t.Errorf("Array value was decoded incorrectly: %#v\n", v)
		}
	}
}

func TestDecodeSerializableObjectValueNoFunc(t *testing.T) {
	decoder := NewPhpDecoder("obj|C:10:\"TestObject\":49:{a:3:{s:1:\"a\";i:5;s:1:\"b\";s:4:\"priv\";s:1:\"c\";i:8;}}")
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode object value %#v \n", err)
	} else {
		if v, ok := (result)["obj"]; !ok {
			t.Errorf("Object value was not decoded \n")
		} else if objValue, ok := v.(*php_serialize.PhpObjectSerialized); ok != true {
			t.Errorf("Object value was decoded incorrectly: %#v \n", v)
		} else if objValue.GetClassName() != "TestObject" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", objValue.GetClassName())
		} else if objValue.GetData() != "a:3:{s:1:\"a\";i:5;s:1:\"b\";s:4:\"priv\";s:1:\"c\";i:8;}" {
			t.Errorf("RawData of object was decoded incorrectly: %#v\n", objValue.GetData())
		}
	}
}

func TestDecodeSerializableObjectValue(t *testing.T) {
	decoder := NewPhpDecoder("object|C:10:\"TestObject\":96:{a:1:{s:4:\"item\";O:8:\"AbcClass\":3:{s:1:\"a\";i:5;s:11:\"\x00AbcClass\x00b\";s:7:\"private\";s:4:\"\x00*\x00c\";i:8;}}}")
	decoder.SetSerializedDecodeFunc(php_serialize.SerializedDecodeFunc(php_serialize.UnSerialize))
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode object value %#v \n", err)
	} else {
		if v, ok := (result)["object"]; !ok {
			t.Errorf("Object value was not decoded \n")
		} else if objValue, ok := v.(*php_serialize.PhpObjectSerialized); ok != true {
			t.Errorf("Object value was decoded incorrectly: %#v \n", v)
		} else if objValue.GetClassName() != "TestObject" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", objValue.GetClassName())
		} else if objValue.GetData() != "a:1:{s:4:\"item\";O:8:\"AbcClass\":3:{s:1:\"a\";i:5;s:11:\"\x00AbcClass\x00b\";s:7:\"private\";s:4:\"\x00*\x00c\";i:8;}}" {
			t.Errorf("RawData of object was decoded incorrectly: %#v\n", objValue.GetData())
		} else if vv := objValue.GetValue(); vv == nil {
			t.Errorf("Object value decoded incorrectly, expected value as PhpArray, have got: %v\n", objValue.GetValue())
		} else if arrVal, ok := vv.(php_serialize.PhpArray); !ok {
			t.Errorf("Unable to convert %v to PhpArray\n", vv)
		} else if v1, ok1 := arrVal["item"]; !ok1 {
			t.Errorf("Array value decoded incorrectly, key `item` doest not exists\n")
		} else if itemObjValue, ok1 := v1.(*php_serialize.PhpObject); !ok1  {
			t.Errorf("Unable to convert %v to int\n", v1)
		} else if itemObjValue.GetClassName() != "AbcClass" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", itemObjValue.GetClassName())
		} else if value1, ok := itemObjValue.GetPublic("a"); !ok || value1 != 5 {
			t.Errorf("Public member of object was decoded incorrectly: %#v\n", itemObjValue.GetMembers())
		} else if value2, ok := itemObjValue.GetPrivate("b"); !ok || value2 != "private" {
			t.Errorf("Private member of object was decoded incorrectly: %#v\n", itemObjValue.GetMembers())
		} else if value3, ok := itemObjValue.GetProtected("c"); !ok || value3 != 8 {
			t.Errorf("Protected member of object was decoded incorrectly: %#v\n", itemObjValue.GetMembers())
		}
	}
}

func TestDecodeSerializableObjectFoo(t *testing.T) {
	decoder := NewPhpDecoder("foo|C:3:\"Foo\":3:{foo}")
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode object value %#v \n", err)
	} else {
		if v, ok := (result)["foo"]; !ok {
			t.Errorf("Object value was not decoded \n")
		} else if objValue, ok := v.(*php_serialize.PhpObjectSerialized); ok != true {
			t.Errorf("Object value was decoded incorrectly: %#v \n", v)
		} else if objValue.GetClassName() != "Foo" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", objValue.GetClassName())
		} else if objValue.GetData() != "foo" {
			t.Errorf("RawData of object was decoded incorrectly: %#v\n", objValue.GetData())
		}
	}
}

func TestDecodeSerializableObjectBar(t *testing.T) {
	var f php_serialize.SerializedDecodeFunc
	f = func(s string) (php_serialize.PhpValue, error) {
		var (
			val map[string]string
			err	error
		)
		err = json.Unmarshal([]byte(s), &val)
		return val, err
	}

	decoder := NewPhpDecoder("bar|C:3:\"Bar\":19:{{\"public\":\"public\"}}")
	decoder.SetSerializedDecodeFunc(f)
	if result, err := decoder.Decode(); err != nil {
		t.Errorf("Can not decode object value %#v \n", err)
	} else {
		if v, ok := (result)["bar"]; !ok {
			t.Errorf("Object value was not decoded \n")
		} else if objValue, ok := v.(*php_serialize.PhpObjectSerialized); ok != true {
			t.Errorf("Object value was decoded incorrectly: %#v \n", v)
		} else if objValue.GetClassName() != "Bar" {
			t.Errorf("Object name was decoded incorrectly: %#v\n", objValue.GetClassName())
		} else if objValue.GetData() != "{\"public\":\"public\"}" {
			t.Errorf("RawData of object was decoded incorrectly: %#v\n", objValue.GetData())
		} else if vv := objValue.GetValue(); vv == nil {
			t.Errorf("Object value decoded incorrectly, expected value as PhpArray, have got: %v\n", objValue.GetValue())
		} else if arrVal, ok := vv.(map[string]string); !ok {
			t.Errorf("Unable to convert %v to map[string]string\n", vv)
		} else if v1, ok1 := arrVal["public"]; !ok1 {
			t.Errorf("Array value decoded incorrectly, key `public` doest not exists\n")
		} else if v1 != "public" {
			t.Errorf("Array value decoded incorrectly, expected: %v, have got: %v\n", "public", v1)
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
