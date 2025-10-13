#!/usr/bin/env bash
set -Eeuo pipefail

# Helper function for error handling
function handle_error {
    echo "Error: $1"
    exit 1
}

store_envfile_in_1password() {
    echo "Storing EnvFile in 1Password..."

    # Validate that .env exists
    if [[ ! -f .env ]]; then
        handle_error ".env file not found. Ensure it exists before running this script."
    fi

    # Check if the item exists; if not, create it
    if ! op item get "${NEW_PROJECT_NAME}" &>/dev/null; then
        # Create the 1Password item with the project name
        op item create --category login --title "${NEW_PROJECT_NAME}" \
            --url "https://github.com/${GITHUB_USERNAME}/${NEW_PROJECT_NAME}" \
            --tags "Projects/${NEW_PROJECT_NAME}" || handle_error "Failed to create 1Password item."
    fi

    # Update the 1Password item with env file
    op item edit "${NEW_PROJECT_NAME}" \
        "EnvFile[file]=.env" \
        || handle_error "Failed to update 1Password item with envfile."

    echo "EnvFile successfully stored in 1Password."
}

# Helper function to store secrets in 1Password
store_secrets_in_1password() {
    echo "Storing secrets in 1Password..."

    # Check if the item exists; if not, create it
    if ! op item get "${NEW_PROJECT_NAME}" &>/dev/null; then
        # Create the 1Password item with the project name
        op item create --category login --title "${NEW_PROJECT_NAME}" \
            --url "https://github.com/${GITHUB_USERNAME}/${NEW_PROJECT_NAME}" \
            --tags "Projects/${NEW_PROJECT_NAME}" || handle_error "Failed to create 1Password item."
    fi

    # Update the 1Password item with generated secrets
    op item edit "${NEW_PROJECT_NAME}" \
        "Cosign.Passphrase[password]=${COSIGN_PASSPHRASE}" \
        "Cosign.Private Key[file]=${NEW_PROJECT_NAME}.key" \
        "Cosign.Public Key[file]=${NEW_PROJECT_NAME}.pub" \
        "GH PAT[password]=${GITHUB_TOKEN}" \
        || handle_error "Failed to update 1Password item with secrets."

    echo "Secrets successfully stored in 1Password."
}


NEW_PROJECT_NAME="${2}"
GITHUB_USERNAME="${3:-toozej}"

if [[ "${1}" == "secrets" ]]; then
    store_secrets_in_1password
elif [[ "${1}" == "envfile" ]]; then
    store_envfile_in_1password
fi
