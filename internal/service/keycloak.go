package service

import (
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"main/pkg"
)

type keycloak struct {
	gocloak      *gocloak.GoCloak
	clientId     string
	clientSecret string
	realm        string
}

func newKeycloak(conf *pkg.Config) *keycloak {
	return &keycloak{
		gocloak:      gocloak.NewClient(fmt.Sprintf("http://%s:%s", conf.KeyCloak.Host, conf.KeyCloak.Port)),
		clientId:     conf.KeyCloak.ClientId,
		clientSecret: conf.KeyCloak.ClientSecret,
		realm:        conf.KeyCloak.Realm,
	}
}
