package server

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
	"review-service/internal/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewRegister, NewGRPCServer, NewHTTPServer)

func NewRegister(conf *conf.Registry) registry.Registrar {
	// new consul client
	c := api.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	reg := consul.New(client, consul.WithHealthCheck(true))
	return reg
}
