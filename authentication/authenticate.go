package authentication

import (
	"errors"
	"fmt"

	configMain "AmbiTradfri/config"
	"AmbiTradfri/discovery"

	tradfri "github.com/eriklupander/tradfri-go/tradfri"
)

// ConnectToClient tries to authenticate and connect to the stored ikea client
func ConnectToClient(clientID string, psk string, reauth bool) (*tradfri.Client, error) {
	cfg, err := configMain.NewConfig()

	if err != nil {
		return nil, err
	}

	if cfg.GeneratedPsk == "" || reauth {
		gateway := discovery.DiscoverGateway()

		if gateway == nil {
			return nil, errors.New("Could not find any gateway")
		}

		address := gateway.V4address + ":5684" // just, why

		generatedPsk, err := authenticateWithClient(address, clientID, psk)

		if err != nil {
			return nil, err
		}

		cfg.Address = address
		cfg.ClientID = clientID
		cfg.StandardPsk = psk
		cfg.GeneratedPsk = *generatedPsk

		err = cfg.Save()

		if err != nil {
			return nil, err
		}
	}

	return connectToClientAfterAuth(cfg.Address, cfg.ClientID, cfg.GeneratedPsk), nil
}

func authenticateWithClient(address string, clientID string, psk string) (*string, error) {
	client := tradfri.NewTradfriClient(address, "Client_identity", psk) // This "Client_identity" is very fucking important, otherwise it wont authenticate

	auth, error := client.AuthExchange(clientID)

	if error != nil {
		return nil, error
	}

	return &auth.Token, nil
}

func connectToClientAfterAuth(address string, clientID string, genPsk string) *tradfri.Client {
	fmt.Println("Connecting to gateway at " + address + ", if this takes too long the gateway probably changed adress then try with -reauth")

	client := tradfri.NewTradfriClient(address, clientID, genPsk)

	return client
}
