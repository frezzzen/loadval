# LOADVAL - Valorant Loadout Manager

A modern Windows desktop application for managing Valorant loadouts, built with Wails v2, Go, and Svelte 5.

## 🎮 Features

- **Real-time Loadout Management**: Automatically sync with your current Valorant loadout
- **Template System**: Create and save custom loadout templates for quick switching
- **Agent-Specific Loadouts**: Set different loadouts for each agent automatically
- **Skin Management**: View and manage your owned weapon skins
- **Modern UI**: Beautiful, responsive interface
- **Auto-Detection**: Automatically detects u picking agent in agent select

## 🚀 Quick Start

### Prerequisites

- Windows 10/11
- Valorant installed and running
- Go 1.24+ (for development)
- Node.js 18+ (for development)

### Installation

1. Download the latest release from the [Releases](https://github.com/frezzzen/loadval/releases) page
2. Extract the executable
3. Run `loadval.exe`

### First Time Setup

1. Launch Valorant and log into your account
2. Start LOADVAL
3. The application will automatically detect your region and authenticate
4. Your current loadout and owned items will be displayed

## 🛠️ Development

### Prerequisites

- Go 1.24+
- Node.js 18+
- Wails v2

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/loadval.git
   cd loadval
   ```

2. Install dependencies:
   ```bash
   # Install Go dependencies
   go mod download
   
   # Install frontend dependencies
   cd frontend
   npm install
   cd ..
   ```

3. Run in development mode:
   ```bash
   wails dev
   ```

### Building

Build for production:
```bash
wails build
```

Build with specific tags:
```bash
wails build -tags webkit2_41
```

## 📖 Usage

### Managing Loadouts

1. **View Current Loadout**: Your current loadout is automatically displayed when the app starts
2. **Create Templates**: 
   - Click "Add Template" to create a new loadout template
   - Customize your weapon skins
   - Save the template with a custom name
3. **Apply Templates**: Select any saved template to instantly apply it to your loadout
4. **Agent-Specific Loadouts**: Enable agent-specific loadouts to automatically switch loadouts based on your selected agent

### Template Management

- **Create**: Design custom loadouts and save them as templates
- **Edit**: Modify existing templates
- **Delete**: Remove unwanted templates
- **Quick Apply**: Instantly apply templates with a single click

### Agent Loadouts

Enable agent-specific loadouts to automatically switch your loadout when you select different agents in Valorant. This feature allows you to have different loadouts for each agent without manual switching.

## 🔧 Technical Details

### Architecture

- **Backend**: Go with Wails v2 framework
- **Frontend**: Svelte 5 with TypeScript
- **Styling**: SCSS
- **Build Tool**: Vite
- **Authentication**: Valorant's local API integration


### Data Flow

1. Application reads Valorant's log files to detect region and shard
2. Connects to Valorant's local API using lockfile authentication
3. Retrieves player loadout and owned items
4. Displays data in the modern Svelte interface
5. Allows real-time loadout modifications

## 🛡️ Security & Privacy

- **Local Only**: All data processing happens locally on your machine
- **No External Servers**: No data is sent to external servers
- **Valorant Integration**: Uses Valorant's official local API
- **Secure Authentication**: Leverages Valorant's built-in authentication system



## 🙏 Acknowledgments

- Built with [Wails v2](https://wails.io/)
- Frontend powered by [Svelte 5](https://svelte.dev/)
- Valorant API integration

## 📞 Support

For support, feature requests, or bug reports, please open an issue on GitHub.

---

**Note**: This application is not affiliated with Riot Games or Valorant. It's a third-party tool for managing your loadouts.
