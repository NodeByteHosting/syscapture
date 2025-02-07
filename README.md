# SysCapture

[![Visit Nodebyte](https://nodebyte.host/banner.png)](https://nodebyte.host)

**SysCapture** is an open source hardware monitoring agent that collects vital system information and exposes it via a RESTful API for easy integration with monitoring services like Prometheus.  

> **Note:** SysCapture is currently available only on **Linux**.

---

## Features

- **Hardware Monitoring:** Captures CPU, memory, disk, and host details.
- **RESTful API:** Retrieve metrics quickly via HTTP endpoints.
- **Lightweight:** Minimal system overhead.
- **Extensible:** Fully open source, allowing for customization.

---

## Installation and Setup

For detailed installation and setup instructions, please refer to the following documents:

- [Setup Guide](./guides/setup.md)
- [NGINX Configuration](./guides/nginx.md)
- [Systemd Service](./guides/systemd.md)

---

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m 'Add feature'`).
4. Push the branch (`git push origin feature/your-feature`).
5. Open a pull request.

Please follow our coding conventions and include tests where applicable.

---