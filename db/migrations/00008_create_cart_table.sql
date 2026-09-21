-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS carts_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS cart_items_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

CREATE TABLE IF NOT EXISTS carts (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('carts_id_seq'),
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    guest_id VARCHAR(6),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    CHECK ( user_id IS NOT NULL OR guest_id IS NOT NULL )
);

CREATE TABLE IF NOT EXISTS cart_items (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('cart_items_id_seq'),
    cart_id TEXT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    item_id TEXT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK ( quantity > 0 ),
    price_at_add NUMERIC(10,2) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);


-- Handles default values on insert
CREATE TRIGGER set_carts_created_at BEFORE INSERT ON carts FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER set_cart_items_created_at BEFORE INSERT ON cart_items FOR EACH ROW EXECUTE PROCEDURE created_at_column();

-- Handles automated timestamp updates on modifications
CREATE TRIGGER update_carts_modtime BEFORE UPDATE ON carts FOR EACH ROW EXECUTE PROCEDURE update_at_column();
CREATE TRIGGER update_cart_items_modtime BEFORE UPDATE ON cart_items FOR EACH ROW EXECUTE PROCEDURE update_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
-- +goose StatementEnd
