package telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nanilabs/hiveclaw/internal/llm"
	"github.com/nanilabs/hiveclaw/internal/session"
)

// Bot represents a Telegram bot
type Bot struct {
	API      *tgbotapi.BotAPI
	Sessions *session.Manager
	LLM      llm.Provider
	Config   Config
}

// Config for Telegram bot
type Config struct {
	Token       string   `json:"token"`
	AllowedIDs  []int64  `json:"allowedIds"`  // Allowed user/chat IDs
	AdminIDs    []int64  `json:"adminIds"`    // Admin user IDs
	SystemPrompt string  `json:"systemPrompt"`
}

// New creates a new Telegram bot
func New(config Config, sessions *session.Manager, llmProvider llm.Provider) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	log.Printf("🤖 Telegram bot authorized as @%s", api.Self.UserName)

	return &Bot{
		API:      api,
		Sessions: sessions,
		LLM:      llmProvider,
		Config:   config,
	}, nil
}

// Start starts the bot
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	log.Println("🐝 Telegram bot listening for messages...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go b.handleMessage(update.Message)
	}

	return nil
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	// Check if user is allowed
	if !b.isAllowed(msg.From.ID, msg.Chat.ID) {
		log.Printf("Unauthorized access attempt from user %d in chat %d", msg.From.ID, msg.Chat.ID)
		return
	}

	// Handle commands
	if msg.IsCommand() {
		b.handleCommand(msg)
		return
	}

	// Handle regular messages
	b.handleChat(msg)
}

func (b *Bot) isAllowed(userID, chatID int64) bool {
	// If no allowlist, allow everyone (not recommended for production)
	if len(b.Config.AllowedIDs) == 0 {
		return true
	}

	for _, id := range b.Config.AllowedIDs {
		if id == userID || id == chatID {
			return true
		}
	}
	return false
}

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		b.sendMessage(msg.Chat.ID, `🐝 *Welcome to HiveClaw!*

I'm your AI assistant powered by the Hive Mind architecture.

*Commands:*
/new - Start a new conversation
/clear - Clear conversation history
/status - Check system status
/help - Show this help message

Just send me a message to chat!`, true)

	case "new":
		sessionKey := b.getSessionKey(msg)
		b.Sessions.Delete(sessionKey)
		b.Sessions.Create(sessionKey)
		b.sendMessage(msg.Chat.ID, "🆕 Started a new conversation!", false)

	case "clear":
		sessionKey := b.getSessionKey(msg)
		b.Sessions.Clear(sessionKey)
		b.sendMessage(msg.Chat.ID, "🧹 Conversation cleared!", false)

	case "status":
		status := fmt.Sprintf(`🐝 *HiveClaw Status*

*Version:* 0.1.0
*Bot:* @%s
*Your ID:* %d
*Chat ID:* %d

System is operational!`, b.API.Self.UserName, msg.From.ID, msg.Chat.ID)
		b.sendMessage(msg.Chat.ID, status, true)

	case "help":
		b.sendMessage(msg.Chat.ID, `🐝 *HiveClaw Help*

I'm an AI assistant that can help you with various tasks.

*Tips:*
• Just type your message to chat
• Use /new to start fresh
• Use /clear to reset context

*About:*
Built with Hive Mind architecture - swarm intelligence meets AI.`, true)

	default:
		b.sendMessage(msg.Chat.ID, "Unknown command. Try /help", false)
	}
}

func (b *Bot) handleChat(msg *tgbotapi.Message) {
	sessionKey := b.getSessionKey(msg)

	// Get or create session
	_, exists := b.Sessions.Get(sessionKey)
	if !exists {
		b.Sessions.Create(sessionKey)
	}

	// Add user message to session
	b.Sessions.AddMessage(sessionKey, "user", msg.Text)

	// Send typing indicator
	typing := tgbotapi.NewChatAction(msg.Chat.ID, tgbotapi.ChatTyping)
	b.API.Send(typing)

	// Build messages for LLM
	messages, _ := b.Sessions.GetMessages(sessionKey)
	llmMessages := make([]llm.Message, len(messages))
	for i, m := range messages {
		llmMessages[i] = llm.Message{Role: m.Role, Content: m.Content}
	}

	// Call LLM
	resp, err := b.LLM.Chat(llmMessages, llm.Options{
		System: b.Config.SystemPrompt,
	})

	if err != nil {
		log.Printf("LLM error: %v", err)
		b.sendMessage(msg.Chat.ID, "❌ Sorry, I encountered an error. Please try again.", false)
		return
	}

	// Add assistant response to session
	b.Sessions.AddMessage(sessionKey, "assistant", resp.Content)

	// Send response
	b.sendMessage(msg.Chat.ID, resp.Content, false)
}

func (b *Bot) getSessionKey(msg *tgbotapi.Message) string {
	// Use chat ID as session key (supports both DMs and groups)
	return fmt.Sprintf("tg_%d", msg.Chat.ID)
}

func (b *Bot) sendMessage(chatID int64, text string, markdown bool) {
	// Split long messages
	const maxLength = 4096
	texts := splitMessage(text, maxLength)

	for _, t := range texts {
		msg := tgbotapi.NewMessage(chatID, t)
		if markdown {
			msg.ParseMode = "Markdown"
		}

		if _, err := b.API.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}

func splitMessage(text string, maxLen int) []string {
	if len(text) <= maxLen {
		return []string{text}
	}

	var parts []string
	lines := strings.Split(text, "\n")
	var current strings.Builder

	for _, line := range lines {
		if current.Len()+len(line)+1 > maxLen {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			// Handle lines longer than maxLen
			for len(line) > maxLen {
				parts = append(parts, line[:maxLen])
				line = line[maxLen:]
			}
		}
		if current.Len() > 0 {
			current.WriteString("\n")
		}
		current.WriteString(line)
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
