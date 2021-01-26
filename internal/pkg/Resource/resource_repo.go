package Resource

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// ResourceRepo interface
type ResourceRepo interface {
	CreateResource(resource *entity.Resource) (*entity.Resource, error)
	DeleteResource(ResourceID uint) error
	GetResourceByID(ResourceID uint) (*entity.Resource, error)
	GetPDF(Offset, Limit uint) (*[]entity.Resource, error)
	GetPicture(Offset, Limit uint) (*[]entity.Resource, error)
	GetFiles(Offset, Limit uint) (*[]entity.Resource, error)
	GetAudios(Offset, Limit uint) (*[]entity.Resource, error)
	GetVideos(Offset, Limit uint) (*[]entity.Resource, error)
	GetResources(Offset, Limit uint) (*[]entity.Resource, error)
	SearchResource(query string, typo uint) ([]entity.Resource, error)
	GetRandomActiveResource() (*entity.Resource, error)
}
