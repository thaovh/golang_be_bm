package service

import (
	"context"
	"time"

	v1 "github.com/go-kratos/kratos-layout/api/province/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
)

type ProvinceService struct {
	v1.UnimplementedProvinceServiceServer

	uc *biz.ProvinceUsecase
}

func NewProvinceService(uc *biz.ProvinceUsecase) *ProvinceService {
	return &ProvinceService{uc: uc}
}

// CreateProvince creates a new province
func (s *ProvinceService) CreateProvince(ctx context.Context, req *v1.CreateProvinceRequest) (*v1.CreateProvinceResponse, error) {
	countryID, err := uuid.FromString(req.CountryId)
	if err != nil {
		return nil, errors.BadRequest("INVALID_COUNTRY_ID", "invalid country ID format")
	}

	province := &biz.Province{
		CountryID:   countryID,
		Code:        req.Code,
		Name:        req.Name,
		NameEn:      req.NameEn,
		Type:        req.Type,
		Area:        req.Area,
		Population:  req.Population,
		Coordinates: req.Coordinates,
		Capital:     req.Capital,
		PostalCode:  req.PostalCode,
		PhonePrefix: req.PhonePrefix,
		SortOrder:   int(req.SortOrder),
	}

	created, err := s.uc.CreateProvince(ctx, province)
	if err != nil {
		return nil, convertProvinceError(err)
	}

	return &v1.CreateProvinceResponse{
		Province: toProtoProvince(created),
	}, nil
}

// UpdateProvince updates an existing province
func (s *ProvinceService) UpdateProvince(ctx context.Context, req *v1.UpdateProvinceRequest) (*v1.UpdateProvinceResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid province ID format")
	}

	province := &biz.Province{
		BaseEntity: biz.BaseEntity{
			ID:     id,
			Status: req.Status,
		},
		Code:        req.Code,
		Name:        req.Name,
		NameEn:      req.NameEn,
		Type:        req.Type,
		Area:        req.Area,
		Population:  req.Population,
		Coordinates: req.Coordinates,
		Capital:     req.Capital,
		PostalCode:  req.PostalCode,
		PhonePrefix: req.PhonePrefix,
		SortOrder:   int(req.SortOrder),
	}

	// Set country ID if provided
	if req.CountryId != "" {
		countryID, err := uuid.FromString(req.CountryId)
		if err != nil {
			return nil, errors.BadRequest("INVALID_COUNTRY_ID", "invalid country ID format")
		}
		province.CountryID = countryID
	}

	updated, err := s.uc.UpdateProvince(ctx, province)
	if err != nil {
		return nil, convertProvinceError(err)
	}

	return &v1.UpdateProvinceResponse{
		Province: toProtoProvince(updated),
	}, nil
}

// DeleteProvince deletes a province
func (s *ProvinceService) DeleteProvince(ctx context.Context, req *v1.DeleteProvinceRequest) (*v1.DeleteProvinceResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid province ID format")
	}

	err = s.uc.DeleteProvince(ctx, id)
	if err != nil {
		return nil, convertProvinceError(err)
	}

	return &v1.DeleteProvinceResponse{Success: true}, nil
}

// GetProvince retrieves a province by ID
func (s *ProvinceService) GetProvince(ctx context.Context, req *v1.GetProvinceRequest) (*v1.GetProvinceResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid province ID format")
	}

	province, err := s.uc.GetProvince(ctx, id)
	if err != nil {
		return nil, convertProvinceError(err)
	}

	return &v1.GetProvinceResponse{
		Province: toProtoProvince(province),
	}, nil
}

// GetProvinceByCode retrieves a province by its code and country ID
func (s *ProvinceService) GetProvinceByCode(ctx context.Context, req *v1.GetProvinceByCodeRequest) (*v1.GetProvinceByCodeResponse, error) {
	countryID, err := uuid.FromString(req.CountryId)
	if err != nil {
		return nil, errors.BadRequest("INVALID_COUNTRY_ID", "invalid country ID format")
	}

	province, err := s.uc.GetProvinceByCode(ctx, req.Code, countryID)
	if err != nil {
		return nil, convertProvinceError(err)
	}

	return &v1.GetProvinceByCodeResponse{
		Province: toProtoProvince(province),
	}, nil
}

