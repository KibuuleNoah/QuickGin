-- +goose Up
WITH new_wishlist AS (
  INSERT INTO wishlists (id, user_id, name, is_public, created_at, updated_at)
  VALUES (
    generate_short_id('wishlists_id_seq'),
    (SELECT id FROM users WHERE username = 'MoxiePlus'),
    'My Wishlist',
    FALSE,
    extract(epoch from now())::bigint,
    extract(epoch from now())::bigint
  )
  RETURNING id
)
INSERT INTO wishlist_items (wishlist_id, item_id, created_at, updated_at)
SELECT new_wishlist.id, items.id, extract(epoch from now())::bigint, extract(epoch from now())::bigint
FROM new_wishlist, items
LIMIT 3;

-- +goose Down
