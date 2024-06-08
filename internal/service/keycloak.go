package service

import (
	"github.com/Nerzal/gocloak/v13"
	"main/pkg"
)

type keycloak struct {
	gocloak      *gocloak.GoCloak
	clientId     string
	clientSecret string
	realm        string
}

func newKeycloak(conf *pkg.KeyCloak) *keycloak {
	return &keycloak{
		gocloak:      pkg.NewKeyCloak(conf),
		clientId:     conf.ClientId,
		clientSecret: conf.ClientSecret,
		realm:        conf.Realm,
	}
}
