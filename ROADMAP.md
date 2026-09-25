# OpenFinality Roadmap

OpenFinality is an open-source settlement coordination layer for financial institutions.

The project is designed to evolve incrementally from a reproducible multi-organization Hyperledger Fabric network into production-grade infrastructure capable of coordinating settlement across banks, payment service providers, mobile-money systems, clearing institutions, central banks, and other regulated financial networks.

The roadmap deliberately separates protocol maturity from commercial adoption.

OpenFinality does not issue money, determine exchange rates, optimize financial markets, or replace existing settlement systems. Its role is to provide a neutral coordination layer between independent institutions and the settlement systems they already use.

---

# Roadmap Principles

OpenFinality development follows several principles.

## Protocol Before Products

The core protocol must remain useful independently of any particular application, financial product, or commercial operator.

## Settlement Coordination, Not Money Movement

OpenFinality coordinates settlement state, authorization, reservations, execution instructions, confirmations, and proofs.

Actual sovereign or regulated money continues to move through external settlement systems.

## Institution-Level Trust

No participant should be required to trust another participant's private database.

Settlement-critical state must be independently verifiable by authorized organizations.

## Incremental Regulatory Integration

Real financial integrations should begin in controlled sandbox environments before progressing to production settlement.

## Fault Tolerance by Design

Institutional settlement infrastructure must remain available through infrastructure failures and must fail safely when consensus cannot be reached.

## Interoperability Over Exclusivity

OpenFinality should integrate with existing financial systems rather than requiring institutions to replace them.

---

# v0.1: Working Settlement Protocol

**Status: Developer Preview**

The objective of v0.1 is to prove that multiple independent organizations can coordinate settlement through a shared, fault-tolerant ledger without trusting one operator.

The release establishes the minimum viable OpenFinality protocol.

## Network

- Three independent participant organizations
- Four-node SmartBFT ordering cluster
- Fabric CA institutional identity
- CouchDB-backed peer state
- Shared settlement channel
- Persistent ledger state
- Byzantine fault tolerance demonstration
- Recovery after node restart

## Settlement Protocol

- Institutional participant identities
- Settlement intents
- Counterparty acceptance and rejection
- Settlement reservations
- Settlement instructions
- Multi-party authorization
- Affected-party state-based endorsement
- Settlement confirmations
- Settlement proofs
- Deterministic settlement state machine
- Idempotent mutation semantics
- Replay protection

## External Settlement

- Settlement adapter interface
- Deterministic simulated settlement adapter
- Simulated institutional balances
- Settlement execution
- External confirmation ingestion
- Final proof commitment

## Application Interfaces

- REST API
- Go Fabric Gateway integration
- Command-line client
- OpenAPI specification
- Structured error model
- Health and readiness endpoints

## Operations

- Local Docker Compose deployment
- One-command bootstrap
- Network health verification
- Automated end-to-end demonstration
- SmartBFT failure demonstration
- Persistent-state restart demonstration
- Structured logging
- Initial metrics

## Security

- MSP-based institutional authorization
- Fabric CA role attributes
- Counterparty authorization
- Settlement-adapter authorization
- State-based endorsement
- Duplicate protection
- Fixed-point monetary representation
- Threat model
- Private data architecture

## Open Source

- Apache 2.0 license
- Contribution guidelines
- Governance model
- Maintainer model
- RFC process
- Security disclosure policy
- Architecture Decision Records
- CI and security scanning
- Release process

### v0.1 Exit Condition

OpenFinality v0.1 is complete when a developer can clone the repository, bootstrap the network, execute an institutional settlement, retrieve its final proof, disable one SmartBFT orderer, and successfully execute another settlement through the normal application API.

---

# v0.2: Institutional Network Foundation

The objective of v0.2 is to move from a fixed demonstration topology to infrastructure that can support independently operated organizations.

The key question becomes:

> Can a new institution join OpenFinality without rebuilding or manually reconfiguring the network?

## Dynamic Organization Onboarding

Introduce a formal lifecycle for adding, suspending, updating, and removing institutional participants.

This includes:

- organization registration
- MSP onboarding
- certificate authority integration
- peer enrollment
- channel participation
- endorsement-policy updates
- organization metadata
- participant lifecycle state
- controlled offboarding

Onboarding procedures should be reproducible and auditable.

## Adapter SDK

Formalize the external settlement adapter interface into a supported integration SDK.

Adapters should expose common operations such as:

- capability discovery
- reserve
- release
- execute
- query status
- cancel where supported
- reconcile

