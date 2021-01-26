package ResourceService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Resource"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// ResourceService struct
type ResourceService struct {
	ResourceRepo Resource.ResourceRepo
}

// NewResourceService function
func NewResourceService(ResourceRepo Resource.ResourceRepo) Resource.ResourceService {
	return &ResourceService{
		ResourceRepo: ResourceRepo,
	}
}

// CreateResource (resource *entity.Resource) (*entity.Resource, error)
func (resser *ResourceService) CreateResource(resource *entity.Resource) *entity.Resource {
	resources, erra := resser.ResourceRepo.CreateResource(resource)
	if erra != nil {
		return nil
	}
	return resources
}

// DeleteResource ( ResourceID uint) error
func (resser *ResourceService) DeleteResource(ResourceID uint) bool {
	erra := resser.ResourceRepo.DeleteResource(ResourceID)
	if erra != nil {
		return false
	}
	return true
}

// GetResourceByID method
func (resser *ResourceService) GetResourceByID(ResourceID uint) *entity.Resource {
	resources, erra := resser.ResourceRepo.GetResourceByID(ResourceID)
	if erra != nil {
		return nil
	}
	return resources
}

// GetResources method
func (resser *ResourceService) GetResources(Offset, Limit uint) *[]entity.Resource {
	resources, erra := resser.ResourceRepo.GetResources(Offset, Limit)
	if erra != nil {
		return nil
	}
	return resources
}

// GetVideos method
func (resser *ResourceService) GetVideos(Offset uint, Limit uint) *[]entity.Resource {
	resources, erra := resser.ResourceRepo.GetVideos(Offset, Limit)
	if erra != nil {
		return nil
	}
	return resources
}

// GetAudios method
func (resser *ResourceService) GetAudios(Offset, Limit uint) *[]entity.Resource {
	resources, erra := resser.ResourceRepo.GetAudios(Offset, Limit)
	if erra != nil {
		return nil
	}
	return resources
}

// GetFiles mthod
func (resser *ResourceService) GetFiles(Offset uint, Limit uint) *[]entity.Resource {
	resources, erra := resser.ResourceRepo.GetFiles(Offset, Limit)
	if erra != nil {
		return nil
	}
	return resources
}

// GetPicture method
func (resser *ResourceService) GetPicture(Offset uint, Limit uint) *[]entity.Resource {
	resources, newErros := resser.ResourceRepo.GetPicture(Offset, Limit)
	if newErros != nil {
		return nil
	}
	return resources
}

// GetPDF mthod
func (resser *ResourceService) GetPDF(Offset uint, Limit uint) *[]entity.Resource {
	resources, newErra := resser.ResourceRepo.GetPDF(Offset, Limit)
	if newErra != nil {
		return nil
	}
	return resources
}

// SearchResource method to search a student
func (resser *ResourceService) SearchResource(query string, typo uint) []entity.Resource {
	res, era := resser.ResourceRepo.SearchResource(query, typo)
	if era != nil {
		return []entity.Resource{}
	}
	return res
}

// GetRandomActiveResource  () (*entity.Resource, error)
func (resser *ResourceService) GetRandomActiveResource() *entity.Resource {
	resource, era := resser.ResourceRepo.GetRandomActiveResource()
	if era != nil {
		return nil
	}
	return resource
}
