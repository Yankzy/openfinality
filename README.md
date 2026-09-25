# Afro-Rail

> Open infrastructure for programmable multi-institution financial settlement.

**License**: Apache-2.0
**Initial Maturity**: Experimental / Developer Preview

> **IMPORTANT**: Afro-Rail is experimental financial infrastructure. The reference implementation uses simulated funds and must not be used to route real customer money.

## Architecture

```text
                Applications / Markets / Optimizers
                            │
                            │
                     X-Unit / Others
                            │
                            │
                       AFRO-RAIL
              open institutional coordination
                            │
               ┌────────────┼────────────┐
               │            │            │
              RTGS       Bank APIs    PSP Rails
               │            │            │
                            ...
```

Broader motivating architecture:

```text
                  TORO
        Intelligence and optimization
                    │
                 X-UNIT
       Common liquidity state space
                    │
                AFRO-RAIL
       Open Settlement Fabric
                    │
          REAL SETTLEMENT SYSTEMS
```

> Toro and X-Unit are examples of systems that may consume Afro-Rail. They are not dependencies of the Afro-Rail protocol.
