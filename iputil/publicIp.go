package iputil

import (
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func GetPublicIP() string {
	url := "https://api.ipify.org?format=text" // we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	log.Info().Msg("Getting Public IP address from  ipify ...\n")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Get Public IP Failed")
		panic(err)
	}
	log.Debug().Msg("Success Get Public IP address from ipify\n")
	return string(ip)
}
