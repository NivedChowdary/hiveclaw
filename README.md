# HiveClaw 🐝

**AGI-native gateway for AI agents with swarm intelligence.**

An OpenClaw alternative built from scratch in Go with Hive Mind architecture.

## Features

- 🚀 **Fast** — Written in Go, 10MB binary, sub-millisecond latency
- 🐝 **Hive Mind Ready** — Architecture for swarm intelligence
- 🔌 **Multi-channel** — Telegram, Discord, WebSocket, REST API
- 🤖 **Multi-agent** — Isolated workspaces per agent
- 🧠 **Multi-LLM** — Claude, OpenRouter, local models
- 📦 **Zero dependencies** — Just run the binary
- ⚡ **60-second setup** — Interactive onboarding wizard

## Quick Start

```bash
# Build from source
git clone https://github.com/nanilabs/hiveclaw
cd hiveclaw
make build

# Run setup wizard (60 seconds)
./hiveclaw onboard

# Start the gateway
./hiveclaw start

# Open dashboard
open http://localhost:8080
```

## Onboarding Wizard

```bash
./hiveclaw onboard
```

The wizard guides you through:
1. **LLM Provider** — Anthropic (Claude) or OpenRouter
2. **Gateway Port** — Default 8080
3. **Telegram Bot** — Optional, get token from @BotFather
4. **Discord Bot** — Optional, create at discord.com/developers
5. **System Prompt** — Customize your AI's personality

Config saves to `~/.hiveclaw/config.json`

## Configuration

```json
{
  "version": "1",
  "gateway": {
    "port": 8080
  },
  "llm": {
    "provider": "anthropic",
    "model": "claude-sonnet-4-20250514",
    "apiKey": "sk-ant-..."
  },
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "allowedIds": [123456789]
    },
    "discord": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN"
    }
  }
}
```

## Architecture

```
hiveclaw/
├── cmd/hiveclaw/          # CLI entry point
├── internal/
│   ├── gateway/           # WebSocket server
│   ├── session/           # Session management
│   ├── llm/               # LLM providers (Claude, OpenRouter)
│   ├── channels/
│   │   ├── telegram/      # Telegram bot
│   │   └── discord/       # Discord bot
│   └── hive/              # Hive Mind (coming soon)
├── web/frontend/          # React dashboard
├── configs/               # Configuration
└── docs/                  # Documentation
```

## Commands

```bash
hiveclaw onboard     # Interactive setup wizard
hiveclaw start       # Start gateway
hiveclaw stop        # Stop gateway
hiveclaw status      # Check status
hiveclaw version     # Print version
```

## REST API

```bash
# Health check
curl http://localhost:8080/api/health

# List sessions
curl http://localhost:8080/api/sessions

# Send message
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"sessionId": "main", "message": "Hello"}'
```

## WebSocket API

Connect to `ws://localhost:8080/ws`

```json
// Send message
{"type": "req", "id": "1", "method": "chat.send", "params": {"message": "Hello"}}

// Response
{"type": "res", "id": "1", "ok": true, "payload": {"response": "Hi!"}}
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ANTHROPIC_API_KEY` | Anthropic API key |
| `OPENROUTER_API_KEY` | OpenRouter API key |
| `HIVECLAW_PORT` | Gateway port (default: 8080) |

## Development

```bash
# Install dependencies
make frontend

# Build binary
make build

# Run in dev mode
make dev

# Build for all platforms
make build-all

# Clean
make clean
```

## Roadmap

- [x] WebSocket gateway
- [x] Session management  
- [x] Claude/OpenRouter LLM
- [x] Telegram bot
- [x] Discord bot
- [x] React dashboard
- [x] Onboarding wizard
- [ ] Tool execution
- [ ] Memory persistence
- [ ] Hive Mind swarm layer
- [ ] Intent Engine integration
- [ ] Multi-agent routing

## Why HiveClaw?

OpenClaw is great, but we wanted:
- **Go performance** — Faster, smaller, single binary
- **AGI-native** — Built for swarm intelligence from day one
- **Simpler onboarding** — 60 seconds to running
- **Open architecture** — Easy to extend and customize

## License

MIT

---

Built by [NaniLabs](https://nanilabs.io) 🐝
