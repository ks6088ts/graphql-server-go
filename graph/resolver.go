package graph

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ks6088ts/graphql-server-go/graph/model"
	"github.com/ks6088ts/graphql-server-go/models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Db *sql.DB
}

// 駅CD検索部分
func (r *Resolver) getStationByCD(ctx context.Context, stationCd *int) (*model.Station, error) {
	stations, err := models.StationByCDsByStationCD(r.Db, *stationCd)
	if err != nil {
		return nil, err
	}
	if len(stations) == 0 {
		return nil, errors.New("not found")
	}
	first := stations[0]

	return &model.Station{
		StationCd:   first.StationCd,
		StationName: first.StationName,
		LineName:    &first.LineName,
		Address:     &first.Address,
	}, nil
}

// 乗り換え駅取得部分
func (r *Resolver) transferStation(ctx context.Context, obj *model.Station) ([]*model.Station, error) {
	stationCd := obj.StationCd

	records, err := models.TransfersByStationCD(r.Db, stationCd)
	if err != nil {
		return nil, err
	}

	resp := make([]*model.Station, 0, len(records))
	for _, v := range records {
		if v.TransferStationName == "" {
			continue
		}
		resp = append(resp, &model.Station{
			StationCd:   v.TransferStationCd,
			StationName: v.TransferStationName,
			LineName:    &v.TransferLineName,
			Address:     &v.TransferAddress,
		})
	}

	return resp, nil
}

// 駅名検索部分
func (r *Resolver) getStationByName(ctx context.Context, stationName *string) ([]*model.Station, error) {
	stations, err := models.StationByNamesByStationName(r.Db, *stationName)
	if err != nil {
		return nil, err
	}
	if len(stations) == 0 {
		return nil, errors.New("not found")
	}

	resp := make([]*model.Station, 0, len(stations))
	for _, v := range stations {
		resp = append(resp, &model.Station{
			StationCd:   v.StationCd,
			StationName: v.StationName,
			LineName:    &v.LineName,
			Address:     &v.Address,
		})
	}

	return resp, nil
}

// 製品ID検索部分
func (r *Resolver) getProductById(ctx context.Context, productID *int) (*model.Product, error) {
	products, err := models.ProductByIdsByProductId(r.Db, *productID)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, errors.New("not found")
	}
	first := products[0]

	return &model.Product{
		ProductID:   first.ProductID,
		CompanyCd:   first.CompanyCd,
		InventoryCd: first.InventoryCd,
		PriceJpy:    first.PriceJpy,
		ProductName: first.ProductName,
		Description: first.Description,
	}, nil
}
