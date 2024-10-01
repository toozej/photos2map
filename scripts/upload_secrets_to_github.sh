#!/usr/bin/env bash
set -Eeuo pipefail

# Helper function for error handling
function handle_error {
    echo "Error: $1"
    exit 1
}

# Validate that .env exists
if [[ ! -f .env ]]; then
    handle_error ".env file not found. Ensure it exists before running this script."
fi

# Read GitHub username and token from the environment
GITHUB_USERNAME="${GITHUB_USERNAME:-}"
GITHUB_TOKEN="${GITHUB_TOKEN:-}"

if [[ -z "$GITHUB_USERNAME" || -z "$GITHUB_TOKEN" ]]; then
    handle_error "GITHUB_USERNAME or GITHUB_TOKEN is not set in the environment. Please set them in .env."
fi

# Helper function to upload secrets to GitHub Actions
upload_secrets_to_github() {
    echo "Pushing .env entries to GitHub Actions secrets for repo: $GITHUB_USERNAME/$REPO_NAME..."

    while IFS='=' read -r key value; do
        if [[ "$key" != "" ]]; then
            response=$(curl -s -X PUT \
              -H "Authorization: token $GITHUB_TOKEN" \
              -H "Content-Type: application/json" \
              -d "{\"encrypted_value\":\"$value\",\"key_id\":\"$key\"}" \
              "https://api.github.com/repos/$GITHUB_USERNAME/$REPO_NAME/actions/secrets/$key")

            if [[ "$response" == *"errors"* ]]; then
                handle_error "Failed to set secret $key in GitHub Actions. Response: $response"
            fi
        fi
    done < .env
    echo "Secrets successfully uploaded to GitHub Actions."
}

# Main script logic
REPO_NAME="$1"

if [[ -z "$REPO_NAME" ]]; then
    handle_error "Usage: $0 <repo_name>"
fi

# Execute the function to upload secrets
upload_secrets_to_github
