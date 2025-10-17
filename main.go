package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	ffmpegPath := os.Getenv("FFMPEG_PATH")
	ytDlpPath := os.Getenv("YTDLP_PATH")

	// yt-dlp escreve no stdout (pipe:1)
	yt := exec.Command(ytDlpPath, "-o", "-", "-f", "best", url)

	// ffmpeg lê do stdin (pipe:0) e escreve pro stdout (pipe:1)
	ffmpeg := exec.Command(ffmpegPath, "-i", "pipe:0", "-c:v", "copy", "-c:a", "copy", "-f", "mp4", "pipe:1")

	ffmpeg.Stdin, _ = yt.StdoutPipe() // conecta saída do yt-dlp na entrada do ffmpeg
	ffmpeg.Stdout = w                 // envia saída final direto pra resposta HTTP
	ffmpeg.Stderr = os.Stderr
	yt.Stderr = os.Stderr

	w.Header().Set("Content-Type", "video/mp4")

	// inicia yt-dlp primeiro, depois ffmpeg
	if err := yt.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Error starting yt-dlp: %v", err), 500)
		return
	}
	if err := ffmpeg.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Error starting ffmpeg: %v", err), 500)
		return
	}

	// espera ambos terminarem
	yt.Wait()
	ffmpeg.Wait()
}

func main() {
	http.HandleFunc("/download", downloadHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
