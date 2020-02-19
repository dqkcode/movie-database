FROM golang:latest AS movies

WORKDIR /home

# COPY ./movies-0.1.bin .
# COPY ./configs/status.yml ./configs/
# RUN chmod +x movies-0.1.bin
# CMD ["./movies-0.1.bin"]


# COPY ./ /home

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -build="go build -v -o movies-0.1.bin" -command=./movies-0.1.bin


