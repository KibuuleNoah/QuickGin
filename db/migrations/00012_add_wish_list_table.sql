-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS wishlists_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE TABLE wishlists (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('wishlists_id_seq'),
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) DEFAULT 'My Wishlist',
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL  
);

CREATE SEQUENCE IF NOT EXISTS wishlist_items_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE TABLE wishlist_items (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('wishlist_items_id_seq'),
    wishlist_id TEXT NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    item_id TEXT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL  

);

CREATE INDEX idx_wishlists_user_id ON wishlists(user_id);
CREATE INDEX idx_wishlist_items_wishlist_id ON wishlist_items(wishlist_id);


CREATE TRIGGER set_wishlists_created_at BEFORE INSERT ON wishlists FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_wishlists_modtime BEFORE UPDATE ON wishlists FOR EACH ROW EXECUTE PROCEDURE update_at_column();

CREATE TRIGGER set_wishlist_items_created_at BEFORE INSERT ON wishlist_items FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_wishlist_items_modtime BEFORE UPDATE ON wishlist_items FOR EACH ROW EXECUTE PROCEDURE update_at_column();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wishlist_items;
DROP TABLE IF EXISTS wishlists;
-- +goose StatementEnd
