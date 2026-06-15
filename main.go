package main

import (
	"fmt"
	"os"
	"os/exec"
)


func main(){
	
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
				cmd := exec.Command("mpv", "--no-video", "--loop-playlist", "https://music.youtube.com/playlist?list=PLoiMA0z1W8dmbPKuu4dErfS4hOFi97M1L")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				
				err := cmd.Run(); 
				if err != nil {
					fmt.Println("mpv exited with error:", err) 
				}
				
			case 2:
				fmt.Println("Add music")
				break
			case 3:
				fmt.Println("Remove music")
				break
			case 4:
				fmt.Println("Exiting...")
				return 
			default:
				fmt.Println("Invalid Option")	
		}
	}	
}
