package rpcs

import (
	"google.golang.org/grpc"

	proto "loki/internal/app/rpcs/proto/sso/v1"
)

type Registry struct {
	permissions proto.PermissionServiceServer
	scopes      proto.ScopeServiceServer
}

func NewRegistry(
	permissions proto.PermissionServiceServer,
	scopes proto.ScopeServiceServer,
) *Registry {
	return &Registry{
		permissions: permissions,
		scopes:      scopes,
	}
}

func (r *Registry) RegisterAll(server *grpc.Server) {
	proto.RegisterPermissionServiceServer(server, r.permissions)
	proto.RegisterScopeServiceServer(server, r.scopes)
}