// ListProvinces lists provinces with pagination and filters
func (s *ProvinceService) ListProvinces(ctx context.Context, req *v1.ListProvincesRequest) (*v1.ListProvincesResponse, error) {
	filter := &biz.ProvinceListFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		Search:   req.Search,
		Type:     req.Type,
		Status:   req.Status,
		Code:     req.Code,
	}

	// Parse country ID if provided
	if req.CountryId != "" {
		countryID, err := uuid.FromString(req.CountryId)
		if err != nil {
			return nil, errors.BadRequest("INVALID_COUNTRY_ID", "invalid country ID format")
		}
		filter.CountryID = countryID
	}

	provinces, total, err := s.uc.ListProvinces(ctx, filter)
	if err != nil {
		return nil, convertProvinceError(err)
	}

	protoProvinces := make([]*v1.Province, 0, len(provinces))
	for _, province := range provinces {
		protoProvinces = append(protoProvinces, toProtoProvince(province))
	}

	return &v1.ListProvincesResponse{
		Provinces: protoProvinces,
		Total:     total,
	}, nil
}

// ListProvincesByCountry lists all provinces of a specific country
func (s *ProvinceService) ListProvincesByCountry(ctx context.Context, req *v1.ListProvincesByCountryRequest) (*v1.ListProvincesByCountryResponse, error) {
	countryID, err := uuid.FromString(req.CountryId)
	if err != nil {
		return nil, errors.BadRequest("INVALID_COUNTRY_ID", "invalid country ID format")
	}

	provinces, err := s.uc.ListProvincesByCountry(ctx, countryID)
	if err != nil {
		return nil, convertProvinceError(err)
	}

	protoProvinces := make([]*v1.Province, 0, len(provinces))
	for _, province := range provinces {
		protoProvinces = append(protoProvinces, toProtoProvince(province))
	}

	return &v1.ListProvincesByCountryResponse{
		Provinces: protoProvinces,
	}, nil
}

// Helper function to convert biz.Province to v1.Province
func toProtoProvince(province *biz.Province) *v1.Province {
	if province == nil {
		return nil
	}

	var createdBy, updatedBy string
	if province.CreatedBy != nil {
		createdBy = province.CreatedBy.String()
	}
	if province.UpdatedBy != nil {
		updatedBy = province.UpdatedBy.String()
	}

	return &v1.Province{
		Id:          province.ID.String(),
		CountryId:   province.CountryID.String(),
		Code:        province.Code,
		Name:        province.Name,
		NameEn:      province.NameEn,
		Type:        province.Type,
		Area:        province.Area,
		Population:  province.Population,
		Coordinates: province.Coordinates,
		Capital:     province.Capital,
		PostalCode:  province.PostalCode,
		PhonePrefix: province.PhonePrefix,
		SortOrder:   int32(province.SortOrder),
		Status:       province.Status,
		CreatedAt:   province.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   province.UpdatedAt.Format(time.RFC3339),
		CreatedBy:   createdBy,
		UpdatedBy:   updatedBy,
	}
}

// Helper function to convert biz errors to v1 errors
func convertProvinceError(err error) error {
	if errors.Is(err, biz.ErrProvinceNotFound) {
		return errors.NotFound(v1.ErrorReason_PROVINCE_NOT_FOUND.String(), err.Error())
	}
	if errors.Is(err, biz.ErrProvinceAlreadyExists) {
		return errors.Conflict(v1.ErrorReason_PROVINCE_ALREADY_EXISTS.String(), err.Error())
	}
	if errors.Is(err, biz.ErrInvalidProvinceCode) {
		return errors.BadRequest(v1.ErrorReason_INVALID_PROVINCE_CODE.String(), err.Error())
	}
	if errors.Is(err, biz.ErrCountryRequired) {
		return errors.BadRequest(v1.ErrorReason_COUNTRY_REQUIRED.String(), err.Error())
	}
	if errors.Is(err, biz.ErrCountryNotFound) {
		return errors.NotFound("COUNTRY_NOT_FOUND", err.Error())
	}
	return err
}

