# System Architecture Strategy & Objectives

## 1. Goal: Architect for Hyper-Scale & Fault Tolerance
**Strategic Intent**: To transition from a monolithic prototype to a distributed system capability of supporting 10k+ concurrent trips. The architecture prioritizes **availability** and **partition tolerance** (AP) for core booking flows, ensuring that temporary downstream failures do not result in revenue loss.

### Objective 1.1: Event-Driven Decoupling
**Architectural Decision**: Adopt an asynchronous Event Sourcing pattern over synchronous REST/gRPC chains for state mutations.
**Impact**:
- **Temporal Decoupling**: The *Booking* domain does not need the *Driver* domain to be online to accept a request. Orders are buffered in **RabbitMQ**.
- **Pressure Release**: During demand spikes (e.g., NYE rain), the message broker acts as a shock absorber. We utilize **fair dispatch QoS** (`internal/messaging/rabbitmq.go`) to prevent consumer overload, ensuring predictable latency variance.

### Objective 1.2: Resilient Error Handling Strategy
**Architectural Decision**: Move beyond "try-catch" to systemic self-healing.
**Impact**:
- **Automatic Recovery**: Implemented **Exponential Backoff** strategies for consumers. Transient network blips don't crash the pipeline; they simply delay processing by milliseconds.
- **Poison Message Containment**: Misconfigured or corrupted payloads are automatically shunted to a **Dead Letter Exchange (DLX)** after strict retry limits. This preserves the main processing queue's throughput and allows for isolated post-mortem analysis without service interruption.

---

## 2. Goal: Minimize Mean Time To Recovery (MTTR)
**Strategic Intent**: In a distributed mesh, "root cause analysis" is the bottleneck. We optimize for **Observability** to reduce the time from *incident* to *diagnosis*.

### Objective 2.1: Distributed Context Propagation
**Architectural Decision**: Implement OpenTelemetry (OTEL) as a first-class citizen in the infrastructure layer.
**Impact**:
- **End-to-End Visibility**: Every request is tagged with a trace ID at the **API Gateway**. This context propagates through RabbitMQ headers and into the `Trip Service` and `Driver Service`.
- **Latency Budgeting**: Using **Jaeger**, we can visualize exactly which span (DB query vs. external OSRM API call) is consuming the latency budget, allowing effectively targeted performance engineering.

### Objective 3. Goal: Enforce Domain Purity & Testability
**Strategic Intent**: Prevent "Spaghetti Code" and ensure that core business rules (Pricing, Matching) are inoculated against infrastructure churn (e.g., switching from MongoDB to Cassandra).

### Objective 3.1: Hexagonal Architecture (Ports & Adapters)
**Architectural Decision**: Strict adherence to Clean Architecture principles.
**Impact**:
- **Dependency Inversion**: The `internal/domain` package has zero external dependencies. The Database is a plugin. The Message Broker is a plugin.
- **Testing Confidence**: We test business logic with 100% coverage using mock adapters. This allows us to refactor complex pricing algorithms safely without spinning up a Docker container.

---

## 4. Goal: Optimize Developer Velocity & Production Parity
**Strategic Intent**: The biggest cost in engineering is idle time. We aim to bring the feedback loop of "code change -> running app" down to < 5 seconds.

### Objective 4.1: Ephemeral Environments
**Architectural Decision**: Kubernetes-native local development using **Tilt**.
**Impact**:
- **Hot Reloading**: Changed a Go struct? **Tilt** syncs the binary into the running container infrastructure immediately. Failed a compiled check? You know in 2 seconds, not after a 10-minute CI pipeline.
- **Production Parity**: Developers run the exact same `deployment.yaml` manifests locally as we run in prod. "It works on my machine" is no longer a valid excuse or a possible bug class.

---

## 5. Goal: Ensure Financial Integrity & Security
**Strategic Intent**: In a transaction-heavy system, data consistency is paramount. We favor **Eventual Consistency** with strong audit barriers.

### Objective 5.1: Idempotent Payment Processing
**Architectural Decision**: Webhook-driven payment finalization with cryptographic verification.
**Impact**:
- **Fraud Prevention**: We do not trust the client state. We trust the Stripe backend signature.
- **Double-Spend Protection**: The `Payment Service` handles duplicate webhooks (at-least-once delivery) by checking transaction state before mutation. This ensures that network retries do not result in double charges.

### Objective 5.2: State Synchronization via Saga Patterns
**Architectural Decision**: Choreography-based Saga for distributed transactions.
**Impact**:
- **Atomic Operations**: A trip is only "Confirmed" when the Saga completes (Trip Created + Driver Locked + Payment Auth). If any step fails, compensating events roll back the state system-wide, preventing "Ghost Trips" (charged but no driver).
