# ADR 0004: No Business Optimization On-Chain

**Status:** Accepted

**Decision:** OpenFinality handles coordination and settlement proof generation only. It does not perform business-level liquidity optimization.

**Alternatives:** Put the optimization logic in chaincode.

**Consequences:** Upstream systems are responsible for submitting optimized settlement instructions to OpenFinality.
