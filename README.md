# 🛡️ etcd-shield

**etcd-shield** is a zero-trust gRPC security proxy and authorization firewall designed to protect **etcd** clusters from unauthorized administrative and maintenance API calls.

It sits in front of etcd using an **Envoy Proxy** frontend coupled with a custom Go-based **External Authorization (`ext_authz`) microservice**.

## 🚀 Features

* **API-Level Access Control**
  Allows safe Key-Value operations such as standard `Put` and `Range` requests while intercepting and blocking dangerous maintenance and administrative commands with a `PermissionDenied` response.

* **Sidecar Architecture**
  Deployable as a lightweight sidecar container inside Kubernetes pods alongside etcd or applications that require protected etcd access.

* **Envoy Powered**
  Uses Envoy Proxy for high-performance gRPC traffic routing and authorization filtering.

* **Zero-Trust Authorization**
  Requests are evaluated by the authorization firewall before they are forwarded to the etcd backend.

## 🏗️ Architecture

```text
                 gRPC Client
                      │
                      ▼
            ┌──────────────────┐
            │   Envoy Proxy    │
            │     :9090        │
            └────────┬─────────┘
                     │
                     │ ext_authz
                     ▼
            ┌──────────────────┐
            │  etcd-shield     │
            │  AuthZ Firewall  │
            │     :50051        │
            └────────┬─────────┘
                     │
              ┌──────┴──────┐
              │             │
          ALLOW │             │ DENY
              │             │
              ▼             ▼
        ┌──────────┐   PermissionDenied
        │   etcd   │
        │  :2379   │
        └──────────┘
```

The authorization layer determines whether an incoming gRPC request is permitted. Safe Key-Value operations can be forwarded to etcd, while restricted administrative or maintenance APIs are denied.

## 📂 Repository Structure

```text
etcd-shield/
├── cmd/
│   └── ...
├── shield/
│   └── ...
├── deploy/
│   └── docker/
│       ├── docker-compose.yml
│       └── envoy.yaml
├── go.mod
└── README.md
```

### Directory Overview

| Directory        | Description                                              |
| ---------------- | -------------------------------------------------------- |
| `cmd/`           | Application entrypoints                                  |
| `shield/`        | Core Go-based authorization firewall                     |
| `deploy/docker/` | Docker Compose and Envoy configuration for local testing |

## 🛠️ Quickstart

### Prerequisites

Make sure you have the following installed:

* Docker
* Docker Compose
* Git

### 1. Clone the repository

```bash
git clone https://github.com/atharvaramawat/etcd-shield.git
cd etcd-shield
```

### 2. Start the stack

```bash
docker compose up --build
```

This starts the local etcd, Envoy Proxy, and `etcd-shield` authorization service.

### 3. Test the proxy

The Envoy frontend is exposed on:

```text
localhost:9090
```

Requests sent through Envoy are evaluated by `etcd-shield` before reaching the etcd backend.

You can use the endpoint to verify that:

* Allowed Key-Value operations are forwarded to etcd.
* Restricted administrative and maintenance operations are rejected.
* Unauthorized requests receive a gRPC `PermissionDenied` response.

## ☸️ Kubernetes Deployment

`etcd-shield` can be deployed using the **sidecar pattern**, placing Envoy and the authorization firewall alongside the etcd container.

A simplified deployment looks like this:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: secured-etcd
spec:
  replicas: 1

  selector:
    matchLabels:
      app: etcd-shield-cluster

  template:
    metadata:
      labels:
        app: etcd-shield-cluster

    spec:
      containers:

        # 1. etcd database
        - name: etcd
          image: quay.io/coreos/etcd:v3.5.9
          ports:
            - containerPort: 2379

        # 2. Envoy Proxy frontend
        - name: envoy-shield-proxy
          image: envoyproxy/envoy:v1.25.0
          ports:
            - containerPort: 9090

        # 3. etcd-shield authorization firewall
        - name: shield-authz
          image: atharvaramawat/etcd-shield:v1.0.0
          ports:
            - containerPort: 50051
```

> **Note:** The Envoy image tag should match the version used by the Envoy configuration shipped with this project. Update the tag if your deployment uses a different Envoy release.

## 🔐 Security Model

The request flow is:

```text
Client
  │
  ▼
Envoy
  │
  ├──► ext_authz ──► etcd-shield
  │                    │
  │                    ├── ALLOW ──► Envoy ──► etcd
  │                    │
  │                    └── DENY ───► PermissionDenied
  │
  ▼
etcd
```

This prevents clients from directly accessing sensitive etcd APIs through the protected endpoint.

The proxy can enforce policies at the **gRPC service and method level**, allowing security rules to distinguish between normal data operations and privileged administrative operations.

## 🧪 Testing

After starting the Docker Compose stack, test requests against:

```text
localhost:9090
```

The expected behavior is:

| Operation                     | Expected Result |
| ----------------------------- | --------------- |
| `Put`                         | ✅ Allowed       |
| `Range`                       | ✅ Allowed       |
| Restricted maintenance API    | ❌ Denied        |
| Restricted administrative API | ❌ Denied        |

## 🐳 Docker Compose

The local development environment is provided under:

```text
deploy/docker/
```

Start the complete stack with:

```bash
docker compose up --build
```

Stop the stack with:

```bash
docker compose down
```

## 🧰 Technology Stack

* **Go** — Authorization firewall
* **Envoy Proxy** — gRPC proxy and request filtering
* **gRPC** — Client-to-proxy communication
* **etcd** — Key-Value database backend
* **Docker / Docker Compose** — Local deployment
* **Kubernetes** — Sidecar deployment

## 🎯 Use Cases

`etcd-shield` is useful in environments where applications need access to etcd data but should not have unrestricted access to its administrative APIs.

Potential use cases include:

* Protecting etcd in Kubernetes environments
* Restricting application access to etcd
* Preventing accidental administrative API calls
* Adding an authorization layer in front of legacy etcd clients
* Enforcing organization-specific gRPC authorization policies

## ⚠️ Project Status

This project is intended as a security proxy and experimental authorization layer for etcd.

Before using it in production, thoroughly review and test the authorization policies, Envoy configuration, TLS configuration, and container security settings for your environment.

## 📄 License

See the repository's license file for licensing information.

## 🔗 Repository

**GitHub:**
https://github.com/atharvaramawat/etcd-shield
