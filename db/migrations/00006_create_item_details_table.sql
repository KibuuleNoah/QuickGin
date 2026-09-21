-- +goose Up
-- +goose StatementBegin
CREATE TABLE item_details (
    item_id TEXT PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    tags JSONB NOT NULL,                -- Stores string[]
    gallery_image_urls JSONB NOT NULL,  -- Stores string[]
    specifications JSONB NOT NULL,     -- Stores ItemSpec[] [{label: "...", value: "..."}]
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL    
);
CREATE INDEX idx_item_details_tags ON item_details USING gin (tags jsonb_path_ops);
-- CREATE INDEX idx_item_details_tags ON item_details USING gin (tags);
-- CREATE INDEX idx_items_category ON item_details(category) WHERE is_available = TRUE;

-- Handles default values on insert
CREATE TRIGGER set_item_details_created_at BEFORE INSERT ON item_details FOR EACH ROW EXECUTE PROCEDURE created_at_column();

-- Handles automated timestamp updates on modifications
CREATE TRIGGER update_item_details_modtime BEFORE UPDATE ON item_details FOR EACH ROW EXECUTE PROCEDURE update_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS item_details;
-- +goose StatementEnd
