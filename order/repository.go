package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	dbTimeout = 10 * time.Second
)

type Repository interface {
	Close() error
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func (p *postgresRepository) Close() error {
	if p.db == nil {
		return nil
	}
	return p.db.Close()
}

func (p *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	if accountID == "" {
		return nil, fmt.Errorf("get orders for account: account_id cannot be empty")
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	rows, err := p.db.QueryContext(
		dbCtx,
		`SELECT
    		o.id, o.account_id, o.created_at, o.total_price::money::numeric::float8,
    		op.id, op.product_id, op.quantity
		FROM orders o
		LEFT JOIN order_products op ON o.id = op.order_id
		WHERE o.account_id = $1
		ORDER BY o.created_at DESC;`,
		accountID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying orders for account %q: %w", accountID, err)
	}
	defer rows.Close()

	orderMap := make(map[string]*Order)

	for rows.Next() {
		var (
			orderID     string
			accID       string
			createdAt   time.Time
			totalPrice  float64
			productID   sql.NullString
			productRefID sql.NullString
			quantity    sql.NullInt64
		)

		if err := rows.Scan(&orderID, &accID, &createdAt, &totalPrice, &productID, &productRefID, &quantity); err != nil {
			return nil, fmt.Errorf("scanning order row: %w", err)
		}

		order, exists := orderMap[orderID]
		if !exists {
			order = &Order{
				ID:         orderID,
				CreatedAt:  createdAt,
				AccountID:  accID,
				TotalPrice: totalPrice,
				Products:   make([]OrderedProduct, 0),
			}
			orderMap[orderID] = order
		}

		if productID.Valid && productRefID.Valid && quantity.Valid {
			order.Products = append(order.Products, OrderedProduct{
				ID:        productID.String,
				ProductID: productRefID.String,
				Quantity:  uint32(quantity.Int64),
				OrderID:   orderID,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating order rows: %w", err)
	}

	orders := make([]Order, 0, len(orderMap))
	for _, o := range orderMap {
		orders = append(orders, *o)
	}

	return orders, nil
}

func (p *postgresRepository) PutOrder(ctx context.Context, o Order) error {
	if o.ID == "" {
		return fmt.Errorf("put order: order id cannot be empty")
	}
	if o.AccountID == "" {
		return fmt.Errorf("put order: account_id cannot be empty")
	}
	if len(o.Products) == 0 {
		return fmt.Errorf("put order: products cannot be empty")
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	tx, err := p.db.BeginTx(dbCtx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rolling back transaction (original error: %v): %w", err, rbErr)
			}
		}
	}()

	_, err = tx.ExecContext(
		dbCtx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	)
	if err != nil {
		return fmt.Errorf("inserting order: %w", err)
	}

	stmt, err := tx.PrepareContext(
		dbCtx,
		pq.CopyIn("order_products", "order_id", "product_id", "quantity"),
	)
	if err != nil {
		return fmt.Errorf("preparing copy statement: %w", err)
	}
	defer stmt.Close()

	for _, product := range o.Products {
		if product.ProductID == "" {
			return fmt.Errorf("put order: product_id cannot be empty")
		}
		if _, err := stmt.ExecContext(dbCtx, o.ID, product.ProductID, product.Quantity); err != nil {
			return fmt.Errorf("executing copy for product %q: %w", product.ProductID, err)
		}
	}

	if _, err := stmt.ExecContext(dbCtx); err != nil {
		return fmt.Errorf("flushing copy statement: %w", err)
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("closing copy statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func NewOrderRepository(url string) (Repository, error) {
	if url == "" {
		return nil, fmt.Errorf("new order repository: database URL cannot be empty")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &postgresRepository{
		db: db,
	}, nil
}
