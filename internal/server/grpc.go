package server

import (
	authv1 "github.com/go-kratos/kratos-layout/api/auth/v1"
	countryv1 "github.com/go-kratos/kratos-layout/api/country/v1"
	helloworldv1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	provincev1 "github.com/go-kratos/kratos-layout/api/province/v1"
	userv1 "github.com/go-kratos/kratos-layout/api/user/v1"
	wardv1 "github.com/go-kratos/kratos-layout/api/ward/v1"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService, user *service.UserService, auth *service.AuthService, country *service.CountryService, province *service.ProvinceService, ward *service.WardService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	helloworldv1.RegisterGreeterServer(srv, greeter)
	userv1.RegisterUserServiceServer(srv, user)
	authv1.RegisterAuthServiceServer(srv, auth)
	countryv1.RegisterCountryServiceServer(srv, country)
	provincev1.RegisterProvinceServiceServer(srv, province)
	wardv1.RegisterWardServiceServer(srv, ward)
	return srv
}
