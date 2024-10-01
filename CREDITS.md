# Credits and Acknowledgements

Below is a list of various projects and code-snippets used in this golang-starter repo, or were otherwise instrumental in the building of this starter. Thanks to everyone involved in these projects!

## Overall
- <https://github.com/MartinHeinz/go-project-blueprint>
- <https://github.com/tslamic/go-starter>

## Golang Project Layout
- <https://github.com/golang-standards/project-layout>

## Golang Module Usage and Upgrading
- <https://github.com/jaegertracing/jaeger>
- <https://stackoverflow.com/a/60675491>
- <https://golang.cafe/blog/how-to-upgrade-golang-dependencies.html>

## SBOM and Signing
- <https://carlosbecker.com/posts/goreleaser-cosign/>
- <https://goreleaser.com/customization/sign/>
- <https://goreleaser.com/customization/docker_sign/>
- <https://github.com/marketplace/actions/cosign-installer>

## Goreleaser
- publish rpm, deb, apk packages to GitHub releases
    - <https://nfpm.goreleaser.com/>
- docker images
    - how to include license in built binary
        - <https://carlosbecker.com/posts/goreleaser-actions-podman/>
- full example
    - <https://github.com/goreleaser/nfpm/blob/main/.goreleaser.yml>
- shell completions and man-page generation
    - <https://carlosbecker.com/posts/golang-completions-cobra/>
- reproducible builds
    - <https://carlosbecker.com/posts/goreleaser-reproducible-buids/> 

## Make
- general Make tips and sane defaults
    - <https://tech.davis-hansson.com/p/make/>
- help text generator
    - <https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html>
- checkmake
	- <https://github.com/mrtazz/checkmake/issues/25>

## GitHub Actions
- generate and auto-commit docs
    - <https://github.com/goreleaser/nfpm/blob/main/.github/workflows/generate.yml>
- build using goreleaser
    - <https://github.com/goreleaser/nfpm/blob/main/.github/workflows/build.yml>
- security scanning
    - <https://github.com/goreleaser/goreleaser/tree/main/.github/workflows>
	- <https://github.com/aquasecurity/trivy-action>

## Pre-Commit
- <https://github.com/dnephin/pre-commit-golang>
- <https://github.com/tekwizely/pre-commit-golang>
- <https://github.com/koalaman/shellcheck-precommit>
- <https://github.com/hadolint/hadolint>
- <https://github.com/mrtazz/checkmake/pull/69>
- <https://github.com/trussworks/pre-commit-hooks#goreleaser-check>
