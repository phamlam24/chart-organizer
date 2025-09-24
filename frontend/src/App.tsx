import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import Layout from './components/Layout/Layout'
import Login from './pages/Login/Login'
import Signup from './pages/Signup/Signup'
import Dashboard from './pages/Dashboard/Dashboard'
import DatasetUpload from './pages/DatasetUpload/DatasetUpload'
import DashboardCreator from './pages/DashboardCreator/DashboardCreator'
import PublicDashboard from './pages/PublicDashboard/PublicDashboard'
import ProtectedRoute from './components/ProtectedRoute/ProtectedRoute'
import './App.css'

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<Signup />} />
          <Route path="/dashboard/:id" element={<PublicDashboard />} />
          <Route element={<Layout />}>
            <Route 
              path="/" 
              element={
                <ProtectedRoute>
                  <Dashboard />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="/upload" 
              element={
                <ProtectedRoute>
                  <DatasetUpload />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="/create-dashboard/:datasetId" 
              element={
                <ProtectedRoute>
                  <DashboardCreator />
                </ProtectedRoute>
              } 
            />
          </Route>
        </Routes>
      </Router>
    </AuthProvider>
  )
}

export default App
