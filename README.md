# ascii-ngin

A terminal-native ASCII animation engine that transforms videos into mesmerizing ASCII art. Convert your own videos or generate new ones with AI.

Built with Go and [Charm Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **Local Video Conversion** - Convert any video (MP4, WebM, AVI, MOV, etc.) to ASCII animation
- **AI Video Generation** - Generate videos from text prompts using Google Veo 3
- **Image Preview** - Preview AI-generated images before video creation (on supported terminals)
- **Animation Library** - Save and manage your ASCII animations
- **Full-Screen Playback** - Immersive looping animations with configurable sizing

## Demo

<!-- TODO: Add demo GIF/video here -->

## Installation

### Prerequisites

- Go 1.21 or higher
- [FFMPEG](https://ffmpeg.org/download.html) installed and in your PATH
- (Optional) [Gemini API key](https://ai.google.dev/) for AI features

### From Source

```bash
git clone https://github.com/yourusername/ascii-ngin.git
cd ascii-ngin
go build -o ascii-ngin ./cmd/ascii-ngin
```

### Verify Installation

```bash
ascii-ngin --version
```

## Quick Start

### Convert a Local Video

```bash
ascii-ngin
# Select "Local Video" → Enter path to your video → Watch the ASCII magic
```

### Generate with AI

```bash
# First, set up your API key
ascii-ngin config set api.gemini_api_key "your-key-here"

# Then generate
ascii-ngin
# Select "AI Generation" → Enter your prompt → Approve preview → Enjoy
```

## Configuration

Configuration is stored at `~/.config/ascii-ngin/config.json`.

### API Key Setup

Choose one of these methods:

**Config file:**
```json
{
  "api": {
    "gemini_api_key": "your-api-key-here"
  }
}
```

**CLI:**
```bash
ascii-ngin config set api.gemini_api_key "your-key-here"
```

**Environment variable:**
```bash
export GEMINI_API_KEY="your-key-here"
```

## How It Works

```
Video → FFMPEG (frame extraction) → ASCII Conversion → JSON Storage → TUI Playback
```

1. **Input**: Provide a video file or text prompt
2. **Frame Extraction**: FFMPEG extracts individual frames
3. **ASCII Conversion**: Each frame is converted to ASCII characters
4. **Storage**: Animations are saved as JSON for instant replay
5. **Playback**: Smooth looping playback in your terminal

## Supported Formats

Any format FFMPEG supports, including:
- MP4, WebM, AVI, MOV, MKV, GIF, and more

## Requirements

| Requirement | Details |
|-------------|---------|
| Terminal | ANSI escape code support (most modern terminals) |
| FFMPEG | Required for video processing |
| Internet | Only for AI generation features |

For the best experience with AI image previews, use a terminal with Kitty or Sixel image protocol support (iTerm2, Kitty, WezTerm, etc.).

## Documentation

See [DOCUMENTATION.md](DOCUMENTATION.md) for detailed architecture, user flows, and technical specifications.

## Roadmap

- [x] Project specification
- [ ] Core TUI application
- [ ] FFMPEG integration
- [ ] ASCII conversion engine
- [ ] Local video workflow
- [ ] Gemini API integration (Nano Banana + Veo 3)
- [ ] Animation library management

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT](LICENSE)

---

Built for the Gemini API Developer Competition
