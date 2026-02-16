<p align="center">
  <img src="https://raw.githubusercontent.com/NivedChowdary/hiveclaw/main/docs/logo.svg" alt="HiveClaw" width="400" />
</p>

<h1 align="center">🐝 HiveClaw</h1>

<p align="center">
  <strong>AGI-native gateway for AI agents with swarm intelligence.</strong>
</p>

<p align="center">
  <a href="#-quick-install">Install</a> •
  <a href="#-features">Features</a> •
  <a href="#%EF%B8%8F-configuration">Config</a> •
  <a href="#-channels">Channels</a> •
  <a href="#-api">API</a> •
  <a href="#-roadmap">Roadmap</a>
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

## 🚀 Quick Install

### One-liner (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/NivedChowdary/hiveclaw/main/install.sh | bash
```

### Windows (PowerShell)

```powershell
# Download binary
Invoke-WebRequest -Uri "https://github.com/NivedChowdary/hiveclaw/raw/main/build/hiveclaw-windows-amd64.exe" -OutFile "hiveclaw.exe"

# Run setup
.\hiveclaw.exe onboard
.\hiveclaw.exe start
```

### With Go

```bash
go install github.com/NivedChowdary/hiveclaw/cmd/hiveclaw@latest
```

### From Source

```bash
git clone https://github.com/NivedChowdary/hiveclaw
cd hiveclaw
go build -o hiveclaw ./cmd/hiveclaw
./hiveclaw onboard
./hiveclaw start
```

## 🎯 60-Second Setup

```bash
hiveclaw onboard
```

The interactive wizard guides you through:

```
🐝 Welcome to HiveClaw!

Step 1: LLM Provider
━━━━━━━━━━━━━━━━━━━━
  1. Anthropic (Claude) - recommended
  2. OpenRouter (multi-model)

Step 2: Gateway Port [8080]

Step 3: Telegram Bot Token (optional)

Step 4: Discord Bot Token (optional)

✅ Setup Complete!
```

Then just:
```bash
hiveclaw start
```

Open **http://localhost:8080** 🎉

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🚀 **Fast** | Written in Go — 7MB binary, sub-millisecond latency |
| 🔌 **Multi-Channel** | Telegram, Discord, WebSocket, REST API |
| 🧠 **Multi-LLM** | Claude, OpenRouter (100+ models) |
| 📱 **Multi-Platform** | Linux, macOS, Windows |
| ⚡ **Simple Setup** | 60-second onboarding wizard |
| 🐝 **Hive Mind Ready** | Architecture for swarm intelligence |

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
3. Run `hiveclaw onboard` or add to config

**Commands:**
- `/start` — Welcome message
- `/new` — Start new conversation
- `/clear` — Clear history
- `/status` — Check status

### Discord

1. Create app at [Discord Developers](https://discord.com/developers/applications)
2. Enable **Message Content Intent**
3. Create bot and get token
4. Invite to server and mention to chat

**Commands:**
- `!help` — Show help
- `!new` — New conversation
- `!clear` — Clear history

## 🔌 API

### REST Endpoints

```bash
# Health check
curl http://localhost:8080/api/health

# Chat
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello!"}'

# Sessions
curl http://localhost:8080/api/sessions
```

### WebSocket

Connect to `ws://localhost:8080/ws`

```javascript
const ws = new WebSocket('ws://localhost:8080/ws')

ws.send(JSON.stringify({
  type: 'req',
  id: '1',
  method: 'chat.send',
  params: { message: 'Hello!' }
}))

ws.onmessage = (e) => console.log(JSON.parse(e.data))
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      HIVECLAW GATEWAY                       │
│  ┌────────────────┐  ┌────────────────┐  ┌──────────────┐   │
│  │   WebSocket    │  │    Channels    │  │   Dashboard  │   │
│  │    Server      │  │  (TG/Discord)  │  │   (React)    │   │
│  └───────┬────────┘  └───────┬────────┘  └──────┬───────┘   │
│          └───────────────────┼───────────────────┘          │
│                              ▼                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                   SESSION MANAGER                    │   │
│  └──────────────────────────┬───────────────────────────┘   │
│                              ▼                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                    LLM PROVIDERS                     │   │
│  │              Claude • OpenRouter • Local             │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
hiveclaw/
├── cmd/hiveclaw/          # CLI & onboarding wizard
├── internal/
│   ├── gateway/           # WebSocket server
│   ├── session/           # Session management
│   ├── llm/               # LLM providers
│   └── channels/          # Telegram, Discord
├── web/frontend/          # React dashboard
├── build/                 # Pre-built binaries
└── configs/               # Configuration
```

## 🛠️ CLI Commands

```bash
hiveclaw onboard     # Interactive setup wizard
hiveclaw start       # Start the gateway
hiveclaw status      # Check gateway status
hiveclaw version     # Print version
```

## 🗺️ Roadmap

- [x] WebSocket gateway
- [x] Session management
- [x] Claude/OpenRouter LLM
- [x] Telegram bot
- [x] Discord bot
- [x] React dashboard
- [x] Onboarding wizard
- [x] Cross-platform binaries
- [ ] Tool execution
- [ ] Memory persistence
- [ ] Hive Mind swarm layer
- [ ] WhatsApp integration
- [ ] Voice support

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push (`git push origin feature/amazing`)
5. Open a Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE)

## 🙏 Credits

- Inspired by [OpenClaw](https://github.com/openclaw/openclaw)
- Built with [Go](https://go.dev), [React](https://react.dev), [Tailwind CSS](https://tailwindcss.com)

---

<p align="center">
  Built with 🐝 by <a href="https://github.com/NivedChowdary">NaniLabs</a>
</p>
