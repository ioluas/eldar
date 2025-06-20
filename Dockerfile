FROM golang:bookworm

RUN apt-get update -qq \
    && apt-get install -y --no-install-recommends gcc libgl1-mesa-dev xorg-dev libxkbcommon-dev \
    && rm -rf /var/apt/list/* \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install golang.org/x/vuln/cmd/govulncheck@latest
