package constant

import (
	"AlexSarva/GophKeeper/models"
	"errors"

	"github.com/sarulabs/di"
)

var GlobalContainer di.Container
var ErrInit = errors.New("error while initializing container")
var ErrBuild = errors.New("error while building container")

func BuildContainer(cfg models.Config) error {
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
