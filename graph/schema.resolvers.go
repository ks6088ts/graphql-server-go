package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/ks6088ts/graphql-server-go/graph/generated"
	"github.com/ks6088ts/graphql-server-go/graph/model"
)

func (r *queryResolver) StationByName(ctx context.Context, stationName *string) ([]*model.Station, error) {
	return r.getStationByName(ctx, stationName)
}

func (r *queryResolver) StationByCd(ctx context.Context, stationCd *int) (*model.Station, error) {
	return r.getStationByCD(ctx, stationCd)
}

func (r *queryResolver) ProductByName(ctx context.Context, productName *string) ([]*model.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductByID(ctx context.Context, productID *int) (*model.Product, error) {
	return r.getProductById(ctx, productID)
}

func (r *stationResolver) BeforeStation(ctx context.Context, obj *model.Station) (*model.Station, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *stationResolver) AfterStation(ctx context.Context, obj *model.Station) (*model.Station, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *stationResolver) TransferStation(ctx context.Context, obj *model.Station) ([]*model.Station, error) {
	return r.transferStation(ctx, obj)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Station returns generated.StationResolver implementation.
func (r *Resolver) Station() generated.StationResolver { return &stationResolver{r} }

type queryResolver struct{ *Resolver }
type stationResolver struct{ *Resolver }
