-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ADD COLUMN user_id UUID REFERENCES users(id);

-- Depending on policy, we might want to make it NOT NULL later, but for existing records it must be nullable initially or we default.
-- For now, we leave it nullable or update existing records if any.
-- Given it's a dev environment with seed data, we could update seed data but let's just add the column.
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN user_id;
-- +goose StatementEnd
