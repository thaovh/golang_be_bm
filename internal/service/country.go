package service

import (
	"context"

	v1 "github.com/go-kratos/kratos-layout/api/country/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
)

type CountryService struct {
	v1.UnimplementedCountryServiceServer

	uc *biz.CountryUsecase
}

func NewCountryService(uc *biz.CountryUsecase) *CountryService {
	return &CountryService{uc: uc}
}

func (s *CountryService) CreateCountry(ctx context.Context, req *v1.CreateCountryRequest) (*v1.CreateCountryResponse, error) {
	country := &biz.Country{
		Code:            req.Code,
		Name:            req.Name,
		NameEn:          req.NameEn,
		Region:          req.Region,
		SubRegion:       req.SubRegion,
		CurrencyCode:    req.CurrencyCode,
		CurrencySymbol:  req.CurrencySymbol,
		PhoneCode:       req.PhoneCode,
		TimeZone:        req.TimeZone,
		Flag:            req.Flag,
		Capital:         req.Capital,
		Population:      req.Population,
		ISO3166Alpha3:   req.Iso3166Alpha3,
		ISO3166Numeric:  req.Iso3166Numeric,
	}

	created, err := s.uc.CreateCountry(ctx, country)
	if err != nil {
		return nil, err
	}

	return &v1.CreateCountryResponse{
		Country: toProtoCountry(created),
	}, nil
}

func (s *CountryService) UpdateCountry(ctx context.Context, req *v1.UpdateCountryRequest) (*v1.UpdateCountryResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid id format")
	}

	country := &biz.Country{
		BaseEntity: biz.BaseEntity{
			ID:     id,
			Status: req.Status,
		},
		Code:            req.Code,
		Name:            req.Name,
		NameEn:          req.NameEn,
		Region:          req.Region,
		SubRegion:       req.SubRegion,
		CurrencyCode:    req.CurrencyCode,
		CurrencySymbol:  req.CurrencySymbol,
		PhoneCode:       req.PhoneCode,
		TimeZone:        req.TimeZone,
		Flag:            req.Flag,
		Capital:         req.Capital,
		Population:      req.Population,
		ISO3166Alpha3:   req.Iso3166Alpha3,
		ISO3166Numeric:  req.Iso3166Numeric,
	}

	updated, err := s.uc.UpdateCountry(ctx, country)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateCountryResponse{
		Country: toProtoCountry(updated),
	}, nil
}

func (s *CountryService) DeleteCountry(ctx context.Context, req *v1.DeleteCountryRequest) (*v1.DeleteCountryResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid id format")
	}

	if err := s.uc.DeleteCountry(ctx, id); err != nil {
		return nil, err
	}

	return &v1.DeleteCountryResponse{Success: true}, nil
}

func (s *CountryService) GetCountry(ctx context.Context, req *v1.GetCountryRequest) (*v1.GetCountryResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid id format")
	}

	country, err := s.uc.GetCountry(ctx, id)
	if err != nil {
		return nil, err
	}

	if country == nil {
		return nil, errors.NotFound("COUNTRY_NOT_FOUND", "country not found")
	}

	return &v1.GetCountryResponse{
		Country: toProtoCountry(country),
	}, nil
}

func (s *CountryService) GetCountryByCode(ctx context.Context, req *v1.GetCountryByCodeRequest) (*v1.GetCountryByCodeResponse, error) {
	country, err := s.uc.GetCountryByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	if country == nil {
		return nil, errors.NotFound("COUNTRY_NOT_FOUND", "country not found")
	}

	return &v1.GetCountryByCodeResponse{
		Country: toProtoCountry(country),
	}, nil
}

func (s *CountryService) ListCountries(ctx context.Context, req *v1.ListCountriesRequest) (*v1.ListCountriesResponse, error) {
	filter := &biz.CountryListFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		Search:   req.Search,
		Region:   req.Region,
		Status:   req.Status,
		Code:     req.Code,
	}

	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100 // Max page size
	}

	countries, total, err := s.uc.ListCountries(ctx, filter)
	if err != nil {
		return nil, err
	}

	protoCountries := make([]*v1.Country, len(countries))
	for i, country := range countries {
		protoCountries[i] = toProtoCountry(country)
	}

	return &v1.ListCountriesResponse{
		Countries: protoCountries,
		Total:     total,
	}, nil
}

func (s *CountryService) SearchCountries(ctx context.Context, req *v1.SearchCountriesRequest) (*v1.SearchCountriesResponse, error) {
	countries, err := s.uc.SearchCountries(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	protoCountries := make([]*v1.Country, len(countries))
	for i, country := range countries {
		protoCountries[i] = toProtoCountry(country)
	}

	return &v1.SearchCountriesResponse{
		Countries: protoCountries,
	}, nil
}

// Helper function to convert domain entity to proto
func toProtoCountry(country *biz.Country) *v1.Country {
	if country == nil {
		return nil
	}

	return &v1.Country{
		Id:             country.ID.String(),
		Code:           country.Code,
		Name:           country.Name,
		NameEn:         country.NameEn,
		Region:         country.Region,
		SubRegion:      country.SubRegion,
		CurrencyCode:   country.CurrencyCode,
		CurrencySymbol: country.CurrencySymbol,
		PhoneCode:      country.PhoneCode,
		TimeZone:       country.TimeZone,
		Flag:           country.Flag,
		Capital:        country.Capital,
		Population:     country.Population,
		Iso3166Alpha3:  country.ISO3166Alpha3,
		Iso3166Numeric: country.ISO3166Numeric,
		Status:         country.Status,
		CreatedAt:      country.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      country.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

