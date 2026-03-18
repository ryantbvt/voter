# Voter
This allows Discord users in a channel to vote and schedule a game/movie night on Discord!

# Bot Framework

This document explains how the bot is structured and how to add new functionality as a contributor.

---

## Project Structure

```
yourapp/
├── main.go                      # Wiring only — connects everything together
└── internal/
    ├── framework/               # Core router — do not edit unless changing routing logic
    │   ├── router.go            # Central dispatcher
    │   ├── commands.go          # Slash command registry + CommandHandler interface
    │   ├── buttons.go           # Button registry + ButtonHandler interface
    │   ├── modals.go            # Modal registry + ModalHandler interface
    │   └── config.go            # Loads environment variables
    ├── commands/                # One file per slash command
    │   ├── ping.go
    │   └── feedback.go
    └── handlers/                # Button and modal handlers
        └── handlers.go
```

As a contributor, you will almost never touch `internal/framework/`. Your work lives in `internal/commands/`, `internal/handlers/`, and one line in `main.go`.

---

## How It Works

Discord sends all user interactions (slash commands, button clicks, modal submits) to a single webhook. The framework receives that event and routes it to the right handler.

```
Discord
  └── router.Handler()           ← single entry point registered with discordgo
        ├── /slash command       → CommandRegistry  → your CommandHandler
        ├── button click         → ButtonRegistry   → your ButtonHandler
        └── modal submit         → ModalRegistry    → your ModalHandler
```

The router dispatches based on the interaction type. Each registry then looks up the handler by ID — command name for slash commands, `custom_id` for buttons and modals — and calls `Handle(s, i)` on it.

---

## Environment Setup

The bot requires a `.env` file in the project root:

```
DISCORD_TOKEN=your_bot_token_here
```

In production, set this as an environment variable directly — the bot will fall back to it automatically if no `.env` file is present.

---

## Adding a Slash Command

**1. Create a new file in `internal/commands/`.**

Every slash command implements the `CommandHandler` interface, which requires two methods:
- `Definition()` — returns the command schema that gets registered with Discord (name, description, options)
- `Handle()` — the logic that runs when a user invokes the command

```go
// internal/commands/greet.go
package commands

import "github.com/bwmarrin/discordgo"

type GreetCommand struct{}

func (c *GreetCommand) Definition() *discordgo.ApplicationCommand {
    return &discordgo.ApplicationCommand{
        Name:        "greet",
        Description: "Sends a greeting",
    }
}

func (c *GreetCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: "Hello! 👋",
        },
    })
}
```

**2. Register it in `main.go`:**

```go
router.Commands().Add(&commands.GreetCommand{})
```

That's it. On next startup, `Sync()` automatically uploads the command definition to Discord.

> **Note:** Synced commands are global by default and can take up to an hour to appear in all servers.
> To test instantly during development, set `guildID` in `internal/framework/commands.go` to your test server's ID — guild-scoped commands update immediately.

---

## Adding a Button Handler

Buttons don't need to be synced with Discord — you define them inline when building a message response, and Discord sends you the interaction when a user clicks.

Buttons are identified by their `custom_id`. The registry supports two matching strategies:

- **Exact match** — for static buttons with a fixed ID
- **Prefix match** — for dynamic buttons where the ID carries embedded state (e.g. `confirm:delete:123`)

**1. Add a handler in `internal/handlers/`:**

```go
type ConfirmDeleteHandler struct{}

func (h *ConfirmDeleteHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
    // Extract the ID suffix if you need it
    // customID := i.MessageComponentData().CustomID  →  "confirm:delete:123"

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: "Deleted!",
            Flags:   discordgo.MessageFlagsEphemeral,
        },
    })
}
```

**2. Register it in `main.go`:**

```go
// Exact match — only fires for custom_id "confirm:delete"
router.Buttons().Add("confirm:delete", &handlers.ConfirmDeleteHandler{})

// Prefix match — fires for "confirm:delete:123", "confirm:delete:456", etc.
router.Buttons().AddPrefix("confirm:delete:", &handlers.ConfirmDeleteHandler{})
```

**3. Send the button in a command response:**

```go
discordgo.Button{
    Label:    "Delete",
    Style:    discordgo.DangerButton,
    CustomID: "confirm:delete:123",
}
```

---

## Adding a Modal Handler

Modals are forms that pop up in the Discord client. A slash command (or button) opens the modal, and when the user submits it, the framework routes the submit interaction to your handler.

Modals follow the same exact/prefix matching as buttons.

**1. Open the modal from a command or button:**

```go
// Inside a command's Handle()
s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
    Type: discordgo.InteractionResponseModal,
    Data: &discordgo.InteractionResponseData{
        CustomID: "modal:report",
        Title:    "Report an Issue",
        Components: []discordgo.MessageComponent{
            discordgo.ActionsRow{
                Components: []discordgo.MessageComponent{
                    discordgo.TextInput{
                        CustomID: "report_text",
                        Label:    "Describe the issue",
                        Style:    discordgo.TextInputParagraph,
                        Required: true,
                    },
                },
            },
        },
    },
})
```

**2. Add a handler in `internal/handlers/`:**

```go
type ReportModalHandler struct{}

func (h *ReportModalHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
    data := i.ModalSubmitData()

    // Pull values out of the submitted components
    var reportText string
    for _, row := range data.Components {
        row := row.(*discordgo.ActionsRow)
        for _, comp := range row.Components {
            input := comp.(*discordgo.TextInput)
            if input.CustomID == "report_text" {
                reportText = input.Value
            }
        }
    }

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: "Thanks for your report!",
            Flags:   discordgo.MessageFlagsEphemeral,
        },
    })
}
```

**3. Register it in `main.go`:**

```go
router.Modals().AddPrefix("modal:report", &handlers.ReportModalHandler{})
```

---

## The One Rule

`main.go` is the **only** file that knows about all the pieces. The framework, commands, and handlers are all independent of each other. When you add a new feature:

1. Create a new file in `internal/commands/` or `internal/handlers/`
2. Add one line to `main.go`

Nothing else needs to change.