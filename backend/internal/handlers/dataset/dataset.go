package dataset

import (
	datasetv1 "chart-organizer/backend/gen/contracts/dataset/v1"
	"chart-organizer/backend/internal/middleware"
	"chart-organizer/backend/internal/repository/dataset"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	// "net/http"

	"connectrpc.com/connect"
)

type DatasetHandler struct {
	DB *sql.DB
}

// GetDataset implements datasetv1connect.DatasetServiceHandler.
// Get the dataset. Remember to get the userId through the authorization header.
// If the dataset does not belong to the user, return not found.
func (h *DatasetHandler) GetDataset(
	ctx context.Context,
	req *connect.Request[datasetv1.GetDatasetRequest],
) (*connect.Response[datasetv1.GetDatasetResponse], error) {
	userId, found := middleware.GetUserId(ctx)
	if !found {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	data, err := dataset.GetDataset(h.DB, userId, req.Msg.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("dataset not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &datasetv1.GetDatasetResponse{
		Data: data,
	}
	return connect.NewResponse(res), nil
}

// UploadDataset implements datasetv1connect.DatasetServiceHandler.
func (h *DatasetHandler) UploadDataset(
	ctx context.Context,
	req *connect.Request[datasetv1.UploadDatasetRequest],
) (*connect.Response[datasetv1.UploadDatasetResponse], error) {
	// Get userId or return error
	userId, found := middleware.GetUserId(ctx)
	slog.Info(fmt.Sprintf("User id: %s", userId))
	if !found {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	id, err := dataset.AddNewDataset(h.DB, userId, req.Msg.Filename, req.Msg.Data)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&datasetv1.UploadDatasetResponse{
		Id: id,
	}), nil
}

// GetAllDatasetsFromUser implements datasetv1connect.DatasetServiceHandler.
func (h *DatasetHandler) GetAllDatasetsFromUser(
	ctx context.Context,
	req *connect.Request[datasetv1.GetAllDatasetsFromUserRequest],
) (*connect.Response[datasetv1.GetAllDatasetsFromUserResponse], error) {
	userId, found := middleware.GetUserId(ctx)
	if !found {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	datasets, err := dataset.GetAllDatasetsFromUser(h.DB, userId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var resDatasets []*datasetv1.GetAllDatasetsFromUser_Dataset
	for _, d := range datasets {
		resDatasets = append(resDatasets, &datasetv1.GetAllDatasetsFromUser_Dataset{
			Id:   d.ID,
			Name: d.Name,
		})
	}

	res := &datasetv1.GetAllDatasetsFromUserResponse{
		Datasets: resDatasets,
	}

	return connect.NewResponse(res), nil
}
