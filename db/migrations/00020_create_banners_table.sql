-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS banners_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

CREATE TABLE IF NOT EXISTS banners (
    id TEXT PRIMARY KEY DEFAULT generate_short_id('banners_id_seq'),
    tag VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    discount VARCHAR(255),
    description TEXT NOT NULL,
    banner_img_url VARCHAR(512) DEFAULT '',
    call_to_action_url VARCHAR(512) DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- PERFORMANCE INDEXES
CREATE INDEX idx_banners_is_active ON banners(is_active);

-- TRIGGER AUTOMATIONS
CREATE TRIGGER set_banners_created_at BEFORE INSERT ON banners FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_banners_modtime BEFORE UPDATE ON banners FOR EACH ROW EXECUTE PROCEDURE update_at_column();

-- SEED MOCK BANNERS
INSERT INTO banners (tag, title, discount, description, banner_img_url, call_to_action_url, is_active)
VALUES 
  (
    '#Kampala Fashion & Style', 
    'Kampala Urban Vibe', 
    'Up to 50% OFF UGX!', 
    'Fresh, premium everyday streetwear and apparel inspired by downtown Kampala style.', 
    'banner1.png', 
    '', 
    TRUE
  ),
  (
    '#Featured Tech Hub Essentials', 
    'Upgrade Your Gadgets', 
    'Save Big on Electronics', 
    'Crush your hustle with high-performance laptops, noise-canceling headphones, and smart devices.', 
    'banner2.png', 
    '', 
    TRUE
  ),
  (
    '#Modern Home Makeover', 
    'Elevate Your Living Space', 
    'Up to 30% Off Decor', 
    'Premium lighting, ergonomic cushions, and local accents designed for beautiful Ugandan homes.', 
    'banner3.png', 
    '', 
    TRUE
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banners;
DROP SEQUENCE IF EXISTS banners_id_seq;
-- +goose StatementEnd
