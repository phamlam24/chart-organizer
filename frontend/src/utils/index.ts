import type { CSVData } from '../types';
import { datasetService } from '../services/datasetService';

/**
 * Parse CSV text content into structured data
 */
export const parseCSV = (text: string): CSVData => {
  const lines = text.trim().split('\n').filter(line => line.trim()); // Filter out empty lines
  
  if (lines.length === 0) {
    throw new Error('CSV file is empty');
  }

  const headers = lines[0].split(',').map(h => h.trim().replace(/^"|"$/g, ''));
  
  if (headers.length === 0) {
    throw new Error('CSV file has no headers');
  }

  const rows = lines.slice(1).map(line => {
    const cells = line.split(',').map(cell => {
      const trimmed = cell.trim().replace(/^"|"$/g, '');
      // Try to parse as number, but keep as string if it fails
      const num = parseFloat(trimmed);
      return isNaN(num) ? trimmed : num;
    });
    
    // Ensure all rows have the same number of columns as headers
    while (cells.length < headers.length) {
      cells.push('');
    }
    
    return cells.slice(0, headers.length); // Trim excess columns
  });

  return { headers, rows };
};

/**
 * Load and parse dataset from the backend
 */
export const loadDatasetCSV = async (datasetId: string): Promise<CSVData> => {
  if (!datasetId) {
    throw new Error('Dataset ID is required');
  }

  const blob = await datasetService.getDataset(datasetId);
  
  if (!blob || blob.size === 0) {
    throw new Error('Dataset file is empty or not found');
  }
  
  const text = await blob.text();
  
  if (!text || text.trim().length === 0) {
    throw new Error('Dataset file contains no data');
  }
  
  return parseCSV(text);
};

/**
 * Get numeric columns from CSV data
 */
export const getNumericColumns = (data: CSVData): string[] => {
  if (data.rows.length === 0) return [];

  return data.headers.filter((_, index) => {
    // Check if most values in this column are numeric
    const numericCount = data.rows.filter(row => 
      typeof row[index] === 'number' && !isNaN(row[index] as number)
    ).length;
    
    return numericCount > data.rows.length * 0.7; // At least 70% numeric
  });
};

/**
 * Validate CSV file
 */
export const validateCSVFile = (file: File): string | null => {
  if (!file) return 'No file selected';
  
  const validTypes = ['text/csv', 'application/vnd.ms-excel'];
  const validExtensions = ['.csv'];
  
  const hasValidType = validTypes.includes(file.type);
  const hasValidExtension = validExtensions.some(ext => file.name.toLowerCase().endsWith(ext));
  
  if (!hasValidType && !hasValidExtension) {
    return 'Please select a valid CSV file';
  }
  
  if (file.size > 10 * 1024 * 1024) { // 10MB limit
    return 'File size must be less than 10MB';
  }
  
  return null;
};

/**
 * Format file size for display
 */
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

/**
 * Generate a shareable dashboard URL
 */
export const generateDashboardURL = (dashboardId: string): string => {
  const baseURL = window.location.origin;
  return `${baseURL}/dashboard/${dashboardId}`;
};

/**
 * Copy text to clipboard
 */
export const copyToClipboard = async (text: string): Promise<boolean> => {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text);
      return true;
    } else {
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = text;
      textArea.style.position = 'fixed';
      textArea.style.left = '-999999px';
      textArea.style.top = '-999999px';
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();
      const success = document.execCommand('copy');
      textArea.remove();
      return success;
    }
  } catch (error) {
    console.error('Failed to copy to clipboard:', error);
    return false;
  }
};