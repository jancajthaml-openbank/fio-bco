package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetConfig(t *testing.T) {
	for _, v := range os.Environ() {
		k := strings.Split(v, "=")[0]
		if strings.HasPrefix(k, "FIO_BCO") {
			os.Unsetenv(k)
		}
	}

	t.Log("has defaults for all values")
	{
		config := LoadConfig()

		if config.Tenant != "" {
			t.Errorf("Tenant default value is not empty")
		}
		if config.SyncRate != 22*time.Second {
			t.Errorf("SyncRate default value is not 22s")
		}
		if config.FioGateway != "https://www.fio.cz/ib_api/rest" {
			t.Errorf("FioGateway default value is not https://www.fio.cz/ib_api/rest")
		}
		if config.LedgerGateway != "https://127.0.0.1:4401" {
			t.Errorf("LedgerGateway default value is not https://127.0.0.1:4401")
		}
		if config.VaultGateway != "https://127.0.0.1:4400" {
			t.Errorf("VaultGateway default value is not https://127.0.0.1:4400")
		}
		if config.RootStorage != "/data/t_/import/fio" {
			t.Errorf("RootStorage default value is not /data/t_/import/fio")
		}
		if config.EncryptionKey != nil {
			t.Errorf("EncryptionKey default value is not empty")
		}
		if config.LakeHostname != "127.0.0.1" {
			t.Errorf("LakeHostname default value is not 127.0.0.1")
		}
		if config.LogLevel != "INFO" {
			t.Errorf("LogLevel default value is not INFO")
		}
		if config.MetricsContinuous != true {
			t.Errorf("MetricsContinuous default value is not true")
		}
		if config.MetricsRefreshRate != time.Second {
			t.Errorf("MetricsRefreshRate default value is not 1s")
		}
		if config.MetricsOutput != "/tmp/fio-bco-import-metrics" {
			t.Errorf("MetricsOutput default value is not /tmp/fio-bco-import-metrics")
		}
	}
}
