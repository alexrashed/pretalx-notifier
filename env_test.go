package main

import (
	"os"
	"testing"
)

func TestGetEnvStr(t *testing.T) {
	err := os.Setenv("TEST_GET_ENV_STR_KEY", "TEST-VALUE")
	if err != nil {
		t.Fail()
	}
	v, err := getEnvStr("TEST_GET_ENV_STR_KEY")
	if v != "TEST-VALUE" || err != nil {
		t.Fail()
	}
	err = os.Unsetenv("TEST_GET_ENV_STR_KEY")
	if err != nil {
		t.Fail()
	}
}

func TestGetEnvStrUnsetFailure(t *testing.T) {
	err := os.Unsetenv("TEST_GET_ENV_STR_KEY")
	if err != nil {
		t.Fail()
	}
	_, err = getEnvStr("TEST_GET_ENV_STR_KEY")
	if err == nil {
		t.Fail()
	}
}

func TestGetEnvInt(t *testing.T) {
	err := os.Setenv("TEST_GET_ENV_INT_KEY", "1337")
	if err != nil {
		t.Fail()
	}
	v, err := getEnvInt("TEST_GET_ENV_INT_KEY")
	if v != 1337 || err != nil {
		t.Fail()
	}
	err = os.Unsetenv("TEST_GET_ENV_INT_KEY")
	if err != nil {
		t.Fail()
	}
}

func TestGetEnvIntUnsetFailure(t *testing.T) {
	err := os.Unsetenv("TEST_GET_ENV_INT_KEY")
	if err != nil {
		t.Fail()
	}
	_, err = getEnvInt("TEST_GET_ENV_INT_KEY")
	if err == nil {
		t.Fail()
	}
}

func TestGetEnvIntInvalidFailure(t *testing.T) {
	err := os.Setenv("TEST_GET_ENV_INT_KEY", "l33t")
	if err != nil {
		t.Fail()
	}
	_, err = getEnvInt("TEST_GET_ENV_INT_KEY")
	if err == nil {
		t.Fail()
	}
}

func TestGetEnvBool(t *testing.T) {
	err := os.Setenv("TEST_GET_ENV_BOOL_KEY", "true")
	if err != nil {
		t.Fail()
	}
	v, err := getEnvBool("TEST_GET_ENV_BOOL_KEY")
	if v != true || err != nil {
		t.Fail()
	}
	err = os.Unsetenv("TEST_GET_ENV_BOOL_KEY")
	if err != nil {
		t.Fail()
	}
}

func TestGetEnvBoolUnsetFailure(t *testing.T) {
	err := os.Unsetenv("TEST_GET_ENV_BOOL_KEY")
	if err != nil {
		t.Fail()
	}
	_, err = getEnvBool("TEST_GET_ENV_BOOL_KEY")
	if err == nil {
		t.Fail()
	}
}

func TestGetEnvBoolInvalidFailure(t *testing.T) {
	err := os.Setenv("TEST_GET_ENV_BOOL_KEY", "1337")
	if err != nil {
		t.Fail()
	}
	_, err = getEnvBool("TEST_GET_ENV_BOOL_KEY")
	if err == nil {
		t.Fail()
	}
}


