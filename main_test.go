package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
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
