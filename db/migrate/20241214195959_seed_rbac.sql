-- +goose Up
INSERT INTO roles (name, description)
VALUES ('admin', 'Administrator role'),
       ('manager', 'Manager role'),
       ('user', 'User role');

INSERT INTO permissions (name, description)
VALUES ('read:self', 'Read own profile'),
       ('write:self', 'Update own profile'),
       ('read:users', 'Read any profile'),
       ('write:users', 'Update any profile');

INSERT INTO scopes (name, description)
VALUES ('self-service', 'self-service application scope'),
        ('sso-service', 'sso-service application scope');

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE (r.name = 'admin' AND p.name IN ('read:self', 'write:self', 'read:users', 'write:users'))
   OR (r.name = 'manager' AND p.name IN ('read:self', 'write:self', 'read:users'))
   OR (r.name = 'user' AND p.name IN ('read:self', 'write:self'));

-- +goose Down
