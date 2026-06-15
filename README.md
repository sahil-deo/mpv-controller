# mpv-cli

A simple command-line tool to manage and play music/video links via [mpv](https://mpv.io/), backed by a SQLite database. Save links with names, list/play/remove them, or quick-play a one-off link without saving it.

## Features

- Store music/video links in a local SQLite database (`db/data.db`)
- Play saved entries by ID, with optional loop and shuffle
- Quick-play any link directly without saving it to the database
- List and remove saved entries

## Requirements

- [Go](https://go.dev/) 1.21+
- [mpv](https://mpv.io/) installed and available in `$PATH`
- SQLite (via `gorm.io/driver/sqlite`, no separate install needed — it's a CGo/pure-Go driver bundled as a dependency)
- **Linux or macOS only** — see [Platform support](#platform-support) below

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

> **Note:** mpv itself installs fine on Windows, but this program **will not run on Windows** — see [Platform support](#platform-support).

## Running

Clone or copy the project, then from the project directory:

```bash
mkdir -p db
go run .
```

The program expects a `db/` directory in your current working directory (it opens `db/data.db`, creating it on first run via GORM's `AutoMigrate`). Run it from the project directory so `db/` is found.

```bash
go run . -add -name "Lofi Mix" -link "https://music.youtube.com/playlist?list=PL..."
go run . -list
go run . -play 1
```

## Building

Once you're happy with it, build a standalone binary:

```bash
go build -o mpv-cli .
```

This produces a binary named `mpv-cli` in the current directory. It still expects a `db/` directory in the current working directory (same as above).

## Platform support

This program uses `syscall.Exec` to replace its own process image with mpv (so mpv inherits the terminal directly). `syscall.Exec` wraps the `execve` system call, which **only exists on Unix-like systems**.

| OS | Supported |
|----|-----------|
| Linux | ✅ Yes |
| macOS | ✅ Yes |
| Windows | ❌ No — `syscall.Exec` is not implemented on Windows and the exec call will fail |

To support Windows, the `-play`/`-quickplay` code paths would need to be rewritten using `os/exec` (`exec.Command(...).Run()`) instead of `syscall.Exec`.

## Adding to `$PATH`

To run `mpv-cli` from anywhere, install the binary somewhere in your `$PATH`.

### Option 1: Use `go install`

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

### Option 2: Move the built binary manually

```bash
go build -o mpv-cli .
sudo mv mpv-cli /usr/local/bin/
```

`/usr/local/bin` is in `$PATH` by default on most Linux/macOS systems, so no further configuration is needed.

### Fixing the database path for `$PATH` usage

As written, `mpv-cli` looks for `db/data.db` relative to the current working directory. So running it from `$PATH` (i.e. from any directory) will try to create/open `db/` wherever you happen to be.

To fix this, edit `getdb()` in `main.go` to use a fixed location in your home directory:

```go
func getdb() *gorm.DB {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error finding home directory", err)
	}

	dbDir := filepath.Join(home, ".local", "share", "mpv-cli")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatal("Error creating db directory", err)
	}

	dbPath := filepath.Join(dbDir, "data.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening db", err)
	}
	return db
}
```

Add `"path/filepath"` to your imports, then rebuild (`go build -o mpv-cli .`). This stores the database at `~/.local/share/mpv-cli/data.db`, so `mpv-cli` works correctly from any directory.

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

## How it works

- Saved entries are stored as `(ID, Name, Link)` rows in the `music` table of `db/data.db` (auto-migrated on startup via GORM).
- `-play` and `-quickplay` build an `mpv` command (`mpv --no-video [--loop-playlist] [--shuffle] <link>`) and replace the current process with it via `syscall.Exec`, so mpv runs directly in your terminal with full TTY control (playback controls, progress bar, etc. work normally).
- `--loop-playlist` is used for `-play` (suited to playlist links); `--loop` (single-file loop) is used for `-quickplay`.

## Troubleshooting

- **`Enable to find mpv path` error**: mpv isn't installed or isn't in `$PATH`. Install it via your package manager (e.g. `sudo pacman -S mpv`, `sudo apt install mpv`, `brew install mpv`).
- **`Error opening db`**: ensure a `db/` directory exists in the current working directory (or update `getdb()` as described above for a fixed path).
- **"No music in db."**: the `music` table is empty — add an entry with `-add` first.