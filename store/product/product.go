package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/billymosis/marketplace-app/model"
	"github.com/billymosis/marketplace-app/service/auth"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ProductStore struct {
	db       *pgxpool.Pool
	Validate *validator.Validate
}

func NewProductStore(db *pgxpool.Pool, validate *validator.Validate) *ProductStore {
	return &ProductStore{
		db:       db,
		Validate: validate,
	}
}

func (ps *ProductStore) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	tagsJSON, err := json.Marshal(product.Tags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal tags to JSON")
	}
	userId, err := auth.GetUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user id")
	}
	query := "INSERT INTO products (name, price, image_url, stock, condition, tags, is_purchasable, user_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	err = ps.db.QueryRow(ctx, query, product.Name, product.Price, product.ImageUrl, product.Stock, product.Condition, tagsJSON, product.IsPurchasable, userId).Scan(&product.Id)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create product")
	}

	return product, nil
}

func (ps *ProductStore) UpdateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	tagsJSON, err := json.Marshal(product.Tags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal tags to JSON")
	}

	query := "UPDATE products SET name=$1, price=$2, image_url=$3, stock=$4, condition=$5, tags=$6, is_purchasable=$7 WHERE id=$8"
	result, err := ps.db.Exec(ctx, query, product.Name, product.Price, product.ImageUrl, product.Stock, product.Condition, tagsJSON, product.IsPurchasable, product.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update product")
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return product, nil
}

func (ps *ProductStore) UpdateProductStock(ctx context.Context, product *model.Product) (*model.Product, error) {
	query := "UPDATE products SET stock=$1 WHERE id=$2"
	result, err := ps.db.Exec(ctx, query, product.Stock, product.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update product")
	}

	rowsAffected  := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return product, nil
}

func (ps *ProductStore) Payment(ctx context.Context, payment *model.Payment) error {

	query := "INSERT INTO payments (account_id, product_id, payment_proof_image_url, quantity) VALUES($1, $2, $3, $4) RETURNING id"
	err := ps.db.QueryRow(ctx, query, payment.AccountId, payment.ProductId, payment.PaymentProofImageUrl, payment.Quantity).Scan(&payment.Id)

	if err != nil {
		return errors.Wrap(err, "failed to create payment")
	} else {
		logrus.Printf("TAMBHAHAHHA")
		query = `
			UPDATE products
			SET purchase_count = purchase_count + 1
			WHERE id = $1;
		`
		_, err := ps.db.Exec(ctx, query, payment.ProductId)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (ps *ProductStore) DeleteProduct(ctx context.Context, id uint) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := ps.db.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete product")
	}

	rowsAffected  := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

func (ps *ProductStore) GetProductById(ctx context.Context, id uint) (*model.Product, error) {
	query := "SELECT * FROM products WHERE id = $1"

	rows := ps.db.QueryRow(ctx, query, id)
	var product model.Product
	var tagsJSON []byte
	if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.ImageUrl, &product.Stock, &product.Condition, &tagsJSON, &product.IsPurchasable, &product.PurchaseCount, &product.UserId); err != nil {
		return nil, errors.Wrap(err, "failed to scan product data")
	}
	if err := json.Unmarshal(tagsJSON, &product.Tags); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal tags JSON")
	}
	return &product, nil
}

func (ps *ProductStore) GetTotalSold(ctx context.Context, productId uint) (int, error) {
	query := `
		SELECT
		    SUM(quantity) AS total_quantity
		FROM
		    payments
		WHERE
		    product_id = $1;
	`

	rows := ps.db.QueryRow(ctx, query, productId)
	var total int
	if err := rows.Scan(&total); err != nil {
		return 0, errors.Wrap(err, "failed to scan product data")
	}
	return total, nil
}

type Query struct {
	b      strings.Builder
	params []interface{}
}

func (q *Query) Query(s string) {
	q.b.WriteString(s)
}

func (q *Query) Param(val interface{}) {
	length := len(q.params)
	q.b.WriteString("$" + strconv.Itoa(length+1))
	q.params = append(q.params, val)
}

func (q *Query) Get() (string, []interface{}) {
	return q.b.String(), q.params
}

type Meta struct {
	Limit  int
	Offset int
	Total  int
}

