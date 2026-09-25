# ADR 0005: Fixed-Point Monetary Values

**Status:** Accepted

**Context:** Floating point precision issues cannot be tolerated in financial applications.

**Decision:** Use deterministic fixed-point monetary values using BigInt structures in all representations (chaincode, API, storage).

**Alternatives:** Float64 (rejected).

**Consequences:** Need a robust math library for safe arithmetic operations (checked add/sub/mul, etc).
