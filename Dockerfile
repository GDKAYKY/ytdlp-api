FROM golang:1.25-alpine

RUN apk add --no-cache ffmpeg python3 py3-pip && \
    pip install yt-dlp

WORKDIR /app
COPY . .

RUN go build -o server ./cmd/api

ENV FFMPEG_PATH=/usr/bin/ffmpeg
ENV YTDLP_PATH=/usr/bin/yt-dlp
ENV PORT=8080

EXPOSE 8080

CMD ["./server"]
