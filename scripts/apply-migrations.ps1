# Load environment variables from .env file
if (Test-Path .env) {
    Get-Content .env | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($key, $value)
        }
    }
}

# Check for required environment variables
if (-not $env:SUPABASE_ACCESS_TOKEN -or -not $env:SUPABASE_PROJECT_ID) {
    Write-Error "Missing required environment variables"
    Write-Error "Please ensure SUPABASE_ACCESS_TOKEN and SUPABASE_PROJECT_ID are set in .env file"
    exit 1
}

# Install Supabase CLI if not present
if (-not (Get-Command supabase -ErrorAction SilentlyContinue)) {
    Write-Output "Installing Supabase CLI..."
    # For Windows, we'll use scoop to install supabase
    if (-not (Get-Command scoop -ErrorAction SilentlyContinue)) {
        Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
        Invoke-RestMethod get.scoop.sh | Invoke-Expression
    }
    scoop bucket add supabase https://github.com/supabase/scoop-bucket.git
    scoop install supabase
}

# Link to your Supabase project
supabase link --project-ref $env:SUPABASE_PROJECT_ID --password $env:SUPABASE_DB_PASSWORD

# Apply migrations
supabase db push

Write-Output "Migrations applied successfully!" 