# Banking Core Mock Server

An open-source mock server for our banking core, built on top of [TigerBeetle](https://docs.tigerbeetle.com/operating/docker/) and designed for fintech startups. This project provides a sandbox environment to simulate banking operations like account management, transfers, and KYC for testing integrations.

## Overview

The goal of this project is to provide a robust, open-source banking core that:
- Leverages TigerBeetle as the high-performance ledger.
- Exposes a REST (or GraphQL) API for fintech partners.
- Supports modules for authentication, entity onboarding, account management, loans, transfers, and more.
- Can eventually be deployed to GCP after local development and testing.

## Roadmap / To-Do List

We use this README as a living to-do list to help guide development and encourage community contributions:

- [ ] **Authentication Module**
  - [ ] Decide between OAuth 2.0 or API key-based auth.
  - [ ] Implement authentication middleware.
  - [ ] Create endpoints for token issuance and validation.
  - [ ] Write unit and integration tests for auth workflows.

- [ ] **Entity Onboarding & KYC**
  - [ ] Design API endpoints for person and business entities.
  - [ ] Integrate KYC fields and document upload flows.
  - [ ] Develop validation and verification mechanisms.
  - [ ] Write tests for onboarding scenarios.

- [ ] **Account Management**
  - [ ] Build CRUD endpoints for bank accounts.
  - [ ] Separate account metadata from account numbers if needed.
  - [ ] Integrate with TigerBeetle’s account operations.
  - [ ] Document the API and add test cases.

- [ ] **Loan Processing**
  - [ ] Create endpoints for loan creation, disbursement, and payment.
  - [ ] Implement business rules (interest, term, collateral, etc.).
  - [ ] Ensure robust error handling and transaction logging.

- [ ] **Transfers Module**
  - [ ] Develop endpoints for ACH, wire, and realtime transfers.
  - [ ] Create multi-step approval workflows for high-value transactions.
  - [ ] Implement cancellation and reversal logic.
  - [ ] Test edge cases such as partial failures and timeouts.

- [ ] **Documents & Reporting**
  - [ ] Provide endpoints for document uploads and secure retrieval.
  - [ ] Build reporting endpoints for compliance and reconciliation.
  - [ ] Generate scheduled and on-demand reports.

- [ ] **Webhooks & Event Logging**
  - [ ] Enable webhook subscriptions for key events (e.g., transaction completed, account created).
  - [ ] Create a robust audit log accessible via an API endpoint.
  - [ ] Document event types and webhook payload formats.

- [ ] **Sandbox & Simulation**
  - [ ] Develop a sandbox mode that mimics production behavior.
  - [ ] Provide test credentials and simulation scenarios for partners.
  - [ ] Implement logging and debugging support in the sandbox.

- [ ] **Deployment & Scalability**
  - [ ] Containerize the application for consistent deployments.
  - [ ] Prepare deployment scripts for GCP (Compute Engine / GKE).
  - [ ] Set up CI/CD pipelines for automated testing and deployment.
  - [ ] Monitor performance and scalability metrics.

- [ ] **Miscellaneous**
  - [ ] Create comprehensive documentation (INSTALL.md, CONTRIBUTING.md).
  - [ ] Establish coding standards and set up code quality tools.
  - [ ] Engage with the community via GitHub Issues and discussions.

## Contributing

We welcome contributions from everyone! To contribute:
1. **Fork the Repository** and clone it locally.
2. Create a branch for your feature or fix.
3. Make your changes and add tests.
4. Submit a pull request with a clear description of your changes.
5. Check out our [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## Getting Started Locally

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-org/banking-core-mock-server.git
   cd banking-core-mock-server ```

2. **Set up TigerBeetle locally:**
Ensure you have TigerBeetle running as described in [their Installation Guide](https://docs.tigerbeetle.com/quick-start/). For local development, you can run it directly as a binary:
bash
```chmod +x ./tigerbeetle
mkdir -p data
./tigerbeetle format --cluster=0 --replica=0 --replica-count=1 ./data/0_0.tigerbeetle
./tigerbeetle start --addresses=0.0.0.0:3000 ./data/0_0.tigerbeetle```