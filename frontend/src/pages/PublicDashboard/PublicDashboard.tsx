import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { dashboardService } from '../../services/dashboardService';
import ParallelCoordinatesChart from '../../components/Charts/ParallelCoordinatesChart';
import ScatterplotChart from '../../components/Charts/ScatterplotChart';
import LinePlotChart from '../../components/Charts/LinePlotChart';
import { loadDatasetCSV } from '../../utils';
import type { CSVData, Visualization, GetDashboardResponse, ParallelCoordinates, Scatterplot, LinePlot } from '../../types';

const PublicDashboard: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [dashboard, setDashboard] = useState<GetDashboardResponse | null>(null);
  const [csvData, setCsvData] = useState<CSVData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (id) {
      loadDashboard();
    }
  }, [id]);

  const loadDashboard = async () => {
    try {
      setIsLoading(true);
      setError('');
      
      if (!id) {
        setError('Dashboard ID is missing');
        return;
      }
      
      const dashboardData = await dashboardService.getDashboard(id);
      setDashboard(dashboardData);

      // Load the associated dataset
      const parsed = await loadDatasetCSV(dashboardData.dataset_id);
      setCsvData(parsed);
    } catch (err) {
      console.error('Error loading dashboard:', err);
      if (err instanceof Error) {
        setError(`Dashboard not found or failed to load: ${err.message}`);
      } else {
        setError('Dashboard not found or failed to load');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const renderChart = (viz: Visualization, index: number) => {
    if (!csvData) return null;

    switch (viz.type) {
      case 'parallel_coordinates':
        return (
          <ParallelCoordinatesChart
            key={index}
            data={csvData}
            config={viz.plot as ParallelCoordinates}
          />
        );
      case 'scatterplot':
        return (
          <ScatterplotChart
            key={index}
            data={csvData}
            config={viz.plot as Scatterplot}
          />
        );
      case 'lineplot':
        return (
          <LinePlotChart
            key={index}
            data={csvData}
            config={viz.plot as LinePlot}
          />
        );
      default:
        return null;
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error || !dashboard || !csvData) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <svg
            className="mx-auto h-12 w-12 text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            aria-hidden="true"
          >
            <path
              vectorEffect="non-scaling-stroke"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
            />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">Dashboard not found</h3>
          <p className="mt-1 text-sm text-gray-500">
            The dashboard you're looking for doesn't exist or has been removed.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="bg-white shadow">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <h1 className="text-2xl font-bold text-gray-900">Data Visualization Dashboard</h1>
          <p className="mt-2 text-sm text-gray-600">
            Shared dashboard with {dashboard.visualizations.length} visualization{dashboard.visualizations.length !== 1 ? 's' : ''}
          </p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="space-y-8">
          {dashboard.visualizations.map((viz, index) => (
            <div key={index} className="bg-white shadow sm:rounded-lg">
              <div className="px-4 py-5 sm:p-6">
                <h3 className="text-lg leading-6 font-medium text-gray-900 mb-4">
                  {(viz.plot as any).title}
                </h3>
                {renderChart(viz, index)}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default PublicDashboard;