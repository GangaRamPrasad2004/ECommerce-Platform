package service

import (
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/dto"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
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
	category := &models.Category{Name: req.Name, Description: req.Description}
	if err := s.productRepo.CreateCategory(category); err != nil {
		return nil, err
	}
	return categoryToDTO(category), nil
}

// GetCategories returns all active product categories.
func (s *ProductService) GetCategories() ([]dto.CategoryResponse, error) {
	categories, err := s.productRepo.GetActiveCategories()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.CategoryResponse, len(categories))
	for i := range categories {
		responses[i] = *categoryToDTO(&categories[i])
	}
	return responses, nil
}

// UpdateCategory updates the category identified by id and returns it.
func (s *ProductService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.productRepo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if err := s.productRepo.UpdateCategory(category); err != nil {
		return nil, err
	}
	return categoryToDTO(category), nil
}

// DeleteCategory removes the category identified by id.
func (s *ProductService) DeleteCategory(id uint) error {
	return s.productRepo.DeleteCategory(id)
}

// CreateProduct creates a product from req and returns its full details.
func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}
	if err := s.productRepo.CreateProduct(product); err != nil {
		return nil, err
	}
	return s.GetProduct(product.ID)
}

// GetProducts returns a paginated list of active products.
func (s *ProductService) GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	products, total, err := s.productRepo.GetProducts(offset, limit)
	if err != nil {
		return nil, nil, err
	}
	responses := make([]dto.ProductResponse, len(products))
	for i := range products {
		responses[i] = *productToDTO(&products[i])
	}
	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}
	return responses, meta, nil
}

// GetProduct retrieves the product identified by id with its category and images.
func (s *ProductService) GetProduct(id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return productToDTO(product), nil
}

// UpdateProduct updates the product identified by id and returns its full details.
func (s *ProductService) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
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
	if err := s.productRepo.UpdateProduct(product); err != nil {
		return nil, err
	}
	return s.GetProduct(id)
}

// DeleteProduct removes the product identified by id.
func (s *ProductService) DeleteProduct(id uint) error {
	return s.productRepo.DeleteProduct(id)
}

// AddProductImage associates an image with a product.
func (s *ProductService) AddProductImage(productID uint, url, altText string) error {
	count, err := s.productRepo.CountProductImages(productID)
	if err != nil {
		return err
	}
	image := &models.ProductImage{
		ProductID: productID,
		URL:       url,
		AltText:   altText,
		IsPrimary: count == 0,
	}
	return s.productRepo.CreateProductImage(image)
}

// SearchProducts returns products matching the supplied search criteria.
func (s *ProductService) SearchProducts(req *dto.SearchProductsRequest) ([]dto.ProductSearchResult, *utils.PaginationMeta, error) {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit
	products, total, err := s.productRepo.SearchProducts(
		req.Query,
		req.CategoryID,
		req.MinPrice,
		req.MaxPrice,
		offset,
		req.Limit,
	)
	if err != nil {
		return nil, nil, err
	}
	results := make([]dto.ProductSearchResult, len(products))
	for i := range products {
		results[i] = dto.ProductSearchResult{
			ProductResponse: *productToDTO(&products[i].Product),
			Rank:            products[i].Rank,
		}
	}
	meta := &utils.PaginationMeta{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      total,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}
	return results, meta, nil
}

func categoryToDTO(category *models.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func productToDTO(product *models.Product) *dto.ProductResponse {
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

	return &dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Category:    *categoryToDTO(&product.Category),
		Images:      images,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
