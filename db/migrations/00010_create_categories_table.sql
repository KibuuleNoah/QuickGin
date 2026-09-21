-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS categories_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE TABLE IF NOT EXISTS categories (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('categories_id_seq'),
    name VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

CREATE TRIGGER set_categories_created_at BEFORE INSERT ON categories FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_categories_modtime BEFORE UPDATE ON categories FOR EACH ROW EXECUTE PROCEDURE update_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
DROP SEQUENCE IF EXISTS categories_id_seq;
-- +goose StatementEnd
