package php_session_decoder

import (
	"strings"
)

const VALUE_NAME_SEPARATOR = '|'
const TYPE_VALUE_SEPARATOR = ':'
const VALUES_SEPARATOR = ';'

type PhpValue interface{}

type PhpSessionData map[string]PhpValue

type PhpObject struct {
	members   PhpSessionData
	className string
}

func NewPhpObject() *PhpObject {
	membersMap := make(PhpSessionData)
	d := &PhpObject{
		members: membersMap,
	}
	return d
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

func (obj *PhpObject) GetPrivateMemberValue(memberName string) (PhpValue, bool) {
	keyParts := [...]string{"\x00", obj.className, "\x00", memberName}
	key := strings.Join(keyParts[:], "")
	v, ok := obj.members[key]
	return v, ok
}

func (obj *PhpObject) SetPrivateMemberValue(memberName string, value PhpValue) {
	keyParts := [...]string{"\x00", obj.className, "\x00", memberName}
	key := strings.Join(keyParts[:], "")
	obj.members[key] = value
}

func (obj *PhpObject) GetProtectedMemberValue(memberName string) (PhpValue, bool) {
	key := "\x00*\x00" + memberName
	v, ok := obj.members[key]
	return v, ok
}

func (obj *PhpObject) SetProtectedMemberValue(memberName string, value PhpValue) {
	key := "\x00*\x00" + memberName
	obj.members[key] = value
}

func (obj *PhpObject) GetPublicMemberValue(memberName string) (PhpValue, bool) {
	v, ok := obj.members[memberName]
	return v, ok
}

func (obj *PhpObject) SetPublicMemberValue(memberName string, value PhpValue) {
	obj.members[memberName] = value
}
