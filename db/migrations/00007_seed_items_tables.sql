-- +goose Up
-- Seed Core Items Table

INSERT INTO items (id, name, rating, discount, price, sold_count, stock_count, thumbnail_url, created_at, updated_at) VALUES
('Xy7Z2a', 'EliteShield Performance Men''s Jackets', 0.00, 51.00, 180000.00, 9, 10, 'images3-removebg-preview.png', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('K9b3Wp', 'Gentlemen''s Summer Gray Hat - Premium Blend', 0.00, 34.00, 45000.00, 9, 10, 'Smoke-Grey-Sleek-Silk-Finish-Fedora-Hat-for-Men-The-Oliver-Agnoulita-Hats-1.webp', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('M4v8Qr', 'OptiZoom Camera Shoulder Bag', 0.00, 41.00, 120000.00, 5, 10, 'png-clipart-handbag-messenger-bags-women-bag-luggage-bags-orange-removebg-preview.png', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('T2n6Fz', 'Cloudy Chic - Grey Peep Toe Heeled Sandals', 0.00, 53.00, 130000.00, 5, 10, '61KzfNJSr-L._AC_UY1000_-removebg-preview.png', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('H5j9Lx', 'Essentials Men''s Long-Sleeve Oxford Shirt', 4.90, NULL, 75000.00, 0, 10, 'images-removebg-preview.png', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Ab8C3d', 'UrbanEdge Men''s Jeans Collection', 4.90, 32.00, 95000.00, 0, 10, '1-removebg-preview.png', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Ef2G4h', 'Essentials Men''s Long-Sleeve Oxford Shirt', 4.90, NULL, 75000.00, 0, 10, 'images-removebg-preview.png', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Jk9L1m', 'StyleHaven Men''s Fashionable Brogues', 4.90, 39.00, 160000.00, 0, 10, 'images2-removebg-preview.png', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Np5Q7r', 'Essential Long-Sleeve Crewneck Shirt for Men', 4.90, NULL, 50000.00, 0, 10, 'Screenshot_From_2026-08-06_13-59-02-removebg-preview.png', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('St3U9v', 'ClassicGent Men''s Formal Shoes', 4.90, NULL, 150000.00, 0, 10, 'images-removebg-preview3.png', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Wx6Y2z', 'UrbanFlex Men''s Short Pants Collection', 4.90, NULL, 65000.00, 0, 10, 'Screenshot_From_2026-08-06_13-58-34-removebg-preview.png', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint);

-- Seed Matching Item Details Table
INSERT INTO item_details (item_id, description, is_available, category, tags, gallery_image_urls, specifications, created_at, updated_at) VALUES
(
    'Xy7Z2a', 
    'High performance all-weather shield protection jacket built for the modern active man.', 
    TRUE, 
    'Apparel', 
    '["jacket", "mens", "outerwear"]'::jsonb, 
    '["images3-removebg-preview.png"]'::jsonb, 
    '[{"label": "Material", "value": "Polyester Blend"}, {"label": "Size", "value": "L"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
),
(
    'K9b3Wp', 
    'A classic timeless gray hat featuring a premium blend design designed for comfortable hot summer days.', 
    TRUE, 
    'Accessories', 
    '["hat", "summer", "mens", "fedora"]'::jsonb, 
    '["Smoke-Grey-Sleek-Silk-Finish-Fedora-Hat-for-Men-The-Oliver-Agnoulita-Hats-1.webp"]'::jsonb, 
    '[{"label": "Style", "value": "Fedora"}, {"label": "Finish", "value": "Silk"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
),
(
    'M4v8Qr', 
    'Ergonomic camera shoulder bag engineered to safely secure cameras, lenses, and accessories on the go.', 
    TRUE, 
    'Bags', 
    '["camera", "bag", "messenger", "orange"]'::jsonb, 
    '["png-clipart-handbag-messenger-bags-women-bag-luggage-bags-orange-removebg-preview.png"]'::jsonb, 
    '[{"label": "Type", "value": "Messenger Bag"}, {"label": "Color", "value": "Orange"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
),
(
    'T2n6Fz', 
    'Elegant peep toe heeled sandals featuring an attractive grey chic coloring and cushioned comfortable sole layout.', 
    TRUE, 
    'Footwear', 
    '["heels", "sandals", "womens", "shoes"]'::jsonb, 
    '["61KzfNJSr-L._AC_UY1000_-removebg-preview.png"]'::jsonb, 
    '[{"label": "Type", "value": "Heeled Sandals"}, {"label": "Heel Height", "value": "3 inches"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
),
(
    'H5j9Lx', 
    'The quintessential classic long-sleeve button-up oxford cotton shirt for corporate daily wear or semi-casual activities.', 
    TRUE, 
    'Apparel', 
    '["shirt", "oxford", "mens", "formal"]'::jsonb, 
    '["images-removebg-preview.png"]'::jsonb, 
    '[{"label": "Material", "value": "100% Cotton"}, {"label": "Fit", "value": "Slim Fit"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
),
(
    'Ab8C3d', 
    'Premium quality tailored denim jeans designed for everyday durability and modern street wear.', 
    TRUE, 
    'Apparel', 
    '["jeans", "mens", "denim", "pants"]'::jsonb, 
    '["1-removebg-preview.png"]'::jsonb, 
    '[{"label": "Material", "value": "Stretch Denim"}, {"label": "Fit", "value": "Slim Fit"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
),
(
    'Ef2G4h', 
    'The quintessential classic long-sleeve button-up oxford cotton shirt for corporate daily wear or semi-casual activities.', 
    TRUE, 
    'Apparel', 
    '["shirt", "oxford", "mens", "formal"]'::jsonb, 
    '["images-removebg-preview.png"]'::jsonb, 
    '[{"label": "Material", "value": "100% Cotton"}, {"label": "Fit", "value": "Regular Fit"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
),
(
    'Jk9L1m', 
    'Exquisitely designed leather brogues tailored perfectly for upscale meetings, business events, or evenings out.', 
    TRUE, 
    'Footwear', 
    '["shoes", "brogues", "mens", "leather"]'::jsonb, 
    '["images2-removebg-preview.png"]'::jsonb, 
    '[{"label": "Material", "value": "Genuine Leather"}, {"label": "Style", "value": "Brogue"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
),
(
    'Np5Q7r', 
    'Comfortable lightweight cotton-blend long sleeve crewneck top perfect for seasonal layering adjustments.', 
    TRUE, 
    'Apparel', 
    '["shirt", "crewneck", "mens", "casual"]'::jsonb, 
    '["Screenshot_From_2026-08-06_13-59-02-removebg-preview.png"]'::jsonb, 
    '[{"label": "Neckline", "value": "Crew Neck"}, {"label": "Sleeve", "value": "Long Sleeve"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
),
(
    'St3U9v', 
    'Elegant executive oxford lace-up formal dress shoes designed to add high confidence to sharp tailored outfits.', 
    TRUE, 
    'Footwear', 
    '["shoes", "formal", "mens", "oxford"]'::jsonb, 
    '["images-removebg-preview3.png"]'::jsonb, 
    '[{"label": "Sole Type", "value": "Anti-Slip Rubber"}, {"label": "Closure", "value": "Lace-up"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
),
(
    'Wx6Y2z', 
    'Breathable casual shorts built from premium flexible textiles for everyday urban utility and leisure settings.', 
    TRUE, 
    'Apparel', 
    '["shorts", "pants", "mens", "summer"]'::jsonb, 
    '["Screenshot_From_2026-08-06_13-58-34-removebg-preview.png"]'::jsonb, 
    '[{"label": "Material", "value": "Cotton Khaki"}, {"label": "Pockets", "value": "4-Pocket Setup"}]'::jsonb,
    EXTRACT(EPOCH FROM NOW())::bigint, 
    EXTRACT(EPOCH FROM NOW())::bigint
);

-- +goose Down
