package Resource

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type ResourceService interface {
	GetPDF(Offset uint, Limit uint) *[]entity.Resource
	GetPicture(Offset uint, Limit uint) *[]entity.Resource
	GetFiles(Offset uint, Limit uint) *[]entity.Resource
	GetAudios(Offset, Limit uint) *[]entity.Resource
	GetVideos(Offset uint, Limit uint) *[]entity.Resource
	GetResources(Offset, Limit uint) *[]entity.Resource
	GetResourceByID(ResourceID uint) *entity.Resource
	DeleteResource(ResourceID uint) bool
	CreateResource(resource *entity.Resource) *entity.Resource
	SearchResource(query string, typo uint) []entity.Resource
	GetRandomActiveResource() *entity.Resource
}
