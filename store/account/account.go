package account

import (
	"context"
	"database/sql"

	"github.com/billymosis/marketplace-app/model"
	"github.com/billymosis/marketplace-app/service/auth"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type AccountStore struct {
	db       *sql.DB
	Validate *validator.Validate
}

func NewAccountStore(db *sql.DB, validate *validator.Validate) *AccountStore {
	return &AccountStore{
		db:       db,
		Validate: validate,
	}
}

func (ps *AccountStore) Create(ctx context.Context, account *model.Account) (*model.Account, error) {
	userId, err := auth.GetUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user id")
	}
	query := "INSERT INTO accounts (bank_name, bank_account_name, bank_account_number, user_id) VALUES($1, $2, $3, $4) RETURNING id"
	err = ps.db.QueryRowContext(ctx, query, account.Name, account.AccountName, account.AccountNumber, userId).Scan(&account.Id)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create account")
	}

	return account, nil
}

func (ps *AccountStore) Update(ctx context.Context, account *model.Account) (*model.Account, error) {
	query := "UPDATE accounts SET bank_name=$1, bank_account_name=$2, bank_account_number=$3 WHERE id=$4"
	result, err := ps.db.ExecContext(ctx, query, account.Name, account.AccountName, account.AccountNumber, account.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update account")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return account, nil
}

func (ps *AccountStore) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM accounts WHERE id = $1"
	result, err := ps.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete account")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (ps *AccountStore) Get(ctx context.Context) ([]*model.Account, error) {
	query := "SELECT * FROM accounts"
	rows, err := ps.db.Query(query)
	var accounts []*model.Account

	for rows.Next() {
		var product model.Account
		if err := rows.Scan(&product.Id, &product.Name, &product.AccountName, &product.AccountNumber, &product.UserId); err != nil {
			return nil, errors.Wrap(err, "failed to scan account data")
		}
		accounts = append(accounts, &product)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete account")
	}
	if accounts == nil {
		accounts = []*model.Account{}
	}

	return accounts, nil
}