package rpcs

import (
    "google.golang.org/grpc"

    proto "loki/internal/app/rpcs/proto/sso/v1"
)

type Registry struct {
    permissions proto.PermissionServiceServer
}

func NewRegistry(
    permissions proto.PermissionServiceServer,
) *Registry {
    return &Registry{
        permissions: permissions,
    }
}

func (r *Registry) RegisterAll(server *grpc.Server) {
    proto.RegisterPermissionServiceServer(server, r.permissions)
}
