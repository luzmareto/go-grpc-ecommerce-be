package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/luzmareto/go-grpc-ecommerce-be/internal/entity"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/common"
)

type IProductRepository interface {
	CreateNewProduct(ctx context.Context, product *entity.Product) error
	GetProductById(ctx context.Context, id string) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error
	GetProductsPagination(ctx context.Context, pagination *common.PaginationRequest) ([]*entity.Product, *common.PaginationResponse, error)
}

type productRepository struct {
	db *sql.DB
}

func (repo *productRepository) CreateNewProduct(ctx context.Context, product *entity.Product) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO product (id, name, description, price, image_file_name, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, is_deleted) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
		product.Id,
		product.Name,
		product.Description,
		product.Price,
		product.ImageFileName,
		product.CreatedAt,
		product.CreatedBy,
		product.UpdatedAt,
		product.UpdatedBy,
		product.DeletedAt,
		product.DeletedBy,
		product.IsDeleted,
	)

	if err != nil {
		return err
	}
	return nil
}

func UUIDOrNil(s string) any {
	if _, err := uuid.Parse(s); err != nil {
		return nil
	}
	return s
}

func (repo *productRepository) GetProductById(ctx context.Context, id string) (*entity.Product, error) {
	idParam := UUIDOrNil(id)
	var productEntity entity.Product
	row := repo.db.QueryRowContext(
		ctx,
		"SELECT id, name, description, price, image_file_name FROM product WHERE id = $1 AND is_deleted = false",
		idParam,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&productEntity.Id,
		&productEntity.Name,
		&productEntity.Description,
		&productEntity.Price,
		&productEntity.ImageFileName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	return &productEntity, nil
}

func (repo *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) error {
	_, err := repo.db.ExecContext(
		ctx,
		"UPDATE product SET name=$1, description=$2, price=$3, image_file_name=$4, updated_at=$5, updated_by=$6 WHERE id =$7",
		product.Name,
		product.Description,
		product.Price,
		product.ImageFileName,
		product.UpdatedAt,
		product.UpdatedBy,
		product.Id,
	)

	if err != nil {
		return err
	}
	return nil
}

func (repo *productRepository) DeleteProduct(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error {
	_, err := repo.db.ExecContext(
		ctx,
		"UPDATE product SET deleted_at = $1, deleted_by = $2, is_deleted = true WHERE id = $3",
		deletedAt,
		deletedBy,
		id,
	)

	if err != nil {
		return err
	}
	return nil
}

func (repo *productRepository) GetProductsPagination(ctx context.Context, pagination *common.PaginationRequest) ([]*entity.Product, *common.PaginationResponse, error) {
 	row := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM product WHERE is_deleted = false")
	if row.Err() != nil {
		return  nil, nil, row.Err()
	}

	var totalCount int
	err :=  row.Scan(&totalCount)
	if err != nil {
		return  nil, nil, err
	}


	offset := (pagination.CurrentPage -1) * pagination.ItemPerPage
	totalPages := (totalCount + int(pagination.ItemPerPage) - 1) / int(pagination.ItemPerPage)
	
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, name, description, price, image_file_name FROM product WHERE is_deleted = false ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		pagination.ItemPerPage,
		offset,

	)

	if err != nil {
		return nil, nil, err
	}	

	var products []*entity.Product = make([]*entity.Product, 0)
	for rows.Next() {
		var product entity.Product

		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.ImageFileName,
		)
		if err != nil {
		return nil, nil, err
		}

		products = append(products, &product)
	}


	paginationResponse := &common.PaginationResponse{
		CurrentPage: pagination.CurrentPage,
		ItemPerPage: pagination.ItemPerPage,
		TotalItemCount: int32(totalCount),
		TotalPageCount: int32(totalPages),
	}
	return products,paginationResponse,nil
}

func NewProductRepository(db *sql.DB) IProductRepository {
	return &productRepository{
		db: db,
	}
}
