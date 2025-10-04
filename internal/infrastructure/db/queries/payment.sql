-- name: GetPaymentByTransaction :one
SELECT 
    transaction,
    request_id,
    currency,
    provider,
    amount,
    payment_dt,
    bank,
    delivery_cost,
    goods_total,
    custom_fee
FROM payment
WHERE transaction = $1;