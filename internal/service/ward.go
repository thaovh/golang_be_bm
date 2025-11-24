package service

import (
	"context"
	"time"

	v1 "github.com/go-kratos/kratos-layout/api/ward/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
)

type WardService struct {
	v1.UnimplementedWardServiceServer

	uc *biz.WardUsecase
}

func NewWardService(uc *biz.WardUsecase) *WardService {
	return &WardService{uc: uc}
}

// CreateWard creates a new ward
func (s *WardService) CreateWard(ctx context.Context, req *v1.CreateWardRequest) (*v1.CreateWardResponse, error) {
	provinceID, err := uuid.FromString(req.ProvinceId)
	if err != nil {
		return nil, errors.BadRequest("INVALID_PROVINCE_ID", "invalid province ID format")
	}

	ward := &biz.Ward{
		ProvinceID:  provinceID,
		Code:        req.Code,
		Name:        req.Name,
		NameEn:      req.NameEn,
		Type:        req.Type,
		Area:        req.Area,
		Population:  req.Population,
		Coordinates: req.Coordinates,
		PostalCode:  req.PostalCode,
		Address:     req.Address,
		SortOrder:   int(req.SortOrder),
	}

	created, err := s.uc.CreateWard(ctx, ward)
	if err != nil {
		return nil, convertWardError(err)
	}

	return &v1.CreateWardResponse{
		Ward: toProtoWard(created),
	}, nil
}

// UpdateWard updates an existing ward
func (s *WardService) UpdateWard(ctx context.Context, req *v1.UpdateWardRequest) (*v1.UpdateWardResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid ward ID format")
	}

	ward := &biz.Ward{
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
		PostalCode:  req.PostalCode,
		Address:     req.Address,
		SortOrder:   int(req.SortOrder),
	}

	// Set province ID if provided
	if req.ProvinceId != "" {
		provinceID, err := uuid.FromString(req.ProvinceId)
		if err != nil {
			return nil, errors.BadRequest("INVALID_PROVINCE_ID", "invalid province ID format")
		}
		ward.ProvinceID = provinceID
	}

	updated, err := s.uc.UpdateWard(ctx, ward)
	if err != nil {
		return nil, convertWardError(err)
	}

	return &v1.UpdateWardResponse{
		Ward: toProtoWard(updated),
	}, nil
}

// DeleteWard deletes a ward
func (s *WardService) DeleteWard(ctx context.Context, req *v1.DeleteWardRequest) (*v1.DeleteWardResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid ward ID format")
	}

	err = s.uc.DeleteWard(ctx, id)
	if err != nil {
		return nil, convertWardError(err)
	}

	return &v1.DeleteWardResponse{Success: true}, nil
}

// GetWard retrieves a ward by ID
func (s *WardService) GetWard(ctx context.Context, req *v1.GetWardRequest) (*v1.GetWardResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid ward ID format")
	}

	ward, err := s.uc.GetWard(ctx, id)
	if err != nil {
		return nil, convertWardError(err)
	}

	return &v1.GetWardResponse{
		Ward: toProtoWard(ward),
	}, nil
}

// GetWardByCode retrieves a ward by its code and province ID
func (s *WardService) GetWardByCode(ctx context.Context, req *v1.GetWardByCodeRequest) (*v1.GetWardByCodeResponse, error) {
	provinceID, err := uuid.FromString(req.ProvinceId)
	if err != nil {
		return nil, errors.BadRequest("INVALID_PROVINCE_ID", "invalid province ID format")
	}

	ward, err := s.uc.GetWardByCode(ctx, req.Code, provinceID)
	if err != nil {
		return nil, convertWardError(err)
	}

	return &v1.GetWardByCodeResponse{
		Ward: toProtoWard(ward),
	}, nil
}

// ListWards lists wards with pagination and filters
func (s *WardService) ListWards(ctx context.Context, req *v1.ListWardsRequest) (*v1.ListWardsResponse, error) {
	filter := &biz.WardListFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		Search:   req.Search,
		Type:     req.Type,
		Status:   req.Status,
		Code:     req.Code,
	}

	// Parse province ID if provided
	if req.ProvinceId != "" {
		provinceID, err := uuid.FromString(req.ProvinceId)
		if err != nil {
			return nil, errors.BadRequest("INVALID_PROVINCE_ID", "invalid province ID format")
		}
		filter.ProvinceID = provinceID
	}

	wards, total, err := s.uc.ListWards(ctx, filter)
	if err != nil {
		return nil, convertWardError(err)
	}

	protoWards := make([]*v1.Ward, 0, len(wards))
	for _, ward := range wards {
		protoWards = append(protoWards, toProtoWard(ward))
	}

	return &v1.ListWardsResponse{
		Wards: protoWards,
		Total: total,
	}, nil
}

// ListWardsByProvince lists all wards of a specific province
func (s *WardService) ListWardsByProvince(ctx context.Context, req *v1.ListWardsByProvinceRequest) (*v1.ListWardsByProvinceResponse, error) {
	provinceID, err := uuid.FromString(req.ProvinceId)
	if err != nil {
		return nil, errors.BadRequest("INVALID_PROVINCE_ID", "invalid province ID format")
	}

	wards, err := s.uc.ListWardsByProvince(ctx, provinceID)
	if err != nil {
		return nil, convertWardError(err)
	}

	protoWards := make([]*v1.Ward, 0, len(wards))
	for _, ward := range wards {
		protoWards = append(protoWards, toProtoWard(ward))
	}

	return &v1.ListWardsByProvinceResponse{
		Wards: protoWards,
	}, nil
}

// Helper function to convert biz.Ward to v1.Ward
func toProtoWard(ward *biz.Ward) *v1.Ward {
	if ward == nil {
		return nil
	}

	var createdBy, updatedBy string
	if ward.CreatedBy != nil {
		createdBy = ward.CreatedBy.String()
	}
	if ward.UpdatedBy != nil {
		updatedBy = ward.UpdatedBy.String()
	}

	return &v1.Ward{
		Id:          ward.ID.String(),
		ProvinceId:  ward.ProvinceID.String(),
		Code:        ward.Code,
		Name:        ward.Name,
		NameEn:      ward.NameEn,
		Type:        ward.Type,
		Area:        ward.Area,
		Population:  ward.Population,
		Coordinates: ward.Coordinates,
		PostalCode:  ward.PostalCode,
		Address:     ward.Address,
		SortOrder:   int32(ward.SortOrder),
		Status:      ward.Status,
		CreatedAt:   ward.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   ward.UpdatedAt.Format(time.RFC3339),
		CreatedBy:   createdBy,
		UpdatedBy:   updatedBy,
	}
}

// Helper function to convert biz errors to v1 errors
func convertWardError(err error) error {
	if errors.Is(err, biz.ErrWardNotFound) {
		return errors.NotFound(v1.ErrorReason_WARD_NOT_FOUND.String(), err.Error())
	}
	if errors.Is(err, biz.ErrWardAlreadyExists) {
		return errors.Conflict(v1.ErrorReason_WARD_ALREADY_EXISTS.String(), err.Error())
	}
	if errors.Is(err, biz.ErrInvalidWardCode) {
		return errors.BadRequest(v1.ErrorReason_INVALID_WARD_CODE.String(), err.Error())
	}
	if errors.Is(err, biz.ErrProvinceRequired) {
		return errors.BadRequest(v1.ErrorReason_PROVINCE_REQUIRED.String(), err.Error())
	}
	if errors.Is(err, biz.ErrProvinceNotFound) {
		return errors.NotFound("PROVINCE_NOT_FOUND", err.Error())
	}
	return err
}

