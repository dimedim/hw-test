package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type UserRole string

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
	testCases := []struct {
		in          any
		expectedErr error
	}{
		{
			in: User{
				ID:     strings.Repeat("x", 36),
				Name:   "Name",
				Age:    30,
				Email:  "email@mail.com",
				Role:   "admin",
				Phones: []string{"01234567890"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "short",
				Name:   "",
				Age:    17,
				Email:  "bad",
				Role:   "unknown",
				Phones: []string{"123"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrLen},
				{Field: "Age", Err: ErrMin},
				{Field: "Email", Err: ErrRegexp},
				{Field: "Role", Err: ErrInt},
				{Field: "Phones[0]", Err: ErrLen},
			},
		},
		{
			in:          App{Version: "1234"},
			expectedErr: ValidationErrors{{Field: "Version", Err: ErrLen}},
		},
		{
			in:          Response{Code: 300, Body: "ok"},
			expectedErr: ValidationErrors{{Field: "Code", Err: ErrInt}},
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in:          "just a string",
			expectedErr: errors.New("expected a struct"),
		},
		{
			in: struct {
				X string `validate:"bad"`
			}{X: "a"},
			expectedErr: errors.New("unknown validation key: bad"),
		},
		{
			in: struct {
				E string `validate:"regexp:([)"`
			}{E: "anything"},
			expectedErr: errors.New("error parsing regexp"),
		},
	}

	for i, tC := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			err := Validate(tC.in)

			if tC.expectedErr == nil {
				if err != nil {
					t.Fatalf("case %d: expected no error, got %v", i, err)
				}
				return
			}

			if err == nil {
				t.Fatalf("case %d: expected error %v, got nil", i, tC.expectedErr)
			}

			if wantVEs, ok := tC.expectedErr.(ValidationErrors); ok {
				var gotVEs ValidationErrors
				if !errors.As(err, &gotVEs) {
					t.Fatalf("case %d: expected ValidationErrors, got %T %v", i, err, err)
				}
				if len(gotVEs) != len(wantVEs) {
					t.Errorf("case %d: expected %d validation errors, got %d", i, len(wantVEs), len(gotVEs))
				}
				for _, exp := range wantVEs {
					found := false
					for _, got := range gotVEs {
						if got.Field == exp.Field && errors.Is(got.Err, exp.Err) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("case %d: missing expected error %q on field %q", i, exp.Err, exp.Field)
					}
				}
				return
			}
			if !strings.Contains(err.Error(), tC.expectedErr.Error()) {
				t.Errorf("case %d: expected error containing %q, got %q", i, tC.expectedErr.Error(), err.Error())
			}
		})
	}
}
