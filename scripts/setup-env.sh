#!/bin/bash

echo "Setting up environment files..."

# Create .env.docker for Docker development
if [ ! -f .env.docker ]; then
    echo "Creating .env.docker for Docker development..."
    cp env.example .env.docker
    # Update for Docker environment
    sed -i '' 's/SMTP_HOST=localhost/SMTP_HOST=mailpit/g' .env.docker
    sed -i '' 's/POSTGRES_HOST=localhost/POSTGRES_HOST=db/g' .env.docker
    echo "✅ .env.docker created with Mailpit configuration"
else
    echo "⚠️  .env.docker already exists"
fi

# Create .env for local development
if [ ! -f .env ]; then
    echo "Creating .env for local development..."
    cp env.example .env
    echo "✅ .env created with local configuration"
    echo "📝 Edit .env to configure your local settings"
else
    echo "⚠️  .env already exists"
fi

echo ""
echo "🎯 Environment setup complete!"
echo ""
echo "📧 Mailpit Configuration:"
echo "   - Web UI: http://localhost:8025"
echo "   - SMTP Server: localhost:1025"
echo ""
echo "🚀 Next steps:"
echo "   1. For Docker: docker compose up"
echo "   2. For local: make dev"
echo "   3. View emails at: http://localhost:8025" 