-- +goose Up
-- +goose StatementBegin
ALTER TABLE items ADD COLUMN category_id TEXT REFERENCES categories(id) ON DELETE SET NULL;

CREATE INDEX idx_items_category_id ON items(category_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_items_category_id;
ALTER TABLE items DROP COLUMN IF EXISTS category_id;
-- +goose StatementEnd
