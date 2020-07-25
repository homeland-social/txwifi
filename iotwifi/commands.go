package iotwifi

import (
	"os/exec"
	"time"

	"github.com/bhoriuchi/go-bunyan/bunyan"
)

// Command for device network commands.
type Command struct {
	Log      bunyan.Logger
	Runner   CmdRunner
	SetupCfg *SetupCfg
}

// RemoveApInterface removes the AP interface.
func (c *Command) RemoveApInterface() {
	cmd := exec.Command("iw", "dev", "uap0", "del")
	cmd.Start()
	cmd.Wait()
}

// ConfigureApInterface configured the AP interface.
func (c *Command) ConfigureApInterface() {
	cmd := exec.Command("ifconfig", "uap0", c.SetupCfg.HostApdCfg.Ip)
	cmd.Start()
	cmd.Wait()
}

// UpApInterface ups the AP Interface.
func (c *Command) UpApInterface() {
	cmd := exec.Command("ifconfig", "uap0", "up")
	cmd.Start()
	cmd.Wait()
}

// AddApInterface adds the AP interface.
func (c *Command) AddApInterface() {
	cmd := exec.Command("iw", "phy", "phy0", "interface", "add", "uap0", "type", "__ap")
	cmd.Start()
	cmd.Wait()
}

// CheckInterface checks the AP interface.
func (c *Command) CheckApInterface() {
	cmd := exec.Command("ifconfig", "uap0")
	go c.Runner.ProcessCmd("ifconfig_uap0", cmd)
}

// EnableAp enables the AP interface.
func (c *Command) EnableAp() {
	cmd := exec.Command("hostapd_cli", "-i", "uap0", "enable")
	cmd.Start()
	cmd.Wait()
}

// DisableAp disables the AP interface.
func (c *Command) DisableAp() {
	cmd := exec.Command("hostapd_cli", "-i", "uap0", "disable")
	cmd.Start()
	cmd.Wait()
}

// StartWpaSupplicant starts wpa_supplicant.
func (c *Command) StartWpaSupplicant() {

	args := []string{
		"-Dnl80211",
		"-iwlan0",
		"-c" + c.SetupCfg.WpaSupplicantCfg.CfgFile,
	}

	cmd := exec.Command("wpa_supplicant", args...)
	go c.Runner.ProcessCmd("wpa_supplicant", cmd)
}

// StartDnsmasq starts dnsmasq.
func (c *Command) StartDnsmasq() {
	// hostapd is enabled, fire up dnsmasq
	args := []string{
		"--no-hosts", // Don't read the hostnames in /etc/hosts.
		"--keep-in-foreground",
		"--log-queries",
		"--no-resolv",
		"--address=" + c.SetupCfg.DnsmasqCfg.Address,
		"--dhcp-range=" + c.SetupCfg.DnsmasqCfg.DhcpRange,
		"--dhcp-vendorclass=" + c.SetupCfg.DnsmasqCfg.VendorClass,
		"--dhcp-authoritative",
		"--log-facility=-",
	}

	cmd := exec.Command("dnsmasq", args...)
	go c.Runner.ProcessCmd("dnsmasq", cmd)
}

// StartHostapd starts hostapd.
func (c *Command) StartHostapd(ssid string, psk string, channel string) {
	args := []string{
		"/dev/stdin",
	}
	cmd := exec.Command("hostapd", args...)

	cfg := `interface=uap0
ssid=` + ssid + `
hw_mode=g
channel=` + channel + `
ctrl_interface=/var/run/hostapd
ctrl_interface_group=0
macaddr_acl=0
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_passphrase=` + psk + `
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP`

	c.Log.Info("Hostapd CFG: %s", cfg)

	// handle in pipe here to pass cfg, out/error handled by Runner
	hostapdPipe, _ := cmd.StdinPipe()
	hostapdPipe.Write([]byte(cfg))

	go c.Runner.ProcessCmd("hostapd", cmd)

	time.Sleep(2)
	hostapdPipe.Close()
}
