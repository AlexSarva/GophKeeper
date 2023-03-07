package constant

import (
	"AlexSarva/GophKeeper/models"
	"errors"

	"github.com/sarulabs/di"
)

// GlobalContainer container for store server config info
var GlobalContainer di.Container
var ErrInit = errors.New("error while initializing container")
var ErrBuild = errors.New("error while building container")

// BuildContainer initializer of GlobalContainer struct
func BuildContainer(cfg models.ServerConfig) error {
	builder, builderErr := di.NewBuilder()
	if builderErr != nil {
		return ErrBuild
	}
	if err := builder.Add(di.Def{
		Name:  "server-config",
		Build: func(ctn di.Container) (interface{}, error) { return cfg, nil }}); err != nil {
		return ErrInit
	}
	GlobalContainer = builder.Build()
	return nil
}
