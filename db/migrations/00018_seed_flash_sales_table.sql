-- +goose Up
-- +goose StatementBegin
WITH new_flashsale1 AS (
  INSERT INTO flashsales (title, start_time, end_time, status, created_at, updated_at)
  VALUES (
    'Weekend Electronics Clearance',
    (extract(epoch from now()) - 3600)::bigint * 1000, -- Started 1 hour ago
    (extract(epoch from now()) + 86400)::bigint * 1000, -- Ends in 24 hours
    'ACTIVE',
    extract(epoch from now())::bigint * 1000,
    extract(epoch from now())::bigint * 1000
  )
  RETURNING id
)
INSERT INTO flashsale_items (flashsale_id, item_id, sale_price, created_at, updated_at)
SELECT 
  new_flashsale1.id, 
  items.id, 
  5000, 
  extract(epoch from now())::bigint * 1000, 
  extract(epoch from now())::bigint * 1000
FROM new_flashsale1, items
LIMIT 6;

WITH new_flashsale2 AS (
  INSERT INTO flashsales (title, start_time, end_time, status, created_at, updated_at)
  VALUES (
    'Black Friday Early Bird Sale',
    (extract(epoch from now()) + 604800)::bigint * 1000, -- Starts in 7 days
    (extract(epoch from now()) + 691200)::bigint * 1000, -- Ends in 8 days
    'SCHEDULED',
    extract(epoch from now())::bigint * 1000,
    extract(epoch from now())::bigint * 1000
  )
  RETURNING id
)
INSERT INTO flashsale_items (flashsale_id, item_id, sale_price, created_at, updated_at)
SELECT 
  new_flashsale2.id, 
  items.id, 
  8000, 
  extract(epoch from now())::bigint * 1000, 
  extract(epoch from now())::bigint * 1000
FROM new_flashsale2, items
LIMIT 7;

WITH new_flashsale3 AS (
  INSERT INTO flashsales (title, start_time, end_time, status, created_at, updated_at)
  VALUES (
    'Flash Fashion Blowout Draft',
    (extract(epoch from now()) + 1209600)::bigint * 1000,
    (extract(epoch from now()) + 1296000)::bigint * 1000,
    'DRAFT',
    extract(epoch from now())::bigint * 1000,
    extract(epoch from now())::bigint * 1000
  )
  RETURNING id
)
INSERT INTO flashsale_items (flashsale_id, item_id, sale_price, created_at, updated_at)
SELECT 
  new_flashsale3.id, 
  items.id, 
  8000, 
  extract(epoch from now())::bigint * 1000, 
  extract(epoch from now())::bigint * 1000
FROM new_flashsale3, items
LIMIT 8;
-- +goose StatementEnd

-- +goose Down
