-- First, drop everything in the correct order
DO $$ 
BEGIN
    -- Drop triggers
    DROP TRIGGER IF EXISTS update_player_timestamp ON players;
    DROP TRIGGER IF EXISTS update_player_stats_timestamp ON player_stats;
    DROP TRIGGER IF EXISTS update_player_affinities_timestamp ON player_affinities;
    DROP TRIGGER IF EXISTS update_player_abilities_timestamp ON player_abilities;

    -- Drop functions
    DROP FUNCTION IF EXISTS update_updated_at_column();

    -- Drop tables
    DROP TABLE IF EXISTS player_abilities CASCADE;
    DROP TABLE IF EXISTS events CASCADE;
    DROP TABLE IF EXISTS team_synergy CASCADE;
    DROP TABLE IF EXISTS match_analysis CASCADE;
    DROP TABLE IF EXISTS match_players CASCADE;
    DROP TABLE IF EXISTS team_compositions CASCADE;
    DROP TABLE IF EXISTS matches CASCADE;
    DROP TABLE IF EXISTS play_styles CASCADE;
    DROP TABLE IF EXISTS position_performance CASCADE;
    DROP TABLE IF EXISTS player_affinities CASCADE;
    DROP TABLE IF EXISTS player_versatility CASCADE;
    DROP TABLE IF EXISTS player_stats CASCADE;
    DROP TABLE IF EXISTS players CASCADE;

    -- Drop types
    DROP TYPE IF EXISTS position_type CASCADE;
    DROP TYPE IF EXISTS event_type CASCADE;
    DROP TYPE IF EXISTS match_result CASCADE;
    DROP TYPE IF EXISTS gender_type CASCADE;
    DROP TYPE IF EXISTS ability_level CASCADE;
END $$;

-- Create types
CREATE TYPE position_type AS ENUM ('forward', 'midfielder', 'defender', 'goalkeeper');
CREATE TYPE event_type AS ENUM ('goal', 'assist', 'pass', 'tackle', 'save', 'interception', 'yellow_card', 'red_card', 'substitution', 'performance');
CREATE TYPE match_result AS ENUM ('home_win', 'away_win', 'draw');
CREATE TYPE gender_type AS ENUM ('male', 'female');
CREATE TYPE ability_level AS ENUM ('basic', 'intermediate', 'advanced');

-- Create tables
CREATE TABLE players (
    id VARCHAR(40) PRIMARY KEY,  -- Changed from UUID
    name VARCHAR(100) NOT NULL,
    position position_type NOT NULL,
    gender gender_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Player stats table
CREATE TABLE player_stats (
    player_id VARCHAR(40) PRIMARY KEY REFERENCES players(id),
    goals_scored INTEGER DEFAULT 0,
    assists INTEGER DEFAULT 0,
    clean_sheets INTEGER DEFAULT 0,
    matches_played INTEGER DEFAULT 0,
    average_rating FLOAT DEFAULT 0,
    win_rate FLOAT DEFAULT 0,
    pass_accuracy FLOAT DEFAULT 0,
    ball_possession FLOAT DEFAULT 0,
    interceptions INTEGER DEFAULT 0,
    tackles INTEGER DEFAULT 0,
    distance_covered FLOAT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Player versatility table
CREATE TABLE player_versatility (
    player_id VARCHAR(40) REFERENCES players(id),
    position position_type,
    PRIMARY KEY (player_id, position)
);

-- Player affinities table
CREATE TABLE player_affinities (
    player_id VARCHAR(40) REFERENCES players(id),
    target_player_id VARCHAR(40) REFERENCES players(id),
    compatibility_score FLOAT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (player_id, target_player_id)
);

-- Position performance table
CREATE TABLE position_performance (
    player_id VARCHAR(40) REFERENCES players(id),
    position position_type,
    performance_score FLOAT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (player_id, position)
);

-- Play styles table
CREATE TABLE play_styles (
    player_id VARCHAR(40) REFERENCES players(id),
    style VARCHAR(50),
    PRIMARY KEY (player_id, style)
);

-- Matches table
CREATE TABLE matches (
    id VARCHAR(40) PRIMARY KEY,  -- Changed from UUID
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    home_score INTEGER DEFAULT 0,
    away_score INTEGER DEFAULT 0,
    home_mvp VARCHAR(40) REFERENCES players(id),
    away_mvp VARCHAR(40) REFERENCES players(id),
    result match_result,
    duration INTEGER, -- in minutes
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Team compositions table
CREATE TABLE team_compositions (
    match_id VARCHAR(40) REFERENCES matches(id),
    is_home BOOLEAN,
    formation VARCHAR(10),
    avg_rating FLOAT,
    chemistry FLOAT,
    PRIMARY KEY (match_id, is_home)
);

-- Match players table
CREATE TABLE match_players (
    match_id VARCHAR(40) REFERENCES matches(id),
    player_id VARCHAR(40) REFERENCES players(id),
    is_home BOOLEAN,
    position position_type,
    role VARCHAR(50),
    PRIMARY KEY (match_id, player_id)
);

-- Match analysis table
CREATE TABLE match_analysis (
    match_id VARCHAR(40) PRIMARY KEY REFERENCES matches(id),
    possession_home FLOAT,
    possession_away FLOAT,
    shots_home INTEGER,
    shots_away INTEGER,
    passes_home INTEGER,
    passes_away INTEGER,
    fouls_home INTEGER,
    fouls_away INTEGER
);

-- Team synergy table
CREATE TABLE team_synergy (
    match_id VARCHAR(40) REFERENCES matches(id),
    player1_id VARCHAR(40) REFERENCES players(id),
    player2_id VARCHAR(40) REFERENCES players(id),
    synergy_score FLOAT,
    PRIMARY KEY (match_id, player1_id, player2_id)
);

-- Events table
CREATE TABLE events (
    id VARCHAR(40) PRIMARY KEY,  -- Changed from UUID
    match_id VARCHAR(40) REFERENCES matches(id),
    type event_type NOT NULL,
    player_id VARCHAR(40) REFERENCES players(id),
    target_player_id VARCHAR(40) REFERENCES players(id),
    description TEXT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    rating FLOAT,
    position position_type,
    x_coord FLOAT,
    y_coord FLOAT,
    success BOOLEAN DEFAULT true,
    impact FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Player abilities table
CREATE TABLE player_abilities (
    player_id VARCHAR(40) REFERENCES players(id),
    ability_name VARCHAR(50),
    ability_level ability_level NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (player_id, ability_name)
);

-- Create the timestamp update function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers after all tables exist
CREATE TRIGGER update_player_timestamp
    BEFORE UPDATE ON players
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_player_stats_timestamp
    BEFORE UPDATE ON player_stats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_player_affinities_timestamp
    BEFORE UPDATE ON player_affinities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_player_abilities_timestamp
    BEFORE UPDATE ON player_abilities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes last
CREATE INDEX idx_events_match_id ON events(match_id);
CREATE INDEX idx_events_player_id ON events(player_id);
CREATE INDEX idx_player_stats_rating ON player_stats(average_rating);
CREATE INDEX idx_matches_date ON matches(date);
CREATE INDEX idx_player_affinities_score ON player_affinities(compatibility_score);