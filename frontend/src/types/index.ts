// Auth types
export interface User {
  id: string;
  username: string;
  created_at: string;
}

export interface SignupRequest {
  username: string;
  password: string;
}

export interface SignupResponse {
  jwtToken: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  jwtToken: string;
}

// Dataset types
export interface Dataset {
  id: string;
  name: string;
}

export interface UploadDatasetRequest {
  filename: string;
  data: File;
}

export interface UploadDatasetResponse {
  id: string;
}

// Visualization types
export interface ParallelCoordinates {
  title: string;
  columns: string[];
}

export interface Scatterplot {
  title: string;
  columnX: string;
  columnY: string;
}

export interface LinePlot {
  title: string;
  columnX: string;
  columnY: string;
}

export type VisualizationType = 'parallel_coordinates' | 'scatterplot' | 'lineplot';

export interface Visualization {
  type: VisualizationType;
  parallelCoordinates?: ParallelCoordinates;
  scatterplot?: Scatterplot;
  lineplot?: LinePlot;
}

export interface CreateDashboardRequest {
  visualizations: Visualization[];
  datasetId: string;
}

export interface CreateDashboardResponse {
  id: string;
}

export interface GetDashboardResponse {
  visualizations: Visualization[];
  datasetId: string;
}

// CSV data structure
export interface CSVData {
  headers: string[];
  rows: (string | number)[][];
}

// Auth context type
export interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (username: string, password: string) => Promise<void>;
  signup: (username: string, password: string) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}