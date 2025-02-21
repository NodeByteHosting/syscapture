# SysCapture

[![Visit Nodebyte](https://cordx.lol/users/510065483693817867/S3XUM4iK.png)](https://nodebytehosting.github.io/syscapture/)

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Installation and Setup](#installation-and-setup)
  - [General Installation](#general-installation)
  - [Systemd Setup Guide](#systemd-setup-guide)
  - [NGINX Setup Guide](#nginx-setup-guide)
- [Configuration Setup](#configuration-setup)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [License](#license)
- [Contact Information](#contact-information)

## Overview
**SysCapture** is an open-source hardware monitoring agent that collects vital system information and exposes it via a RESTful API for easy integration with monitoring services like Prometheus.

> **Note:** SysCapture is currently available only on **Linux**.

## Features
- **Hardware Monitoring:** Captures CPU, memory, disk, and host details.
- **RESTful API:** Retrieve metrics quickly via HTTP endpoints.
- **Lightweight:** Minimal system overhead.
- **Extensible:** Fully open source, allowing for customization.

## Installation and Setup
For detailed installation and setup instructions, please refer to the following documents:

### General Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/nodebytehosting/syscapture.git
   cd syscapture
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build the application:
   ```bash
   make build
   ```
4. Run the application:
   ```bash
   make dev
   ```

### Systemd Setup Guide
To run SysCapture as a systemd service, create a service file:
```bash
sudo nano /etc/systemd/system/syscapture.service
```
Add the following content:
```ini
[Unit]
Description=SysCapture Service
After=network.target

[Service]
Type=simple
User=your_username
ExecStart=/path/to/syscapture
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
Enable and start the service:
```bash
sudo systemctl enable syscapture
sudo systemctl start syscapture
```

### NGINX Setup Guide
To set up NGINX as a reverse proxy for SysCapture, create a configuration file:
```bash
sudo nano /etc/nginx/sites-available/syscapture
```
Add the following content:
```nginx
server {
    listen 80;
    server_name your_domain_or_IP;

    location / {
        proxy_pass http://localhost:42000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```
Enable the configuration and restart NGINX:
```bash
sudo ln -s /etc/nginx/sites-available/syscapture /etc/nginx/sites-enabled/
sudo systemctl restart nginx
```

## Configuration Setup

To set up the configuration for this project, please refer to the example files available in the `temp` directory. You can find:
- `.env.example`: This file contains environment variables needed for the application.
- `config.example.yml`: This file contains configuration settings for the application.

You can move **one** of these files to your project's root directory and edit it to customize the configuration.

## Usage
After installation, you can access SysCapture by navigating to `http://localhost:42000/` in your web browser. The application will provide real-time monitoring data.

## API Documentation
SysCapture provides a RESTful API for accessing system metrics. The following endpoints are available:
- `GET /api/metrics`: Retrieve all system metrics.
- `GET /api/metrics/cpu`: Retrieve CPU usage metrics.
- `GET /api/metrics/memory`: Retrieve memory usage metrics.
- `GET /api/metrics/disk`: Retrieve disk usage metrics.

## Contributing
We welcome contributions! Please refer to the [Contributing Guide](CONTRIBUTING.md) for more information on how to get involved. Follow our coding conventions and include tests where applicable.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.