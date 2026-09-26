package env_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/andreimerlescu/goenv/env"
)

func TestParse(t *testing.T) {
	input := strings.Join([]string{
		"# a comment",
		"",
		"PLAIN=value",
		"export EXPORTED=yes",
		`DOUBLE="hello world\nnext"`,
		`SINGLE='literal \n # not a comment'`,
		"INLINE=value # comment",
		"HASH=abc#def",
		"EMPTY=",
		"SPACED = padded ",
		"URL=postgres://u:p@host/db?sslmode=disable",
	}, "\n")
	got, err := env.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	want := map[string]string{
		"PLAIN":    "value",
		"EXPORTED": "yes",
		"DOUBLE":   "hello world\nnext",
		"SINGLE":   `literal \n # not a comment`,
		"INLINE":   "value",
		"HASH":     "abc#def",
		"EMPTY":    "",
		"SPACED":   "padded",
		"URL":      "postgres://u:p@host/db?sslmode=disable",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %#v, want %#v", got, want)
	}

	if _, err := env.Parse(strings.NewReader("NOT_A_PAIR")); err == nil {
		t.Error("Parse() expected error for line without '='")
	}
}

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("GOENV_TEST_NEW=from_file\nGOENV_TEST_EXISTING=from_file\n"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Unsetenv("GOENV_TEST_NEW") })

	withEnv(t, "GOENV_TEST_EXISTING", "from_env", func() {
		if err := env.LoadFile(path); err != nil {
			t.Fatalf("LoadFile() returned error: %v", err)
		}
		if got := env.String("GOENV_TEST_NEW", ""); got != "from_file" {
			t.Errorf("GOENV_TEST_NEW = %q, want %q", got, "from_file")
		}
		if got := env.String("GOENV_TEST_EXISTING", ""); got != "from_env" {
			t.Errorf("LoadFile() overwrote GOENV_TEST_EXISTING = %q, want %q", got, "from_env")
		}

		if err := env.OverloadFile(path); err != nil {
			t.Fatalf("OverloadFile() returned error: %v", err)
		}
		if got := env.String("GOENV_TEST_EXISTING", ""); got != "from_file" {
			t.Errorf("OverloadFile() GOENV_TEST_EXISTING = %q, want %q", got, "from_file")
		}
	})

	if err := env.LoadFile(filepath.Join(dir, "missing.env")); err != nil {
		t.Errorf("LoadFile() on missing file returned error: %v", err)
	}
}

func TestLoadFileDefaultPath(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("GOENV_TEST_DEFAULT=cwd\n"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	t.Cleanup(func() { _ = os.Unsetenv("GOENV_TEST_DEFAULT") })

	if env.EnvFile != env.DefaultEnvFile {
		t.Fatalf("EnvFile = %q, want %q", env.EnvFile, env.DefaultEnvFile)
	}
	if err := env.LoadFile(""); err != nil {
		t.Fatalf("LoadFile(\"\") returned error: %v", err)
	}
	if got := env.String("GOENV_TEST_DEFAULT", ""); got != "cwd" {
		t.Errorf("GOENV_TEST_DEFAULT = %q, want %q", got, "cwd")
	}
}