func (ps *ProductStore) GetProducts(ctx context.Context, queryParams url.Values) ([]*model.Product, Meta, error) {
	var meta Meta
	userId, err := auth.GetUserId(ctx)
	q := Query{}
	q.Query("SELECT * FROM products WHERE")
	userOnlyStr := queryParams.Get("userOnly")
	if err != nil {
		userOnlyStr = "false"
	}
	tags := queryParams["tags"]
	condition := queryParams.Get("condition")
	showEmptyStockStr := queryParams.Get("showEmptyStock")
	maxPriceStr := queryParams.Get("maxPrice")
	minPriceStr := queryParams.Get("minPrice")
	search := queryParams.Get("search")
	hasParams := false
	userOnly, err := strconv.ParseBool(userOnlyStr)
	if err != nil {
		userOnly = false
	}
	if userOnly || len(tags) > 0 || condition != "" || showEmptyStockStr != "" || maxPriceStr != "" || minPriceStr != "" || search != "" {
		hasParams = true
	}

	if userOnly {
		if userOnly {
			q.Query(" AND user_id = ")
			q.Param(userId)
		}
	}

	if condition != "" {
		if strings.ToLower(condition) == "new" {
			q.Query(" AND condition = ")
			q.Param("new")
		}
		if strings.ToLower(condition) == "second" {
			q.Query(" AND condition = ")
			q.Param("second")
		}
	}

	if showEmptyStockStr != "" {
		showEmptyStock, err := strconv.ParseBool(showEmptyStockStr)
		if err != nil {
			return nil, meta, err
		}
		if showEmptyStock {
			q.Query(" AND stock > ")
			q.Param(-1)
		} else {
			q.Query(" AND stock > ")
			q.Param(0)
		}

	}

	if maxPriceStr != "" {
		maxPrice, err := strconv.Atoi(maxPriceStr)
		if err != nil {
			return nil, meta, err
		}
		q.Query(" AND price < ")
		q.Param(maxPrice)
	}

	if minPriceStr != "" {
		minPrice, err := strconv.Atoi(minPriceStr)
		if err != nil {
			return nil, meta, err
		}
		q.Query(" AND price > ")
		q.Param(minPrice)
	}
	if len(tags) > 0 {
		for _, element := range tags {
			q.Query(" AND tags @> ")
			q.Param("[\"" + element + "\"]")
		}

	}

	if search != "" {
		q.Query(" AND name LIKE ")
		q.Param("%" + search + "%")
	}

	limitStr := queryParams.Get("limit")
	limit := 10
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, meta, err
		}
	}
	q.Query(" LIMIT ")
	q.Param(limit)

	offsetStr := queryParams.Get("offset")
	offset := 0
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return nil, meta, err
		}
	}
	q.Query(" OFFSET ")
	q.Param(offset)

	query, params := q.Get()
	logrus.Printf("BEFORE: %+v\n", query)
	if !hasParams {
		query = strings.Replace(query, "WHERE", "", 1)
	} else {
		if strings.Count(query, "AND") == 0 {
			query = strings.Replace(query, "WHERE", "", 1)
		}
		query = strings.Replace(query, "WHERE AND", "WHERE", 1)
	}
	logrus.Printf("%+v\n", query)
	logrus.Printf("%+v\n", params)

	rows, err := ps.db.Query(ctx, query, params...)
	if err != nil {
		return nil, meta, errors.Wrap(err, "failed to get product")
	}
	var products []*model.Product
	for rows.Next() {
		var product model.Product
		var tagsJSON []byte
		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.ImageUrl, &product.Stock, &product.Condition, &tagsJSON, &product.IsPurchasable, &product.PurchaseCount, &product.UserId); err != nil {
			return nil, meta, errors.Wrap(err, "failed to scan product data")
		}
		if err := json.Unmarshal(tagsJSON, &product.Tags); err != nil {
			return nil, meta, errors.Wrap(err, "failed to unmarshal tags JSON")
		}
		products = append(products, &product)
	}

	countQuery := strings.Replace(query, "SELECT *", "SELECT COUNT(*)", 1)
	countQuery = strings.Split(countQuery, "LIMIT")[0]
	params = params[:len(params)-2]
	var count int
	err = ps.db.QueryRow(ctx, countQuery, params...).Scan(&count)

	if err != nil {
		return nil, meta, errors.Wrap(err, "failed to scan product data")
	}
	meta = Meta{
		Limit:  limit,
		Offset: offset,
		Total:  count,
	}

	logrus.Printf("TOTAL LEN: %+v\n", count)
	if products == nil {
		products = []*model.Product{}
	}
	return products, meta, nil
}
