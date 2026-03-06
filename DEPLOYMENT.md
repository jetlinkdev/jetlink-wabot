# Deployment Guide

## Local Development

### Prerequisites
- Go 1.21+
- SQLite3
- Make (optional)
- Air (optional, for hot reload)

### Running Locally

1. **Install dependencies:**
```bash
go mod tidy
```

2. **Configure environment:**
```bash
cp .env.example .env
# Edit .env with your credentials
```

3. **Run the application:**

**Standard run:**
```bash
# Using go run
go run ./cmd/server

# Or using make
make run

# Or build first
make build
./bin/bot-wa
```

**With hot reload (recommended for development):**
```bash
# Install air first
make air-install

# Run with hot reload
make watch
# or simply
air
```

4. **Test the webhook:**
```bash
# Health check
curl http://localhost:8080/health

# Send test webhook (replace with actual payload)
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{"sender":"6281234567890","message":"Hello"}'
```

## Docker Deployment

### Prerequisites
- Docker
- Docker Compose

### Using Docker Compose (Recommended)

1. **Configure environment:**
```bash
cp .env.example .env
# Edit .env with your credentials
```

2. **Build and start:**
```bash
docker-compose up -d --build
```

3. **Check logs:**
```bash
docker-compose logs -f
```

4. **Stop:**
```bash
docker-compose down
```

### Using Docker directly

1. **Build image:**
```bash
docker build -t whatsapp-bot .
```

2. **Run container:**
```bash
docker run -d \
  --name whatsapp-bot \
  -p 8080:8080 \
  -v bot-data:/root/ \
  -e FOONTE_TOKEN=your_token \
  -e GROQ_API_KEY=your_api_key \
  whatsapp-bot
```

## Production Deployment

### System Requirements
- 1 CPU core minimum
- 512MB RAM minimum
- 1GB disk space
- Public IP or domain for webhook

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `FOONTE_TOKEN` | Yes | - | Your Foonte API token |
| `GROQ_API_KEY` | Yes | - | Your Groq API key |
| `WEBHOOK_PORT` | No | 8080 | Server port |
| `GROQ_MODEL` | No | qwen/qwen3-32b | LLM model |
| `DATABASE_PATH` | No | bot.db | SQLite database path |

### Nginx Reverse Proxy Setup

1. **Install Nginx:**
```bash
sudo apt-get install nginx
```

2. **Configure Nginx:**
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /webhook {
        proxy_pass http://localhost:8080/webhook;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_read_timeout 30s;
    }

    location /health {
        proxy_pass http://localhost:8080/health;
        access_log off;
    }
}
```

3. **Enable SSL (recommended):**
```bash
sudo certbot --nginx -d your-domain.com
```

### Systemd Service

1. **Create service file:**
```bash
sudo nano /etc/systemd/system/whatsapp-bot.service
```

2. **Add configuration:**
```ini
[Unit]
Description=WhatsApp Bot Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/bot
ExecStart=/path/to/bot/bot-wa
Restart=always
RestartSec=10
Environment="FOONTE_TOKEN=your_token"
Environment="GROQ_API_KEY=your_api_key"
Environment="WEBHOOK_PORT=8080"

[Install]
WantedBy=multi-user.target
```

3. **Enable and start:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable whatsapp-bot
sudo systemctl start whatsapp-bot
sudo systemctl status whatsapp-bot
```

## Foonte Webhook Configuration

1. Login to Foonte dashboard
2. Go to Devices
3. Select your device
4. Configure webhook URL: `https://your-domain.com/webhook`
5. Save changes

## Monitoring

### Health Check
```bash
curl https://your-domain.com/health
```

Expected response:
```json
{"status":"healthy"}
```

### Logs

**Docker:**
```bash
docker-compose logs -f
```

**Systemd:**
```bash
sudo journalctl -u whatsapp-bot -f
```

**Direct:**
```bash
# Application logs to stdout by default
```

## Backup Database

```bash
# Copy SQLite database
cp bot.db bot.db.backup

# Or using sqlite3
sqlite3 bot.db ".backup 'bot.db.backup'"
```

## Troubleshooting

### Bot not responding
1. Check if server is running: `curl http://localhost:8080/health`
2. Check Foonte webhook configuration
3. Check logs for errors
4. Verify API tokens are correct

### Database errors
1. Check file permissions
2. Ensure disk space is available
3. Try removing database file (will be recreated)

### High memory usage
1. Reduce context window size
2. Limit concurrent connections
3. Consider using a smaller LLM model

## Performance Tuning

### Recommended Settings

**For low traffic (< 100 messages/hour):**
- Context window: 10 messages
- Default model

**For medium traffic (100-1000 messages/hour):**
- Context window: 5-10 messages
- Consider caching frequently used prompts

**For high traffic (> 1000 messages/hour):**
- Context window: 3-5 messages
- Use connection pooling
- Consider horizontal scaling

## Security Considerations

1. **Use HTTPS** for webhook endpoint
2. **Validate webhook source** (implement signature verification if available)
3. **Rate limiting** to prevent abuse
4. **Regular backups** of database
5. **Keep dependencies updated**
6. **Use environment variables** for sensitive data
7. **Restrict database file permissions**

## Scaling

### Horizontal Scaling
- Use load balancer
- Shared database (consider PostgreSQL for production)
- Session affinity not required (stateless design)

### Vertical Scaling
- Increase CPU for faster LLM responses
- Increase RAM for larger context windows
- Use SSD for faster database operations
