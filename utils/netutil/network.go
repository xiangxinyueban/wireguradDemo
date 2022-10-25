package netutil

import (
	"net"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogNetworkStats logs network information to the Trace log level.
var LogNetworkStats = defaultLogNetworkStats

// AddDefaultRoute adds default VPN tunnel route.
func AddDefaultRoute(iface string) error {
	return addDefaultRoute(iface)
}

func IptablesInit(iface string, listenPort int, subnet net.IPNet) error {
	return iptablesinit(iface, listenPort, subnet)
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
