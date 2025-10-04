-- name: GetLastOrders :many
SELECT 
    order_uid,
    track_number,
    entry,
    locale,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    date_created,
    oof_shard
FROM orders
ORDER BY date_created DESC
LIMIT $1;