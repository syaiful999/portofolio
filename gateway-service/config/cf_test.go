package config

import (
	"os"
	"testing"
)

func TestConfigInit(t *testing.T) {
	os.Setenv("SERVICE_NAME", "test-gateway")
	os.Setenv("SERVICE_VERSION", "1.0.0")
	os.Setenv("SERVICE_PORT", "8080")
	os.Setenv("SERVICE_JWT_SECRET", "secret")
	os.Setenv("MASTER_DATA_SERVICE_HOST", "http://localhost:8096")
	os.Setenv("CACHE_DRIVER", "default")

	defer func() {
		os.Unsetenv("SERVICE_NAME")
		os.Unsetenv("SERVICE_VERSION")
		os.Unsetenv("SERVICE_PORT")
		os.Unsetenv("SERVICE_JWT_SECRET")
		os.Unsetenv("MASTER_DATA_SERVICE_HOST")
		os.Unsetenv("CACHE_DRIVER")
	}()

	conf := Config{}
	conf.Init()

	if conf.Service.Name != "test-gateway" {
		t.Errorf("expected 'test-gateway', got '%s'", conf.Service.Name)
	}
	if conf.Service.Version != "1.0.0" {
		t.Errorf("expected '1.0.0', got '%s'", conf.Service.Version)
	}
	if conf.Service.Port != 8080 {
		t.Errorf("expected 8080, got %d", conf.Service.Port)
	}
	if conf.Services.MasterDataURL != "http://localhost:8096" {
		t.Errorf("expected 'http://localhost:8096', got '%s'", conf.Services.MasterDataURL)
	}
	if conf.Hosts.Cache.Driver != "default" {
		t.Errorf("expected 'default', got '%s'", conf.Hosts.Cache.Driver)
	}
}

func TestConfigInit_EmptyEnv(t *testing.T) {
	conf := Config{}
	conf.Init()

	if conf.Service.Name != "" {
		t.Errorf("expected empty name, got '%s'", conf.Service.Name)
	}
	if conf.Service.Port != 0 {
		t.Errorf("expected 0, got %d", conf.Service.Port)
	}
}
