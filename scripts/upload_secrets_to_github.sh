#!/usr/bin/env bash

# Helper function for error handling
function handle_error {
    echo "Error: $1"
    exit 1
}

# Validate that .env exists
if [[ ! -f .env ]]; then
    handle_error ".env file not found. Ensure it exists before running this script."
fi

# Source .env file's environment variables
# shellcheck disable=SC2046
export $(cat .env)

# Read GitHub username and token from the environment
GITHUB_USERNAME="${GITHUB_USERNAME:-}"
GITHUB_TOKEN="${GH_TOKEN:-}"

if [[ -z "$GITHUB_USERNAME" ]]; then
    handle_error "GITHUB_USERNAME is not set in the environment. Please set it in .env."
elif [[  -z "$GITHUB_TOKEN" ]]; then
    handle_error "GITHUB_TOKEN is not set in the environment. Please set it in .env."
fi

# Helper function to upload secrets to GitHub Actions
upload_secrets_to_github() {
    echo "Pushing .env entries to GitHub Actions secrets for repo: $GITHUB_USERNAME/$REPO_NAME..."
    gh secret set --repo "$GITHUB_USERNAME"/"$REPO_NAME" --app actions --env-file .env
    gh secret set COSIGN_PRIVATE_KEY --repo "$GITHUB_USERNAME"/"$REPO_NAME" --app actions < "$REPO_NAME.key"
    echo "Secrets successfully uploaded to GitHub Actions."
}

# Helper function to upload secrets to GitHub secrets for use by Dependabot
upload_secrets_to_dependabot() {
    echo "Pushing .env entries to GitHub secrets for use by Dependabot for repo: $GITHUB_USERNAME/$REPO_NAME..."
    gh secret set --repo "$GITHUB_USERNAME"/"$REPO_NAME" --app dependabot --env-file .env
    echo "Secrets successfully uploaded to GitHub Dependabot."
}

# Main script logic
REPO_NAME="$1"

if [[ -z "$REPO_NAME" ]]; then
    handle_error "Usage: $0 <repo_name>"
fi

# Execute the functions to upload secrets
upload_secrets_to_github
upload_secrets_to_dependabot
