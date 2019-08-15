package liguetaxi

import (
	"testing"
)

func TestUserStatusUnmarshalJSON(t *testing.T) {
	testCases := []struct{
		b []byte
		want UserStatus
	}{
		{[]byte(`24`), UserStatusActive},
		{[]byte(`25`), UserStatusInactive},
	}

	for _, tc := range testCases {
		var status userStatus

		status.UnmarshalJSON(tc.b)

		if status != tc.want {
			t.Errorf("got userStatus.UnmarshalJSON(%s): %v; want %v.", tc.b, status, tc.want)
		}
	}
}

func TestEmptyObjToStrUnmarshalJSON(t *testing.T) {
	testCases := []struct{
		b []byte
		want string
	}{
		{[]byte(`{}`), ""},
	}

	for _, tc := range testCases {
		var str emptyObjToStr

		str.UnmarshalJSON(tc.b)

		if string(str) != tc.want {
			t.Errorf("got emptyObjToStr.UnmarshalJSON(%s): %v; want %s.", tc.b, str, tc.want)
		}
	}
}
