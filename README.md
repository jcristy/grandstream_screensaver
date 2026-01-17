# Grandstream Screensaver

Easily serve images with the necessary XML file for a Grandstream phone's screensaver.

## Testing

I have run this on a Grandstream GRP2613W served from:
- Ubuntu Server on a Raspberry Pi with Docker
- Ubuntu Desktop 24.04 on a PC without docker

## Build

Docker:

    make build

Standalone: (I should have the make file support either, but this for now)

    go build main.go

## Running

### Docker Compose

Put the jpegs (I believe they must be <500kB) such that they will be in the /app/images directory

      grandstream_screensaver:
        image: "grandstream_screensaver:latest"
        ports:
          - 8080:8080
        volumes:
          - <your volume>:/app/images

### Standalone

Put the jpegs in ./images.

    ./grandstream_screensaver

## Configuration

If the host is messed up in the XML, you can set `SERVER_EXTERNAL_HOST`.
