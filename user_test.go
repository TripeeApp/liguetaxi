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
	}

	for _, tc := range testCases {
		var status UserStatus

		status.UnmarshalJSON(tc.b)

		if status != tc.want {
			t.Errorf("got UserStatus.UnmarshalJSON(%s): %v; want %v.", tc.b, status, tc.want)
		}
	}
}
