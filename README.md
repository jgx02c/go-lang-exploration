# 🧠 Pan Pacific Internship Prep Plan (2 Weeks)

### 🗓️ Duration: 2 Weeks

### 🏢 Focus Areas: Internal Tools, Security, Go + Echo Backend, GitHub Actions, Cloud Deployments, Terraform

---

## ✅ Week 1: Core Stack & Systems

### 1. **Go (Golang) + Echo Framework**

* [ ] Complete [A Tour of Go](https://tour.golang.org/)
* [ ] Learn Echo basics: routing, middleware, handlers
  🔗 [Echo Quickstart](https://echo.labstack.com/guide)
* [ ] Practice JSON parsing and building REST endpoints
* [ ] Write a simple Echo server (CRUD + middleware)

---

### 2. **GitHub Actions + Git Workflow**

* [ ] Learn GitHub Actions structure: `jobs`, `steps`, `runs-on`
* [ ] Understand workflow triggers (push, PR, schedule)
* [ ] Use GitHub Secrets in a workflow
* [ ] Create a basic workflow for:

  * Linting a Go project
  * Running `go test ./...`
    🔗 [Actions Docs](https://docs.github.com/en/actions)

---

### 3. **Caching & Validation**

* [ ] Understand when and how to use Redis
* [ ] Learn Redis TTL and basic caching patterns
* [ ] Build a Go service that:

  * Caches a value (e.g., ad-rate)
  * Validates it via mock endpoint fallback
* [ ] Use a Go Redis client like `go-redis`

---

## ✅ Week 2: Cloud Infra + Security + Observability

### 4. **Terraform + Cloud Automation**

* [ ] Learn basics: resources, variables, modules, state
* [ ] Create a simple Terraform script to:

  * Provision a VM or S3 bucket (can use `localstack` or mock provider)
    🔗 [Terraform Getting Started](https://developer.hashicorp.com/terraform/tutorials)
* [ ] Understand what "Versioned Cloud Control" could imply:

  * Versioning infrastructure code
  * Managing resource drift

---

### 5. **Security + Internal Tooling Concepts**

* [ ] Study RBAC patterns and access control
* [ ] Understand JWT, OAuth2, session-based auth
* [ ] Review best practices for secrets handling
* [ ] Explore security concerns in internal tools (e.g., audit trails, input validation)

---

### 6. **Logging & Monitoring**

* [ ] Learn structured logging in Go using:

  * `zap`, `logrus`, or built-in logging
* [ ] Review Azure Monitor basics
  🔗 [Azure Monitor Docs](https://learn.microsoft.com/en-us/azure/azure-monitor/)
* [ ] Add logging & request tracking middleware to Echo project

---

## 🧩 Bonus (Optional)

### LLMs + Data Warehouse (Context Prep)

* [ ] Skim overview of LLM use in procurement (RAG, document processing)
* [ ] Light intro to Data Warehouse systems: Redshift, BigQuery, or Snowflake





## Access the docker w/ postgres

docker exec -it echo_postgres psql -U postgres -d echo_api

## Exit postgres Shell

\q


## Run the docker

docker-compose up -d
