package ResourceRepo

import (
	"fmt"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Resource"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// ResourceRepo struct
type ResourceRepo struct {
	DB *gorm.DB
}

// NewResourceRepo function
func NewResourceRepo(db *gorm.DB) Resource.ResourceRepo {
	return &ResourceRepo{
		DB: db,
	}
}

// CreateResource (resource *entity.Resource) (*entity.Resource, error)
func (resrepo *ResourceRepo) CreateResource(resource *entity.Resource) (*entity.Resource, error) {
	newErro := resrepo.DB.Create(resource).Error
	return resource, newErro
}

// DeleteResource ( ResourceID uint) error
func (resrepo *ResourceRepo) DeleteResource(ResourceID uint) error {
	resource := &entity.Resource{}
	resource.ID = ResourceID
	newError := resrepo.DB.Delete(resource, "id=?", ResourceID).Error
	return newError
}

// GetResourceByID method
func (resrepo *ResourceRepo) GetResourceByID(ResourceID uint) (*entity.Resource, error) {
	resource := &entity.Resource{}
	resource.ID = ResourceID
	newError := resrepo.DB.Where("id=?", ResourceID).First(resource).Error
	return resource, newError
}

// GetResources method
func (resrepo *ResourceRepo) GetResources(Offset, Limit uint) (*[]entity.Resource, error) {
	resources := &[]entity.Resource{}
	newErro := resrepo.DB.Table("resources").Offset(Offset).Limit(Limit).Find(resources).Error
	for k := 0; k < len(*resources); k++ {
		resource := &(*resources)[k]
		resrepo.DB.Model(resource).Related(&resource.UploadDate, "UDID")
	}
	return resources, newErro
}

// GetVideos method
func (resrepo *ResourceRepo) GetVideos(Offset, Limit uint) (*[]entity.Resource, error) {
	resources := &[]entity.Resource{}
	newErroa := resrepo.DB.Where("type=?", entity.VIDEO).Offset(Offset).Limit(Limit).Find(resources).Error
	return resources, newErroa
}

// GetAudios method
func (resrepo *ResourceRepo) GetAudios(Offset, Limit uint) (*[]entity.Resource, error) {
	resources := &[]entity.Resource{}
	newErroa := resrepo.DB.Where("type=?", entity.AUDIO).Offset(Offset).Limit(Limit).Find(resources).Error
	return resources, newErroa
}

// GetFiles mthod
func (resrepo *ResourceRepo) GetFiles(Offset, Limit uint) (*[]entity.Resource, error) {
	resources := &[]entity.Resource{}
	newErroa := resrepo.DB.Where("type=?", entity.FILES).Offset(Offset).Limit(Limit).Find(resources).Error
	return resources, newErroa
}

// GetPicture method
func (resrepo *ResourceRepo) GetPicture(Offset, Limit uint) (*[]entity.Resource, error) {
	resources := &[]entity.Resource{}
	newErroa := resrepo.DB.Where("type=?", entity.IMAGES).Offset(Offset).Limit(Limit).Find(resources).Error
	return resources, newErroa
}

// GetPDF mthod
func (resrepo *ResourceRepo) GetPDF(Offset, Limit uint) (*[]entity.Resource, error) {
	resources := &[]entity.Resource{}
	newErroa := resrepo.DB.Where("type=?", entity.PDFS).Offset(Offset).Limit(Limit).Find(resources).Error
	return resources, newErroa
}

// SearchResource (query string, typo uint) []entity.Resource
func (resrepo *ResourceRepo) SearchResource(query string, typo uint) ([]entity.Resource, error) {
	resources := []entity.Resource{}
	var newEra error
	if typo == 0 {
		newEra = resrepo.DB.Table("resources").Find(&resources, " title LIKE ? ", "%"+query+"%").Error
		newEra = resrepo.DB.Table("resources").Find(&resources, " title LIKE ? ", query+"%").Error
	} else {
		newEra = resrepo.DB.Table("resources").Find(&resources, " title LIKE ? and type=? ", "%"+query+"%", typo).Error
	}
	if len(resources) > 0 {
		for h := 0; h < len(resources); h++ {
			resource := &resources[h]
			resrepo.DB.Model(resource).Related(&(resource.UploadDate), "UDID")
		}
	}
	fmt.Println(resources)
	return resources, newEra
}

// GetRandomActiveResource  () (*entity.Resource, error)
func (resrepo *ResourceRepo) GetRandomActiveResource() (*entity.Resource, error) {
	resource := &entity.Resource{}
	newEra := resrepo.DB.First(resource, "type in (?)", []uint{entity.VIDEO, entity.AUDIO}).Error
	resrepo.DB.Model(resource).Related(&resource.UploadDate, "UploadDate")
	return resource, newEra
}
