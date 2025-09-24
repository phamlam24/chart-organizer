package viz

import (
	"context"
	"database/sql"
	"errors"

	"connectrpc.com/connect"

	vizv1 "chart-organizer/backend/gen/contracts/viz/v1"
	"chart-organizer/backend/internal/interceptors"
	"chart-organizer/backend/internal/repository/viz"
)

type VisualizationHandler struct {
	DB *sql.DB
}

// CreateDashboard implements vizv1connect.DashboardServiceHandler.
func (h *VisualizationHandler) CreateDashboard(
	ctx context.Context,
	req *connect.Request[vizv1.CreateDashboardRequest],
) (*connect.Response[vizv1.CreateDashboardResponse], error) {
	userId, found := interceptors.GetUserId(ctx)
	if !found {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	id, err := viz.AddNewDashboard(h.DB, userId, req.Msg.DatasetId, req.Msg.Visualizations)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &vizv1.CreateDashboardResponse{
		Id: id,
	}
	return connect.NewResponse(res), nil
}

// GetDashboard implements vizv1connect.DashboardServiceHandler.
func (h *VisualizationHandler) GetDashboard(
	ctx context.Context,
	req *connect.Request[vizv1.GetDashboardRequest],
) (*connect.Response[vizv1.GetDashboardResponse], error) {
	visualizations, datasetId, err := viz.GetDashboard(h.DB, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if visualizations == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("dashboard not found"))
	}

	res := &vizv1.GetDashboardResponse{
		Visualizations: visualizations,
		DatasetId:      datasetId,
	}
	return connect.NewResponse(res), nil
}
