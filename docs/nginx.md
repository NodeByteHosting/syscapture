## Set Up Guide with NGINX
Deploying SysCapture with NGINX is simple and straightforward

## Setting Up NGINX

1. **Install NGINX**

    ```shell
    sudo apt update
    sudo apt install nginx
    ```

2. **Configure NGINX**

    Create a new configuration file for SysCapture:

    ```shell
    sudo nano /etc/nginx/sites-available/syscapture
    ```

    Add the following content to the file (update paths and domains as needed):

    ```nginx
    server {
        listen 80;
        server_name your_domain_or_ip;

        location / {
            proxy_pass http://localhost:42000;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
    ```

3. **Enable the Configuration**

    ```shell
    sudo ln -s /etc/nginx/sites-available/syscapture /etc/nginx/sites-enabled/
    ```

4. **Test the Configuration**

    ```shell
    sudo nginx -t
    ```

5. **Restart NGINX**

    ```shell
    sudo systemctl restart nginx
    ```

6. **Access SysCapture**

    Open your browser and navigate to `http://your_domain_or_ip` to access SysCapture.

---

By following these steps, you can set up SysCapture with NGINX when the `ALLOW_PUBLIC_API` environment variable is set to true. If you encounter any further issues, please provide the error messages for further assistance.