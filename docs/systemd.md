## Running as a systemd Service

For a robust and continuously running setup, deploy SysCapture as a systemd service.

### Sample systemd Service File

Create a file at `/etc/systemd/system/syscapture.service` with the following content (update paths and usernames as needed):

```ini
[Unit]
Description=SysCapture Monitoring Agent
After=network.target

[Service]
User=your_username       # Replace with your local username
Group=your_username      # Replace with your local group
ExecStart=/full/path/to/syscapture   # Full path to the SysCapture binary
Restart=always
Environment="API_SECRET=your_secret"
Environment="GIN_MODE=release"
Environment="PORT=59232"

[Install]
WantedBy=multi-user.target
```

### Setup Steps

1. **Copy the Service File:**  
   Save the above content to `/etc/systemd/system/syscapture.service`.

2. **Reload systemd:**

    ```shell
    sudo systemctl daemon-reload
    ```

3. **Enable and Start the Service:**

    ```shell
    sudo systemctl enable syscapture
    sudo systemctl start syscapture
    ```

4. **Check Service Status:**

    ```shell
    sudo systemctl status syscapture
    ```

This configuration ensures SysCapture runs continuously in the background and restarts automatically if it fails or after a reboot.

---