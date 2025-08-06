#!/bin/bash

# Environment Setup Script for PDF Emailer
# This script helps set up environment files for different deployment scenarios

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -l, --local             Set up local development environment"
    echo "  -d, --docker            Set up Docker environment"
    echo "  -a, --all               Set up both local and Docker environments"
    echo "  -f, --force             Overwrite existing files"
    echo ""
    echo "Examples:"
    echo "  $0 -d                   # Set up Docker environment only"
    echo "  $0 -l -f                # Set up local environment, overwrite existing"
    echo "  $0 -a                   # Set up both environments"
}

# Default values
SETUP_LOCAL=false
SETUP_DOCKER=false
FORCE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -l|--local)
            SETUP_LOCAL=true
            shift
            ;;
        -d|--docker)
            SETUP_DOCKER=true
            shift
            ;;
        -a|--all)
            SETUP_LOCAL=true
            SETUP_DOCKER=true
            shift
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# If no options specified, show usage
if [ "$SETUP_LOCAL" = false ] && [ "$SETUP_DOCKER" = false ]; then
    print_error "No setup option specified"
    show_usage
    exit 1
fi

# Function to setup local environment
setup_local() {
    print_status "Setting up local development environment..."
    
    if [ -f ".env" ] && [ "$FORCE" = false ]; then
        print_warning ".env file already exists. Use -f to overwrite."
        return
    fi
    
    if [ -f "env.example" ]; then
        cp env.example .env
        print_success "Created .env from env.example"
        print_warning "Please edit .env with your local database credentials"
    else
        print_error "env.example not found"
        return 1
    fi
}

# Function to setup Docker environment
setup_docker() {
    print_status "Setting up Docker environment..."
    
    if [ -f ".env.docker" ] && [ "$FORCE" = false ]; then
        print_warning ".env.docker file already exists. Use -f to overwrite."
        return
    fi
    
    if [ -f "env.docker.example" ]; then
        cp env.docker.example .env.docker
        print_success "Created .env.docker from env.docker.example"
        print_warning "Please edit .env.docker with your Docker database credentials"
    else
        print_error "env.docker.example not found"
        return 1
    fi
}

# Main execution
if [ "$SETUP_LOCAL" = true ]; then
    setup_local
fi

if [ "$SETUP_DOCKER" = true ]; then
    setup_docker
fi

print_success "Environment setup completed!"
echo ""
print_status "Next steps:"
if [ "$SETUP_LOCAL" = true ]; then
    echo "  1. Edit .env with your local database credentials"
    echo "  2. Run: make run"
fi
if [ "$SETUP_DOCKER" = true ]; then
    echo "  1. Edit .env.docker with your Docker database credentials"
    echo "  2. Run: docker-compose up -d"
fi
echo ""
print_warning "Remember: .env and .env.docker files are excluded from git for security" 