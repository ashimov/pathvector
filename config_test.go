package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configFile := `
asn: 34553
router-id: 192.0.2.1
prefixes:
  - 192.0.2.0/24
  - 2001:db8::/48
augments:
  statics:
    "203.0.113.0/24" : "192.0.2.10"
    "2001:db8:2::/64" : "2001:db8::1"
vrrp:
  - state: primary
    interface: eth0
    priority: 255
    vips:
      - 192.0.2.1/24
      - 2001:db8::1/64
  - state: backup
    interface: eth1
    priority: 255
    vips:
      - 192.0.2.2/24
      - 2001:db8::2/64
peers:
  Example:
    asn: 65530
    neighbors:
      - 203.0.113.25
      - 2001:db8:2::25
`

	globalConfig, err := loadConfig([]byte(configFile))
	if err != nil {
		t.Error(err)
	}

	if globalConfig.Asn != 34553 {
		t.Errorf("expected asn 34553 got %d", globalConfig.Asn)
	}
	if globalConfig.RouterId != "192.0.2.1" {
		t.Errorf("expected router-id 192.0.2.1 got %s", globalConfig.RouterId)
	}
	if len(globalConfig.Peers) != 1 {
		t.Errorf("expected 1 peer, got %d", len(globalConfig.Peers))
	}
	if globalConfig.Peers["Example"].Asn != 65530 {
		t.Errorf("expected peer asn 34553 got %d", globalConfig.Peers["Example"].Asn)
	}
	if !reflect.DeepEqual(globalConfig.Peers["Example"].NeighborIPs, []string{"203.0.113.25", "2001:db8:2::25"}) {
		t.Errorf("expected neighbor ips [203.0.113.25 2001:db8:2::25] got %v", globalConfig.Peers["Example"].NeighborIPs)
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	configFile := "INVALID YAML"
	_, err := loadConfig([]byte(configFile))
	if err == nil || !strings.Contains(err.Error(), "yaml unmarshal") {
		t.Errorf("expected yaml unmarshal error, got %+v", err)
	}
}

func TestLoadConfigValidationError(t *testing.T) {
	configFile := "router-id: foo"
	_, err := loadConfig([]byte(configFile))
	if err == nil || !strings.Contains(err.Error(), "validation") {
		t.Errorf("expected validation error, got %+v", err)
	}
}

func TestLoadConfigInvalidOriginPrefix(t *testing.T) {
	configFile := `
asn: 34553
router-id: 192.0.2.1
prefixes:
  - foo/24
  - 2001:db8::/48`
	_, err := loadConfig([]byte(configFile))
	if err == nil || !strings.Contains(err.Error(), "invalid origin prefix") {
		t.Errorf("expected invalid origin prefix error, got %+v", err)
	}
}

func TestLoadConfigInvalidVRRPState(t *testing.T) {
	configFile := `
asn: 34553
router-id: 192.0.2.1
vrrp:
  - state: invalid
    interface: eth1
    priority: 255
    vips:
      - 192.0.2.2/24
      - 2001:db8::2/64`
	_, err := loadConfig([]byte(configFile))
	if err == nil || !strings.Contains(err.Error(), "VRRP state must be") {
		t.Errorf("expected VRRP state error, got %+v", err)
	}
}

func TestLoadConfigInvalidStaticPrefix(t *testing.T) {
	configFile := `
asn: 34553
router-id: 192.0.2.1
augments:
  statics:
    "foo/24" : "192.0.2.10"
    "2001:db8:2::/64" : "2001:db8::1"
`
	_, err := loadConfig([]byte(configFile))
	if err == nil || !strings.Contains(err.Error(), "invalid static prefix") {
		t.Errorf("expected invalid static prefix error, got %+v", err)
	}
}

func TestLoadConfigInvalidVIP(t *testing.T) {
	configFile := `
asn: 34553
router-id: 192.0.2.1
vrrp:
  - state: invalid
    interface: eth1
    priority: 255
    vips:
      - foo/24
      - 2001:db8::2/64`

	_, err := loadConfig([]byte(configFile))
	if err == nil || !strings.Contains(err.Error(), "invalid VIP") {
		t.Errorf("expected invalid VIP error, got %+v", err)
	}
}

func TestDocumentCLIFlags(t *testing.T) {
	documentCliFlags()
}

func TestDocumentConfig(t *testing.T) {
	documentConfig()
}
