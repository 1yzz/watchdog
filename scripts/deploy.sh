#!/bin/bash

# Watchdog gRPC Service Deployment Script
# Usage: ./scripts/deploy.sh [production|development]

set -e

DEPLOY_ENV="${1:-development}"
SERVICE_NAME="watchdog"
BINARY_NAME="watchdog-server"
SERVICE_USER="watchdog"

echo "🚀 Deploying Watchdog gRPC Service (${DEPLOY_ENV})"

# Check if running as root for production deployment
if [[ "$DEPLOY_ENV" == "production" && $EUID -ne 0 ]]; then
   echo "❌ Production deployment must be run as root (use sudo)"
   exit 1
fi

# Build the binary
echo "📦 Building binary..."
make build

# For production deployment
if [[ "$DEPLOY_ENV" == "production" ]]; then
    echo "🏭 Setting up production deployment..."
    
    # Create service user if doesn't exist
    if ! id "$SERVICE_USER" &>/dev/null; then
        echo "👤 Creating service user: $SERVICE_USER"
        useradd --system --no-create-home --shell /bin/false $SERVICE_USER
    fi
    
    # Create directories
    echo "📁 Creating directories..."
    mkdir -p /etc/watchdog
    mkdir -p /var/log/watchdog
    
    # Copy binary
    echo "📋 Copying binary to /usr/local/bin/"
    cp bin/$BINARY_NAME /usr/local/bin/
    chmod +x /usr/local/bin/$BINARY_NAME
    
    # Copy configuration if .env exists
    if [[ -f .env ]]; then
        echo "⚙️  Copying configuration..."
        cp .env /etc/watchdog/
    else
        echo "⚠️  No .env file found. Creating template..."
        cp .env.default /etc/watchdog/.env
        echo "📝 Edit /etc/watchdog/.env with your configuration"
    fi
    
    # Set permissions
    chown -R $SERVICE_USER:$SERVICE_USER /etc/watchdog
    chown -R $SERVICE_USER:$SERVICE_USER /var/log/watchdog
    
    # Create systemd service
    echo "🔧 Creating systemd service..."
    cat > /etc/systemd/system/$SERVICE_NAME.service <<EOF
[Unit]
Description=Watchdog gRPC Service
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=$SERVICE_USER
Group=$SERVICE_USER
WorkingDirectory=/etc/watchdog
ExecStart=/usr/local/bin/$BINARY_NAME
Restart=always
RestartSec=5
StandardOutput=append:/var/log/watchdog/watchdog.log
StandardError=append:/var/log/watchdog/watchdog.log

# Security settings
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/var/log/watchdog

[Install]
WantedBy=multi-user.target
EOF
    
    # Reload systemd and enable service
    echo "🔄 Enabling systemd service..."
    systemctl daemon-reload
    systemctl enable $SERVICE_NAME
    
    echo "✅ Production deployment complete!"
    echo ""
    echo "Next steps:"
    echo "1. Edit /etc/watchdog/.env with your database configuration"
    echo "2. Test database connection: cd /etc/watchdog && /usr/local/bin/$BINARY_NAME (Ctrl+C to stop)"
    echo "3. Start service: systemctl start $SERVICE_NAME"
    echo "4. Check status: systemctl status $SERVICE_NAME"
    echo "5. View logs: journalctl -u $SERVICE_NAME -f"

else
    # Development deployment
    echo "🛠️  Development deployment complete!"
    echo ""
    echo "Available commands:"
    echo "• Start server: make run"
    echo "• Test database: make db-test"
    echo "• View help: make help"
fi

echo ""
echo "🎉 Deployment finished successfully!"