package wireguard

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
	"vpn/utils/cmdutil"
	"vpn/utils/iputil"

	"github.com/pkg/errors"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/mysteriumnetwork/node/services/wireguard/connection/dns"
	"vpn/utils"
	"vpn/utils/actionstack"
	"vpn/utils/key"
	"vpn/utils/netutil"
)

// Peer represents wireguard peer.
type Peer struct {
	PublicKey              string       `json:"public_key"`
	Endpoint               *net.UDPAddr `json:"endpoint"`
	AllowedIPs             []string     `json:"allowed_i_ps"`
	KeepAlivePeriodSeconds int          `json:"keep_alive_period_seconds"`
}

// Stats represents wireguard peer statistics information.
type Stats struct {
	BytesSent     uint64    `json:"bytes_sent"`
	BytesReceived uint64    `json:"bytes_received"`
	LastHandshake time.Time `json:"last_handshake"`
}

// DeviceConfig describes wireguard device configuration.
type DeviceConfig struct {
	IfaceName  string    `json:"iface_name"`
	Subnet     net.IPNet `json:"subnet"`
	PrivateKey string    `json:"private_key"`
	ListenPort int       `json:"listen_port"`
	DNSPort    int       `json:"dns_port,omitempty"`
	DNS        []string  `json:"dns"`
	// Used only for unix.
	DNSScriptDir string `json:"dns_script_dir"`

	Peer         Peer `json:"peer"`
	ReplacePeers bool `json:"replace_peers,omitempty"`

	ProxyPort int `json:"proxy_port,omitempty"`
}

type client struct {
	iface      string
	wgClient   *wgctrl.Client
	dnsManager dns.Manager
}

// NewWireguardClient creates new wireguard kernel space client.
func NewWireguardClient() (*client, error) {
	wgClient, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	return &client{
		wgClient:   wgClient,
		dnsManager: dns.NewManager(),
	}, nil
}

func NewWireguardEndpoint() (*client, error) {
	wgClient, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	return &client{
		wgClient:   wgClient,
		dnsManager: dns.NewManager(),
	}, nil
}

func (c *client) ReConfigureDevice(config DeviceConfig) error {
	err := c.configureDevice(config)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) ConfigureDevice(config DeviceConfig) error {
	rollback := actionstack.NewActionStack()

	if err := c.up(config.IfaceName); err != nil {
		return err
	}
	rollback.Push(func() {
		_ = c.DestroyDevice(config.IfaceName)
	})

	if config.Peer.Endpoint != nil {
		if err := netutil.AddDefaultRoute(config.IfaceName); err != nil {
			rollback.Run()
			return err
		}
	}

	err := c.configureDevice(config)
	if err != nil {
		rollback.Run()
		return err
	}

	return nil
}

func (c *client) configureDevice(config DeviceConfig) error {
	if err := cmdutil.SudoExec("ip", "address", "replace", "dev", config.IfaceName, config.Subnet.String()); err != nil {
		return err
	}

	peer, err := peerConfig(config.Peer)
	if err != nil {
		return err
	}

	privateKey, err := stringToKey(config.PrivateKey)
	if err != nil {
		return err
	}

	c.iface = config.IfaceName
	deviceConfig := wgtypes.Config{
		PrivateKey:   &privateKey,
		ListenPort:   &config.ListenPort,
		Peers:        []wgtypes.PeerConfig{peer},
		ReplacePeers: true,
	}

	if err := c.wgClient.ConfigureDevice(c.iface, deviceConfig); err != nil {
		return fmt.Errorf("could not configure kernel space device: %w", err)
	}

	if err := c.dnsManager.Set(dns.Config{
		ScriptDir: config.DNSScriptDir,
		IfaceName: config.IfaceName,
		DNS:       config.DNS,
	}); err != nil {
		return fmt.Errorf("could not set DNS: %w", err)
	}
	return nil
}

