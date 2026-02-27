# DD-02: Codegen Strategy — Adapt Stigmer's Pipeline

**Date**: 2026-02-26
**Status**: PROPOSED

## Context

Stigmer has a two-stage codegen pipeline (proto → schema → Go) that generates typed MCP input structs with `ToProto()` methods. Should we build the same for Planton?

## Options

### Option A: Full codegen pipeline from day one
- Pros: Future-proof, consistent with Stigmer, automated maintenance
- Cons: Significant upfront investment, only 2 tools initially, may be over-engineering

### Option B: Hand-written types now, codegen later
- Pros: Faster to ship, only 2 simple input types needed, iterate on design first
- Cons: Tech debt if not addressed, diverges from Stigmer pattern initially

### Option C: Adapt Stigmer's codegen directly
- Pros: Proven system, minimal design work
- Cons: Planton's CloudResource wrapping differs from Stigmer's flat resources

## Decision

**Option B initially, with Option C as follow-up.**

Start with hand-written `CloudResourceApplyInput` and `DeleteCloudResourceInput` types. Once we add more domains (resourcemanager, servicehub, connect), build the codegen pipeline adapted from Stigmer's.

## Consequences

- Fast initial delivery
- Must plan for codegen migration when adding 3rd+ domain
- Hand-written types serve as reference for what codegen should produce
