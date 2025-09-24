# Chart Organizer

This document provides comprehensive instructions for running the Chart Organizer application using Docker and docker-compose.

## Prerequisites

- Docker (20.10+)
- Docker Compose (2.0+)
- Git (to clone the repository)

## Quick Start

1. **Clone the repository** (if you haven't already):
   ```bash
   git clone <repository-url>
   cd chart-organizer
   ```

2. **Create storage directory**:
   ```bash
   mkdir -p storage
   ```

3. **Run the application**:
   ```bash
   docker-compose up --build
   ```

4. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## Services

### Backend
- **Port**: 8080
- **Technology**: Go 1.24.2 with SQLite database
- **API**: Connect/gRPC-Web protocol
- **Authentication**: JWT-based
- **Storage**: SQLite database + file system for CSV uploads

### Frontend  
- **Port**: 3000 (served via Nginx)
- **Technology**: React 19 with TypeScript, Vite build tool
- **Features**: Interactive data visualization dashboard creator
- **Charts**: Parallel coordinates, scatter plots, line plots

## Environment Variables

### Backend
- `JWT_KEY`: Secret key for JWT token signing (default: auto-generated)
- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database path (default: ./storage/chart-organizer.db)
- `DATASET_STORAGE_PATH`: Directory for uploaded CSV files (default: ./storage/datasets)

### Frontend
- `VITE_API_BASE_URL`: Backend API URL (default: http://localhost:8080)

## Data Persistence

- **Database**: SQLite database stored in `./storage/chart-organizer.db`
- **Uploads**: CSV files stored in `./storage/datasets/`
- **Volumes**: The `./storage` directory is mounted to ensure data persistence

## Development Mode

For development with hot reloading:

```bash
docker-compose -f docker-compose.dev.yml up --build
```

This provides:
- Go backend with source watching
- Vite dev server for frontend (port 5173)
- Hot reloading for both frontend and backend

### Manual Development Setup

**Backend** (Go):
```bash
cd backend
go mod download
go run cmd/main.go
```

**Frontend** (React):
```bash
cd frontend
npm install
npm run dev
```

## Production Deployment

### Security Considerations

1. **Change the JWT secret**:
   ```bash
   export JWT_KEY="your-super-secure-production-key"
   docker-compose up -d
   ```

2. **Use environment file**:
   Create `.env` file:
   ```env
   JWT_KEY=your-super-secure-production-key
   VITE_API_BASE_URL=https://your-api-domain.com
   ```

### Deployment Options

1. **Single server deployment**: Use the provided docker-compose.yml
2. **Reverse proxy setup**: Configure nginx/Apache to proxy to containers
3. **Container orchestration**: Deploy with Kubernetes, Docker Swarm, etc.
4. **Cloud deployment**: Use AWS ECS, Google Cloud Run, Azure Container Instances

### Database Migration

For production, consider migrating from SQLite to PostgreSQL:

1. Update backend to support PostgreSQL
2. Add PostgreSQL service to docker-compose.yml
3. Update environment variables accordingly

## Stopping the Application

**Normal stop:**
```bash
docker-compose down
```

**Stop and remove volumes (⚠️ deletes all data):**
```bash
docker-compose down -v
```

## Troubleshooting

### Common Issues

1. **Port conflicts**: 
   - Modify ports in `docker-compose.yml`
   - Kill processes using ports 3000/8080

2. **Database issues**: 
   ```bash
   rm -rf storage/
   mkdir storage
   docker-compose up --build
   ```

3. **Build issues**:
   ```bash
   docker-compose build --no-cache
   docker-compose up
   ```

4. **Permission issues (Linux/macOS)**:
   ```bash
   sudo chown -R $USER:$USER storage/
   ```

### Debugging

**View logs:**
```bash
docker-compose logs backend
docker-compose logs frontend
```

**Execute into containers:**
```bash
docker-compose exec backend sh
docker-compose exec frontend sh
```

## Application Features

### Core Functionality
- ✅ User authentication (signup/login with JWT)
- ✅ CSV dataset upload and management
- ✅ Interactive data visualization creation
- ✅ Dashboard builder with drag-and-drop interface
- ✅ Shareable dashboard links
- ✅ Public dashboard viewing (no auth required)

### Visualization Types
- **Parallel Coordinates**: Multi-dimensional data exploration
- **Scatter Plots**: Correlation analysis between variables
- **Line Plots**: Time-series and trend visualization

### Technical Features
- React-based responsive UI
- Real-time chart updates
- Data export capabilities
- RESTful API with Connect protocol
- File upload with validation
- Session management

## API Documentation

The backend provides these main endpoint groups:

### Authentication (`/contracts.auth.v1.AuthService/`)
- `Signup` - Create new user account
- `Login` - Authenticate existing user
- JWT token-based session management

### Dataset Management (`/contracts.dataset.v1.DatasetService/`)
- `UploadDataset` - Upload CSV files
- `GetAllDatasetsFromUser` - List user's datasets  
- `GetDataset` - Retrieve specific dataset

### Dashboard & Visualization (`/contracts.viz.v1.DashboardService/`)
- `CreateDashboard` - Create new dashboard with charts
- `GetDashboard` - Retrieve dashboard configuration

Refer to the `contracts/` directory for detailed protobuf specifications.

## File Structure

```
chart-organizer/
├── docker-compose.yml          # Production setup
├── docker-compose.dev.yml      # Development setup  
├── backend/
│   ├── Dockerfile             # Backend container config
│   ├── cmd/main.go           # Application entry point
│   └── internal/             # Application logic
├── frontend/
│   ├── Dockerfile            # Frontend container config
│   ├── nginx.conf           # Production nginx config
│   └── src/                 # React application
└── storage/                 # Persistent data (created automatically)
    ├── chart-organizer.db   # SQLite database
    └── datasets/            # Uploaded CSV files
```