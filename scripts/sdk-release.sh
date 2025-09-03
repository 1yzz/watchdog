#!/bin/bash

# SDK Release Script
# This script handles version bumping and publishing of the JavaScript SDK

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if we're in a git repository
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_error "Not in a git repository"
        exit 1
    fi
}

# Function to check if there are uncommitted changes
check_clean_working_dir() {
    if ! git diff-index --quiet HEAD --; then
        print_error "Working directory is not clean. Please commit or stash your changes."
        exit 1
    fi
}

# Function to check if we're on main branch
check_main_branch() {
    current_branch=$(git branch --show-current)
    if [ "$current_branch" != "main" ]; then
        print_warning "Not on main branch (currently on $current_branch)"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# Function to bump version
bump_version() {
    local version_type=$1
    local sdk_dir="sdk/javascript"
    
    print_status "Bumping version ($version_type)..."
    
    cd "$sdk_dir"
    
    case $version_type in
        "patch")
            npm version patch --no-git-tag-version
            ;;
        "minor")
            npm version minor --no-git-tag-version
            ;;
        "major")
            npm version major --no-git-tag-version
            ;;
        *)
            print_error "Invalid version type: $version_type"
            print_error "Valid types: patch, minor, major"
            exit 1
            ;;
    esac
    
    # Get the new version
    new_version=$(node -p "require('./package.json').version")
    print_success "Version bumped to $new_version"
    
    cd - > /dev/null
}

# Function to build the SDK
build_sdk() {
    print_status "Building SDK..."
    cd sdk/javascript
    npm run build
    cd - > /dev/null
    print_success "SDK built successfully"
}

# Function to run tests
run_tests() {
    print_status "Running tests..."
    cd sdk/javascript
    npm test
    cd - > /dev/null
    print_success "Tests passed"
}

# Function to publish to npm
publish_to_npm() {
    local sdk_dir="sdk/javascript"
    local test_mode=${1:-false}
    
    if [[ "$test_mode" == true ]]; then
        print_warning "TEST MODE: Skipping npm publish"
        return 0
    fi
    
    print_status "Publishing to npm..."
    
    cd "$sdk_dir"
    
    # Check if user is logged in to npm
    if ! npm whoami > /dev/null 2>&1; then
        print_error "Not logged in to npm. Please run 'npm login' first."
        exit 1
    fi
    
    # Publish to npm
    npm publish
    
    cd - > /dev/null
    print_success "Published to npm successfully"
}

# Function to create git tag
create_git_tag() {
    local sdk_dir="sdk/javascript"
    local version=$(node -p "require('./$sdk_dir/package.json').version")
    local test_mode=${1:-false}
    
    if [[ "$test_mode" == true ]]; then
        print_warning "TEST MODE: Would create git tag sdk-v$version"
        return 0
    fi
    
    print_status "Creating git tag v$version..."
    
    # Add the updated package.json
    git add "$sdk_dir/package.json"
    
    # Commit the version bump
    git commit -m "chore: bump SDK version to $version"
    
    # Create and push the tag
    git tag "sdk-v$version"
    git push origin main
    git push origin "sdk-v$version"
    
    print_success "Git tag sdk-v$version created and pushed"
}

# Main execution
main() {
    local version_type=${1:-patch}
    local dry_run=false
    local test_mode=false
    
    # Check for dry-run flag
    if [[ "$1" == "--dry-run" ]]; then
        dry_run=true
        version_type=${2:-patch}
    fi
    
    # Check for test mode flag
    if [[ "$1" == "--test" ]]; then
        test_mode=true
        version_type=${2:-patch}
    fi
    
    print_status "Starting SDK release process..."
    print_status "Version type: $version_type"
    if [[ "$dry_run" == true ]]; then
        print_warning "DRY RUN MODE - No actual changes will be made"
    fi
    if [[ "$test_mode" == true ]]; then
        print_warning "TEST MODE - Will skip npm publish and git operations"
    fi
    
    # Pre-flight checks
    check_git_repo
    if [[ "$dry_run" == false && "$test_mode" == false ]]; then
        check_clean_working_dir
        check_main_branch
    fi
    
    # Build and test
    build_sdk
    run_tests
    
    if [[ "$dry_run" == true ]]; then
        print_status "DRY RUN: Would bump version ($version_type)..."
        print_success "DRY RUN completed successfully!"
        return 0
    fi
    
    # Bump version
    bump_version "$version_type"
    
    # Build again with new version
    build_sdk
    
    # Run tests again
    run_tests
    
    # Publish to npm
    publish_to_npm "$test_mode"
    
    # Create git tag
    create_git_tag "$test_mode"
    
    print_success "SDK release completed successfully!"
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
