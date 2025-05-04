package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			// Place your code here.
		},
		// ...
		// Place your code here.
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.
			_ = tt
		})
	}
}

func GetValidCheckList(len, min, max int, reg string, in []string) CheckList {
	return CheckList{
		Len:    len,
		Regexp: reg,
		In:     in,
		Min:    min,
		Max:    max,
	}
}
func TestNewValidators(t *testing.T) {
	tests := []struct {
		validate          string
		expectedCheckList CheckList
	}{
		{
			validate:          "in:200,404,500",
			expectedCheckList: GetValidCheckList(0, 0, 0, "", []string{"200", "404", "500"}),
		},
		{
			validate:          "in:admin,stuff",
			expectedCheckList: GetValidCheckList(0, 0, 0, "", []string{"admin", "stuff"}),
		},
		{
			validate:          "min:18|max:50",
			expectedCheckList: GetValidCheckList(0, 18, 50, "", nil),
		},
		{
			validate:          "regexp:^\\w+@\\w+\\.\\w+$",
			expectedCheckList: GetValidCheckList(0, 0, 0, "^\\w+@\\w+\\.\\w+$", nil),
		},
		{
			validate:          "len:50",
			expectedCheckList: GetValidCheckList(50, 0, 0, "", nil),
		},
		{
			validate:          "regexp:awdawdaw:awdwad:awdawd:123",
			expectedCheckList: GetValidCheckList(0, 0, 0, "awdawdaw:awdwad:awdawd:123", nil),
		},
	}

	for _, tc := range tests {
		vc, err := GetCheckListFromStructTag(tc.validate)
		require.NoError(t, err)
		require.Equal(t, tc.expectedCheckList, *vc)
	}
}
