# ⚽ Soccer Team Management Service

A robust REST API service for managing soccer teams, players, matches, and performance tracking.

## 🏗️ Architecture

The project follows a clean architecture pattern with the following structure:

```
soccer-service/
├── cmd/
│ └── api/
│ └── main.go # Application entry point
├── internal/
│ ├── api/
│ │ ├── handlers/ # HTTP request handlers
│ │ └── routes/ # Route definitions
│ ├── core/
│ │ ├── models/ # Domain models
│ │ └── services/ # Business logic
│ └── repository/ # Data persistence layer
├── pkg/
│ └── utils/ # Shared utilities
└── README.md
```

### 📦 Components

- **cmd/api**: Contains the main application entry point with Echo server setup
- **internal/api**: HTTP handlers and route definitions
  - `handlers/`: Request handlers for each domain entity
  - `routes/`: API route configuration
- **internal/core**: Business logic and domain models
  - `models/`: Data structures for Player, Match, and Event
  - `services/`: Business logic implementation
- **internal/repository**: Data persistence layer with thread-safe operations
- **pkg/utils**: Shared utilities and helpers

## 🎯 Current Implementation

### Models
- **Player**: Manages player information including stats
- **Match**: Handles match details and team information
- **Event**: Tracks match events (goals, cards, etc.)

### Repositories
- Thread-safe implementations using sync.RWMutex
- In-memory storage with CRUD operations
- Specialized queries (e.g., GetByMatch for events)

### Services
- Player management with CRUD operations
- Input validation and business logic

### API Endpoints

#### Players
- `POST /api/v1/players` - Create a new player
- `GET /api/v1/players` - List all players
- `GET /api/v1/players/{id}` - Get player details
- `PUT /api/v1/players/{id}` - Update player
- `DELETE /api/v1/players/{id}` - Delete player

### 🚀 Getting Started

1. Clone the repository
2. Initialize the Go module:
```bash
go mod init soccer-service
```

3. Install dependencies:
```bash
go mod tidy
```

4. Run the service:
```bash
go run cmd/api/main.go
```

### 📝 API Usage Examples

#### Create a Player
```bash
curl -X POST http://localhost:8080/api/v1/players \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "name": "John Doe",
    "position": "forward",
    "stats": {
      "goals_scored": 0,
      "assists": 0,
      "clean_sheets": 0,
      "matches_played": 0,
      "average_rating": 0
    }
  }'
```

#### Get Player Details
```bash
curl http://localhost:8080/api/v1/players/1
```

## 🛠️ Technical Details

### Dependencies
- Echo v4: High performance web framework
- UUID: For generating unique identifiers

### Thread Safety
- All repositories implement thread-safe operations using sync.RWMutex
- Concurrent access to data is properly handled

## 🎮 Fun Note

Just like in LoL, may your plays be as perfect as a pentakill on Summoner's Rift! 😎

## 🚧 Coming Soon
- Match management endpoints
- Event tracking system
- Performance ranking
- MVP validation
- Unit tests