package api

import (
	"main/pkg"
)

type Api interface {
	InitRoutes(*pkg.Config)
}
