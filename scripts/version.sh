#!/bin/bash
# Version management script for semantic versioning

set -e

VERSION_FILE="VERSION"
CHANGELOG_FILE="CHANGELOG.md"

function get_current_version() {
    cat "$VERSION_FILE" | sed 's/-dev$//'
}

function bump_version() {
    local current=$1
    local part=$2

    local major minor patch
    IFS='.' read -r major minor patch <<< "$current"

    case $part in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            echo "Invalid version part: $part (use major, minor, or patch)"
            exit 1
            ;;
    esac

    echo "${major}.${minor}.${patch}"
}

function update_version_file() {
    local version=$1
    echo "$version" > "$VERSION_FILE"
    echo "Updated $VERSION_FILE to $version"
}

function update_changelog() {
    local version=$1
    local date=$(date +%Y-%m-%d)

    # Replace [Unreleased] with the new version
    sed -i.bak "s/## \[Unreleased\]/## [Unreleased]\n\n## [$version] - $date/" "$CHANGELOG_FILE"

    # Update comparison links at bottom
    sed -i.bak "s|\[Unreleased\]:.*|[Unreleased]: https://github.com/scttfrdmn/nrp-aws-kip/compare/v$version...HEAD\n[$version]: https://github.com/scttfrdmn/nrp-aws-kip/releases/tag/v$version|" "$CHANGELOG_FILE"

    rm -f "${CHANGELOG_FILE}.bak"
    echo "Updated $CHANGELOG_FILE with version $version"
}

function create_git_tag() {
    local version=$1
    git add "$VERSION_FILE" "$CHANGELOG_FILE"
    git commit -m "chore: bump version to $version"
    git tag -a "v$version" -m "Release v$version"
    echo "Created git tag v$version"
    echo "Push with: git push && git push --tags"
}

# Main script
if [ $# -lt 1 ]; then
    echo "Usage: $0 <major|minor|patch> [--no-tag]"
    echo "Current version: $(get_current_version)"
    exit 1
fi

PART=$1
NO_TAG=false

if [ "${2:-}" = "--no-tag" ]; then
    NO_TAG=true
fi

current=$(get_current_version)
new_version=$(bump_version "$current" "$PART")

echo "Bumping version from $current to $new_version"

update_version_file "$new_version"
update_changelog "$new_version"

if [ "$NO_TAG" = false ]; then
    create_git_tag "$new_version"
else
    echo "Skipping git tag creation (--no-tag specified)"
    echo "Don't forget to commit and tag manually!"
fi

echo "Version bump complete!"
