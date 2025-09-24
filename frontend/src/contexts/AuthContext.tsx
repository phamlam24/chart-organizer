import React, { createContext, useContext, useState, useEffect } from 'react';
import { authService } from '../services/authService';
import type { User, AuthContextType } from '../types';

const AuthContext = createContext<AuthContextType | null>(null);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: React.ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const storedToken = authService.getAuthToken();
    if (storedToken) {
      setToken(storedToken);
      // In a real app, you might want to validate the token or fetch user data
    }
    setIsLoading(false);
  }, []);

  const login = async (username: string, password: string) => {
    setIsLoading(true);
    try {
      const response = await authService.login({ username, password });
      const { jwtToken } = response;
      
      authService.setAuthToken(jwtToken);
      setToken(jwtToken);
      
      // Decode JWT to get user info (simplified - in production use a proper JWT library)
      const payload = JSON.parse(atob(jwtToken.split('.')[1]));
      setUser({ id: payload.user_id, username, created_at: new Date().toISOString() });
    } catch (error) {
      throw error; // Re-throw the error so the Login component can handle it
    } finally {
      setIsLoading(false);
    }
  };

  const signup = async (username: string, password: string) => {
    setIsLoading(true);
    try {
      const response = await authService.signup({ username, password });
      const { jwtToken } = response;
      
      authService.setAuthToken(jwtToken);
      setToken(jwtToken);
      
      // Decode JWT to get user info (simplified - in production use a proper JWT library)
      const payload = JSON.parse(atob(jwtToken.split('.')[1]));
      setUser({ id: payload.user_id, username, created_at: new Date().toISOString() });
    } catch (error) {
      throw error; // Re-throw the error so the Signup component can handle it
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    authService.removeAuthToken();
    setToken(null);
    setUser(null);
  };

  const value: AuthContextType = {
    user,
    token,
    login,
    signup,
    logout,
    isLoading,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};