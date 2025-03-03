#!/bin/bash

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | awk '/=/ {print $1}')
fi

# Check for required environment variables
if [ -z "$SUPABASE_ACCESS_TOKEN" ] || [ -z "$SUPABASE_PROJECT_ID" ]; then
    echo "Error: Missing required environment variables"
    echo "Please ensure SUPABASE_ACCESS_TOKEN and SUPABASE_PROJECT_ID are set in .env file"
    exit 1
fi

# Install Supabase CLI if not present
if ! command -v supabase &> /dev/null; then
    echo "Installing Supabase CLI..."
    brew install supabase/tap/supabase
fi

# Link to your Supabase project
supabase link --project-ref "$SUPABASE_PROJECT_ID" --password "$SUPABASE_DB_PASSWORD"

# Apply migrations
supabase db push

echo "Migrations applied successfully!" 