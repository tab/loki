-- +goose Up
ALTER TABLE users ADD CONSTRAINT users_identity_number_not_empty
  CHECK (length(trim(identity_number)) > 10);

ALTER TABLE users ADD CONSTRAINT users_personal_code_not_empty
  CHECK (length(trim(personal_code)) > 10);

-- +goose Down
ALTER TABLE users DROP CONSTRAINT users_identity_number_not_empty;
ALTER TABLE users DROP CONSTRAINT users_personal_code_not_empty;
