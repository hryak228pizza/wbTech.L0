-- name: GetItemsByTrackNumber :many
SELECT 
    id,
    order_uid,
    chrt_id,
    track_number,
    price,
    rid,
    name,
    sale,
    size,
    total_price,
    nm_id,
    brand,
    status
FROM items
WHERE track_number = $1;