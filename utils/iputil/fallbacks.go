package iputil

import (
	"errors"
	"github.com/mysteriumnetwork/node/utils/random"
	"io"
	"net"
	"strings"
	"vpn/requests"
)

var rng = random.NewTimeSeededRand()

func shuffleStringSlice(slice []string) []string {
	tmp := make([]string, len(slice))
	copy(tmp, slice)
	rng.Shuffle(len(tmp), func(i, j int) {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	})
	return tmp
}

// RequestAndParsePlainIPResponse requests and parses a plain IP response.
func RequestAndParsePlainIPResponse(c *requests.HTTPClient, url string) (string, error) {
	req, err := requests.NewGetRequest(url, "", nil)
	if err != nil {
		return "", err
	}

	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	r, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	ipv4addr := net.ParseIP(strings.TrimSpace(string(r)))
	if ipv4addr == nil {
		return "", errors.New("could not parse iputil response")
	}
	return ipv4addr.String(), err
}
