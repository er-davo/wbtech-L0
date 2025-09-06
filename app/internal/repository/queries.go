package repository

const (
	selectExtendedOrderWithoutItemsQuery = `
	SELECT
		o.id, o.order_uid, o.track_number,
		o.entry, o.delivery_id, o.payment_id,
		o.locale, o.internal_signature,
		o.customer_id, o.delivery_service,
		o.shardkey, o.sm_id, o.date_created, o.oof_shard,

		d.id, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,

		p.id, p.transaction, p.request_id,
		p.currency, p.provider, p.amount,
		p.payment_dt, p.bank, p.delivery_cost,
		p.goods_total, p.custom_fee

	FROM orders AS o
	INNER JOIN delivery AS d ON o.delivery_id = d.id
	INNER JOIN payment AS p ON o.payment_id = p.id
`

	insertOrderQuery = `
	INSERT INTO orders (
			order_uid, track_number, entry,
			delivery_id, payment_id, locale,
			internal_signature, customer_id,
			delivery_service, shardkey,	sm_id,
			date_created, oof_shard
		) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (order_uid)
		DO UPDATE SET order_uid = EXCLUDED.order_uid
		RETURNING id;
	`

	insertDeliveryQuery = `
	INSERT INTO delivery (
			name, phone, zip, city, address, region, email
		) VALUES ( $1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	insertPaymentQuery = `
	INSERT INTO payment (
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
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id;
	`

	insertItemQuery = `
	INSERT INTO items (
		order_id,
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
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id;
	`

	selectItemsWitoutWhereQuery = `
	SELECT
		id,
		order_id,
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
	`
)
