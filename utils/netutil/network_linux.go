package netutil

import (
	"net"
	"os/exec"
	"strings"
	"vpn/utils/cmdutil"
	"vpn/utils/iputil"

	"github.com/rs/zerolog/log"
)

func assignIP(iface string, subnet net.IPNet) error {
	if err := cmdutil.SudoExec("ip", "address", "replace", "dev", iface, subnet.String()); err != nil {
		return err
	}
	return cmdutil.SudoExec("ip", "link", "set", "dev", iface, "up")
}

func excludeRoute(ip, gw net.IP) error {
	return cmdutil.SudoExec("ip", "route", "add", ip.String(), "via", gw.String())
}

func deleteRoute(ip, gw string) error {
	return cmdutil.SudoExec("ip", "route", "delete", ip, "via", gw)
}

func addDefaultRoute(iface string) error {
	if err := cmdutil.SudoExec("ip", "route", "add", "0.0.0.0/1", "dev", iface); err != nil {
		return err
	}

	if err := cmdutil.SudoExec("ip", "route", "add", "128.0.0.0/1", "dev", iface); err != nil {
		return err
	}

	if ipv6Enabled() {
		if err := cmdutil.SudoExec("ip", "-6", "route", "add", "::/1", "dev", iface); err != nil {
			return err
		}

		if err := cmdutil.SudoExec("ip", "-6", "route", "add", "8000::/1", "dev", iface); err != nil {
			return err
		}
	}
	return nil
}

func iptablesinit(iface string, port int, subnet net.IPNet) error {
	outface := iputil.GetOutboundInterface()
	if err := cmdutil.SudoExec("iptables", "-t", "net", "-I", "POSTROUTING", "1", "-s", subnet.String(), "-o", outface, "-j", "MASQUERADE"); err != nil {
		return err
	}

	if err := cmdutil.SudoExec("iptables", "-I", "INPUT", "1", "-i", iface, "-j", "ACCEPT"); err != nil {
		return err
	}

	if err := cmdutil.SudoExec("iptables", "-I", "FORWARD", "1", "-i", iface, "-o", outface, "-j", "ACCEPT"); err != nil {
		return err
	}

	if err := cmdutil.SudoExec("iptables", "-I", "FORWARD", "1", "-i", outface, "-o", iface, "-j", "ACCEPT"); err != nil {
		return err
	}

	if err := cmdutil.SudoExec("iptables", "-I", "INPUT", "1", "-i", outface, "-p", "udp", "--dport", port, "-j", "ACCEPT"); err != nil {
		return err
	}
	return nil
}

func logNetworkStats() {
	for _, args := range [][]string{{"iptables", "-L", "-n"}, {"iptables", "-L", "-n", "-t", "nat"}, {"ip", "route", "list"}, {"ip", "address", "list"}} {
		out, err := exec.Command("sudo", args...).CombinedOutput()
		logOutputToTrace(out, err, args...)
	}
}

func ipv6Enabled() bool {
	out, err := cmdutil.ExecOutput("sysctl", "net.ipv6.conf.all.disable_ipv6")
	if err != nil {
		log.Error().Err(err).Msg("Failed to detect if IPv6 disabled on the host")

		return true
	}

	return strings.Contains(string(out), "net.ipv6.conf.all.disable_ipv6 = 0")
}
