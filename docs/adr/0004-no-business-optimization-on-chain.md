# ADR 0004: No Business Optimization On-Chain

**Status:** Accepted

**Context:** Complex liquidity optimization and netting (e.g. Toro/X-Unit) can require huge data inputs and frequent state updates.

**Decision:** Afro-Rail handles coordination and settlement proof generation only. It does not perform business-level liquidity optimization.

**Alternatives:** Put the optimization logic in chaincode.

**Consequences:** Upstream systems are responsible for submitting optimized settlement instructions to Afro-Rail.
