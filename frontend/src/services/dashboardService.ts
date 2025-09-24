import apiClient from './apiClient';
import type { CreateDashboardRequest, CreateDashboardResponse, GetDashboardResponse } from '../types';

export const dashboardService = {
  async createDashboard(request: CreateDashboardRequest): Promise<CreateDashboardResponse> {
    const response = await apiClient.post('/contracts.viz.v1.DashboardService/CreateDashboard', request);
    return response.data;
  },

  async getDashboard(id: string): Promise<GetDashboardResponse> {
    const response = await apiClient.post('/contracts.viz.v1.DashboardService/GetDashboard', { id });
    return response.data;
  }
};