# Basic Monitoring Examples

## Quick Start

### Using cURL

```bash
# Get all system metrics
curl -H "Authorization: Bearer your-secret-here" http://localhost:42000/api/metrics

# Get specific metrics
curl -H "Authorization: Bearer your-secret-here" http://localhost:42000/api/metrics/cpu
curl -H "Authorization: Bearer your-secret-here" http://localhost:42000/api/metrics/memory
```

### Using PowerShell

```powershell
$headers = @{
    Authorization = "Bearer your-secret-here"
}

# Get all metrics
Invoke-RestMethod -Uri "http://localhost:42000/api/metrics" -Headers $headers

# Get specific metrics
Invoke-RestMethod -Uri "http://localhost:42000/api/metrics/cpu" -Headers $headers
```

## Code Examples

### Go Client

```go
package main

import (
    "fmt"
    "net/http"
    "encoding/json"
)

func main() {
    client := &http.Client{}
    req, _ := http.NewRequest("GET", "http://localhost:42000/api/metrics", nil)
    req.Header.Add("Authorization", "Bearer your-secret-here")

    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    var metrics map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&metrics)
    fmt.Printf("CPU Usage: %.2f%%\n", metrics["cpu"].(map[string]interface{})["usage_percent"])
}
```

### Python Client

```python
import requests

headers = {
    'Authorization': 'Bearer your-secret-here'
}

# Get all metrics
response = requests.get('http://localhost:42000/api/metrics', headers=headers)
metrics = response.json()

print(f"CPU Usage: {metrics['cpu']['usage_percent']}%")
print(f"Memory Usage: {metrics['memory']['usage_percent']}%")
```