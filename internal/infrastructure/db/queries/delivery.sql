-- name: GetDeliveryByOrderUID :one
SELECT 
    order_uid,
    name,
    phone,
    zip,
    city,
    address,
    region,
    email
FROM delivery
WHERE order_uid = $1;