-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types for positions and roles
CREATE TYPE player_position AS ENUM ('Forward', 'Midfielder', 'Defender', 'Goalkeeper');
CREATE TYPE player_role AS ENUM ('Starter', 'Substitute');

-- Add updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create validate_match_player trigger function
CREATE OR REPLACE FUNCTION validate_match_player()
RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM matches m
        WHERE m.id = NEW.match_id 
        AND (m.home_team = NEW.team_id OR m.away_team = NEW.team_id)
    ) THEN
        RAISE EXCEPTION 'team_id must match either home_team or away_team of the match';
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create players table
CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create matches table
CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    home_team UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    away_team UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    home_score INTEGER NOT NULL DEFAULT 0,
    away_score INTEGER NOT NULL DEFAULT 0,
    home_mvp UUID REFERENCES players(id) ON DELETE SET NULL,
    away_mvp UUID REFERENCES players(id) ON DELETE SET NULL,
    result VARCHAR(10), -- WIN, LOSS, DRAW
    duration INTEGER, -- Match duration in minutes
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT different_teams CHECK (home_team != away_team)
);

-- Create match_players junction table
CREATE TABLE IF NOT EXISTS match_players (
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
    player_id UUID REFERENCES players(id) ON DELETE CASCADE,
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    position player_position NOT NULL,
    role player_role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (match_id, player_id)
);

-- Create trigger to validate team_id in match_players
CREATE TRIGGER validate_match_player_trigger
    BEFORE INSERT OR UPDATE ON match_players
    FOR EACH ROW
    EXECUTE FUNCTION validate_match_player();

-- Create indexes for better query performance
CREATE INDEX idx_matches_date ON matches(date DESC);
CREATE INDEX idx_match_players_player ON match_players(player_id);
CREATE INDEX idx_match_players_team ON match_players(team_id);

-- Add updated_at triggers
CREATE TRIGGER update_teams_updated_at
    BEFORE UPDATE ON teams
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_players_updated_at
    BEFORE UPDATE ON players
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_matches_updated_at
    BEFORE UPDATE ON matches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Enable Row Level Security
ALTER TABLE teams ENABLE ROW LEVEL SECURITY;
ALTER TABLE players ENABLE ROW LEVEL SECURITY;
ALTER TABLE matches ENABLE ROW LEVEL SECURITY;
ALTER TABLE match_players ENABLE ROW LEVEL SECURITY;

-- Create RLS policies
CREATE POLICY "Enable read access for all users" ON teams FOR SELECT USING (true);
CREATE POLICY "Enable read access for all users" ON players FOR SELECT USING (true);
CREATE POLICY "Enable read access for all users" ON matches FOR SELECT USING (true);
CREATE POLICY "Enable read access for all users" ON match_players FOR SELECT USING (true);

CREATE POLICY "Enable insert for authenticated users only" ON teams FOR INSERT WITH CHECK (auth.role() = 'authenticated');
CREATE POLICY "Enable insert for authenticated users only" ON players FOR INSERT WITH CHECK (auth.role() = 'authenticated');
CREATE POLICY "Enable insert for authenticated users only" ON matches FOR INSERT WITH CHECK (auth.role() = 'authenticated');
CREATE POLICY "Enable insert for authenticated users only" ON match_players FOR INSERT WITH CHECK (auth.role() = 'authenticated');

-- Grant permissions to authenticated users
GRANT SELECT, INSERT ON teams TO authenticated;
GRANT SELECT, INSERT ON players TO authenticated;
GRANT SELECT, INSERT ON matches TO authenticated;
GRANT SELECT, INSERT ON match_players TO authenticated; 