import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { dashboardService } from '../../services/dashboardService';
import { ParallelCoordinatesChart, ScatterplotChart, LinePlotChart } from '../../components/Charts';
import { loadDatasetCSV } from '../../utils';
import type { CSVData, Visualization, VisualizationType, ParallelCoordinates, Scatterplot, LinePlot } from '../../types';

const DashboardCreator: React.FC = () => {
  const { datasetId } = useParams<{ datasetId: string }>();
  const navigate = useNavigate();
  
  const [csvData, setCsvData] = useState<CSVData | null>(null);
  const [visualizations, setVisualizations] = useState<Visualization[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [showAddForm, setShowAddForm] = useState(false);
  const [newVizType, setNewVizType] = useState<VisualizationType>('scatterplot');
  const [newVizConfig, setNewVizConfig] = useState({
    title: '',
    columns: [] as string[],
    column_x: '',
    column_y: '',
  });

  useEffect(() => {
    if (datasetId) {
      loadDataset();
    }
  }, [datasetId]);

  const loadDataset = async () => {
    try {
      setIsLoading(true);
      setError('');
      
      if (!datasetId) {
        setError('Dataset ID is missing');
        return;
      }
      
      const parsed = await loadDatasetCSV(datasetId);
      setCsvData(parsed);
    } catch (err) {
      console.error('Error loading dataset:', err);
      if (err instanceof Error) {
        setError(`Failed to load dataset: ${err.message}`);
      } else {
        setError('Failed to load dataset');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const addVisualization = () => {
    if (!newVizConfig.title.trim()) {
      setError('Please enter a title for the visualization');
      return;
    }

    let plot: ParallelCoordinates | Scatterplot | LinePlot;

    switch (newVizType) {
      case 'parallel_coordinates':
        if (newVizConfig.columns.length < 2) {
          setError('Please select at least 2 columns for parallel coordinates');
          return;
        }
        plot = {
          title: newVizConfig.title,
          columns: newVizConfig.columns,
        };
        break;
      case 'scatterplot':
      case 'lineplot':
        if (!newVizConfig.column_x || !newVizConfig.column_y) {
          setError('Please select both X and Y columns');
          return;
        }
        plot = {
          title: newVizConfig.title,
          column_x: newVizConfig.column_x,
          column_y: newVizConfig.column_y,
        };
        break;
    }

    const newVisualization: Visualization = {
      type: newVizType,
      plot,
    };

    setVisualizations([...visualizations, newVisualization]);
    setShowAddForm(false);
    setNewVizConfig({ title: '', columns: [], column_x: '', column_y: '' });
    setError('');
  };

  const removeVisualization = (index: number) => {
    setVisualizations(visualizations.filter((_, i) => i !== index));
  };

  const saveDashboard = async () => {
    if (visualizations.length === 0) {
      setError('Please add at least one visualization');
      return;
    }

    try {
      setIsLoading(true);
      const response = await dashboardService.createDashboard({
        visualizations,
        dataset_id: datasetId!,
      });
      
      alert(`Dashboard created! Share this link: ${window.location.origin}/dashboard/${response.id}`);
      navigate('/');
    } catch (err) {
      setError('Failed to create dashboard');
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
      <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-center min-h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
      </div>
    );
  }

  if (!csvData) {
    return (
      <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="text-red-600">Failed to load dataset</div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <div className="sm:flex sm:items-center">
        <div className="sm:flex-auto">
          <h1 className="text-xl font-semibold text-gray-900">Create Dashboard</h1>
          <p className="mt-2 text-sm text-gray-700">
            Add and configure visualizations for your dataset. You can reorder them by dragging.
          </p>
        </div>
        <div className="mt-4 sm:mt-0 sm:ml-16 sm:flex-none space-x-2">
          <button
            onClick={() => setShowAddForm(true)}
            className="inline-flex items-center justify-center rounded-md border border-transparent bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
          >
            Add Visualization
          </button>
          <button
            onClick={saveDashboard}
            disabled={visualizations.length === 0}
            className="inline-flex items-center justify-center rounded-md border border-transparent bg-green-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 disabled:opacity-50"
          >
            Save Dashboard
          </button>
        </div>
      </div>

      {error && (
        <div className="mt-4 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      {/* Add Visualization Form */}
      {showAddForm && (
        <div className="mt-6 bg-white shadow sm:rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900">
              Add New Visualization
            </h3>
            <div className="mt-2 max-w-xl text-sm text-gray-500">
              <p>Configure your visualization settings.</p>
            </div>
            <div className="mt-5 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Visualization Type
                </label>
                <select
                  value={newVizType}
                  onChange={(e) => setNewVizType(e.target.value as VisualizationType)}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                >
                  <option value="scatterplot">Scatter Plot</option>
                  <option value="lineplot">Line Plot</option>
                  <option value="parallel_coordinates">Parallel Coordinates</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Title
                </label>
                <input
                  type="text"
                  value={newVizConfig.title}
                  onChange={(e) => setNewVizConfig({ ...newVizConfig, title: e.target.value })}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  placeholder="Enter visualization title"
                />
              </div>

              {newVizType === 'parallel_coordinates' ? (
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    Columns (select multiple)
                  </label>
                  <select
                    multiple
                    value={newVizConfig.columns}
                    onChange={(e) => setNewVizConfig({
                      ...newVizConfig,
                      columns: Array.from(e.target.selectedOptions, option => option.value)
                    })}
                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    size={Math.min(csvData.headers.length, 6)}
                  >
                    {csvData.headers.map(header => (
                      <option key={header} value={header}>{header}</option>
                    ))}
                  </select>
                  <p className="mt-1 text-sm text-gray-500">Hold Ctrl/Cmd to select multiple columns</p>
                </div>
              ) : (
                <>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      X-axis Column
                    </label>
                    <select
                      value={newVizConfig.column_x}
                      onChange={(e) => setNewVizConfig({ ...newVizConfig, column_x: e.target.value })}
                      className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    >
                      <option value="">Select X column</option>
                      {csvData.headers.map(header => (
                        <option key={header} value={header}>{header}</option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Y-axis Column
                    </label>
                    <select
                      value={newVizConfig.column_y}
                      onChange={(e) => setNewVizConfig({ ...newVizConfig, column_y: e.target.value })}
                      className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    >
                      <option value="">Select Y column</option>
                      {csvData.headers.map(header => (
                        <option key={header} value={header}>{header}</option>
                      ))}
                    </select>
                  </div>
                </>
              )}
            </div>
            <div className="mt-5 flex justify-end space-x-3">
              <button
                onClick={() => {
                  setShowAddForm(false);
                  setError('');
                }}
                className="bg-white py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                onClick={addVisualization}
                className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
              >
                Add Visualization
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Visualizations */}
      <div className="mt-8 space-y-8">
        {visualizations.map((viz, index) => (
          <div key={index} className="bg-white shadow sm:rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg leading-6 font-medium text-gray-900">
                  {(viz.plot as any).title} ({viz.type.replace('_', ' ')})
                </h3>
                <button
                  onClick={() => removeVisualization(index)}
                  className="text-red-600 hover:text-red-900"
                >
                  Remove
                </button>
              </div>
              {renderChart(viz, index)}
            </div>
          </div>
        ))}
      </div>

      {visualizations.length === 0 && (
        <div className="text-center py-12">
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
              d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
            />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">No visualizations</h3>
          <p className="mt-1 text-sm text-gray-500">
            Get started by adding your first visualization.
          </p>
          <div className="mt-6">
            <button
              onClick={() => setShowAddForm(true)}
              className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
            >
              Add Visualization
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default DashboardCreator;