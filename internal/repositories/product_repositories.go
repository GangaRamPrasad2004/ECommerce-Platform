package repositories

import (
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"gorm.io/gorm"
)

var _ ProductRepositoryInterface = (*ProductRepository)(nil)

// ProductRepository stores and retrieves products and categories.
type ProductRepository struct {
	db *gorm.DB
}

// ProductWithRank represents a product and its full-text search rank.
type ProductWithRank struct {
	models.Product
	Rank float32 `gorm:"column:rank"`
}

// NewProductRepository creates a product repository backed by db.
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// CreateCategory persists a category.
func (repo *ProductRepository) CreateCategory(category *models.Category) error {
	return repo.db.Create(category).Error
}

// GetActiveCategories retrieves all active categories.
func (repo *ProductRepository) GetActiveCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := repo.db.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID retrieves a category by its ID.
func (repo *ProductRepository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := repo.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory persists changes to a category.
func (repo *ProductRepository) UpdateCategory(category *models.Category) error {
	return repo.db.Save(category).Error
}

// DeleteCategory removes the category with the specified ID.
func (repo *ProductRepository) DeleteCategory(id uint) error {
	return repo.db.Delete(&models.Category{}, id).Error
}

// CreateProduct persists a product.
func (repo *ProductRepository) CreateProduct(product *models.Product) error {
	return repo.db.Create(product).Error
}

// GetProducts retrieves a page of active products and its total count.
func (repo *ProductRepository) GetProducts(offset, limit int) ([]models.Product, int64, error) {
	var total int64
	if err := repo.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []models.Product
	if err := repo.db.Preload("Category").Preload("Images").
		Where("is_active = ?", true).
		Offset(offset).Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

// GetProductByID retrieves the product with the specified ID.
func (repo *ProductRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := repo.db.Preload("Category").Preload("Images").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateProduct persists changes to a product.
func (repo *ProductRepository) UpdateProduct(product *models.Product) error {
	return repo.db.Save(product).Error
}

// DeleteProduct removes the product with the specified ID.
func (repo *ProductRepository) DeleteProduct(id uint) error {
	return repo.db.Delete(&models.Product{}, id).Error
}

// CountProductImages returns the number of images associated with a product.
func (repo *ProductRepository) CountProductImages(productID uint) (int64, error) {
	var count int64
	if err := repo.db.Model(&models.ProductImage{}).Where("product_id = ?", productID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CreateProductImage persists an image associated with a product.
func (repo *ProductRepository) CreateProductImage(image *models.ProductImage) error {
	return repo.db.Create(image).Error
}

// SearchProducts searches active products and returns ranked results and their total count.
func (repo *ProductRepository) SearchProducts(query string, categoryID *uint, minPrice, maxPrice *float64, offset, limit int) ([]ProductWithRank, int64, error) {
	searchQuery := func() *gorm.DB {
		query := repo.db.Model(&models.Product{}).
			Where("search_vector @@ plainto_tsquery('english', ?)", query).
			Where("is_active = ?", true)
		if categoryID != nil {
			query = query.Where("category_id = ?", *categoryID)
		}
		if minPrice != nil {
			query = query.Where("price >= ?", *minPrice)
		}
		if maxPrice != nil {
			query = query.Where("price <= ?", *maxPrice)
		}
		return query
	}

	var total int64
	if err := searchQuery().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []ProductWithRank
	if err := searchQuery().
		Select("products.*, ts_rank(search_vector, plainto_tsquery('english', ?)) as rank", query).
		Order("rank DESC, created_at DESC").
		Preload("Category").
		Preload("Images").
		Offset(offset).
		Limit(limit).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
