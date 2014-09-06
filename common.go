package php_session_decoder

const (
	SEPARATOR_VALUE_NAME	= '|'
	SEPARATOR_VALUE_TYPE	= ':'
	SEPARATOR_VALUES		= ';'

	DELIMITER_STRING_LEFT	= '"'
	DELIMITER_STRING_RIGHT	= '"'
	DELIMITER_OBJECT_LEFT	= '{'
	DELIMITER_OBJECT_RIGHT	= '}'

	FORMATTER_FLOAT			= 'f'
)

type PhpValue interface{}

type PhpSessionData map[string]PhpValue

type PhpObject struct {
	RawData		string
	members		PhpSessionData
	className	string
	custom		bool
}

func NewPhpObject() *PhpObject {
	membersMap := make(PhpSessionData)
	d := &PhpObject{
		members: membersMap,
	}
	return d
}

func (obj *PhpObject) Custom(value bool) {
	obj.custom = value
}

func (obj *PhpObject) IsCustom() bool {
	return obj.custom
}

func (obj *PhpObject) GetClassName() string {
	return obj.className
}

func (obj *PhpObject) SetClassName(cName string) {
	obj.className = cName
}

func (obj *PhpObject) GetMembers() PhpSessionData {
	return obj.members
}

func (obj *PhpObject) GetPrivateMemberValue(memberName string) (v PhpValue, ok bool) {
	v, ok = obj.members["\x00" + obj.className + "\x00" + memberName]
	return
}

func (obj *PhpObject) SetPrivateMemberValue(memberName string, value PhpValue) {
	obj.members["\x00" + obj.className + "\x00" + memberName] = value
}

func (obj *PhpObject) GetProtectedMemberValue(memberName string) (v PhpValue, ok bool) {
	v, ok = obj.members["\x00*\x00" + memberName]
	return
}

func (obj *PhpObject) SetProtectedMemberValue(memberName string, value PhpValue) {
	obj.members["\x00*\x00" + memberName] = value
}

func (obj *PhpObject) GetPublicMemberValue(memberName string) (v PhpValue, ok bool) {
	v, ok = obj.members[memberName]
	return
}

func (obj *PhpObject) SetPublicMemberValue(memberName string, value PhpValue) {
	obj.members[memberName] = value
}
