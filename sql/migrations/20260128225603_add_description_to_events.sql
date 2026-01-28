-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ADD COLUMN description TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN description;
-- +goose StatementEnd
