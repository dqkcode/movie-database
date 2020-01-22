FROM alpine:3.8

WORKDIR /home
COPY ./movies-0.1.bin .
COPY ./configs/status.yml ./configs/
RUN chmod +x movies-0.1.bin


CMD ["./movies-0.1.bin"]