# ⚽ Soccer Team Builder Service

A data-driven REST API service for tracking player performance and generating balanced teams. Unlike traditional team management systems, this service focuses on individual player tracking across matches to optimize team composition using AI.

## 🎯 Core Concepts

- **Dynamic Team Formation**: Teams are not persistent entities but rather generated compositions based on player statistics and compatibility
- **Player Performance Tracking**: Comprehensive tracking of player stats, performance metrics, and interaction data across all matches
- **AI-Ready Data Structure**: All data is structured to facilitate machine learning analysis for team balancing
- **Player Affinity**: Tracks how well players perform together to identify optimal team compositions

### 📊 Key Metrics Tracked

- Individual performance stats (goals, assists, etc.)
- Player-to-player interaction success rates
- Position-based performance metrics
- Historical team composition effectiveness
- Player versatility across positions
- Match outcome influence factors

### 🤖 Advanced Matchmaking System

The service uses a sophisticated AI-driven matchmaking system that considers multiple factors:

#### Player Ability Scoring
- **Word-Based Levels**: Players' abilities are rated on a 5-level scale:
  - Beginner (0.2)
  - Developing (0.4)
  - Intermediate (0.6)
  - Advanced (0.8)
  - Expert (1.0)

#### Position-Specific Weighting
Each ability is weighted differently based on position importance:
```
Example: Ball Control
- Forward:    1.0  (Critical for control under pressure)
- Midfielder: 0.95 (Essential for midfield play)
- Defender:   0.7  (Important but not critical)
- Goalkeeper: 0.3  (Less important)
```

#### Team Generation Algorithm
1. **Gender Balance** (for mixed matches)
   - Ensures equal distribution of male/female players
   - Maintains position balance within gender groups

2. **Formation Selection**
   - Analyzes available players' positions
   - Considers player versatility (primary/secondary positions)
   - Adapts to player pool composition

3. **Player Distribution**
   - Primary position matching (100% effectiveness)
   - Secondary position consideration (80% effectiveness)
   - Out-of-position placement (50% effectiveness)

4. **Team Chemistry Calculation**
   - Historical player partnerships
   - Complementary abilities analysis
   - Position-specific synergies

5. **Balancing Factors**
   - Overall team rating
   - Position coverage
   - Skill distribution
   - Physical/tactical balance

#### Complementary Abilities System
The system identifies and rewards complementary skill pairs:
- Speed + Ball Control
- Game Reading + Positioning
- Communication + Teamwork
- Heading + Jump/Strength
- etc.

#### Versatility Bonus
Players receive bonuses for:
- Multiple position capabilities
- Balanced ability sets
- Leadership skills in key positions
- Tactical flexibility

This advanced system ensures balanced, competitive matches while considering both individual skills and team dynamics.

## 🏗️ Architecture

The project follows a clean architecture pattern optimized for data collection and analysis:

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

- **cmd/api**: Main application with Echo server setup
- **internal/api**: HTTP handlers and routes
- **internal/core**: Business logic and domain models
  - `models/`: Data structures optimized for AI processing
  - `services/`: Performance analysis and team generation logic
- **internal/repository**: PostgreSQL persistence layer
- **pkg/utils**: Shared utilities and analysis helpers

## 🎮 Features

### Player Tracking
- Comprehensive performance metrics
- Player interaction success rates
- Position effectiveness mapping
- Historical team composition data

### Team Building
- AI-ready data structures for team optimization
- Player compatibility scoring
- Position-based team balancing
- Performance prediction modeling

### Match Analysis
- Detailed event tracking
- Team composition effectiveness
- Player contribution metrics
- Interaction pattern analysis

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

## 🤖 AI Integration Points

The service is designed to provide rich data for AI analysis:

- **Player Compatibility**: Track successful player combinations
- **Team Balance**: Historical performance of different team compositions
- **Position Optimization**: Player effectiveness in various positions
- **Performance Prediction**: Data structure ready for ML model training

## 🚧 Coming Soon
- Advanced player affinity metrics
- ML-based team composition suggestions
- Performance prediction endpoints
- Player chemistry analysis
- Automated team balancing recommendations

## 🗄️ Database Setup

1. Copy `.env.example` to `.env`:
```powershell
# PowerShell
Copy-Item .env.example .env

# Or using Command Prompt
copy .env.example .env
```

2. Fill in your Supabase credentials in `.env`:
- SUPABASE_PROJECT_ID: Found in Project Settings > General
- SUPABASE_DB_PASSWORD: Database password
- SUPABASE_ACCESS_TOKEN: Found in Project Settings > API

3. Apply migrations:
```powershell
# If using PowerShell script
.\scripts\apply-migrations.ps1

# If using Git Bash
./scripts/apply-migrations.sh
```