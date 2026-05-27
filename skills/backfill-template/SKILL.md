---
name: backfill-template
description: Back-fills template updates and improvements from the golang-starter template repository into target projects, keeping them consistent, updated, and passing builds.
---

# Back-filling Template Updates from golang-starter

This skill guides the AI agent on how to back-fill scaffolding, build pipelines, and configuration updates from the `golang-starter` template repository into downstream target repositories. The goal is to keep the downstream projects consistent and updated with the latest optimizations, security practices, and workflow changes, without breaking any downstream-specific functionality.

## Core Guidelines

1. **Intelligent Merging, Not Overwriting**: Always inspect diffs carefully. Do not copy files wholesale if the target repository has customized behavior. Merge the template improvements (e.g. version bumps, extra linting checks, newer workflow actions) into the existing files.
2. **Dockerfile Similarity Matching**: Target repos might have renamed, added, or deleted Dockerfiles. Never rename target Dockerfiles. Identify the closest matching template Dockerfile by content structure, review the differences, and apply template improvements (such as updated Go version tags, nonroot user configurations, or shell flags) without wiping out custom packages, ports, or entrypoints.
3. **Build and Test Verification**: After applying changes, the project MUST build and pass checks. Validate using `make local` (local Go toolchain checks) and `make all` (Docker-based workflow checks).
4. **Local Staging and Commit**: Stage the changes (`git add`) and write a clear, descriptive commit message summarizing the backfilled changes. Commit locally but **do NOT push** to the remote.

---

## Step-by-Step Execution Plan

### Step 1: Run the Comparison Script
Run the helper Python script to fetch the latest `golang-starter` template and compare it against the current target repository. The script will automatically substitute placeholders (`golang-starter` and `toozej` with the target's project name and owner) and compute Dockerfile structural similarities.

Run the following command from the root of the target repository:
```bash
python3 <path_to_skill_folder>/scripts/compare.py
```
*(Note: If you are testing the template or have a local clone of the template repository, you can pass `--template-dir <path_to_template_dir>` to the script.)*

The script outputs:
- **[ADD]**: Files that exist in the template but are missing in the target (e.g. new linter configuration or pre-commit hooks).
- **[MODIFIED]**: Unified diffs showing the difference between the target file and the customized template file.
- **Dockerfile Matches**: Mappings showing which target Dockerfile is most similar to which template Dockerfile (along with the similarity ratio and diff).

---

### Step 2: Apply Scaffolding & Configuration Diffs
For each modified file, carefully apply the diffs:

- **`.github/workflows/` (GitHub Actions)**:
  - Back-fill new steps, updated action versions (e.g., `actions/checkout@v4`), and optimizations.
  - Preserve target-specific environment variables, secrets, and repository-specific steps.
  - *Note*: Ignore minor cron minute variations (which are randomized to distribute load).
- **`Makefile`**:
  - Merge new targets, build arguments, or linter integrations.
  - Preserve custom target-specific make targets, build tags, or special compiler flags.
- **`.goreleaser.yml`**:
  - Merge improved release stages, SBOM generation, or signing logic.
  - Preserve project-specific binary names, custom archive formats, or tap configurations.
- **Config & Version Boilerplate (`pkg/version/`, `pkg/config/`)**:
  - Update version parsing logic or standard config loaders.
  - Ensure the Go package paths match the target's module path.

---

### Step 3: Back-fill Dockerfile Updates
For each target Dockerfile matched by the script:
1. Identify the template file it was paired with (e.g., `Dockerfile.distroless`).
2. Review the diff. Typical template updates include:
   - Bumping the builder Go version (e.g., `FROM golang:1.26-bookworm`).
   - Bumping the runtime base images (e.g., `FROM gcr.io/distroless/static-debian13`).
   - Security configurations (e.g., `SHELL ["/bin/bash", "-o", "pipefail", "-c"]` or `USER nonroot`).
   - Better caching/multi-stage builds.
3. Apply these updates into the target Dockerfile.
4. **CRITICAL**: Do NOT change the target Dockerfile's filename, and do NOT remove target-specific logic (like custom `apt-get` packages, custom ports, volume mounts, or customized entrypoints/arguments).

---

### Step 4: Verification
Verify that the target repository builds, tests, and passes all checks. Run the following targets from the Makefile:

1. **Local Toolchain Verification**:
   ```bash
   make local
   ```
   This will run `go mod tidy`, `go mod vendor`, `go vet`, `pre-commit` hooks, run tests with race detection, generate coverage, build the local binary, and check goreleaser snapshots. Ensure everything compiles and all tests pass.

2. **Docker Workflow Verification**:
   ```bash
   make all
   ```
   This runs Docker-based vet, test, build, image signing verification with Cosign, and launches the container.
   
   > [!NOTE]
   > Running `make all` locally builds a fresh, unsigned Docker image, which will cause the `cosign verify` step to fail with `Error: no signatures found`. This is expected. As long as the `vet`, `pre-commit`, `test`, and `build` stages of the Docker build complete successfully, the Docker verification is considered successful. You can test running the container manually if needed with:
   > `docker run --rm --name <project_name> --env-file .env <owner>/<project_name>:latest`

If any check fails, troubleshoot the build/test issues and refine the backfilled configurations.

---

### Step 5: Git Commit
Once everything builds and passes perfectly:
1. Stage the modified files:
   ```bash
   git add .
   ```
2. Commit the changes locally:
   ```bash
   git commit -m "chore: backfill updates from golang-starter template

- Bumped Go and base image versions in Dockerfiles
- Updated GitHub Action workflows to latest versions
- Applied Makefile improvements and linter updates
- Synced pre-commit hooks and goreleaser configuration"
   ```
3. **DO NOT PUSH** the changes to the remote repository. Leave it staged for manual review and pushing by the developer.
