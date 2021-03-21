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
