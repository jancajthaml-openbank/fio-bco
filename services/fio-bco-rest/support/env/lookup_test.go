package env

import (
	"io/ioutil"
	"os"
	"testing"
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

func TestEnvInteger(t *testing.T) {

	t.Log("TEST_INT missing")
	{
		if Int("TEST_INT", 0) != 0 {
			t.Errorf("Int did not provide default value")
		}
	}

	t.Log("TEST_INT present and valid")
	{
		os.Setenv("TEST_INT", "1")
		defer os.Unsetenv("TEST_INT")

		if Int("TEST_INT", 0) != 1 {
			t.Errorf("Int did not obtain env value")
		}
	}

	t.Log("TEST_INT present and invalid")
	{
		os.Setenv("TEST_INT", "x")
		defer os.Unsetenv("TEST_INT")

		if Int("TEST_INT", 0) != 0 {
			t.Errorf("Int did not fallback to default value")
		}
	}
}

func TestEnvUnsignedInteger(t *testing.T) {

	t.Log("TEST_UINT missing")
	{
		if Uint64("TEST_UINT", 0) != 0 {
			t.Errorf("Uint64 did not provide default value")
		}
	}

	t.Log("TEST_UINT present and valid")
	{
		os.Setenv("TEST_UINT", "1")
		defer os.Unsetenv("TEST_UINT")

		if Uint64("TEST_UINT", 2) != 1 {
			t.Errorf("Int did not obtain env value")
		}
	}

	t.Log("TEST_UINT present and invalid")
	{
		os.Setenv("TEST_UINT", "x")
		defer os.Unsetenv("TEST_UINT")

		if Uint64("TEST_UINT", 0) != 0 {
			t.Errorf("Int did not fallback to default value")
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
