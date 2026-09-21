-- +goose Up
WITH new_collection1 AS (
  INSERT INTO collections (user_id, name, is_public, created_at, updated_at)
  VALUES (
    (SELECT id FROM users WHERE username = 'MoxiePlus'),
    'Fit Check',
    TRUE,
    extract(epoch from now())::bigint,
    extract(epoch from now())::bigint
  )
  RETURNING id
)
INSERT INTO collection_items (collection_id, item_id, created_at, updated_at)
SELECT new_collection1.id, items.id, extract(epoch from now())::bigint, extract(epoch from now())::bigint
FROM new_collection1, items
LIMIT 3;

WITH new_collection2 AS (
  INSERT INTO collections (user_id, name, is_public, created_at, updated_at)
  VALUES (
    (SELECT id FROM users WHERE username = 'MoxiePlus'),
    'My Dream',
    TRUE,
    extract(epoch from now())::bigint,
    extract(epoch from now())::bigint
  )
  RETURNING id
)
INSERT INTO collection_items (collection_id, item_id, created_at, updated_at)
SELECT new_collection2.id, items.id, extract(epoch from now())::bigint, extract(epoch from now())::bigint
FROM new_collection2, items
LIMIT 4;

-- +goose Down
