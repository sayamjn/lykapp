# Video Player with Dynamic Ad Overlays

A web application featuring a video player with dynamically rotating ad overlays, click tracking, and real-time analytics.

## Features

- Video player with custom controls
- Dynamic ad overlays with configurable positions
- Real-time click tracking and analytics
- Comprehensive logging system
- SQLite database for persistence
- Docker support

## Tech Stack

### Backend
- Go 1.21
- SQLite for data persistence
- Unsplash API for dynamic ad images

### Frontend
- React + Vite
- TailwindCSS for styling

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- npm or yarn
- Docker and Docker Compose (for production deployment)
- Unsplash API key

## Local Development Setup

1. Clone the repository:
```bash
git clone https://github.com/sayamjn/lykapp.git
cd lykapp
```

2. Set up environment variables:
```bash
# Backend (.env)
PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:5173
UNSPLASH_ACCESS_KEY=your_access_key
AD_CACHE_TIMEOUT=5m
AD_REFRESH_ENABLED=true

# Frontend (.env.local)
VITE_API_BASE_URL=http://localhost:8080/api
```

3. Start the backend server:
```bash
cd backend
go mod download
go run cmd/server/main.go
```

4. Start the frontend development server:
```bash
cd frontend
npm install
npm run dev
```

## Production Deployment

1. Build and run using Docker Compose:
```bash
docker-compose up --build -d
```

2. Access the Frontend at `http://localhost:5173`

## System Design

### Architecture

```
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│   Frontend  │ ←── │   Backend    │ ←── │  Unsplash  │
│   (React)   │     │     (Go)     │     │    API     │
└─────────────┘     └──────────────┘     └────────────┘
      ↓                    ↓                    
┌─────────────┐     ┌──────────────┐     
│  Analytics  │     │    SQLite    │     
│  Dashboard  │     │   Database   │     
└─────────────┘     └──────────────┘     
```

### Components

#### Backend

1. **API Server**
   - RESTful endpoints for ads and click tracking
   - CORS configuration
   - Request rate limiting
   - Error handling middleware

2. **Store Layer**
   - SQLite database integration
   - In-memory caching for ads
   - Thread-safe operations
   - Analytics data persistence

3. **Logging System**
   - Structured logging
   - Request/response logging
   - Error tracking
   - Analytics events

#### Frontend

1. **Video Player**
   - Custom playback controls
   - Ad overlay management
   - Click tracking

2. **Analytics Dashboard**
   - Real-time metrics
   - Click visualization
   - User activity tracking
   - Auto-refresh functionality

## API Endpoints

### GET /api/ads
Returns list of current ads

### POST /api/ads/click
Records ad click with tracking data

### GET /api/ads/clicks
Returns click analytics data

Run frontend tests:
```bash
cd frontend
npm test
```

## Logging

Logs are written to `logs/app.log` by default. The logging system captures:
- HTTP requests/responses
- Ad click events
- Error events
- Performance metrics
- System events

## Monitoring

The application provides:
- Real-time analytics dashboard
- Structured logs for monitoring
- Error tracking and reporting

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.