## The Natakos Microservice Engineering Retrospective (learn.md)

### 1. The Core Philosophy: The Hierarchy of Software Engineering

Software engineering is not about typing code blindly; it is an interconnected chain of structural hierarchies:

* Product Requirements dictate System Designs.
* System Designs anchor Data Contracts.
* Data Contracts inform Function Definitions.
* Coding is simply the final step of translating these structural blueprints into physical software adapters.

### 2. Microservice Decoupling & The Edge Webhook Boundary

A robust backend architecture operates like a modular lego chain. The Ingress Gatekeeper (Go Gateway) focuses strictly on structural security, token validation, and network protocol primitives (HTTP/REST).
By completely decoupling the edge ingress from core heavy processing using an asynchronous message broker (RabbitMQ Alpine), the edge gate can accept inbound WhatsApp webhooks and immediately return a lightning-fast HTTP 200 response to the client. The user never experiences bottleneck waiting delays, while downstream processing layers focus purely on the core business state logic.

### 3. Evolutionary Schema Design & RBAC Anchors

System design is fundamentally evolutionary. Instead of writing brittle code logic to handle edge actions, shifts must be caught and anchored at the database schema layer first. By pivoting our data layout early to natively support Role-Based Access Control (RBAC), we protected our internal runtime cluster from executing unauthenticated data tasks before any core backend code was ever invoked.

### 4. Dual-State Performance: Postgres Relational vs. Rolling Redis Caches

To balance high durability with high performance, state architectures must be split based on execution requirements:

* Persistent Storage (PostgreSQL Relational Tables) holds long-term transactional records such as user credentials, roles, and static room assignments.
* Volatile Session Caches (Optimized Redis Cache Layer) manage short-term session context limits.
By utilizing native Redis Lists with RPUSH and LRANGE operations combined with rolling size trimmings (LTRIM), we engineered a sliding-window memory management channel that feeds optimized, historical chat turn sequences to the AI without blowing past our model context tokens.

### 5. Structured Data vs. Unstructured Data (Function Calling vs. RAG)

Enterprise marketing overhypes complex standalone vector infrastructure. For 95% of operational business logic, a hybrid database like PostgreSQL with extensions like pgvector is bulletproof. Data processing must follow its native texture:

* Unstructured Information (House Rules) uses semantic lookup pipelines. Raw text strings are transformed into 1536-dimensional math arrays (embeddings) and matched using pgvector's optimized native cosine similarity distance calculations (<=> operator).
* Structured Information (Live Room Availability) must never rely on RAG, as language models will pull stale cached context data and hallucinate. Instead, we implement a two-pass Tool Calling (Function Calling) handshake. DeepSeek parses the intent, flags that it needs live status data, requests a structured JSON execution tool parameter block, and allows Python to query the real relational database values dynamically before synthesizing the final natural response.

### 6. Graceful Degradation & Global Failure Domain Protection

A production-ready cluster must expect third-party networks, external APIs, and hardware nodes to break. By enclosing our background AMQP transaction lifecycle inside a Global Failure Domain Defense Ring, any downstream syntax, database deadlock, or cloud credential error is caught gracefully. Instead of crashing the script and trapping the system in an infinite RabbitMQ loop, the circuit breaker triggers a Tier-2 hardcoded safe-mode fallback string, preserves the cache history state, and safely acknowledges the message token to clear the broker pipeline highway.

### 7. Absolute Observability: Structured JSON & Correlation IDs

Emojis and plain-text console logging do not scale. True cluster visibility requires transforming all operational statements into structured JSON logs injected with unique, trace-tracking Correlation IDs. From the exact microsecond the Go Gateway mints a request identification token, that ID travels inside the message metadata suitcase across every queue, caching layer, and AI model pass. Tracking down a multi-language distributed exception is reduced to a single search query, establishing complete engineering transparency.