func peerConfig(peer Peer) (wgtypes.PeerConfig, error) {
	endpoint := peer.Endpoint
	publicKey, err := stringToKey(peer.PublicKey)
	if err != nil {
		return wgtypes.PeerConfig{}, errors.Wrap(err, "could not convert string key to wgtypes.Key")
	}

	// Apply keep alive interval
	var keepAliveInterval *time.Duration
	if peer.KeepAlivePeriodSeconds > 0 {
		interval := time.Duration(peer.KeepAlivePeriodSeconds) * time.Second
		keepAliveInterval = &interval
	}

	// Apply allowed IPs network
	var allowedIPs []net.IPNet
	for _, ip := range peer.AllowedIPs {
		_, network, err := net.ParseCIDR(ip)
		if err != nil {
			return wgtypes.PeerConfig{}, fmt.Errorf("could not parse IP %q: %v", ip, err)
		}
		allowedIPs = append(allowedIPs, *network)
	}

	return wgtypes.PeerConfig{
		Endpoint:                    endpoint,
		PublicKey:                   publicKey,
		AllowedIPs:                  allowedIPs,
		PersistentKeepaliveInterval: keepAliveInterval,
	}, nil
}

func (c *client) PeerStats(string) (Stats, error) {
	d, err := c.wgClient.Device(c.iface)
	if err != nil {
		return Stats{}, err
	}

	if len(d.Peers) != 1 {
		return Stats{}, errors.New("kernelspace: exactly 1 peer expected")
	}

	return Stats{
		BytesReceived: uint64(d.Peers[0].ReceiveBytes),
		BytesSent:     uint64(d.Peers[0].TransmitBytes),
		LastHandshake: d.Peers[0].LastHandshakeTime,
	}, nil
}

func (c *client) DestroyDevice(name string) error {
	return cmdutil.SudoExec("ip", "link", "del", "dev", name)
}

func (c *client) up(iface string) error {
	rollback := actionstack.NewActionStack()
	if d, err := c.wgClient.Device(iface); err != nil || d.Name != iface {
		if err := cmdutil.SudoExec("ip", "link", "add", "dev", iface, "type", "wireguard"); err != nil {
			return err
		}
	}
	rollback.Push(func() {
		_ = c.DestroyDevice(iface)
	})

	if err := cmdutil.SudoExec("ip", "link", "set", "dev", iface, "up"); err != nil {
		rollback.Run()
		return err
	}

	return nil
}

func (c *client) Close() (err error) {
	errs := utils.ErrorCollection{}
	if err := c.DestroyDevice(c.iface); err != nil {
		errs.Add(err)
	}
	if err := c.wgClient.Close(); err != nil {
		errs.Add(err)
	}
	if err := c.dnsManager.Clean(); err != nil {
		errs.Add(err)
	}
	if err := errs.Error(); err != nil {
		return fmt.Errorf("could not close client: %w", err)
	}
	return nil
}

func stringToKey(key string) (wgtypes.Key, error) {
	k, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return wgtypes.Key{}, err
	}
	return wgtypes.NewKey(k)
}

func getValidIfaceName() string {
	var generateRandName func() string
	generateRandName = func() string {
		index := rand.Intn(10)
		ifaceName := "wg" + strconv.Itoa(index)
		return ifaceName
	}
	var isNameValid func(string) bool
	isNameValid = func(name string) bool {
		ifaces, _ := net.Interfaces()
		for _, v := range ifaces {
			if v.Name == name {
				if name == v.Name {
					return false
				}
			}
		}
		return true
	}
	ifaceName := generateRandName()
	for ; !isNameValid(ifaceName); ifaceName = generateRandName() {
	}
	return ifaceName
}

var allocator *iputil.Allocator = iputil.NewAllocator(net.IPNet{
	IP:   net.ParseIP("10.182.0.0").To4(),
	Mask: net.IPv4Mask(255, 255, 0, 0),
})

var ifaceMap map[string]*DeviceConfig

func ConfigFactory() (*DeviceConfig, error) {
	var config DeviceConfig
	config.IfaceName = getValidIfaceName()
	subnet, err := allocator.AllocateIPNet()
	if err != nil {
		return nil, err
	}
	config.Subnet = subnet
	privatekey, err := key.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	config.PrivateKey = privatekey
	config.DNSScriptDir = "/etc/openvpn/"
	config.DNS = []string{"8.8.8.8"}
	if ifaceMap == nil {
		ifaceMap = make(map[string]*DeviceConfig, 10)
	}
	ifaceMap[config.IfaceName] = &config
	return &config, nil
}

func DestroyConfig(ifaceName string) {
	config := ifaceMap[ifaceName]
	allocator.ReleaseIPNet(config.Subnet)
}
