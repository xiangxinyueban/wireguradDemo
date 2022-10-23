package iputil

import (
	"github.com/rs/zerolog"
	"net"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/mysteriumnetwork/node/utils/cmdutil"
)

var LogNetworkStats = defaultLogNetworkStats

// AddDefaultRoute adds default VPN tunnel route.
func AddDefaultRoute(iface string) error {
	return addDefaultRoute(iface)
}

// AssignIP assigns subnet to given interface.
func AssignIP(iface string, subnet net.IPNet) error {
	return assignIP(iface, subnet)
}

func defaultLogNetworkStats() {
	if log.Logger.GetLevel() != zerolog.TraceLevel {
		return
	}

	logNetworkStats()
}

func logOutputToTrace(out []byte, err error, args ...string) {
	logSkipFrame := log.With().CallerWithSkipFrameCount(3).Logger()

	if err != nil {
		(&logSkipFrame).Trace().Msgf("Failed to get %s error: %v", strings.Join(args, " "), err)
	} else {
		(&logSkipFrame).Trace().Msgf("%q output:\n%s", strings.Join(args, " "), out)
	}
}

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
