# Distributed Ride-Sharing System (Uber Clone)

![Architecture](https://mermaid.ink/img/pako:eNqNVt9v2jAQ_lcsP21qGvGjZSEPlSpaTX1YxWDVpAmpMvZBIkicOQ6UVf3fd4mdEgcKzQOK47v7vjt_d-aVcimAhnSW5vC3gJTDXcyWiiWzlOCTMaVjHmcs1eQpB3X49Xb88J1p2LLd4d4vFWdTUJuYw-HmnYo3oM5sH34fs10Cqf7Qb6oRFL-bnZLz5c3NnmRIRgrwteJGJmXOuTa2eyP0aFAPyXIyHtWO5YaxN7-PEoNJpNrM1nOSCw3Y_QuPWLq0nBvWl4h30fIos_Bhg5n6vMIVDTgVLyNN5II4LO9L65A8wrbyJimAyAkjwlbSBHBwENisQ2vl80T4pfezapbGRW1RHckkYaloILMNi9dsvgYXtHUQv2E-lXwF-hCccQ5ZE3sNiwZ0A_O2sjSw6oPTbB_nWbT3TB3dGEiykGrLlABBtCQTNp_H-sdP8sUwK3EmkGcS2-lnAQV8rUtwVCchGSvJIc8tJ2KoMGxDjxSZKIVapZZrpovc9_2j4nF7whGPifvM8jxepp8XkW2SzAQmOVKMZXoyF6_Nwq5bunetkLxp2HfIUQR8JQts5Bqz9DJGx3K1ZtZdHAVpCc9mZStkc3s-0WZtTFukOsGnB5QeE7u6PM4gKScQsozktmF_TmuN1qidUHZJa6q9V04m2RowlrV1SnbQc5GUq33YacFL_R2faHtPzxXt0ZN1O-6iXbRW1Zu4HxeiqjTJivk6ziPTcp-U1aXD-NPgZ46a21KL061Qj6kzc38_facarzCyv1tOdGgVk7OUpCipOSxjZEI9moBKWCzwKn8tQ8yojiCBGQ3xVTC1muEV_4Z2rNByuks5DbUqwKNKFsuIhgu2znFlZo79C1Cb4L36R8rmkoav9IWGvW_-1XVn0O_1-kE3GAyHgUd3-Lnb8fu9frc_xKfbvQ6CN4_-qyJ0_KDX7Q86QTDoDAfD66ve23_1IPGQ)

A high-performance, fault-tolerant microservices platform built with **Golang** and **Kubernetes**, designed to handle 5k+ concurrent ride requests. This project implements advanced distributed systems patterns like **Hexagonal Architecture**, **Saga Pattern**, and **Event-Driven Architecture**.

## 👥 Team and my role

Built by **four engineers**. I co-architected the service decomposition and the
gRPC contracts between services, and I owned two pieces end to end:

- **The matching engine** — the in-memory geo-hash index and its concurrency.
  `pprof` attributed the throughput ceiling to contention on a single mutex;
  moving to read-write locks over a sharded index lifted throughput **200%**
  from that baseline.
- **The load-test and profiling campaign** — the k6 scenarios and the
  Prometheus/Grafana instrumentation that produced every number in the table
  below.

The Saga/Outbox consistency work and the payment and dispatch services were my
teammates'. Numbers below are for the system as a whole.

## 🧪 Performance Benchmarks
> Benchmarks run on a 3-node GKE Cluster (4 vCPU/16GB RAM) using **k6** and **Prometheus**.

| Metric | Result | Target / SLA |
| :--- | :--- | :--- |
| **Throughput** | **5,200 RPS** | 5,000 RPS |
| **P99 Latency** | **42ms** | < 50ms |
| **Error Rate** | **0.01%** | < 0.1% |
| **Matching Speed**| **12ms** (avg) | < 100ms |
| **Unit Tests** | **PASS** (Covered Core Logic) | 100% Pass |

## 🚀 Key Features

### 🏗️ Distributed Architecture & Design
*   **Microservices Ecosystem**: Decomposed into 4 bounded contexts: `Trip`, `Driver`, `Payment`, and `API Gateway`.
*   **Hexagonal Architecture (Ports & Adapters)**: Business logic is isolated from infrastructure concerns (RabbitMQ, MongoDB, gRPC), enabling independent testing and swapping of adapters.
*   **Event-Driven Communication**: Asynchronous data flow using **RabbitMQ** with **Topic Exchanges** and **Fair Dispatch** (QoS 1) to ensure high availability.
*   **High Performance**: Critical paths optimized to **<50ms P99 latency** using pprof profiling and efficient concurrency patterns.

### 🌐 Real-Time & Geospatial
*   **Driver Matching Engine**: Implements **In-Memory GeoHash Indexing** to perform $O(1)$ proximity queries, significantly faster than traditional spatial database lookups.
*   **Reactive Dashboard**: **Next.js** frontend with **Bi-directional WebSockets** replacing polling, providing sub-100ms live fleet tracking.

### 🛡️ Reliability & Data Consistency
*   **Saga Pattern (Choreography)**: Manages distributed transactions across Trip and Payment services.
*   **Outbox Pattern**: Ensures "at-least-once" delivery guarantee for critical events (e.g., Trip Created, Payment Success) by persisting events atomically with state changes.
*   **Fault Tolerance**: Implements **Dead Letter Queues (DLQ)** and automatic retry policies for transient failures.

---

## 🛠️ Tech Stack

| Category | Technologies |
|----------|--------------|
| **Backend** | Golang 1.23, gRPC, Protobuf |
| **Architecture** | Hexagonal, DDD, Event-Driven (RabbitMQ) |
| **Database** | MongoDB (Geospatial & ACID Transactions) |
| **Frontend** | Next.js, TypeScript, TailwindCSS, WebSockets |
| **DevOps** | Docker, Kubernetes (GKE), Tilt, GitHub Actions |
| **Observability** | Jaeger (Distributed Tracing), Prometheus, Grafana, pprof |
| **Testing** | k6 (Load Testing), Go Test |

---

## 🏃‍♂️ Getting Started

### Prerequisites
*   Docker & Kubernetes (Minikube or Docker Desktop)
*   Go 1.23+
*   Tilt (for local dev orchestration)

### One-Command Startup
We use **Tilt** to orchestrate the entire development environment (Service mesh, hot-reloading, logs):

```bash
tilt up
```

Access the dashboard at `http://localhost:10350`.

### Load Testing (Simulate 5,000 RPS)
Verify system stability under high concurrency using our calibrated **k6** script:

```bash
# Install k6
brew install k6

# Run 5k RPS Stress Test
k6 run load_tests/k6_script.js
```

### Profiling
The Driver Service exposes a **pprof** debug endpoint to analyze thread contention and memory usage:

```bash
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile
```

---

## 📂 Project Structure

```
├── services/
│   ├── driver-service/    # In-memory Matching, pprof enabled
│   ├── trip-service/      # Trip Lifecycle, Outbox Pattern
│   ├── payment-service/   # Stripe Integration, Saga Consumer
│   └── api-gateway/       # Rest -> gRPC, WebSocket Aggregator
├── shared/
│   ├── messaging/         # RabbitMQ Adapters (Fair Dispatch, DLQ)
│   ├── proto/             # gRPC Service Contracts
│   └── db/                # MongoDB Drivers
├── infra/                 # Kubernetes Manifests & Terraform
└── load_tests/            # k6 Stress Testing Scripts
```

## 📜 Deployment (GKE)
Production-ready deployment scripts are located in `infra/production`. The CI/CD pipeline (`.github/workflows/deploy.yml`) handles:
1.  Automated testing & linting.
2.  Building optimized Docker images to Google Artifact Registry.
3.  Zero-downtime rolling updates to GKE.

---
