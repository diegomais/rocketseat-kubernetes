# 🚀 Kubernetes Challenge: API + PostgreSQL (with Kind)

This project demonstrates deploying a simple API connected to a PostgreSQL database on a Kubernetes cluster using **Kind**. It includes configurations for Deployments, Services, Persistent Volumes, Secrets, ConfigMaps, Probes, and Horizontal Pod Autoscaler (HPA).

---

## 📁 Project Structure

```
├── go-api
│   ├── create_items_table.sql
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── k8s-api
│   ├── config-map.yaml
│   ├── deployment.yaml
│   ├── horizontal-pod-autoscaler.yaml
│   ├── namespace.yaml
│   ├── secret.yaml
│   └── service.yaml
├── k8s-cluster
│   ├── kind.yaml
│   └── metrics-server.yaml
├── k8s-db
│   ├── config-map.yaml
│   ├── deployment.yaml
│   ├── namespace.yaml
│   ├── persistent-volume-claim.yaml
│   ├── secret.yaml
│   └── service.yaml
├── CHALLENGE.md
├── LICENSE
└── README.md
```

---

## 🛠️ Prerequisites

- [Docker](https://www.docker.com/)
- [Kind](https://kind.sigs.k8s.io/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [Lens](https://k8slens.dev/) (Optional)

---

## :bricks: Cluster Setup

Navigate to root directory and run the following commands to proceed.

### 1. Create the Kind Cluster

```bash
kind create cluster --config=k8s-cluster/kind.yaml --name=challenge-k8s
```

---

## :whale: API Image

Navigate to `go-api` directory and run the following commands to proceed.

### 2. Build Container Images

```bash
docker build -t username/go-api:acf6ec8 .
```

### 3. Push to Docker Hub

```bash
docker push username/go-api:acf6ec8
```

---

## :rocket: Deployment Steps

Navigate to root directory and run the following commands to proceed.

### 4. Create Namespaces

```bash
kubectl apply -f k8s-api/namespace.yaml
kubectl apply -f k8s-db/namespace.yaml
```

### 5. Create Secrets

```bash
kubectl apply -f k8s-api/secret.yaml
kubectl apply -f k8s-db/secret.yaml
```

### 6. Create ConfigMap

```bash
kubectl apply -f k8s-api/config-map.yaml
kubectl apply -f k8s-db/config-map.yaml
```

### 7. Create PersistentVolumeClaim

```bash
kubectl apply -f k8s-db/persistent-volume-claim.yaml
```

### 8. Deploy Database (PostgreSQL)

```bash
kubectl apply -f k8s-db/deployment.yaml
```

### 9. Create Database Service

```bash
kubectl apply -f k8s-db/service.yaml
```

### 10. Deploy API

```bash
kubectl apply -f k8s-api/deployment.yaml
```

### 11. Create API Service

```bash
kubectl apply -f k8s-api/service.yaml
```

### 12. Download Metrics Server YAML manifest

```bash
curl --output k8s-cluster/metrics-server.yaml -L https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

### 13. Configure Metrics Server to not verify the CA

Add arg `--kubelet-insecure-tls` to k8s-cluster/metrics-server.yaml > Deployment.

### 14. Install Metrics Server

```bash
kubectl apply -f k8s-cluster/metrics-server.yaml
```

### 15. Create Horizontal Pod Autoscaler

```bash
kubectl apply -f k8s-api/horizontal-pod-autoscaler.yaml
```

---

## :test_tube: Testing

### 16. Accessing the API

To access the API, you can use `kubectl port-forward` to forward the service port to your local machine:

```bash
kubectl port-forward service/api-svc 8080:80 -n challenge-api
```

## :electric_plug: API Endpoints

After forwarding the port, access your API via:

```bash
curl http://localhost:8080/status
```

### Routes

* `GET /status` → Checks DB connectivity (`Connection OK`)
* `POST /data` → Inserts data into the database
* `GET /data` *(optional)* → Lists stored data

---

## :balance_scale: HPA Auto-Scaling

Check HPA status:

```bash
kubectl get hpa -n challenge-api
```

Simulate CPU load to test scaling behavior (e.g. with a stress route or external tool).

### 17. Add mock data

```bash
for i in $(seq 1 1000); do        
  curl -X POST -H "Content-Type: application/json" -d "{\"name\": \"Item $i\"}" http://localhost:8080/data
  echo "Sent request $i"
done
```

### 18. Load Testing

```bash
kubectl run -it fortio -n challenge-api --rm --image=fortio/fortio -- load -qps 6000 -t 120s -c 50 "http://api-svc/data"
```

---

## :chart_with_upwards_trend: Observability

### Logs

```bash
kubectl logs -n challenge-api deploy/go-api
kubectl logs -n challenge-db deploy/postgres
```

### Resource Usage

```bash
kubectl top pods -n challenge-api
kubectl top pods -n challenge-db
```

---

## :package: Kubernetes Features Used

* **Namespaces** (`challenge-api`, `challenge-db`)
* **Deployments** for API and DB
* **ClusterIP Services** for internal communication
* **Persistent Volumes/Claims** for DB persistence
* **Secrets/ConfigMaps** for environment management
* **Probes** (liveness/readiness) for health checks
* **HPA** with Metrics Server

---

## :white_check_mark: Outcome

You will:

* Run a full app (API + PostgreSQL) inside Kubernetes with **Kind**
* Scale the API with HPA
* Use probes to manage readiness and health
* Persist data across pod restarts
* Observe and monitor the cluster’s health and behavior

---

## :technologist: Author

Made with :heart: by [Diego Mais](https://diegomais.github.io/) :wave:.

Created as part of a Kubernetes infrastructure challenge to build, deploy, and scale real-world microservices using best practices.
