package env

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestEnvString(t *testing.T) {

	t.Log("TEST_STR missing")
	{
		if String("TEST_STR", "x") != "x" {
			t.Errorf("String did not provide default value")
		}
	}

	t.Log("TEST_STR present")
	{
		os.Setenv("TEST_STR", "y")
		defer os.Unsetenv("TEST_STR")

		if String("TEST_STR", "x") != "y" {
			t.Errorf("String did not obtain env value")
		}
	}
}

func TestEnvDuration(t *testing.T) {

	t.Log("TEST_DUR missing")
	{
		if Duration("TEST_DUR", time.Second) != time.Second {
			t.Errorf("Duration did not provide default value")
		}
	}

	t.Log("TEST_DUR present and valid")
	{
		os.Setenv("TEST_DUR", "2s")
		defer os.Unsetenv("TEST_DUR")

		if Duration("TEST_DUR", time.Second) != 2*time.Second {
			t.Errorf("Duration did not obtain env value")
		}
	}

	t.Log("TEST_DUR present and invalid")
	{
		os.Setenv("TEST_DUR", "x")
		defer os.Unsetenv("TEST_DUR")

		if Duration("TEST_DUR", time.Second) != time.Second {
			t.Errorf("Duration did not obtain fallback to default value")
		}
	}
}

func TestEnvHexadecimalFile(t *testing.T) {

	t.Log("TEST_HEX missing")
	{
		if string(HexFile("TEST_HEX", []byte("x"))) != "x" {
			t.Errorf("HexFile did not provide default value")
		}
	}

	t.Log("TEST_HEX present and valid")
	{
		file, err := ioutil.TempFile(os.TempDir(), "hexfile")
		if err != nil {
			t.Fatalf(err.Error())
		}
		defer os.Remove(file.Name())
		_, err = file.WriteString("61626364")
		if err != nil {
			t.Fatalf(err.Error())
		}

		os.Setenv("TEST_HEX", file.Name())
		defer os.Unsetenv("TEST_HEX")

		if string(HexFile("TEST_HEX", []byte("x"))) != "abcd" {
			t.Errorf("HexFile did not obtain env value")
		}
	}

	t.Log("TEST_HEX present and invalid")
	{
		os.Setenv("TEST_HEX", "/dev/null")
		defer os.Unsetenv("TEST_HEX")

		if string(HexFile("TEST_HEX", []byte("x"))) != "x" {
			t.Errorf("HexFile did not obtain fallback to default value")
		}
	}
}
