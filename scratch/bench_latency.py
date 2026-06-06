import time
import urllib.request
import json
import statistics

BASE_URL = "http://localhost:4300"

def measure_latency(endpoint, method="GET", payload=None):
    latencies = []
    url = f"{BASE_URL}{endpoint}"

    for _ in range(20):
        start = time.perf_counter()
        if method == "GET":
            with urllib.request.urlopen(url) as response:
                response.read()
        else:
            data = json.dumps(payload).encode('utf-8')
            req = urllib.request.Request(url, data=data, method='POST', headers={'Content-Type': 'application/json'})
            with urllib.request.urlopen(req) as response:
                response.read()
        end = time.perf_counter()
        latencies.append((end - start) * 1000)

    print(f"Endpoint: {endpoint}")
    print(f"  Min: {min(latencies):.2f}ms")
    print(f"  Max: {max(latencies):.2f}ms")
    print(f"  Avg: {statistics.mean(latencies):.2f}ms")
    if len(latencies) >= 20:
        print(f"  P95: {statistics.quantiles(latencies, n=20)[18]:.2f}ms")

if __name__ == "__main__":
    print("Go Sidecar REST API Latency Benchmarks")
    try:
        measure_latency("/api/index")
        measure_latency("/api/native/tools/list")
        measure_latency("/api/agent/tool", method="POST", payload={"name": "prompt_list", "arguments": {}})
        measure_latency("/api/agent/tool", method="POST", payload={"name": "ls", "arguments": {"path": "."}})
    except Exception as e:
        print(f"Error: {e}")
