package discovery

import (
	"regexp"
	"strings"

	mdns "github.com/hashicorp/mdns"
)

// TradfriGatewayAdress describes the exact location on the network of a tradfri hub
type TradfriGatewayAdress struct {
	Name      string
	Host      string
	Version   string
	V4address string
	V6address string
}

// DiscoverGateway uses the mDNS service to locate the first gateway it finds on the network
func DiscoverGateway() *TradfriGatewayAdress {

	tradfriRegex := regexp.MustCompile("^gw\\-[0-9a-f]{12}")

	entriesCh := make(chan *mdns.ServiceEntry, 4)

	var gateway *TradfriGatewayAdress = nil

	go func() {
		for entry := range entriesCh {
			if tradfriRegex.MatchString(entry.Name) {

				var version string

				for _, field := range entry.InfoFields {
					split := strings.Split(field, "=")
					if split[0] == "version" {
						version = split[1]
					}
				}

				gateway = &TradfriGatewayAdress{
					Name:    entry.Name,
					Host:    entry.Host,
					Version: version,

					V4address: entry.AddrV4.String(),
					V6address: entry.AddrV6.String(),
				}
			}
		}
	}()

	mdns.Lookup("_coap._udp", entriesCh)
	close(entriesCh)

	return gateway
}
