package pkg

import (
	"fmt"
	"github.com/Nerzal/gocloak/v13"
)

func NewKeyCloak(conf *KeyCloak) *gocloak.GoCloak {
	return gocloak.NewClient(fmt.Sprintf("%s:%s", conf.Host, conf.Port))
}
