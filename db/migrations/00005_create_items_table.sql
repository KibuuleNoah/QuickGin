-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE SEQUENCE IF NOT EXISTS items_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

CREATE TABLE items (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('items_id_seq'),
    name VARCHAR(255) NOT NULL,
    rating DECIMAL(3, 2) NOT NULL DEFAULT 0.00,
    discount DECIMAL(5, 2) DEFAULT NULL,
    price DECIMAL(10, 2) NOT NULL,
    sold_count INT NOT NULL DEFAULT 0,
    stock_count INT NOT NULL DEFAULT 0,
    thumbnail_url VARCHAR(512) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,

    search_vector tsvector GENERATED ALWAYS AS (
        to_tsvector('english', coalesce(name, ''))
    ) STORED
);


CREATE INDEX idx_items_search_vector ON items USING GIN (search_vector);

CREATE INDEX idx_items_name_trgm ON items USING GIN (name gin_trgm_ops);

CREATE INDEX idx_items_stock_price ON items (stock_count, price) WHERE stock_count > 0;

CREATE INDEX idx_items_rating_popularity ON items (rating DESC, sold_count DESC);


CREATE TRIGGER set_items_created_at BEFORE INSERT ON items FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_items_modtime BEFORE UPDATE ON items FOR EACH ROW EXECUTE PROCEDURE update_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
