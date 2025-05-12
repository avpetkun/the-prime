package dbx

import (
	"context"

	"github.com/avpetkun/the-prime/internal/common"
)

func (db *DB) SaveProduct(ctx context.Context, p *common.Product) error {
	if p.ID == 0 {
		const q = `
			INSERT INTO products (
				"name", "type", amount, price, badge
			) VALUES ($1,$2,$3,$4,$5)
			RETURNING id, created_at
		`
		return db.c.QueryRow(ctx, q, p.Name, p.Type, p.Amount, p.Price, p.Badge).Scan(&p.ID, &p.CreatedAt)
	}
	const q = `
		UPDATE products SET
			"name" = $1, "type" = $2,
			amount = $3, price = $4, badge = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`
	return db.c.QueryRow(ctx, q, p.Name, p.Type, p.Amount, p.Price, p.Badge, p.ID).Scan(&p.UpdatedAt)
}

func (db *DB) DeleteProduct(ctx context.Context, productID int64) error {
	const q = `UPDATE products SET deleted_at = NOW() WHERE id = $1`
	_, err := db.c.Exec(ctx, q, productID)
	return err
}

func (db *DB) GetProduct(ctx context.Context, productID int64) (*common.Product, error) {
	const q = `
		SELECT
			created_at, updated_at, deleted_at,
			"name", "type",
			amount, price, badge
		FROM products
		WHERE id = $1
	`
	p := common.Product{ID: productID}
	err := db.c.QueryRow(ctx, q, productID).Scan(
		&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		&p.Name, &p.Type,
		&p.Amount, &p.Price, &p.Badge,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (db *DB) GetAllProducts(ctx context.Context) (products []*common.Product, err error) {
	const q = `
		SELECT
			id, created_at, updated_at,
			"name", "type",
			amount, price, badge
		FROM products
		WHERE deleted_at is null
		ORDER BY price ASC
	`
	rows, err := db.c.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p common.Product
		err = rows.Scan(
			&p.ID, &p.CreatedAt, &p.UpdatedAt,
			&p.Name, &p.Type,
			&p.Amount, &p.Price, &p.Badge,
		)
		if err != nil {
			break
		}
		products = append(products, &p)
	}
	if products == nil {
		products = []*common.Product{}
	}
	return products, err
}
