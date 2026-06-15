package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Music struct {
	ID         int `gorm:"primarykey"`
	Name       string
	Link       string
	IsPlaylist bool
}

func main() {

	db, err := gorm.Open(sqlite.Open("db/data.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		return
	}

	err = db.AutoMigrate(&Music{})

	if err != nil {
		log.Fatal(err)
		return
	}

	for {

		fmt.Println("\nWelcome mpv control\n")

		fmt.Println("Options:")

		fmt.Println("1. Play music.")
		fmt.Println("2. Add music.")
		fmt.Println("3. Remove music.")
		fmt.Println("4. Exit music")

		var entry int
		fmt.Scan(&entry)

		switch entry {

		case 1:
			fmt.Println("Play music")

			var music []Music
			result := db.Find(&music)

			if result.Error != nil {
				log.Fatal(result.Error)
			}

			log.Println("Count: ", result.RowsAffected)

			if result.RowsAffected == 0 {
				fmt.Println("0 entires.\nAdd Music first.")
				break
			}

			for i := 0; i < len(music); i++ {
				fmt.Println(music[i].ID, music[i].Name)
			}

			var id int
			fmt.Println("Enter ID to play: ")
			fmt.Scan(&id)

			var musicRow Music
			db.First(&musicRow, "id=?", id)

			var loop string
			if musicRow.IsPlaylist {
				loop = "--loop-playlist"
			} else {
				loop = "--loop"
			}

			fmt.Println("Executing: mpv --no-video", "--shuffle", loop, musicRow.Link)
			cmd := exec.Command("mpv", "--no-video", "--shuffle", loop, musicRow.Link)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				fmt.Println("mpv exited with error:", err)
			}

		case 2:
			fmt.Println("Add music")
			var name, link, isPlaylist string
			var music Music
			fmt.Println("Name: ")
			fmt.Scan(&name)
			fmt.Println("Link: ")
			fmt.Scan(&link)
			fmt.Println("Playlist(P)/Single(S)[Default=P]: ")
			fmt.Scan(&isPlaylist)

			if isPlaylist == "S" || isPlaylist == "s" {
				fmt.Println("Single")

				music.Name = name
				music.Link = link
				music.IsPlaylist = false
			} else {
				fmt.Println("Playlist")

				music.Name = name
				music.Link = link
				music.IsPlaylist = true
			}

			result := db.Create(&music)
			if result.Error != nil {
				log.Fatal(result.Error)
				return
			}

			log.Println("Inserted id: ", music.ID)

		case 3:
			fmt.Println("Remove music")

		case 4:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid Option")
		}
	}
}
