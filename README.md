# frigate2pushover
A real tiny golang based Docker container to sent Frigate NVR Messages/Snapshots as Pushover notifications. 
This is my first golang project so stay tuned

## How to use
You need Frigate and a MQTT Broker running and configured. In my case I have a RasPi 4 with Docker Compose.

1. Copy the config.yaml-template to config.yaml
2. Adopt the values in config.yaml
3. Build the container and run the container see build-and-run.sh
```bash
#!/bin/bash
docker build -t "frigate2pushover" .
docker run --rm -it -v ./config.yaml:/app/config.yaml frigate2pushover:latest
```

##  TODO
- ~~subscribe to any topic snapshot via config~~
- ~~Decode Pictures and send Pictures~~
- Map Topics to messages templates via config.yaml
- Make Messages more configurable
- more hardening
- use Event Topic instead of snapshots as recommended here https://github.com/blakeblackshear/frigate/issues/8397#issuecomment-1787116272

## References
- Very small Docker Container Build for golang https://klotzandrew.com/blog/smallest-golang-docker-image
- Pushover golang API 
