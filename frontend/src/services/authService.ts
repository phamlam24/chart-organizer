import apiClient from './apiClient';
import type { LoginRequest, LoginResponse, SignupRequest, SignupResponse } from '../types';

export const authService = {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post('/contracts.auth.v1.AuthService/Login', credentials);
    return response.data;
  },

  async signup(credentials: SignupRequest): Promise<SignupResponse> {
    const response = await apiClient.post('/contracts.auth.v1.AuthService/Signup', credentials);
    return response.data;
  },

  setAuthToken(token: string) {
    localStorage.setItem('auth_token', token);
  },

  getAuthToken(): string | null {
    return localStorage.getItem('auth_token');
  },

  removeAuthToken() {
    localStorage.removeItem('auth_token');
  }
};