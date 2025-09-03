#!/bin/bash

# Smart development startup script
# Automatically handles port conflicts and process cleanup

PORT=50051
BINARY_PATH="./bin/watchdog-server"

echo "🚀 Starting watchdog in development mode..."

# Function to cleanup
cleanup() {
    echo "🛑 Cleaning up..."
    pkill -f "watchdog-server" 2>/dev/null || true
    pkill -f "tmp/main" 2>/dev/null || true
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Function to check if port is in use
check_port() {
    if command -v ss >/dev/null 2>&1; then
        ss -tlnp | grep ":$PORT " >/dev/null 2>&1
    elif command -v netstat >/dev/null 2>&1; then
        netstat -tlnp | grep ":$PORT " >/dev/null 2>&1
    else
        # Fallback: try to bind to the port
        timeout 1 bash -c "echo >/dev/tcp/localhost/$PORT" 2>/dev/null
        return $?
    fi
}

# Kill existing processes and free port
echo "🔍 Checking for existing processes..."
pkill -f "watchdog-server" 2>/dev/null || echo "No existing watchdog-server process"
pkill -f "tmp/main" 2>/dev/null || echo "No existing tmp/main process"

# Wait a moment for processes to fully terminate
sleep 1

# Check if port is still in use
if check_port; then
    echo "⚠️  Port $PORT is still in use. Attempting to kill process..."
    lsof -ti:$PORT | xargs kill -9 2>/dev/null || echo "Could not kill process on port $PORT"
    sleep 2
fi

# Final port check
if check_port; then
    echo "❌ Port $PORT is still in use. Please manually stop the process using:"
    echo "   lsof -ti:$PORT | xargs kill -9"
    exit 1
else
    echo "✅ Port $PORT is free"
fi

# Start air for development
echo "🎯 Starting air development server..."
air
