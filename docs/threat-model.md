# Threat Model: Config Tampering

dlinter-go's levers (see the [discipline harness section](../README.md#dlinter-go-as-a-discipline-harness) of the README) are machine-checked, but a machine-checked rule only constrains code as long as the rule itself stays in place. This document states the threat plainly, explains why the mechanism is deferred, and records what a consumer should do about it.

## 1. The threat

An AI agent — or any contributor under time pressure — can satisfy a rule by removing it instead of fixing the code it flagged:

- Editing `.golangci.yml` to widen a threshold (`gocognit.min-complexity: 15` → `50`).
- Deleting a linter from `linters.enable`.
- Adding a `//nolint` directive to silence a specific finding.
- Excluding a path via `linters.exclusions.rules` that was not excluded before.

Every one of these is a one-line change, cheaper than the refactor the rule was trying to force. An agent optimizing for "make the lint pass" rather than "make the code good" will find this path if nothing pushes back.

## 2. Why the linter cannot fix this

A linter cannot police the config that configures it. Any check we add — a meta-rule that flags threshold changes, a hash-pinned config file, a self-referential analyzer — lives in the same mutable tree as the config it would be checking, and is editable by exactly the same actor with exactly the same one-line change. Shipping a partial mechanism here would imply a guarantee dlinter-go cannot actually make. We would rather state the limitation than ship decoration.

This is the same reasoning already applied to the package-cohesion analyzer's deferral: a mechanism that can be trivially defeated does not belong in a project whose entire premise is gaming resistance.

## 3. What dlinter-go ships today

- **`nolintlint`**, enabled in both `recommended.golangci.yml` and this repo's self-applied `.golangci.yml`, with `require-explanation: true`, `require-specific: true`, and `allow-unused: false`. It cannot stop a `//nolint` directive from being added, but it forces every one to name a specific linter and explain itself in English — silent suppression becomes a reviewable sentence in the diff instead of a bare comment or nothing at all.
- Every threshold lives in a single, small, reviewable file (`.golangci.yml`), not scattered across build scripts or CI YAML. A reviewer checking "did anyone weaken the cage" has exactly one file and a handful of lines to look at.

## 4. Governance baseline the consumer owns

The credible controls here are external to dlinter-go and are the consumer's repository settings to configure, not something this project can enforce from inside the linter:

- **CODEOWNERS** on `.golangci.yml` (and `.custom-gcl.yml`), so a threshold change requires sign-off from someone accountable for the harness, not just whoever's PR happens to touch it.
- **Branch protection** requiring review before merge, applied to the same files.
- **A CI diff check** that fails a PR when it detects a threshold weakening (a `min-complexity`/`lines`/`max` value increasing) or a new `//nolint` directive without a corresponding, human-written justification in the PR description. This is not something dlinter-go ships — it is a project-specific CI job the consumer writes for their own repository, informed by whatever review discipline they already run.

None of these live inside dlinter-go's binary. That is intentional: they are decisions about human process, and human process is where this kind of guarantee actually has to live.

## 5. Non-guarantee statement

**dlinter-go raises the cost and visibility of cheating. It does not make cheating impossible.** Every lever in this harness can still be defeated by someone willing to edit the config and get that edit past review. What the harness buys is that defeating it is no longer free or silent — it is a small, legible diff that a CODEOWNERS-gated review can catch, instead of an invisible act embedded in a thousand-line feature PR.
