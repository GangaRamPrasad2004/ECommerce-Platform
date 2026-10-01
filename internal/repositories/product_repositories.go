package repositories

import (
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/dto"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/utils"
	"gorm.io/gorm"
)

var _ ProductRepositoryInterface = (*ProductRepository)(nil)

// ProductRepository stores and retrieves products and categories.
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a product repository backed by db.
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// CreateCategory creates a category from req.
func (repo *ProductRepository) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := repo.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return categoryToResponse(&category), nil
}

// GetCategories retrieves all active categories.
func (repo *ProductRepository) GetCategories() ([]dto.CategoryResponse, error) {
	var categories []models.Category
	if err := repo.db.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		return nil, err
	}

	response := make([]dto.CategoryResponse, len(categories))
	for i := range categories {
		response[i] = *categoryToResponse(&categories[i])
	}
	return response, nil
}

// UpdateCategory updates the category with the specified ID.
func (repo *ProductRepository) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	var category models.Category
	if err := repo.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if err := repo.db.Save(&category).Error; err != nil {
		return nil, err
	}
	return categoryToResponse(&category), nil
}

// DeleteCategory removes the category with the specified ID.
func (repo *ProductRepository) DeleteCategory(id uint) error {
	return repo.db.Delete(&models.Category{}, id).Error
}

// CreateProduct creates a product from req.
func (repo *ProductRepository) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}
	if err := repo.db.Create(&product).Error; err != nil {
		return nil, err
	}
	return repo.GetProduct(product.ID)
}

// GetProducts retrieves a page of active products and its pagination metadata.
func (repo *ProductRepository) GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	if err := repo.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, nil, err
	}

	var products []models.Product
	if err := repo.db.Preload("Category").Preload("Images").
		Where("is_active = ?", true).
		Offset(offset).Limit(limit).
		Find(&products).Error; err != nil {
		return nil, nil, err
	}

	response := make([]dto.ProductResponse, len(products))
	for i := range products {
		response[i] = productToResponse(&products[i])
	}

	return response, paginationMeta(page, limit, total), nil
}

// GetProduct retrieves the product with the specified ID.
func (repo *ProductRepository) GetProduct(id uint) (*dto.ProductResponse, error) {
	var product models.Product
	if err := repo.db.Preload("Category").Preload("Images").First(&product, id).Error; err != nil {
		return nil, err
	}
	response := productToResponse(&product)
	return &response, nil
}

// UpdateProduct updates the product with the specified ID.
func (repo *ProductRepository) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	var product models.Product
	if err := repo.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if err := repo.db.Save(&product).Error; err != nil {
		return nil, err
	}
	return repo.GetProduct(id)
}

// DeleteProduct removes the product with the specified ID.
func (repo *ProductRepository) DeleteProduct(id uint) error {
	return repo.db.Delete(&models.Product{}, id).Error
}

// AddProductImage associates an image with a product.
func (repo *ProductRepository) AddProductImage(productID uint, url, altText string) error {
	var count int64
	if err := repo.db.Model(&models.ProductImage{}).Where("product_id = ?", productID).Count(&count).Error; err != nil {
		return err
	}

	image := models.ProductImage{
		ProductID: productID,
		URL:       url,
		AltText:   altText,
		IsPrimary: count == 0,
	}
	return repo.db.Create(&image).Error
}

// SearchProducts searches active products and returns ranked paginated results.
func (repo *ProductRepository) SearchProducts(req *dto.SearchProductsRequest) ([]dto.ProductSearchResult, *utils.PaginationMeta, error) {
	page, limit := req.Page, req.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	searchQuery := func() *gorm.DB {
		query := repo.db.Model(&models.Product{}).
			Where("search_vector @@ plainto_tsquery('english', ?)", req.Query).
			Where("is_active = ?", true)
		if req.CategoryID != nil {
			query = query.Where("category_id = ?", *req.CategoryID)
		}
		if req.MinPrice != nil {
			query = query.Where("price >= ?", *req.MinPrice)
		}
		if req.MaxPrice != nil {
			query = query.Where("price <= ?", *req.MaxPrice)
		}
		return query
	}

	var total int64
	if err := searchQuery().Count(&total).Error; err != nil {
		return nil, nil, err
	}

	type productWithRank struct {
		models.Product
		Rank float32 `gorm:"column:rank"`
	}
	var rows []productWithRank
	if err := searchQuery().
		Select("products.*, ts_rank(search_vector, plainto_tsquery('english', ?)) as rank", req.Query).
		Order("rank DESC, created_at DESC").
		Preload("Category").
		Preload("Images").
		Offset(offset).
		Limit(limit).
		Find(&rows).Error; err != nil {
		return nil, nil, err
	}

	results := make([]dto.ProductSearchResult, len(rows))
	for i := range rows {
		results[i] = dto.ProductSearchResult{
			ProductResponse: productToResponse(&rows[i].Product),
			Rank:            rows[i].Rank,
		}
	}
	return results, paginationMeta(page, limit, total), nil
}

func categoryToResponse(category *models.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func productToResponse(product *models.Product) dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range product.Images {
		images[i] = dto.ProductImageResponse{
			ID:        product.Images[i].ID,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
			CreatedAt: product.Images[i].CreatedAt,
		}
	}

	return dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Category:    *categoryToResponse(&product.Category),
		Images:      images,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func paginationMeta(page, limit int, total int64) *utils.PaginationMeta {
	return &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}
}
