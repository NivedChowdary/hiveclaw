<p align="center">
  <img src="https://raw.githubusercontent.com/NivedChowdary/hiveclaw/main/docs/logo.png" alt="HiveClaw" width="400" />
</p>

<h1 align="center">🐝 HiveClaw</h1>

<p align="center">
  <strong>AGI-native gateway for AI agents with swarm intelligence.</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#channels">Channels</a> •
  <a href="#api">API</a> •
  <a href="#roadmap">Roadmap</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/version-0.1.0-blue" alt="Version" />
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8?logo=go" alt="Go" />
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License" />
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-lightgrey" alt="Platform" />
</p>

---

> *"One hive. Infinite intelligence."*

HiveClaw is a **self-hosted AI gateway** that connects Telegram, Discord, and more to LLMs like Claude. Built in Go for blazing performance, with architecture ready for swarm intelligence.

## ✨ Features

<table>
<tr>
<td width="50%">

### 🚀 Performance
- Written in **Go** — 10MB binary
- Sub-millisecond latency
- Zero runtime dependencies
- Single binary deployment

</td>
<td width="50%">

### 🔌 Multi-Channel
- **Telegram** bot integration
- **Discord** bot integration  
- **WebSocket** real-time API
- **REST** API endpoints

</td>
</tr>
<tr>
<td width="50%">

### 🧠 Multi-LLM
- **Anthropic** Claude
- **OpenRouter** (100+ models)
- Local models (coming soon)
- Automatic failover

</td>
<td width="50%">

### ⚡ Developer Experience
- 60-second onboarding wizard
- Hot config reload
- Session management
- React dashboard

</td>
</tr>
</table>

## 🚀 Quick Start

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/NivedChowdary/hiveclaw
cd hiveclaw

# Build
make build

# Run the setup wizard
./hiveclaw onboard

# Start the gateway
./hiveclaw start
```

### Option 2: Go Install

```bash
go install github.com/NivedChowdary/hiveclaw/cmd/hiveclaw@latest
hiveclaw onboard
hiveclaw start
```

### Open Dashboard

```
http://localhost:8080
```

## 🎯 Onboarding Wizard

The interactive wizard gets you running in 60 seconds:

```bash
./hiveclaw onboard
```

```
🐝 Welcome to HiveClaw!

Step 1: LLM Provider
━━━━━━━━━━━━━━━━━━━━
  1. Anthropic (Claude) - recommended
  2. OpenRouter (multi-model)

Step 2: Gateway Settings
━━━━━━━━━━━━━━━━━━━━━━━
Port [8080]: 

Step 3: Telegram Bot (optional)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Token: 

Step 4: Discord Bot (optional)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Token: 

✅ Setup Complete!
```

## ⚙️ Configuration

Config lives at `~/.hiveclaw/config.json`:

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

### Environment Variables

| Variable | Description |
|----------|-------------|
| `ANTHROPIC_API_KEY` | Anthropic API key |
| `OPENROUTER_API_KEY` | OpenRouter API key |
| `HIVECLAW_PORT` | Gateway port (default: 8080) |

## 📱 Channels

### Telegram

1. Create a bot via [@BotFather](https://t.me/BotFather)
2. Get your token
3. Add to config or use `hiveclaw onboard`

**Commands:**
- `/start` — Welcome message
- `/new` — Start new conversation
- `/clear` — Clear history
- `/status` — Check status

### Discord

1. Create app at [Discord Developers](https://discord.com/developers/applications)
2. Create a bot and get token
3. Invite to your server with Message Content intent
4. Mention the bot or DM it to chat

**Commands:**
- `!help` — Show help
- `!new` — New conversation
- `!clear` — Clear history
- `!status` — Check status

## 🔌 API

### REST Endpoints

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

### WebSocket

Connect to `ws://localhost:8080/ws`

```javascript
const ws = new WebSocket('ws://localhost:8080/ws')

// Send message
ws.send(JSON.stringify({
  type: 'req',
  id: '1',
  method: 'chat.send',
  params: { message: 'Hello!' }
}))

// Receive response
ws.onmessage = (event) => {
  const data = JSON.parse(event.data)
  console.log(data.payload.response)
}
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      HIVECLAW GATEWAY                       │
│  ┌────────────────┐  ┌────────────────┐  ┌──────────────┐   │
│  │   WebSocket    │  │    Channels    │  │   Dashboard  │   │
│  │    Server      │  │  (TG/Discord)  │  │   (React)    │   │
│  └───────┬────────┘  └───────┬────────┘  └──────┬───────┘   │
│          │                   │                  │           │
│  ┌───────▼───────────────────▼──────────────────▼─────────┐ │
│  │                   SESSION MANAGER                      │ │
│  └───────────────────────────┬────────────────────────────┘ │
│                              │                              │
│  ┌───────────────────────────▼────────────────────────────┐ │
│  │                    LLM PROVIDERS                       │ │
│  │            Claude • OpenRouter • Local                 │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
hiveclaw/
├── cmd/hiveclaw/          # CLI entry point
│   └── main.go            # Commands & onboarding wizard
├── internal/
│   ├── gateway/           # WebSocket server
│   ├── session/           # Session management
│   ├── llm/               # LLM providers
│   └── channels/
│       ├── telegram/      # Telegram bot
│       └── discord/       # Discord bot
├── web/frontend/          # React dashboard
├── configs/               # Configuration
├── docs/                  # Documentation
├── Makefile               # Build automation
└── README.md
```

## 🛠️ Development

```bash
# Install dependencies & build frontend
make frontend

# Build binary
make build

# Run in dev mode
make dev

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

## 🗺️ Roadmap

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
- [ ] Intent Engine
- [ ] Multi-agent routing
- [ ] WhatsApp integration
- [ ] Voice support

## 🤝 Contributing

Contributions welcome! Please read our contributing guidelines first.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing`)
5. Open a Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- Inspired by [OpenClaw](https://github.com/openclaw/openclaw)
- Built with [Go](https://go.dev), [React](https://react.dev), [Tailwind CSS](https://tailwindcss.com)
- LLM support via [Anthropic](https://anthropic.com) and [OpenRouter](https://openrouter.ai)

---

<p align="center">
  Built with 🐝 by <a href="https://github.com/NivedChowdary">NaniLabs</a>
</p>
