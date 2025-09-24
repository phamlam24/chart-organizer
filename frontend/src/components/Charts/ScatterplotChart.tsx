import React from 'react';
import Plot from 'react-plotly.js';
import type { CSVData, Scatterplot } from '../../types';

interface ScatterplotChartProps {
  data: CSVData;
  config: Scatterplot;
}

const ScatterplotChart: React.FC<ScatterplotChartProps> = ({ data, config }) => {
  const xIndex = data.headers.indexOf(config.column_x);
  const yIndex = data.headers.indexOf(config.column_y);

  if (xIndex === -1 || yIndex === -1) {
    return (
      <div className="text-red-600 p-4">
        Error: Selected columns not found in dataset
      </div>
    );
  }

  const xValues = data.rows.map(row => {
    const value = row[xIndex];
    return typeof value === 'string' ? parseFloat(value) : value;
  }).filter(val => !isNaN(val));

  const yValues = data.rows.map(row => {
    const value = row[yIndex];
    return typeof value === 'string' ? parseFloat(value) : value;
  }).filter(val => !isNaN(val));

  const plotData = [{
    x: xValues,
    y: yValues,
    mode: 'markers' as const,
    type: 'scatter' as const,
    marker: {
      color: 'rgba(54, 162, 235, 0.7)',
      size: 8,
    },
  }];

  const layout = {
    title: { text: config.title },
    xaxis: { title: { text: config.column_x } },
    yaxis: { title: { text: config.column_y } },
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

export default ScatterplotChart;