# ascii-ngin Documentation

> A terminal-native ASCII animation engine built with Go and Charm Bubble Tea

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Core Components](#core-components)
- [User Flows](#user-flows)
- [Technical Specifications](#technical-specifications)
- [Configuration](#configuration)
- [Storage Format](#storage-format)
- [Dependencies](#dependencies)
- [Development Roadmap](#development-roadmap)

---

## Overview

### What is ascii-ngin?

ascii-ngin is a full-featured TUI (Terminal User Interface) application that converts videos into ASCII art animations. Users can either provide their own video files or leverage Google's AI models to generate videos from text prompts.

### Goals

**Primary Goals:**
- Create a performant ASCII animation engine
- Support two input modes: local video files and AI-generated videos
- Deliver a polished, terminal-native experience using Charm Bubble Tea
- Provide flexible rendering options (full-screen and custom sizes)

**Secondary Goals:**
- Submit for the Gemini API Developer Competition (Due: February 10, 2026)

### Key Features

- **Local Video Conversion**: Convert any video format supported by FFMPEG to ASCII animation
- **AI Video Generation**: Generate videos from text prompts using Google Veo 3
- **AI Image Preview**: Preview generated images using Google Nano Banana before video generation (on supported terminals)
- **File-based Storage**: Save and manage multiple ASCII animations
- **Flexible Display**: Full-screen default with configurable sizing options
- **Loop Playback**: Continuous animation playback

---

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         ascii-ngin                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │   TUI Layer  │    │ Processing   │    │   Storage    │      │
│  │ (Bubble Tea) │◄──►│    Engine    │◄──►│    Layer     │      │
│  └──────────────┘    └──────────────┘    └──────────────┘      │
│         │                   │                                   │
│         ▼                   ▼                                   │
│  ┌──────────────┐    ┌──────────────┐                          │
│  │   Renderer   │    │  External    │                          │
│  │              │    │  Services    │                          │
│  └──────────────┘    └──────────────┘                          │
│                            │                                    │
│              ┌─────────────┼─────────────┐                     │
│              ▼             ▼             ▼                     │
│        ┌─────────┐   ┌─────────┐   ┌─────────┐                │
│        │ FFMPEG  │   │Nano     │   │  Veo 3  │                │
│        │         │   │Banana   │   │         │                │
│        └─────────┘   └─────────┘   └─────────┘                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Component Breakdown

```
ascii-ngin/
├── cmd/
│   └── ascii-ngin/
│       └── main.go              # Application entry point
├── internal/
│   ├── tui/
│   │   ├── app.go               # Main Bubble Tea application
│   │   ├── views/               # Different TUI views/screens
│   │   ├── components/          # Reusable UI components
│   │   └── styles/              # Lip Gloss styling
│   ├── engine/
│   │   ├── converter.go         # ASCII conversion logic
│   │   ├── renderer.go          # Animation playback
│   │   └── processor.go         # Frame processing pipeline
│   ├── video/
│   │   ├── ffmpeg.go            # FFMPEG integration
│   │   └── extractor.go         # Frame extraction
│   ├── ai/
│   │   ├── client.go            # Gemini API client
│   │   ├── nanoBanana.go        # Image generation
│   │   └── veo.go               # Video generation
│   ├── storage/
│   │   ├── manager.go           # Animation file management
│   │   └── animation.go         # Animation data structures
│   └── config/
│       └── config.go            # Configuration management
├── pkg/
│   └── ascii/
│       └── converter.go         # [TO BE IMPLEMENTED] ASCII conversion algorithm
├── configs/
│   └── default.json             # Default configuration
└── go.mod
```

---

## Core Components

### 1. TUI Layer (Bubble Tea)

The user interface is built using [Charm Bubble Tea](https://github.com/charmbracelet/bubbletea), providing a rich terminal experience.

**Views:**
- **Main Menu**: Choose between local video or AI generation
- **File Browser**: Navigate and select local video files
- **Prompt Input**: Enter text prompts for AI generation
- **Image Preview**: Display generated image for approval (if terminal supports image protocol)
- **Processing View**: Show conversion progress
- **Animation Player**: Full-screen ASCII animation playback
- **Library View**: Browse and select saved animations

### 2. Processing Engine

Handles the core logic of converting video frames to ASCII art.

**Pipeline:**
```
Video Input → Frame Extraction → ASCII Conversion → Animation Assembly → Storage
```

### 3. ASCII Converter

> **Note**: This section will be implemented separately. The converter transforms image frames into ASCII character representations.

<!-- [TO BE FILLED: ASCII conversion algorithm details] -->

### 4. External Services Integration

**FFMPEG Integration:**
- Frame extraction from video files
- Support for all FFMPEG-compatible formats
- Configurable extraction rate

**Gemini API Integration:**
- Nano Banana for image generation (preview)
- Veo 3 for video generation
- API key management via configuration

### 5. Storage Layer

File-based storage system for managing ASCII animations.

---

## User Flows

### Flow 1: Local Video Conversion

```
┌─────────┐    ┌─────────────┐    ┌─────────────┐    ┌──────────┐
│  Start  │───►│ Select Mode │───►│ Enter Video │───►│ Validate │
└─────────┘    │   (Local)   │    │    Path     │    │   Path   │
               └─────────────┘    └─────────────┘    └────┬─────┘
                                                          │
       ┌──────────────────────────────────────────────────┘
       ▼
┌──────────────┐    ┌─────────────┐    ┌─────────────┐    ┌──────────┐
│   Extract    │───►│  Convert    │───►│    Save     │───►│  Preview │
│   Frames     │    │  to ASCII   │    │  Animation  │    │  & Loop  │
│  (FFMPEG)    │    │             │    │   (JSON)    │    │          │
└──────────────┘    └─────────────┘    └─────────────┘    └──────────┘
```

**Step-by-step:**

1. User launches ascii-ngin
2. User selects "Local Video" mode from main menu
3. User provides path to video file (via file browser or direct input)
4. System validates the file exists and is a supported format
5. FFMPEG extracts frames from the video
6. Each frame is converted to ASCII representation
7. Animation is saved to storage as JSON
8. Animation plays in loop (full-screen by default)

### Flow 2: AI-Generated Video

```
┌─────────┐    ┌─────────────┐    ┌─────────────┐    ┌──────────────┐
│  Start  │───►│ Select Mode │───►│ Enter Text  │───►│   Generate   │
└─────────┘    │    (AI)     │    │   Prompt    │    │    Image     │
               └─────────────┘    └─────────────┘    │(Nano Banana) │
                                                     └──────┬───────┘
                                                            │
       ┌────────────────────────────────────────────────────┘
       ▼
┌──────────────────┐    ┌─────────────┐    ┌─────────────┐
│  Preview Image   │───►│   Approve   │───►│  Generate   │
│ (if supported)   │    │  or Retry?  │    │   Video     │
└──────────────────┘    └──────┬──────┘    │   (Veo 3)   │
                               │           └──────┬──────┘
                               │                  │
       ┌───────────────────────┘                  │
       │ (Retry with new prompt)                  │
       ▼                                          ▼
┌─────────────┐                          ┌──────────────┐
│ Enter New   │                          │   Extract    │
│   Prompt    │                          │   Frames     │
└─────────────┘                          └──────┬───────┘
                                                │
       ┌────────────────────────────────────────┘
       ▼
┌──────────────┐    ┌─────────────┐    ┌──────────┐
│  Convert     │───►│    Save     │───►│  Preview │
│  to ASCII    │    │  Animation  │    │  & Loop  │
└──────────────┘    └─────────────┘    └──────────┘
```

**Step-by-step:**

1. User launches ascii-ngin
2. User selects "AI Generation" mode from main menu
3. User enters a text prompt describing desired animation
4. System calls Nano Banana API to generate a preview image
5. **If terminal supports image protocol:**
   - Image is displayed for user review
   - User can approve or request regeneration with a new prompt
6. **If terminal does not support image protocol:**
   - This preview step is skipped
7. Upon approval (or skip), system calls Veo 3 API to generate video
8. FFMPEG extracts frames from the generated video
9. Each frame is converted to ASCII representation
10. Animation is saved to storage as JSON
11. Animation plays in loop (full-screen by default)

---

## Technical Specifications

### Performance Targets

| Metric | Target | Acceptable |
|--------|--------|------------|
| Frame Rate | 30 FPS | 15 FPS |
| Startup Time | < 1s | < 2s |
| Frame Conversion | < 50ms/frame | < 100ms/frame |

### Supported Video Formats

All formats supported by FFMPEG, including but not limited to:
- MP4 (H.264, H.265)
- WebM (VP8, VP9)
- AVI
- MOV
- MKV
- GIF

### Terminal Requirements

**Minimum:**
- Terminal with ANSI escape code support
- Minimum 80x24 character dimensions

**Recommended:**
- Modern terminal emulator (iTerm2, Kitty, Alacritty, WezTerm, etc.)
- Support for Kitty or Sixel image protocol (for AI image preview)
- 256+ color support

### System Requirements

- Go 1.21 or higher
- FFMPEG installed and available in PATH
- Internet connection (for AI features)
- Google Cloud API key with Gemini API access (for AI features)

---

## Configuration

### Configuration File Location

```
~/.config/ascii-ngin/config.json
```

### Configuration Schema

```json
{
  "api": {
    "gemini_api_key": "your-api-key-here"
  },
  "playback": {
    "default_fps": 30,
    "fullscreen": true,
    "loop": true
  },
  "conversion": {
    "target_width": 120,
    "target_height": 40,
    "charset": "standard"
  },
  "storage": {
    "animations_dir": "~/.config/ascii-ngin/animations/"
  }
}
```

### Setting the API Key

**Option 1: Configuration File**

Edit `~/.config/ascii-ngin/config.json` and add your API key.

**Option 2: Command Line**

```bash
ascii-ngin config set api.gemini_api_key "your-api-key-here"
```

**Option 3: Environment Variable**

```bash
export GEMINI_API_KEY="your-api-key-here"
```

---

## Storage Format

### Animation File Structure

Animations are stored as JSON files with the following structure:

```json
{
  "metadata": {
    "name": "my-animation",
    "created_at": "2025-01-15T10:30:00Z",
    "source": "local|ai",
    "original_prompt": "a cat dancing in the rain",
    "fps": 30,
    "total_frames": 150,
    "width": 120,
    "height": 40
  },
  "frames": {
    "1": "ASCII content for frame 1...",
    "2": "ASCII content for frame 2...",
    "3": "ASCII content for frame 3...",
    "...": "..."
  }
}
```

### Storage Location

```
~/.config/ascii-ngin/
├── config.json
└── animations/
    ├── my-animation.json
    ├── dancing-cat.json
    └── ...
```

---

## Dependencies

### Go Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/bubbletea` | TUI framework |
| `github.com/charmbracelet/bubbles` | TUI components |
| `github.com/charmbracelet/lipgloss` | TUI styling |
| `google.golang.org/genai` | Gemini API client |

### System Dependencies

| Dependency | Purpose | Required |
|------------|---------|----------|
| FFMPEG | Video frame extraction | Yes |
| Go 1.21+ | Runtime | Yes |

---

## Development Roadmap

### Phase 1: Foundation
- [ ] Project setup and structure
- [ ] Basic Bubble Tea application shell
- [ ] Configuration management
- [ ] Storage layer implementation

### Phase 2: Local Video Flow
- [ ] FFMPEG integration for frame extraction
- [ ] ASCII conversion algorithm
- [ ] Animation playback renderer
- [ ] File browser component

### Phase 3: AI Integration
- [ ] Gemini API client setup
- [ ] Nano Banana image generation
- [ ] Image preview (terminal image protocol detection)
- [ ] Veo 3 video generation

### Phase 4: Polish
- [ ] Error handling and user feedback
- [ ] Performance optimization
- [ ] Animation library management
- [ ] Sizing options for playback

### Phase 5: Release
- [ ] Documentation
- [ ] Testing
- [ ] Gemini hackathon submission preparation

---

## License

MIT License - See [LICENSE](LICENSE) for details.
