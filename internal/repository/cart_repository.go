package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/luzmareto/go-grpc-ecommerce-be/internal/entity"
)

type ICartRepository interface {
	GetCartByProductAndUserId(ctx context.Context, ProductId, userID string) (*entity.UserCart,error)
	CreateNewCart(ctx context.Context, cart *entity.UserCart) error
	UpdateCart(ctx context.Context, cart *entity.UserCart) error
	GetlistCart (ctx context.Context, userId string) ([]*entity.UserCart, error)
	GetCartById (ctx context.Context, cartId string) (*entity.UserCart, error)
	DeleteCart (ctx context.Context, cartId string) error
}

type cartRepository struct {
	db *sql.DB
}

func (cr *cartRepository) GetCartByProductAndUserId(ctx context.Context, ProductId, userId string) (*entity.UserCart,error) {
	row := cr.db.QueryRowContext(
		ctx,
		"SELECT id, product_id, user_id, quantity, created_at, created_by, updated_at, updated_by FROM user_cart WHERE product_id = $1 AND user_id = $2",
		ProductId,
		userId,
	)
	if row.Err() != nil {
		return  nil, row.Err()
	}

	var cartEntity entity.UserCart
	err :=  row.Scan(
		&cartEntity.Id,
		&cartEntity.ProductId,
		&cartEntity.UserId,
		&cartEntity.Quantity,
		&cartEntity.CreatedAt,
		&cartEntity.CreatedBy,
		&cartEntity.UpdatedAt,
		&cartEntity.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, nil
		}
		return nil, err
	}

	return  &cartEntity, nil
}

func (cs *cartRepository) CreateNewCart(ctx context.Context, cart *entity.UserCart) error {
	_, err :=cs.db.ExecContext(
		ctx,
		"INSERT INTO user_cart (id, product_id, user_id, quantity, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		cart.Id,
		cart.ProductId,
		cart.UserId,
		cart.Quantity,
		cart.CreatedAt,
		cart.CreatedBy,
		cart.UpdatedAt,
		cart.UpdatedBy,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cs *cartRepository) UpdateCart(ctx context.Context, cart *entity.UserCart) error{
	_, err :=cs.db.ExecContext(
		ctx,
		"UPDATE user_cart SET product_id = $1, user_id= $2, quantity = $3, updated_at = $4, updated_by = $5 WHERE id = $6",
		cart.ProductId,
		cart.UserId,
		cart.Quantity,
		cart.UpdatedAt,
		cart.UpdatedBy,
		cart.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) GetlistCart (ctx context.Context, userId string) ([]*entity.UserCart, error){
	 rows, err := cr.db.QueryContext(
		ctx,
		"SELECT uc.id, uc.product_id, uc.user_id, uc.quantity , uc.created_at, uc.created_by, uc.updated_at, uc.updated_by, p.id, p.name, p.image_file_name, p.price FROM user_cart uc JOIN product p ON uc.product_id = p.id WHERE uc.user_id = $1 AND p.is_deleted = false",
		userId,
	)
	if err != nil {
		return nil, err
	}

	var carts []*entity.UserCart = make([]*entity.UserCart, 0)
	for rows.Next(){
		var cart entity.UserCart
		cart.Product = &entity.Product{}

		err = rows.Scan(
			&cart.Id,
			&cart.ProductId,
			&cart.UserId,
			&cart.Quantity,
			&cart.CreatedAt,
			&cart.CreatedBy,
			&cart.UpdatedAt,
			&cart.UpdatedBy,
			&cart.Product.Id,
			&cart.Product.Name,
			&cart.Product.ImageFileName,
			&cart.Product.Price,
		)
		if err != nil {
			return nil, err
		}

		carts = append(carts, &cart)
	}

	return  carts, nil
}

func (cr *cartRepository) GetCartById (ctx context.Context, cartId string) (*entity.UserCart, error){
	row := cr.db.QueryRowContext(
		ctx,
		"SELECT id, product_id, user_id, quantity, created_at, created_by, updated_at, updated_by FROM user_cart WHERE id = $1",
		cartId, 
	)
	if row.Err() != nil {
		return  nil, row.Err()
	}

	var cart entity.UserCart
	err := row.Scan(
		&cart.Id,
		&cart.ProductId,
		&cart.UserId,
		&cart.Quantity,
		&cart.CreatedAt,
		&cart.CreatedBy,
		&cart.CreatedAt,
		&cart.CreatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return  nil, nil
		}
		return nil, err
	}

	return  &cart, nil
}
func (cr *cartRepository) DeleteCart (ctx context.Context, cartId string) error{
	_, err := cr.db.ExecContext(
		ctx,
		"DELETE FROM user_cart WHERE id = $1",
		cartId, 
	)
	if err != nil {
		return  err
	}

	return  nil
}


func NewCartRepository(db *sql.DB) ICartRepository  {
	return  &cartRepository{
		db: db,
	}
}