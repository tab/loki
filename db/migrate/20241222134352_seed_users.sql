-- +goose Up
INSERT INTO users (identity_number, personal_code, first_name, last_name)
VALUES
  ('PNOEE-50001029996', '50001029996', 'TESTNUMBER', 'ADULT'),
  ('PNOBE-00010299944', '00010299944', 'TESTNUMBER', 'OK')
ON CONFLICT (identity_number) DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON
  (u.identity_number = 'PNOEE-50001029996' AND r.name = 'admin') OR
  (u.identity_number = 'PNOBE-00010299944' AND r.name = 'manager');

INSERT INTO user_scopes (user_id, scope_id)
SELECT u.id, s.id
FROM users u
JOIN scopes s ON s.name = 'sso-service'
WHERE u.identity_number IN ('PNOEE-50001029996', 'PNOBE-00010299944');

-- +goose Down
DELETE FROM users WHERE identity_number IN ('PNOEE-50001029996', 'PNOBE-00010299944');
