# Eldar

[![Audit](https://github.com/ioluas/eldar/actions/workflows/audit.yml/badge.svg)](https://github.com/ioluas/eldar/actions/workflows/audit.yml)

<img src="PalestinanTux.png" alt="Palestinian Tux" width="256" height="256" style="border-radius: 25%;">

Eldar is a cross-platform tasks management application built with Go and the Fyne UI toolkit. It provides a simple interface for managing group tasks across different platforms including desktop and mobile.

## Features

- Cross-platform support (Linux, macOS, Windows, Android, iOS)
- Simple and intuitive user interface
- ...

## Installation

### Desktop

#### System-wide installation

```bash
sudo make install
```

#### User-specific installation

```bash
make user-install
```

### Mobile

Android APK is available in the repository as `eldar.apk`.

## Usage

1. Launch the application
2. ...

## Development

### Requirements

- Go 1.24 or later
- Fyne dependencies (for UI)

### Building from source

```bash
# Clone the repository
git clone https://github.com/ioluas/eldar.git
cd eldar

# Build the application
go build

# For Android build
fyne package -os android
```

### Using Docker for development

A Dockerfile is provided for development purposes:

```bash
# Build the development container
docker build -t eldar-dev .

# Run the container
docker run -it --rm -v $(pwd):/app eldar-dev
```
