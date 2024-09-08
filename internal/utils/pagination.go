package utils

import (
	"fmt"
	"os"
	"os/exec"
	"ozon1/internal/entity/item"
	"runtime"
	"time"

	"github.com/eiannone/keyboard"
)

func clearTerminal() {

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}


func Paginate(items []item.Item) {

	pageSize := 5
	offset := 0
	totalItems := len(items)

	if err := keyboard.Open(); err != nil {
		fmt.Println("Error initializing keyboard:", err)
		return
	}

	defer keyboard.Close()

	for {

		clearTerminal()

		fmt.Printf("Items %d to %d:\n", offset+1, offset+pageSize)

		for i := offset; i < offset+pageSize && i < totalItems; i++ {
			fmt.Println(items[i])
		}

		fmt.Println("\nUse arrow keys to scroll (q to quit)")

		char, key, err := keyboard.GetKey()

		if err != nil {
			fmt.Println("Error reading keyboard input:", err)
			return
		}

		switch key {

		case keyboard.KeyArrowDown:
			if offset+pageSize < totalItems {
				offset++
			}
		case keyboard.KeyArrowUp:
			if offset > 0 {
				offset--
			}
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return
		case 0: 
			if char == 'q' {
				return
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}

}