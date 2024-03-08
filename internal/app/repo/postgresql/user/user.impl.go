package user

import (
	"context"
	"fmt"
	"time"

	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"github.com/fjrid/parking/pkg/dbtxn"
	"github.com/fjrid/parking/pkg/sqkit"

	sq "github.com/Masterminds/squirrel"
)

// @ctor
func NewUserRepository(impl UserRepositoryImpl) UserRepository {
	return &impl
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, data dto.CreateUserRequest) (res entity.User, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	res = entity.User{
		Email:      data.Email,
		Password:   data.Password,
		Role:       data.Role,
		CreatedAt:  time.Now(),
		CreatedBy:  data.CreatedBy,
		ModifiedAt: time.Now(),
		ModifiedBy: data.CreatedBy,
	}

	builder := sq.
		Insert(entity.UserTableName).
		Columns(
			entity.UserTable.Email,
			entity.UserTable.Password,
			entity.UserTable.Role,
			entity.UserTable.CreatedBy,
			entity.UserTable.ModifiedBy,
		).
		Suffix(
			fmt.Sprintf("RETURNING \"%s\"", entity.UserTable.ID),
		).
		PlaceholderFormat(sq.Dollar).
		Values(
			res.Email,
			res.Password,
			res.Role,
			res.CreatedBy,
			res.ModifiedBy,
		)

	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	if err = scanner.Scan(&res.ID); err != nil {
		txn.AppendError(err)
		return
	}

	return
}

func (r *UserRepositoryImpl) Find(ctx context.Context, opts ...sqkit.SelectOption) (results []entity.User, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	builder := sq.
		Select(
			"*",
		).
		From(entity.UserTableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	results = make([]entity.User, 0)
	for rows.Next() {
		user := entity.User{}

		if err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.CreatedBy,
			&user.ModifiedAt,
			&user.ModifiedBy,
			&user.DeletedAt,
			&user.DeletedBy,
		); err != nil {
			return
		}

		results = append(results, user)
	}

	return
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (results []entity.User, err error) {
	condition := []sqkit.SelectOption{
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.UserTableName, entity.UserTable.Email): email},
	}

	return r.Find(ctx, condition...)
}
