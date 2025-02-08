## Set Up Guide
Deploying SysCapture is simple and straightforward

### Go Package
You can install Syscapture using the `go install` command.

```shell
go install github.com/nodebytehosting/syscapture@latest
```

### Build from Source
You can build SysCapture from its source code.

1. **Clone the Repository**

    ```shell
    git clone https://github.com/NodeByteHosting/syscapture
    ```

2. **Change to the Project Directory**

    ```shell
    cd syscapture
    ```

3. **Install Dependencies**

    ```shell
    go mod download
    ```

4. **Build the Project**

    ```shell
    go build -o dist/syscapture ./cmd/syscapture/
    ```

5. **Run the Project**

    Directly:

    ```shell
    ./dist/syscapture
    ```

    or using `go run`:

    ```shell
    go run ./cmd/syscapture/
    ```

6. **Environment Variables**

    If you want to change the Default Port or API secret, you can use the following environment variables:

   > **NOTE**: an api secret is required to interact with `Syscapture's` "API".

   | Variable         | Description                                      | Example Value          | Required |
   |------------------|--------------------------------------------------|------------------------|----------|
   | `PORT`           | Port on which the server will run (def: 42000)   | `8080`                 | No       |
   | `API_SECRET`     | Secret key for API authentication (required)     | `your_secret`          | Yes      |
   | `GIN_MODE`       | Mode in which Gin will run (release/debug)       | `release`              | No       |

   > **INFO**: Your API Secret can be used to authenticate requests to the server from services like Prometheus.

   - **Example Usage**:
   ```shell
     PORT=8080 API_SECRET=your_secret ./dist/syscapture
   ```