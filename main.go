package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}
	ffmpegPath := os.Getenv("FFMPEG_PATH")
	ytDlpPath := os.Getenv("YTDLP_PATH")
	w.Header().Set("Content-Type", "video/mp4")
	yt := exec.Command(ytDlpPath, "-o", "-", url)
	ffmpeg := exec.Command(ffmpegPath, "-i", "pipe:0", "-c:v", "copy", "-c:a", "copy", "pipe:1")

	pipe, err := ffmpeg.StdoutPipe()
	if err != nil {
		http.Error(w, "Error creating pipe", 500)
		return
	}

	ffmpeg.Stdin = pipe
	ffmpeg.Stdout = w

	if err := yt.Start(); err != nil {
		http.Error(w, "Error starting yt-dlp", 500)
		return
	}

	if err := ffmpeg.Start(); err != nil {
		http.Error(w, "Error starting ffmpeg", 500)
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
	log.Printf("API Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
