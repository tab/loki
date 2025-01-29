-- +goose Up
INSERT INTO roles (name, description)
VALUES ('admin', 'Administrator role'),
       ('manager', 'Manager role'),
       ('user', 'User role');

INSERT INTO permissions (name, description)
VALUES ('read:self', 'Read own profile'),
       ('write:self', 'Update own profile'),
       ('read:permissions', 'Read permissions'),
       ('write:permissions', 'Update permissions'),
        ('read:roles', 'Read roles'),
        ('write:roles', 'Update roles'),
        ('read:scopes', 'Read scopes'),
        ('write:scopes', 'Update scopes'),
       ('read:users', 'Read any profile'),
       ('write:users', 'Update any profile'),
        ('read:tokens', 'Read tokens'),
        ('write:tokens', 'Update tokens');

INSERT INTO scopes (name, description)
VALUES ('self-service', 'self-service application scope'),
        ('sso-service', 'sso-service application scope');

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE (r.name = 'admin' AND p.name IN (
  'read:self',
  'write:self',
  'read:users',
  'write:users',
  'read:roles',
  'write:roles',
  'read:scopes',
  'write:scopes',
  'read:permissions',
  'write:permissions',
  'read:tokens',
  'write:tokens'
))
OR (r.name = 'manager' AND p.name IN ('read:self', 'write:self', 'read:users'))
OR (r.name = 'user' AND p.name IN ('read:self', 'write:self'));

-- +goose Down
