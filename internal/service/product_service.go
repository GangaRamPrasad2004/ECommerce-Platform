package service

import (
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/dto"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/repositories"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/utils"
)

var _ ProductServiceInterface = (*ProductService)(nil)

// ProductService manages product catalog data.
type ProductService struct {
	productRepo repositories.ProductRepositoryInterface
}

// NewProductService creates a ProductService backed by a product repository.
func NewProductService(productRepo repositories.ProductRepositoryInterface) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// CreateCategory creates a new product category from req.
func (s *ProductService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	return s.productRepo.CreateCategory(req)
}

// GetCategories returns all active product categories.
func (s *ProductService) GetCategories() ([]dto.CategoryResponse, error) {
	return s.productRepo.GetCategories()
}

// UpdateCategory updates the category identified by id and returns it.
func (s *ProductService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	return s.productRepo.UpdateCategory(id, req)
}

// DeleteCategory removes the category identified by id.
func (s *ProductService) DeleteCategory(id uint) error {
	return s.productRepo.DeleteCategory(id)
}

// CreateProduct creates a product from req and returns its full details.
func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	return s.productRepo.CreateProduct(req)
}

// GetProducts returns a paginated list of active products.
func (s *ProductService) GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	return s.productRepo.GetProducts(page, limit)
}

// GetProduct retrieves the product identified by id with its category and images.
func (s *ProductService) GetProduct(id uint) (*dto.ProductResponse, error) {
	return s.productRepo.GetProduct(id)
}

// UpdateProduct updates the product identified by id and returns its full details.
func (s *ProductService) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	return s.productRepo.UpdateProduct(id, req)
}

// DeleteProduct removes the product identified by id.
func (s *ProductService) DeleteProduct(id uint) error {
	return s.productRepo.DeleteProduct(id)
}

// AddProductImage associates an image with a product.
func (s *ProductService) AddProductImage(productID uint, url, altText string) error {
	return s.productRepo.AddProductImage(productID, url, altText)
}

// SearchProducts returns products matching the supplied search criteria.
func (s *ProductService) SearchProducts(req *dto.SearchProductsRequest) ([]dto.ProductSearchResult, *utils.PaginationMeta, error) {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 10
	}

	return s.productRepo.SearchProducts(req)
}