The SDK should define consistent:

- request schemas
- response schemas
- error semantics
- idempotency behavior
- timeout behavior
- retry behavior
- settlement finality semantics

Reference adapters should remain simulated unless integrated under an appropriate regulated environment.

## Institutional Privacy

Expand Fabric private-data usage to support more realistic institutional requirements.

Areas include:

- confidential account references
- bilateral settlement metadata
- institution-specific limits
- private reservation details
- restricted settlement instructions
- selective disclosure
- hashed public commitments

The shared network should expose only the minimum information required for consensus and auditability.

## Federated Policy Configuration

Introduce protocol-level mechanisms for institutions and network governance to publish and version operational policies.

Potential policy domains include:

- permitted participant types
- transaction limits
- required endorsements
- supported settlement rails
- reservation requirements
- jurisdictional restrictions
- settlement deadlines
- emergency restrictions

Policy changes must be signed, versioned, auditable, and reproducible.

## Disaster Recovery

Move from development-grade recovery toward documented institutional recovery procedures.

Include:

- peer recovery
- orderer recovery
- CA recovery
- ledger snapshot restoration
- CouchDB restoration
- cryptographic material restoration
- configuration backup
- recovery drills
- recovery-point objectives
- recovery-time objectives

## Observability

Expand telemetry into an operational monitoring model.

Include:

- settlement throughput
- commit latency
- endorsement latency
- adapter latency
- orderer health
- peer health
- channel state
- pending settlements
- failed settlements
- reservation utilization
- settlement-finality time
- error rates

The system should expose sufficient telemetry for independent institutional operators.

### v0.2 Exit Condition

A new institution should be able to join an existing OpenFinality network using documented onboarding procedures, operate its own infrastructure, integrate a settlement adapter, participate in settlement, and independently monitor its node without requiring repository-level code changes.

---

# v0.3: Regulated Settlement Pilot

The objective of v0.3 is to move OpenFinality from simulated settlement into controlled interaction with real financial infrastructure.

This phase must occur within appropriate legal, regulatory, and institutional environments.

## Regulated Sandbox Adapter

Implement the first adapter connected to an actual financial system in a controlled environment.

Potential initial integrations may include:

- bank sandbox APIs
- payment-provider sandbox APIs
- test RTGS environments
- regulated payment sandboxes
- central-bank experimentation environments

No production customer money should be routed until the participating institutions and relevant authorities approve the deployment.

## Shadow Settlement

Before OpenFinality controls real settlement execution, support shadow operation.

In shadow mode:

1. Real transactions continue through existing infrastructure.
2. OpenFinality receives permitted representations of those transactions.
3. OpenFinality independently models the settlement lifecycle.
4. Results are compared against actual institutional settlement.

Shadow mode allows measurement of:

- state consistency
- reconciliation accuracy
- settlement timing
- failure handling
- operational reliability
- protocol correctness

without placing customer funds at risk.

## Multi-Rail Settlement Interface

Extend the adapter architecture so a settlement instruction can identify and coordinate with multiple classes of settlement infrastructure.

The interface should support differences in:

- finality models
- settlement windows
- synchronous versus asynchronous execution
- reservation semantics
- cancellation behavior
- retry semantics
- confirmation formats

OpenFinality itself should remain neutral regarding which rail an upstream application chooses.

## Institutional Conformance Suite

Create a formal test suite that institutions and adapter developers can run before joining a network.

The suite should validate:

- identity behavior
- protocol compatibility
- settlement state transitions
- idempotency
- authorization
- endorsement
- reservation behavior
- failure handling
- adapter conformance
- confirmation integrity
- settlement-proof generation

Conformance results should be machine-readable and reproducible.

### v0.3 Exit Condition

At least one regulated or institutional sandbox should be capable of participating through a real OpenFinality adapter, while the same protocol remains interoperable with the existing simulated environment and conformance suite.

---

# v0.4: Multi-Institution Settlement Network

The objective of v0.4 is to prepare OpenFinality for sustained operation across independently administered financial institutions.

## Network Governance

Formalize governance for shared infrastructure, including:

- membership decisions
- protocol upgrades
- emergency actions
- channel governance
- certificate revocation
- participant suspension
- dispute procedures
- protocol voting
- maintainer authority
- security response

No single commercial participant should be able to unilaterally alter settlement-critical protocol behavior.

## Production Identity Architecture

Support institutional certificate lifecycle requirements including:

