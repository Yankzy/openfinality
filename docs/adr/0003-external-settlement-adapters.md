# ADR 0003: External Settlement Adapters

**Status:** Accepted

**Context:** Afro-Rail coordinates settlement but does not itself issue or move sovereign fiat currency directly within its chaincode.

**Decision:** Use external settlement adapters to interact with existing real-world payment rails, RTGS systems, APIs, or mobile money operators.

**Alternatives:** Native on-chain tokens for all assets.

**Consequences:** Clean architectural boundary, but requires careful idempotency and state synchronization.
