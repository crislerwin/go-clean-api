-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD COLUMN reason VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP COLUMN reason;
-- +goose StatementEnd
