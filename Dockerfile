FROM golang:1.21-alpine

WORKDIR /app

# Copia código e mod
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Baixa yt-dlp e ffmpeg direto no container
RUN apk add --no-cache curl tar xz \
 && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /app/bin/yt-dlp \
 && chmod +x /app/bin/yt-dlp \
 && wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz \
 && tar -xvf ffmpeg-release-amd64-static.tar.xz \
 && cp ffmpeg-*/ffmpeg /app/bin/ffmpeg \
 && chmod +x /app/bin/ffmpeg

ENV FFMPEG_PATH=/app/bin/ffmpeg
ENV YTDLP_PATH=/app/bin/yt-dlp
ENV PORT=8080

RUN go build -o main .

EXPOSE 8080
CMD ["./main"]
