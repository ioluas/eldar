FROM golang:bookworm

RUN apt-get update -qq \
    && apt-get install -y --no-install-recommends gcc libgl1-mesa-dev xorg-dev libxkbcommon-dev \
    && rm -rf /var/apt/list/*