- hardware-backed keys
- HSM integration
- certificate rotation
- automated renewal
- emergency revocation
- key compromise recovery
- separation of administrative and transaction identities

## Production Deployment Profiles

Provide deployment guidance and supported profiles for:

- institutional data centers
- private cloud
- public cloud
- hybrid deployments
- geographically distributed peers
- geographically distributed orderers

Docker Compose remains a developer environment, not the production architecture.

## Protocol Upgrade Framework

Define safe mechanisms for:

- chaincode upgrades
- protocol-version negotiation
- schema migration
- rolling infrastructure upgrades
- backward compatibility
- deprecation
- network capability upgrades

## Settlement Reconciliation

Introduce formal reconciliation interfaces between:

- OpenFinality ledger state
- external settlement system state
- institutional accounting systems

Reconciliation discrepancies must become explicit protocol events rather than silent operational errors.

### v0.4 Exit Condition

Independent organizations should be able to operate OpenFinality continuously with defined governance, security, deployment, upgrade, reconciliation, and incident-management procedures.

---

# v1.0: Production-Grade Open Settlement Fabric

OpenFinality v1.0 represents the point at which the protocol is considered suitable for production evaluation by regulated institutions.

Production use remains subject to the laws, licenses, approvals, and risk requirements of each deployment.

The v1.0 objective is not universal adoption.

It is protocol maturity.

## Required Characteristics

OpenFinality should provide:

- independently operated institutional nodes
- Byzantine fault-tolerant ordering
- mature identity lifecycle
- hardware-backed key support
- dynamic membership
- settlement adapter standard
- deterministic protocol semantics
- private institutional data
- policy federation
- complete settlement auditability
- robust disaster recovery
- protocol-version negotiation
- institutional conformance testing
- reconciliation interfaces
- operational observability
- secure upgrade procedures
- external security audits
- documented production threat model

The protocol should be deployable without dependency on any specific commercial application.

---

# Future Integration Tracks

After the core protocol reaches sufficient maturity, development can expand through independent integration tracks.

These are not assumptions about which financial systems will adopt OpenFinality. They are technical interoperability targets.

## RTGS Adapters

Develop standardized connectors for real-time gross settlement systems.

The adapter specification must account for:

- settlement finality
- operating windows
- queued payments
- liquidity requirements
- acknowledgement formats
- failure semantics
- reconciliation

---

## Bank Adapters

Allow participating banks to connect internal settlement infrastructure to OpenFinality.

Possible integration surfaces include:

- core banking systems
- treasury systems
- internal ledgers
- correspondent banking systems
- settlement accounts

---

## PSP Adapters

Support payment service providers and payment switches through standardized settlement interfaces.

This can include:

- payment processors
- mobile-money operators
- wallet providers
- merchant acquirers
- national payment switches

---

## Central-Bank Integration

Explore integration models in cooperation with monetary authorities.

Possible roles include:

- network observer
- participant authorization
- settlement operator
- policy authority
- RTGS adapter operator
- institutional identity authority
- emergency network governance participant

OpenFinality does not require a central bank to operate the protocol, but the architecture should support direct central-bank participation where appropriate.

---

## Settlement Interoperability Standards

Develop open specifications that allow implementations beyond the reference OpenFinality software to interoperate.

Potential standards include:

- settlement-intent schema
- reservation protocol
- settlement-instruction schema
- confirmation schema
- settlement-proof format
- participant identity model
- adapter interface
- error taxonomy
- finality semantics

The long-term goal is for OpenFinality to become more than a single codebase.

It should become an **open settlement protocol** that independent software implementations can speak.

---

# Long-Term Direction

The long-term goal of OpenFinality is not to replace existing financial infrastructure.

It is to make existing infrastructure interoperable.

Banks should continue operating their ledgers.

Central banks should continue operating sovereign settlement systems.

Payment providers should continue operating their networks.

OpenFinality should provide the neutral coordination layer between them.

The desired architecture is:

```text
        Financial Applications
                 │
                 │
            OPENFINALITY
     Open Settlement Protocol
                 │
     ┌───────────┼───────────┐
     │           │           │
    RTGS       Banks        PSPs
     │           │           │
     └───────────┼───────────┘
                 │
        Existing Financial
        Settlement Systems
````

The project succeeds when independent institutions can coordinate settlement securely without being forced onto one company's proprietary financial network.

---

## North Star

> **Any authorized financial institution should be able to connect to OpenFinality, coordinate settlement with another institution, execute through an appropriate external rail, and independently verify the resulting finality.**

That is the protocol OpenFinality is building toward.
