#!/bin/bash

# Deploy script untuk Jetlink WhatsApp Bot
# Usage: ./deploy.sh [vps-ip] [vps-port] [vps-user]

set -e

# Configuration
VPS_HOST="${1:-151.243.222.93}"
VPS_PORT="${2:-34589}"
VPS_USER="${3:-root}"
REMOTE_PATH="/root/jetlink-wabot"

echo "🚀 Deploying Jetlink WhatsApp Bot..."
echo "   Host: $VPS_HOST"
echo "   Port: $VPS_PORT"
echo "   User: $VPS_USER"
echo "   Remote Path: $REMOTE_PATH"
echo ""

# Build binary for Linux (no CGO needed with modernc.org/sqlite)
echo "📦 Building binary for Linux..."
CGO_ENABLED=0 GOOS=linux go build -o bin/bot-wa-linux ./cmd/server

# Create remote directory
echo "📁 Creating remote directory..."
ssh -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" "mkdir -p $REMOTE_PATH"

# Sync files (exclude unnecessary files)
echo "📤 Syncing files to VPS..."
rsync -avz --progress \
    -e "ssh -p $VPS_PORT" \
    --exclude='.git' \
    --exclude='tmp/' \
    --exclude='*.db' \
    --exclude='.env' \
    --exclude='bot-wa' \
    --exclude='*.test' \
    ./ \
    "$VPS_USER@$VPS_HOST:$REMOTE_PATH"

# Upload binary
echo "⬆️  Uploading binary..."
scp -P "$VPS_PORT" bin/bot-wa-linux "$VPS_USER@$VPS_HOST:$REMOTE_PATH/bot-wa"

# Setup on VPS
echo "🔧 Setting up on VPS..."
ssh -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" << 'ENDSSH'
    cd /root/jetlink-wabot
    
    # Make binary executable
    chmod +x bot-wa
    
    # Create .env if not exists
    if [ ! -f .env ]; then
        echo "⚠️  .env file not found! Please create it manually:"
        echo "    cp .env.example .env"
        echo "    nano .env"
        exit 1
    fi
    
    # Stop existing process if running
    pkill -f bot-wa || true
    
    echo "✅ Deploy completed!"
ENDSSH

echo ""
echo "🎉 Deployment successful!"
echo ""
echo "Next steps on VPS:"
echo "  1. SSH to VPS: ssh -p $VPS_PORT $VPS_USER@$VPS_HOST"
echo "  2. cd $REMOTE_PATH"
echo "  3. cp .env.example .env (if not exists)"
echo "  4. nano .env (edit with your credentials)"
echo "  5. ./bot-wa"
echo ""
echo "Or use systemd service:"
echo "  sudo systemctl enable jetlink-bot"
echo "  sudo systemctl start jetlink-bot"
