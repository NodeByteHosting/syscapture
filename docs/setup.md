## Set Up Guide
Deploying SysCapture is simple and straightforward

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
    go build -o syscapture ./cmd/api/
    ```

5. **Run the Project**

    Directly:

    ```shell
    ./syscapture
    ```

    or using `go run`:

    ```shell
    go run ./cmd/api/
    ```

6. **Environment Variables**

    If you want to change the port, API secret, allow public API, or set allowed origins, you can use the following environment variables:

    | Variable         | Description                                      | Example Value          |
    |------------------|--------------------------------------------------|------------------------|
    | `PORT`           | Port on which the server will run                | `8080`                 |
    | `API_SECRET`     | Secret key for API authentication                | `your_secret`          |
    | `ALLOW_PUBLIC_API` | Allow public access to the API                  | `true` or `false`      |
    | `GIN_MODE`       | Mode for Gin framework (`release` or `debug`)    | `release`              |
    | `ALLOWED_ORIGINS`| Comma-separated list of allowed origins for CORS | `origin1,origin2` or `*` |

    Example Usage:

    ```shell
    PORT=8080 API_SECRET=your_secret ALLOW_PUBLIC_API=true GIN_MODE=release ALLOWED_ORIGINS=* ./syscapture
    ```