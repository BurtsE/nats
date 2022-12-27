SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
FROM items 
INNER JOIN (
    SELECT item_id FROM order_item WHERE order_id=$1
    ) AS its
    ON(item_id=chrt_id)
    ;