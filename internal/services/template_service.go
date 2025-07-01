package services

import (
	"gobi/internal/models"
	"gobi/pkg/errors"
)

// TemplateService handles template-related business logic
type TemplateService struct {
	templateRepo      TemplateRepository
	permissionService PermissionService
}

// NewTemplateService creates a new TemplateService instance
func NewTemplateService(
	templateRepo TemplateRepository,
	permissionService PermissionService,
) *TemplateService {
	return &TemplateService{
		templateRepo:      templateRepo,
		permissionService: permissionService,
	}
}

// CreateTemplate creates a new template
func (s *TemplateService) CreateTemplate(template *models.ExcelTemplate, userID uint) error {
	template.UserID = userID

	if err := s.templateRepo.Create(template); err != nil {
		return errors.WrapError(err, "Could not create template")
	}

	return nil
}

// ListTemplates retrieves templates based on user permissions
func (s *TemplateService) ListTemplates(userID uint, isAdmin bool) ([]models.ExcelTemplate, error) {
	templates, err := s.templateRepo.FindByUser(userID, isAdmin)
	if err != nil {
		return nil, errors.WrapError(err, "Could not fetch templates")
	}

	// Do not return template content in list view
	for i := range templates {
		templates[i].Template = nil
	}

	return templates, nil
}

// GetTemplate retrieves a specific template
func (s *TemplateService) GetTemplate(templateID uint, userID uint, isAdmin bool) (*models.ExcelTemplate, error) {
	template, err := s.templateRepo.FindByID(templateID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	if !s.permissionService.CanAccess(userID, templateID, "template", isAdmin) {
		return nil, errors.ErrForbidden
	}

	// Do not return template content in detail view
	template.Template = nil

	return template, nil
}

// UpdateTemplate updates a template
func (s *TemplateService) UpdateTemplate(templateID uint, updates *models.ExcelTemplate, userID uint, isAdmin bool) (*models.ExcelTemplate, error) {
	template, err := s.templateRepo.FindByID(templateID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	if !s.permissionService.CanAccess(userID, templateID, "template", isAdmin) {
		return nil, errors.ErrForbidden
	}

	// Update allowed fields
	if updates.Name != "" {
		template.Name = updates.Name
	}
	if updates.Description != "" {
		template.Description = updates.Description
	}

	if err := s.templateRepo.Update(template); err != nil {
		return nil, errors.WrapError(err, "Could not update template")
	}

	// Do not return template content
	template.Template = nil

	return template, nil
}

// DeleteTemplate deletes a template
func (s *TemplateService) DeleteTemplate(templateID uint, userID uint, isAdmin bool) error {
	_, err := s.templateRepo.FindByID(templateID)
	if err != nil {
		return errors.ErrNotFound
	}

	if !s.permissionService.CanAccess(userID, templateID, "template", isAdmin) {
		return errors.ErrForbidden
	}

	if err := s.templateRepo.Delete(templateID); err != nil {
		return errors.WrapError(err, "Could not delete template")
	}

	return nil
}

// DownloadTemplate retrieves a template for download
func (s *TemplateService) DownloadTemplate(templateID uint, userID uint, isAdmin bool) (*models.ExcelTemplate, error) {
	template, err := s.templateRepo.FindByID(templateID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	if !s.permissionService.CanAccess(userID, templateID, "template", isAdmin) {
		return nil, errors.ErrForbidden
	}

	return template, nil
}

// GetDashboardStats retrieves dashboard statistics
func (s *TemplateService) GetDashboardStats() (map[string]interface{}, error) {
	stats, err := s.templateRepo.GetStats()
	if err != nil {
		return nil, errors.WrapError(err, "Could not fetch dashboard stats")
	}

	return stats, nil
}
