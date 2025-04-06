package main

import "testing"

func TestRunCmd(t *testing.T) {
	t.Run("Empty Cmd", func(t *testing.T) {
		env := Environment{}
		code := RunCmd([]string{}, env)
		if code != 1 {
			t.Fatalf("expected exit code 1, got %d", code)
		}
	})

	t.Run("With Environment", func(t *testing.T) {
		env := Environment{
			"VAR1": {
				Value:      "123",
				NeedRemove: false,
			},
		}

		cmd := []string{"sh", "-c", `test "$VAR1" = "123" && exit 42 || exit 1`}
		code := RunCmd(cmd, env)
		if code != 42 {
			t.Fatalf("expected exit code 42, got %d", code)
		}
	})

	t.Run("Need Remove Environment", func(t *testing.T) {
		env := Environment{
			"REMOVEME": {
				Value:      "some_value",
				NeedRemove: true,
			},
		}
		cmd := []string{"sh", "-c", `test "$REMOVEME" = "" && exit 43 || exit 1`}
		code := RunCmd(cmd, env)
		if code != 43 {
			t.Fatalf("expected exit code 43, got %d", code)
		}
	})
}
