package rpcs

import (
	"google.golang.org/grpc"

	proto "loki/internal/app/rpcs/proto/sso/v1"
)

type Registry struct {
	permissions proto.PermissionServiceServer
	roles       proto.RoleServiceServer
	scopes      proto.ScopeServiceServer
	tokens      proto.TokenServiceServer
	users       proto.UserServiceServer
}

func NewRegistry(
	permissions proto.PermissionServiceServer,
	roles proto.RoleServiceServer,
	scopes proto.ScopeServiceServer,
	tokens proto.TokenServiceServer,
	users proto.UserServiceServer,
) *Registry {
	return &Registry{
		permissions: permissions,
		roles:       roles,
		scopes:      scopes,
		tokens:      tokens,
		users:       users,
	}
}

func (r *Registry) RegisterAll(server *grpc.Server) {
	proto.RegisterPermissionServiceServer(server, r.permissions)
	proto.RegisterRoleServiceServer(server, r.roles)
	proto.RegisterScopeServiceServer(server, r.scopes)
	proto.RegisterTokenServiceServer(server, r.tokens)
	proto.RegisterUserServiceServer(server, r.users)
}
