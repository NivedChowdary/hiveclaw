package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nanilabs/hiveclaw/internal/llm"
	"github.com/nanilabs/hiveclaw/internal/session"
)

// Bot represents a Discord bot
type Bot struct {
	Session  *discordgo.Session
	Sessions *session.Manager
	LLM      llm.Provider
	Config   Config
}

// Config for Discord bot
type Config struct {
	Token        string   `json:"token"`
	GuildID      string   `json:"guildId"`      // Optional: limit to specific guild
	AllowedRoles []string `json:"allowedRoles"` // Optional: role-based access
	Prefix       string   `json:"prefix"`       // Command prefix (default: !)
	SystemPrompt string   `json:"systemPrompt"`
}

// New creates a new Discord bot
func New(config Config, sessions *session.Manager, llmProvider llm.Provider) (*Bot, error) {
	if config.Prefix == "" {
		config.Prefix = "!"
	}

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	bot := &Bot{
		Session:  dg,
		Sessions: sessions,
		LLM:      llmProvider,
		Config:   config,
	}

	// Register handlers
	dg.AddHandler(bot.messageCreate)
	dg.AddHandler(bot.ready)

	// Set intents
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	return bot, nil
}

// Start starts the Discord bot
func (b *Bot) Start() error {
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord connection: %w", err)
	}

	log.Println("🎮 Discord bot is running...")
	return nil
}

// Stop stops the Discord bot
func (b *Bot) Stop() error {
	return b.Session.Close()
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("🎮 Discord bot logged in as %s#%s", event.User.Username, event.User.Discriminator)

	// Set status
	s.UpdateGameStatus(0, "🐝 HiveClaw | !help")
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot's own messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if it's a command
	if strings.HasPrefix(m.Content, b.Config.Prefix) {
		b.handleCommand(s, m)
		return
	}

	// Check if bot is mentioned or it's a DM
	mentioned := false
	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			mentioned = true
			break
		}
	}

	// In DMs, always respond
	isDM := m.GuildID == ""

	if mentioned || isDM {
		b.handleChat(s, m)
	}
}

func (b *Bot) handleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.TrimPrefix(m.Content, b.Config.Prefix)
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "help":
		embed := &discordgo.MessageEmbed{
			Title:       "🐝 HiveClaw Help",
			Description: "I'm an AI assistant powered by Hive Mind architecture.",
			Color:       0xFFD700, // Gold
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "💬 Chat",
					Value:  "Mention me or DM me to chat",
					Inline: true,
				},
				{
					Name:   "🆕 !new",
					Value:  "Start a new conversation",
					Inline: true,
				},
				{
					Name:   "🧹 !clear",
					Value:  "Clear conversation history",
					Inline: true,
				},
				{
					Name:   "📊 !status",
					Value:  "Check system status",
					Inline: true,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Built by NaniLabs 🐝",
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)

	case "new":
		sessionKey := b.getSessionKey(m)
		b.Sessions.Delete(sessionKey)
		b.Sessions.Create(sessionKey)
		s.ChannelMessageSend(m.ChannelID, "🆕 Started a new conversation!")

	case "clear":
		sessionKey := b.getSessionKey(m)
		b.Sessions.Clear(sessionKey)
		s.ChannelMessageSend(m.ChannelID, "🧹 Conversation cleared!")

	case "status":
		embed := &discordgo.MessageEmbed{
			Title: "🐝 HiveClaw Status",
			Color: 0x00FF00, // Green
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Version",
					Value:  "0.1.0",
					Inline: true,
				},
				{
					Name:   "Status",
					Value:  "✅ Operational",
					Inline: true,
				},
				{
					Name:   "Your ID",
					Value:  m.Author.ID,
					Inline: true,
				},
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)

	case "ping":
		s.ChannelMessageSend(m.ChannelID, "🏓 Pong!")

	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unknown command: `%s`. Try `%shelp`", cmd, b.Config.Prefix))
	}
}

func (b *Bot) handleChat(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Remove bot mention from message
	content := m.Content
	for _, mention := range m.Mentions {
		content = strings.ReplaceAll(content, "<@"+mention.ID+">", "")
		content = strings.ReplaceAll(content, "<@!"+mention.ID+">", "")
	}
	content = strings.TrimSpace(content)

	if content == "" {
		s.ChannelMessageSend(m.ChannelID, "Hey! How can I help you?")
		return
	}

	sessionKey := b.getSessionKey(m)

	// Get or create session
	_, exists := b.Sessions.Get(sessionKey)
	if !exists {
		b.Sessions.Create(sessionKey)
	}

	// Add user message
	b.Sessions.AddMessage(sessionKey, "user", content)

	// Show typing indicator
	s.ChannelTyping(m.ChannelID)

	// Build messages for LLM
	messages, _ := b.Sessions.GetMessages(sessionKey)
	llmMessages := make([]llm.Message, len(messages))
	for i, msg := range messages {
		llmMessages[i] = llm.Message{Role: msg.Role, Content: msg.Content}
	}

	// Call LLM
	resp, err := b.LLM.Chat(llmMessages, llm.Options{
		System: b.Config.SystemPrompt,
	})

	if err != nil {
		log.Printf("LLM error: %v", err)
		s.ChannelMessageSend(m.ChannelID, "❌ Sorry, I encountered an error. Please try again.")
		return
	}

	// Add assistant response
	b.Sessions.AddMessage(sessionKey, "assistant", resp.Content)

	// Send response (split if too long)
	b.sendMessage(s, m.ChannelID, resp.Content, m.Reference())
}

func (b *Bot) getSessionKey(m *discordgo.MessageCreate) string {
	// Use channel ID for guilds, user ID for DMs
	if m.GuildID == "" {
		return fmt.Sprintf("discord_dm_%s", m.Author.ID)
	}
	return fmt.Sprintf("discord_%s_%s", m.GuildID, m.ChannelID)
}

func (b *Bot) sendMessage(s *discordgo.Session, channelID, content string, ref *discordgo.MessageReference) {
	const maxLength = 2000 // Discord's message limit

	if len(content) <= maxLength {
		s.ChannelMessageSendReply(channelID, content, ref)
		return
	}

	// Split long messages
	for len(content) > 0 {
		end := maxLength
		if end > len(content) {
			end = len(content)
		}

		// Try to split at newline
		if end < len(content) {
			if idx := strings.LastIndex(content[:end], "\n"); idx > 0 {
				end = idx + 1
			}
		}

		s.ChannelMessageSend(channelID, content[:end])
		content = content[end:]
	}
}
