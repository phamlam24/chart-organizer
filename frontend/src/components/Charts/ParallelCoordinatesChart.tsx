import React from 'react';
import Plot from 'react-plotly.js';
import type { CSVData, ParallelCoordinates } from '../../types';

interface ParallelCoordinatesChartProps {
  data: CSVData;
  config: ParallelCoordinates;
}

const ParallelCoordinatesChart: React.FC<ParallelCoordinatesChartProps> = ({ data, config }) => {
  // Prepare data for parallel coordinates plot
  const dimensions = config.columns.map(column => {
    const columnIndex = data.headers.indexOf(column);
    if (columnIndex === -1) return null;

    const values = data.rows.map(row => {
      const value = row[columnIndex];
      return typeof value === 'string' ? parseFloat(value) : value;
    }).filter(val => !isNaN(val));

    return {
      label: column,
      values: values,
    };
  }).filter(dim => dim !== null);

  const plotData = [{
    type: 'parcoords' as const,
    dimensions: dimensions,
    line: {
      color: 'blue',
      colorscale: 'Viridis',
    },
  }];

  const layout = {
    title: { text: config.title },
    margin: { t: 50, r: 50, b: 50, l: 50 },
  };

  return (
    <div className="w-full">
      <Plot
        data={plotData}
        layout={layout}
        style={{ width: '100%', height: '400px' }}
        config={{ responsive: true }}
      />
    </div>
  );
};

export default ParallelCoordinatesChart;