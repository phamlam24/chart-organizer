import React from 'react';
import Plot from 'react-plotly.js';
import type { CSVData, LinePlot } from '../../types';

interface LinePlotChartProps {
  data: CSVData;
  config: LinePlot;
}

const LinePlotChart: React.FC<LinePlotChartProps> = ({ data, config }) => {
  const xIndex = data.headers.indexOf(config.columnX);
  const yIndex = data.headers.indexOf(config.columnY);

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
    mode: 'lines+markers' as const,
    type: 'scatter' as const,
    line: {
      color: 'rgba(75, 192, 192, 1)',
      width: 2,
    },
    marker: {
      color: 'rgba(75, 192, 192, 0.7)',
      size: 6,
    },
  }];

  const layout = {
    title: { text: config.title },
    xaxis: { title: { text: config.columnX } },
    yaxis: { title: { text: config.columnY } },
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

export default LinePlotChart;