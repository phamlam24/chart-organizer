import apiClient from './apiClient';
import type { Dataset, UploadDatasetResponse } from '../types';

export const datasetService = {
  async uploadDataset(file: File): Promise<UploadDatasetResponse> {
    // Convert file to base64
    const fileBuffer = await file.arrayBuffer();
    const uint8Array = new Uint8Array(fileBuffer);
    const base64Data = btoa(String.fromCharCode(...uint8Array));

    const requestData = {
      filename: file.name,
      data: base64Data
    };

    const response = await apiClient.post('/contracts.dataset.v1.DatasetService/UploadDataset', requestData, {
      headers: {
        'Content-Type': 'application/json',
      },
    });
    return response.data;
  },

  async getAllDatasets(): Promise<Dataset[]> {
    const response = await apiClient.post('/contracts.dataset.v1.DatasetService/GetAllDatasetsFromUser', {});
    return response.data.datasets || [];
  },

  async getDataset(id: string): Promise<Blob> {
    const response = await apiClient.post('/contracts.dataset.v1.DatasetService/GetDataset', { id });
    
    // The response.data contains { data: "base64-encoded-content" }
    const base64Data = response.data.data;
    
    if (!base64Data) {
      throw new Error('No data field in response');
    }
    
    try {
      // Decode base64 to binary
      const binaryString = atob(base64Data);
      const bytes = new Uint8Array(binaryString.length);
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
      }
      
      // Create blob from binary data
      return new Blob([bytes], { type: 'text/csv' });
    } catch (error) {
      console.error('Error decoding base64:', error);
      throw new Error('Failed to decode dataset data');
    }
  }
};