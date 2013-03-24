package php_session_decoder

import (
	"testing"
)

func TestPhpObjectFabrica(t *testing.T) {
	obj := NewPhpObject()
	if obj == nil {
		t.Error("Can not create PHP Object instance\n")
	}
}

func TestClassName(t *testing.T) {
	obj := NewPhpObject()
	cName := "SClass"

	obj.SetClassName(cName)
	if obj.GetClassName() != cName {
		t.Error("Class name getter or setter is broken \n")
	}
}

func TestPublicMembers(t *testing.T) {
	obj := NewPhpObject()
	name := "wwwwww"
	value := 34

	obj.SetPublicMemberValue(name, value)
	if v, ok := obj.GetPublicMemberValue(name); v != value || !ok {
		t.Error("Public class members  getter or setter is broken \n")
	}
}

func TestProtectedMembers(t *testing.T) {
	obj := NewPhpObject()
	name := "wwwwww11"
	value := 23232323

	obj.SetProtectedMemberValue(name, value)
	if v, ok := obj.GetProtectedMemberValue(name); v != value || !ok {
		t.Error("Protected class members  getter or setter is broken \n")
	}
}

func TestPrivateMembers(t *testing.T) {
	obj := NewPhpObject()
	name := "wwwwww112"
	value := 232323234

	obj.SetPrivateMemberValue(name, value)
	if v, ok := obj.GetPrivateMemberValue(name); v != value || !ok {
		t.Error("Private class members  getter or setter is broken \n")
	}
}
