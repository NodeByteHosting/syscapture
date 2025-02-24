# SysCapture

[![Visit Nodebyte](https://cordx.lol/users/510065483693817867/S3XUM4iK.png)](https://nodebytehosting.github.io/syscapture/)

## Overview
**SysCapture** is a powerful system monitoring agent that provides real-time hardware metrics, system alerts, and monitoring capabilities through a RESTful API. It's designed for easy integration with existing monitoring stacks and supports multiple notification channels.

### Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Quick Start](#quick-start)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Features

### Core Features
- **Hardware Monitoring:** Real-time CPU, memory, disk, and network metrics
- **RESTful API:** Comprehensive HTTP endpoints with authentication
- **Alert System:** Configurable thresholds and notifications
- **Multi-Channel Notifications:** Support for Discord, Slack, and Email

### Monitoring Capabilities
- CPU usage and temperature
- Memory utilization and swap
- Disk usage and I/O metrics
- Network bandwidth and connections
- DDOS detection and prevention

### Notification Channels
- Discord webhooks with rich embeds
- Email support (SMTP, SendGrid, Postmark, Resend)
- Slack webhooks
- Custom notification templates

## Quick Start

### Installation
```bash
# Clone the repository
git clone https://github.com/nodebytehosting/syscapture.git
cd syscapture

# Install dependencies
make install

# Start the service
make start
```

### Basic Usage
```bash
# Get all metrics
curl http://localhost:42000/api/metrics

# Get specific metrics
curl http://localhost:42000/api/metrics/cpu
curl http://localhost:42000/api/metrics/memory
curl http://localhost:42000/api/metrics/disk
```

## Documentation

### Setup Guides
- [Available Commands](guides/MAKEFILE.md)
- [Configuration Guide](guides/setup/CONFIG.md)
- [Systemd Service](guides/setup/SYSTEMD.md)
- [NGINX Configuration](guides/setup/NGINX.md)

### Feature Documentation
- [Plugin System](guides/features/PLUGINS.md)
- [Monitoring System](guides/features/MONITORS.md)
- [Notification System](guides/features/NOTIFICATIONS.md)
- [Authentication](guides/features/AUTHENTICATION.md)
- [API Reference](guides/api/README.md)

### Examples
- [Basic Monitoring](guides/examples/BASIC-MONITORING.md)
- [Alert Configuration](guides/examples/ALERT-CONFIGURATION.md)
- [Custom Notifications](guides/examples/CUSTOM-NOTIFICATIONS.md)
- [Custom Plugins](guides/examples/PLUGIN-SYSTEM.md)

For detailed configuration options, see the [Configuration Guide](guides/setup/CONFIG.md).

## Contributing
We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support
- [GitHub Issues](https://github.com/nodebytehosting/syscapture/issues)
- [Discord Community](https://discord.gg/f99Rr9UybB)