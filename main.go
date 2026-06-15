package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Music struct {
	ID   int `gorm:"primarykey"`
	Name string
	Link string
}

func main() {

	db := getdb()
	err := db.AutoMigrate(&Music{})

	if err != nil {
		log.Fatal(err)
		return
	}

	if len(os.Args) < 2 {
		help()
		return
	}

	args := os.Args[1:]

	var tree = make(map[string]int)

	for i, arg := range args {
		switch arg {
		case "-a", "-add":
			tree["ADD"] = i

		case "-r", "-remove":
			tree["REMOVE"] = i

		case "-p", "-play":
			tree["PLAY"] = i

		case "-l", "-ls", "-list":
			tree["LIST"] = i

		case "-loop":
			tree["LOOP"] = i

		case "-n", "-name":
			tree["NAME"] = i

		case "-lk", "-link":
			tree["LINK"] = i

		case "-sh", "-shuffle":
			tree["SHUFFLE"] = i

		case "-h", "-help":
			tree["HELP"] = i

		case "-q", "-quickplay":
			tree["QUICK"] = i

		case "-d", "-default":
			tree["SETDEFAULT"] = i
		}

	}

	// precedence:
	// list -> play -> add -> remove -> quickplay -> default -> help

	if _, ok := tree["LIST"]; ok {

		listMusic()

	} else if val, ok := tree["PLAY"]; ok {
		id, err := strconv.Atoi(args[val+1])

		if err != nil {
			log.Fatal("Error parsing id", err)
			return
		}

		var cmdargs []string = []string{"mpv", "--no-video"}
		if _, ok := tree["LOOP"]; ok {
			cmdargs = append(cmdargs, "--loop-playlist")
		}
		if _, ok := tree["SHUFFLE"]; ok {
			cmdargs = append(cmdargs, "--shuffle")
		}

		playMusic(id, cmdargs)

	} else if _, ok := tree["ADD"]; ok {

		var name, link string

		if i, ok := tree["NAME"]; ok {
			name = args[i+1]
		} else {
			log.Fatal("name not provided\nUse -h or -help to see usage")
			return
		}

		if i, ok := tree["LINK"]; ok {
			link = args[i+1]
		} else {
			log.Fatal("link not provided\nUse -h or -help to see usage")
			return
		}

		addMusic(name, link)

	} else if val, ok := tree["REMOVE"]; ok {

		id, err := strconv.Atoi(args[val+1])

		if err != nil {
			log.Fatal("Error parsing id", err)
			return
		}
		removeMusic(id)

	} else if val, ok := tree["QUICK"]; ok {

		link := args[val+1]

		var cmdargs []string = []string{"mpv", "--no-video"}

		if _, ok := tree["LOOP"]; ok {
			cmdargs = append(cmdargs, "--loop-playlist")
		}
		if _, ok := tree["SHUFFLE"]; ok {
			cmdargs = append(cmdargs, "--shuffle")
		}

		cmdargs = append(cmdargs, link)

		playQuickMusic(cmdargs)

	} else if _, ok := tree["HELP"]; ok {

		help()

	} else {
		fmt.Println("Invalid options, use -h or -help to check usage")
		return
	}
}

func listMusic() {
	db := getdb()

	var entries []Music
	result := db.Find(&entries)

	if len(entries) == 0 {
		log.Fatal("No music in db.")
		return
	}

	log.Println("Total Count:", result.RowsAffected)

	for _, entry := range entries {
		log.Println(entry.ID, entry.Name)
	}
}

func playMusic(id int, cmdargs []string) {
	db := getdb()

	var entry Music
	db.First(&entry, "id=?", id)

	link := entry.Link
	cmdargs = append(cmdargs, link)
	env := os.Environ()

	// cmd := exec.Command("mpv", cmdargs, link)
	mpv, err := exec.LookPath("mpv")

	if err != nil {
		log.Fatal("Enable to find mpv path", err)
	}

	syscall.Exec(mpv, cmdargs, env)

}

func addMusic(name, link string) {
	db := getdb()

	var entry Music
	entry.Name = name
	entry.Link = link

	db.Create(&entry)

	log.Println("Music Added: ", entry.ID, entry.Name, entry.Link)
}

func removeMusic(id int) {
	db := getdb()
	db.Delete(&Music{}, id)
}

func playQuickMusic(cmdargs []string) {

	env := os.Environ()
	mpv, err := exec.LookPath("mpv")

	if err != nil {
		log.Fatal("Enable to find mpv path", err)
	}

	syscall.Exec(mpv, cmdargs, env)

}

func help() {
	fmt.Println(`
	mpv-cli — manage and play music links via mpv

	Usage:
	mpv-cli [options]

	Options:
	-a,  -add              Add a new music entry
	-r,  -remove <id>      Remove music entry by ID
	-p,  -play <id>        Play music entry by ID
	-l,  -ls, -list        List all music entries
	-q,  -quickplay <link> Play a link directly without saving
	-n,  -name <name>      Name of the music (used with -add)
	-lk, -link <link>      Link of the music (used with -add)
	-loop                  Loop the track/playlist
	-sh, -shuffle          Shuffle playback
	-h,  -help             Show this help message

	Examples:
	mpv-cli -add -name "Lofi Mix" -link "https://youtube.com/watch?v=..."
	mpv-cli -play 3 -loop -shuffle
	mpv-cli -quickplay "https://youtube.com/watch?v=..." -loop
	mpv-cli -list
	mpv-cli -remove 2
	`)
}
func getdb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db/data.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Error opening db", err)
	}
	return db
}
