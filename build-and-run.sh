#!/bin/bash
docker build -t "frigate2pushover" .
docker run --rm -it -v ./config.yaml:/app/config.yaml frigate2pushover:latest
