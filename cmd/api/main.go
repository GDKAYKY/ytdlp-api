package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "url parameter required", http.StatusBadRequest)
		return
	}

	ffmpegPath := os.Getenv("FFMPEG_PATH")
	ytdlpPath := os.Getenv("YTDLP_PATH")

	if ffmpegPath == "" || ytdlpPath == "" {
		http.Error(w, "missing FFMPEG_PATH or YTDLP_PATH", http.StatusInternalServerError)
		return
	}

	log.Printf("Starting download for: %s", url)
	w.Header().Set("Content-Type", "video/mp4")

	// yt-dlp envia o vídeo cru para stdout
	yt := exec.Command(ytdlpPath, "-o", "-", "-f", "best", url)
	// ffmpeg lê da entrada padrão e replica em mp4 via stdout
	ffmpeg := exec.Command(ffmpegPath, "-i", "pipe:0", "-c:v", "copy", "-c:a", "copy", "-f", "mp4", "pipe:1")

	ytStdout, err := yt.StdoutPipe()
	if err != nil {
		http.Error(w, "yt-dlp pipe error", 500)
		return
	}
	ffmpeg.Stdin = ytStdout
	ffmpeg.Stdout = w

	if err := yt.Start(); err != nil {
		http.Error(w, "yt-dlp start error", 500)
		return
	}
	if err := ffmpeg.Start(); err != nil {
		http.Error(w, "ffmpeg start error", 500)
		return
	}

	yt.Wait()
	ffmpeg.Wait()
}

func main() {
	http.HandleFunc("/download", downloadHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
