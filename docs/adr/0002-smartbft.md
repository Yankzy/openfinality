# ADR 0002: SmartBFT

**Status:** Accepted

**Context:** Financial networks require crash and Byzantine fault tolerance.

**Decision:** Use SmartBFT for the ordering service.

**Alternatives:** Raft (crash fault tolerance only).

**Consequences:** Tolerates up to f Byzantine faults with n=3f+1. Requires minimum 4 orderers.
