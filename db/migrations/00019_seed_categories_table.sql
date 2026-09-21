-- +goose Up
-- +goose StatementBegin
INSERT INTO categories (name, created_at, updated_at) VALUES
('Electronics', EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('Clothing', EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('Home & Kitchen', EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('Books', EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
('Sports & Outdoors', EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT)
ON CONFLICT (id) DO NOTHING;

UPDATE items 
SET category_id = (
    SELECT id 
    FROM categories 
    ORDER BY RANDOM() 
    LIMIT 1
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE categories RESTART IDENTITY CASCADE;
-- +goose StatementEnd
