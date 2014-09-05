package php_serialize

const (
	TOKEN_NULL				= 'N'
	TOKEN_BOOL				= 'b'
	TOKEN_INT				= 'i'
	TOKEN_FLOAT				= 'd'
	TOKEN_STRING			= 's'
	TOKEN_ARRAY				= 'a'
	TOKEN_OBJECT			= 'O'
	TOKEN_OBJECT_SERIALIZED	= 'C'

	SEPARATOR_VALUE_NAME	= '|'
	SEPARATOR_VALUE_TYPE	= ':'
	SEPARATOR_VALUES		= ';'

	DELIMITER_STRING_LEFT	= '"'
	DELIMITER_STRING_RIGHT	= '"'
	DELIMITER_OBJECT_LEFT	= '{'
	DELIMITER_OBJECT_RIGHT	= '}'

	FORMATTER_FLOAT			= 'f'
)

type SerializedDecodeFunc func(string) (PhpValue, error)

type SerializedEncodeFunc func(PhpValue) (string, error)

type PhpValue interface{}

type PhpArray map[string]PhpValue

type PhpObject struct {
	className	string
	members		PhpArray
}

func (self *PhpObject) GetClassName() string {
	return self.className
}

func (self *PhpObject) SetClassName(name string) *PhpObject {
	self.className = name
	return self
}

func (self *PhpObject) GetMembers() PhpArray {
	return self.members
}

func (self *PhpObject) SetMembers(members PhpArray) *PhpObject {
	self.members = members
	return self
}

func (self *PhpObject) GetPrivate(name string) (v PhpValue, ok bool) {
	v, ok = self.members["\x00" + self.className + "\x00" + name]
	return
}

func (self *PhpObject) SetPrivate(name string, value PhpValue) *PhpObject {
	self.members["\x00" + self.className + "\x00" + name] = value
	return self
}

func (self *PhpObject) GetProtected(name string) (v PhpValue, ok bool) {
	v, ok = self.members["\x00*\x00" + name]
	return
}

func (self *PhpObject) SetProtected(name string, value PhpValue) *PhpObject {
	self.members["\x00*\x00" + name] = value
	return self
}

func (self *PhpObject) GetPublic(name string) (v PhpValue, ok bool) {
	v, ok = self.members[name]
	return
}

func (self *PhpObject) SetPublic(name string, value PhpValue) *PhpObject {
	self.members[name] = value
	return self
}

type PhpObjectSerialized struct {
	className	string
	data		string
	value		PhpValue
}

func (self *PhpObjectSerialized) GetClassName() string {
	return self.className
}

func (self *PhpObjectSerialized) SetClassName(name string) *PhpObjectSerialized {
	self.className = name
	return self
}

func (self *PhpObjectSerialized) GetData() string {
	return self.data
}

func (self *PhpObjectSerialized) SetData(data string) *PhpObjectSerialized {
	self.data = data
	return self
}

func (self *PhpObjectSerialized) GetValue() PhpValue {
	return self.value
}

func (self *PhpObjectSerialized) SetValue(value PhpValue) *PhpObjectSerialized {
	self.value = value
	return self
}
