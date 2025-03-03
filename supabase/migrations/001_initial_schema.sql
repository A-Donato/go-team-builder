CREATE TABLE players (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    position VARCHAR(20) NOT NULL,
    goals_scored INT DEFAULT 0,
    assists INT DEFAULT 0,
    clean_sheets INT DEFAULT 0,
    matches_played INT DEFAULT 0,
    average_rating DECIMAL(3,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE matches (
    id VARCHAR(36) PRIMARY KEY,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) DEFAULT 'completed',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE match_players (
    match_id VARCHAR(36) REFERENCES matches(id),
    player_id VARCHAR(36) REFERENCES players(id),
    team VARCHAR(10) NOT NULL, -- 'home' or 'away'
    is_mvp BOOLEAN DEFAULT FALSE,
    performance_rating DECIMAL(3,2),
    PRIMARY KEY (match_id, player_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE match_scores (
    match_id VARCHAR(36) REFERENCES matches(id) PRIMARY KEY,
    home_score INT DEFAULT 0,
    away_score INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events (
    id VARCHAR(36) PRIMARY KEY,
    match_id VARCHAR(36) REFERENCES matches(id),
    type VARCHAR(20) NOT NULL,
    player_id VARCHAR(36) REFERENCES players(id),
    description TEXT,
    rating DECIMAL(3,2),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_event_type CHECK (
        type IN ('goal', 'yellow_card', 'red_card', 'substitution', 'performance')
    )
);

-- Player performance history for AI analysis
CREATE TABLE player_performance_history (
    id VARCHAR(36) PRIMARY KEY,
    player_id VARCHAR(36) REFERENCES players(id),
    match_id VARCHAR(36) REFERENCES matches(id),
    played_with_ids TEXT[], -- Array of player IDs they played with
    played_against_ids TEXT[], -- Array of player IDs they played against
    team VARCHAR(10) NOT NULL, -- 'home' or 'away'
    match_result VARCHAR(10), -- 'win', 'loss', 'draw'
    individual_rating DECIMAL(3,2),
    goals_scored INT DEFAULT 0,
    assists INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (player_id, match_id)
);

-- Indexes for better query performance
CREATE INDEX idx_match_players_match ON match_players(match_id);
CREATE INDEX idx_match_players_player ON match_players(player_id);
CREATE INDEX idx_events_match ON events(match_id);
CREATE INDEX idx_events_player ON events(player_id);
CREATE INDEX idx_events_type ON events(type);
CREATE INDEX idx_performance_history_player ON player_performance_history(player_id);
CREATE INDEX idx_performance_history_match ON player_performance_history(match_id);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_players_updated_at
    BEFORE UPDATE ON players
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_matches_updated_at
    BEFORE UPDATE ON matches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_match_scores_updated_at
    BEFORE UPDATE ON match_scores
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add more tables for matches and events here 