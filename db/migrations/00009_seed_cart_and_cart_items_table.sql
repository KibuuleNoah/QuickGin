-- +goose Up
INSERT INTO carts (id, guest_id, created_at, updated_at) VALUES 
( 'Hj58W1','Hj58Pg', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint ),
( 'Hj58X2','Hj58Og', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint ),
( 'Hj58M1','Hj58Mg', EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint );

INSERT INTO cart_items (cart_id, item_id, quantity, price_at_add, created_at, updated_at) VALUES
('Hj58W1','Xy7Z2a', 10, 50000, EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Hj58W1','K9b3Wp', 5, 38900, EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Hj58X2','T2n6Fz', 12, 7000, EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Hj58X2','H5j9Lx', 45, 10000, EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Hj58X2','Ab8C3d', 5, 59000, EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Hj58M1','Ef2G4h', 2, 50000, EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Hj58M1','Jk9L1m', 9, 25000, EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Hj58M1','Np5Q7r', 7, 10500, EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Hj58M1','St3U9v', 90, 12000, EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint),
('Hj58M1','Wx6Y2z', 1, 10000, EXTRACT(EPOCH FROM NOW())::bigint, EXTRACT(EPOCH FROM NOW())::bigint);

-- +goose Down
