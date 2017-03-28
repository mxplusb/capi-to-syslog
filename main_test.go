package main

import (
	"os"
	"testing"
)

func TestSetEnvVars(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.FailNow()
		}
	}()

	if err := os.Setenv("CAPI_CLIENT_ID", "idontcare"); err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := os.Setenv("CAPI_CLIENT_SECRET", "whatisfordinner"); err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := os.Setenv("CAPI_SYSTEM_URI", "breakfast.is.awesome"); err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := os.Setenv("CAPI_EVENTS", "type:audit.app.ssh-authorized,type:audit.app.ssh-unauthorized,type:audit.service_key.create,type:audit.service_key.delete,type:audit.space.create,type:audit.app.create"); err != nil {
		t.Log(err)
		t.Fail()
	}

	SetEnvVars()

	if CapiClientID != "idontcare" {
		t.Fail()
	}

	if CapiClientSecret != "whatisfordinner" {
		t.Fail()
	}

	if CapiSystemURI != "breakfast.is.awesome" {
		t.Fail()
	}

	if len(AuditableEvents) == 0 {
		t.Log("no events parsed!")
		t.Fail()
	}

	if err := os.Setenv("CAPI_CLIENT_ID", "uaa-to-syslog"); err != nil {
		t.Fail()
	}

	if err := os.Setenv("CAPI_CLIENT_SECRET", "pivotal123!"); err != nil {
		t.Fail()
	}

	if err := os.Setenv("CAPI_SYSTEM_URI", "run.haas-88.pez.pivotal.io"); err != nil {
		t.Fail()
	}

	if err := os.Setenv("CAPI_EVENTS", "type:audit.app.ssh-authorized,type:audit.app.ssh-unauthorized,type:audit.service_key.create,type:audit.service_key.delete,type:audit.space.create,type:audit.app.create"); err != nil {
		t.Fail()
	}
}
