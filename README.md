# mpv-cli

A simple command-line tool to manage and play music/video links via [mpv](https://mpv.io/), backed by a SQLite database. Save links with names, list/play/remove them, or quick-play a one-off link without saving it.

## Features

- Stores music/video links in a SQLite database
- Play saved entries by ID, with optional loop and shuffle
- Quick-play any link directly without saving it to the database
- Works on Linux, macOS, and Windows

## Requirements

- [Go](https://go.dev/) 1.21+
- [mpv](https://mpv.io/) installed and available in `$PATH`
- SQLite (via `gorm.io/driver/sqlite`, no separate install needed — it's a CGo/pure-Go driver bundled as a dependency)

## Installing mpv

### Linux

```bash
# Arch / Manjaro
sudo pacman -S mpv

# Debian / Ubuntu
sudo apt install mpv

# Fedora
sudo dnf install mpv
```

### macOS

```bash
brew install mpv
```

### Windows

```powershell
winget install mpv-player.mpv
```
or via [Chocolatey](https://chocolatey.org/):
```powershell
choco install mpv
```

## Installation

### Option 1: Download a prebuilt binary (recommended)

Prebuilt binaries for Linux, macOS, and Windows are available on the [Releases](../../releases) page. Download the one for your platform, then:

**Linux / macOS:**
```bash
chmod +x mpv-cli
./mpv-cli -help
```

**Windows:**
Just run `mpv-cli.exe` from PowerShell or Command Prompt.

To run it from anywhere, see [Adding to `$PATH`](#adding-to-path) below.

### Option 2: Build from source

Requires [Go](https://go.dev/) 1.21+.

```bash
git clone <repo-url>
cd mpv-cli
go build -o mpv-cli .
```

On Windows this produces `mpv-cli.exe`.


## Running


```bash
./mpv-cli -add -name "Lofi Mix" -link "https://music.youtube.com/playlist?list=PL..."
./mpv-cli -list
./mpv-cli -play 1
```

## Adding to `$PATH`

To run `mpv-cli` from anywhere, install the binary somewhere in your `$PATH`.

### Linux / macOS

**Using `go install`** (if built from source):

```bash
go install .
```

This places the binary in `$GOPATH/bin` (commonly `~/go/bin`). Make sure that directory is in your `$PATH`:

```bash
# Add to ~/.bashrc, ~/.zshrc, etc.
export PATH="$HOME/go/bin:$PATH"
```

Then reload your shell config:

```bash
source ~/.bashrc   # or ~/.zshrc
```

**Moving a binary manually** (prebuilt or self-built):

```bash
sudo mv mpv-cli /usr/local/bin/
```

`/usr/local/bin` is in `$PATH` by default on most Linux/macOS systems.

### Windows

Move `mpv-cli.exe` to a folder of your choice (e.g. `C:\Tools\`), then add that folder to your `PATH`:

1. Search "Environment Variables" in the Start menu → "Edit the system environment variables"
2. Click "Environment Variables..."
3. Under "User variables", select `Path` → "Edit" → "New"
4. Add the folder path (e.g. `C:\Tools`)
5. Click OK on all dialogs, then open a new terminal

## Usage

```
mpv-cli [options]
```

### Options

| Flag | Aliases | Description |
|------|---------|-------------|
| `-a` | `-add` | Add a new music entry |
| `-r <id>` | `-remove <id>` | Remove music entry by ID |
| `-p <id>` | `-play <id>` | Play music entry by ID |
| `-l` | `-ls`, `-list` | List all saved entries |
| `-q <link>` | `-quickplay <link>` | Play a link directly without saving |
| `-n <name>` | `-name <name>` | Name for the entry (used with `-add`) |
| `-lk <link>` | `-link <link>` | Link for the entry (used with `-add`) |
| `-loop` | | Loop playback (playlist loop for `-play`, file loop for `-quickplay`) |
| `-sh` | `-shuffle` | Shuffle playback |
| `-h` | `-help` | Show help message |

### Precedence

If multiple action flags are given, only one runs, in this order:

```
list > play > add > remove > quickplay > help
```

### Examples

**Add a new entry:**
```bash
mpv-cli -add -name "Lofi Mix" -link "https://music.youtube.com/playlist?list=PL..."
```

**List all saved entries:**
```bash
mpv-cli -list
```

**Play entry #3, looping the playlist with shuffle:**
```bash
mpv-cli -play 3 -loop -shuffle
```

**Quick-play a link without saving, looping the single track:**
```bash
mpv-cli -quickplay "https://youtube.com/watch?v=..." -loop
```

**Remove entry #2:**
```bash
mpv-cli -remove 2
```

**Show help:**
```bash
mpv-cli -help
```
Running `mpv-cli` with no arguments also shows help.

## Project structure

```
.
├── main.go              # CLI parsing, DB models, command logic
├── exec_unix.go         # //go:build !windows — syscall.Exec implementation
└── exec_windows.go      # //go:build windows — os/exec implementation
```

## Troubleshooting

- **`Enable to find mpv path` error**: mpv isn't installed or isn't in `$PATH`. Install it via your package manager (see [Installing mpv](#installing-mpv)).
- **`Error opening db`**: check that `~/.local/share/mpv-cli/` exists and is writable.
- **"No music in db."**: the `music` table is empty — add an entry with `-add` first.