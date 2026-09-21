-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS collections_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE TABLE collections (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('collections_id_seq'),
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) DEFAULT 'My Collection',
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL  
);

CREATE SEQUENCE IF NOT EXISTS collection_items_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE TABLE collection_items (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('collection_items_id_seq'),
    collection_id TEXT NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    item_id TEXT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL 
);

CREATE SEQUENCE IF NOT EXISTS collection_contributors_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE TABLE collection_contributors (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('collection_contributors_id_seq'),
    collection_id TEXT NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL 
);


CREATE INDEX idx_collections_user_id ON collections(user_id);
CREATE INDEX idx_collection_items_collection_id ON collection_items(collection_id);


CREATE TRIGGER set_collections_created_at BEFORE INSERT ON collections FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_collections_modtime BEFORE UPDATE ON collections FOR EACH ROW EXECUTE PROCEDURE update_at_column();

CREATE TRIGGER set_collection_items_created_at BEFORE INSERT ON collection_items FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_collection_items_modtime BEFORE UPDATE ON collection_items FOR EACH ROW EXECUTE PROCEDURE update_at_column();

CREATE TRIGGER set_collection_contributors_created_at BEFORE INSERT ON collection_contributors FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_collection_contributors_modtime BEFORE UPDATE ON collection_contributors FOR EACH ROW EXECUTE PROCEDURE update_at_column();



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS collection_contributors;
DROP TABLE IF EXISTS collection_items;
DROP TABLE IF EXISTS collections;
-- +goose StatementEnd

