import apiClient from './apiClient';
import type { Dataset, UploadDatasetResponse } from '../types';

export const datasetService = {
  async uploadDataset(file: File): Promise<UploadDatasetResponse> {
    const formData = new FormData();
    formData.append('filename', file.name);
    formData.append('data', file);

    const response = await apiClient.post('/contracts.dataset.v1.DatasetService/UploadDataset', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  async getAllDatasets(): Promise<Dataset[]> {
    const response = await apiClient.post('/contracts.dataset.v1.DatasetService/GetAllDatasetsFromUser', {});
    return response.data.datasets || [];
  },

  async getDataset(id: string): Promise<Blob> {
    const response = await apiClient.post('/contracts.dataset.v1.DatasetService/GetDataset', { id }, {
      responseType: 'blob',
    });
    return response.data;
  }
};