-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS flashsales_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE TABLE flashsales (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('flashsales_id_seq'),
    title VARCHAR(255) NOT NULL,
    start_time BIGINT NOT NULL,
    end_time BIGINT NOT NULL CHECK (end_time > start_time),
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'SCHEDULED', 'ACTIVE', 'ENDED')),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

CREATE SEQUENCE IF NOT EXISTS flashsale_items_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE TABLE flashsale_items (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('flashsale_items_id_seq'),
    flashsale_id TEXT NOT NULL REFERENCES flashsales(id) ON DELETE CASCADE,
    item_id TEXT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    sale_price NUMERIC(10, 2) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL  

);


CREATE INDEX idx_flashsale_items_flashsale_id ON flashsale_items(flashsale_id);


CREATE TRIGGER set_flashsales_created_at BEFORE INSERT ON flashsales FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_flashsales_modtime BEFORE UPDATE ON flashsales FOR EACH ROW EXECUTE PROCEDURE update_at_column();

CREATE TRIGGER set_flashsale_items_created_at BEFORE INSERT ON flashsale_items FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_flashsale_items_modtime BEFORE UPDATE ON flashsale_items FOR EACH ROW EXECUTE PROCEDURE update_at_column();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS flashsale_items;
DROP TABLE IF EXISTS flashsales;
-- +goose StatementEnd

