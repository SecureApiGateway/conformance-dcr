# DCR spec-version analysis (OPENIG-10765) — context record

## Decision (2026-09-28): do nothing

**No DCR v3.4 (or v3.3) adoption in this fork, no fix ported.** Reason: the team is focusing on FAPI, not Open Banking — this repo and its consumer (`openig-pyforge`'s DCR tests) are both Open Banking–specific, so upstream `conformance-dcr` changes are currently out of scope. This is a prioritization call, not a disagreement with the analysis below — see the decision comments on [OPENIG-10765](https://pingidentity.atlassian.net/browse/OPENIG-10765) for the record.

**Current state as a result:** this fork stays on `sbat-master` at `bd317c72` (or wherever it's since moved to for *unrelated* reasons), still targeting DCR spec v3.2 only. No `dcr34.go`, no DN-escaping fix, no `spec_version` change. `openig-pyforge`'s `registration_response_schema.py` (Core's Python re-implementation of this same schema, v3.2-only) is likewise untouched.

**Purpose of the rest of this file:** context for a human or an AI agent picking up work on this repo later, so the analysis behind the decision above doesn't have to be redone from scratch. If you are an AI agent asked to work on DCR version support, `spec_version`, or "sync with upstream" in this repo, **read this file first.**

- **Jira ticket:** [OPENIG-10765](https://pingidentity.atlassian.net/browse/OPENIG-10765) — "Analyze upstream conformance-dcr commits from 2026 for relevance to our FAPI DCR testing"
- **Analysis performed:** 2026-09-21 to 2026-09-23
- **This file is self-contained** — a more detailed task-by-task breakdown exists as local working notes from the analysis, but those aren't published alongside this repo, so everything needed to pick this back up is captured here instead.

## Exact reference points — use these, don't re-derive them from scratch

| What | Value |
|---|---|
| This fork's branch analyzed | `sbat-master` |
| **This fork's HEAD at analysis time** | `bd317c722372cfc5a4d325615097b0319c0560af` (`bd317c7`), dated 2024-12-09 — commit message: "1556: Fix error messages to be compatible with openig-fapi module (#66)" |
| Upstream repo analyzed | `https://github.com/OpenBankingUK/conformance-dcr` |
| Upstream branch analyzed | `develop` |
| **Last upstream commit actually reviewed** | `c11b78a` ("Update README.md"), dated 2026-05-20 |
| Upstream release covering the reviewed range | `v1.4.0`, tagged at commit `cc00a00`, dated 2026-05-12 (the `v1.4.0` release notes cover everything up to `cc00a00`; `c11b78a` is one untagged README commit on top of it that was also reviewed) |
| Number of upstream commits reviewed | 26 (all commits on `develop` between the fork's sync point and `c11b78a`) |

**If this gets revisited later:** don't re-diff from the beginning of history. Fetch upstream `develop`, and diff from `c11b78a` forward (`git log c11b78a..upstream/develop`) to see only what's genuinely new since this analysis, then combine with the findings below rather than repeating this work.

## What was analyzed

All 26 upstream commits since the fork's last sync point were grouped into 5 logical change units and individually reviewed via `git show`/diff (not just commit titles):

| # | Change | Verdict |
|---|---|---|
| 1 | DCR **v3.4** spec support (new manifest `dcr34.go` + schema validator `version34.go`) | Relevant — but only 3 concrete schema-rule changes, no new test scenarios (see below) |
| 2 | `disableKeepAlives` CLI option for mTLS clients | Not relevant (no known handshake/connection issue on record for us) |
| 3 | `dgrijalva/jwt-go` → `golang-jwt/jwt` v4 dependency swap | Not relevant — confirmed pure import swap, no functional change |
| 4 | `use_oid` option + RFC 2253 subject-DN escaping fix in `pkg/compliant/auth/signer.go` | **Relevant, high confidence** — see below, this is a real bug in this fork's current code |
| 5 | Go version / CI / linter / dependency bumps | Not relevant to test behavior |

### The DCR v3.4 schema delta, precisely

`dcr34.go`'s `NewDCR34` reuses the exact same 10 scenarios as `NewDCR32`/`NewDCR33` — **no new test scenarios**, only a different response-schema validator (`responseValidator34`). Verified by diffing `version33.go` against upstream's `version34.go`, the only 3 rule changes are:

1. `application_type`: allowed values change from `"web","mobile"` → `"web","native"` — **breaking rename**, not additive.
2. `grant_types`: gains `urn:ietf:params:oauth:grant-type:jwt-bearer` as newly allowed — purely additive.
3. `tls_client_auth_subject_dn`: max length raised `128` → `512`, **and** a new structural rule is added requiring the DN to start with `CN=` (comma-separated `attr=value` pairs) — v3.2/v3.3 never checked this.

For reference, **v3.3 is already fully implemented in this fork** (`dcr33.go`/`version33.go`, ported since 2020, well before this analysis) — its only delta from v3.2 is `grant_types` gaining `urn:openid:params:grant-type:ciba`, and `client_secret` max length tightening from `256` to `36`. Both are already correctly present in this fork's `version33.go`. **v3.4 is not a strict superset of v3.3** (it drops `"mobile"` and adds the DN pattern requirement, which v3.3 doesn't have) but in practice covers everything v3.3 needs, since this fork's own `signer.go` never sends `application_type: "mobile"` in the first place.

### The real bug found (item 4) — independent of the version question

This fork's current `pkg/compliant/auth/signer.go` builds the `tls_client_auth_subject_dn` claim like this:

```go
// Instead of potentially custom ASN/OID parsing to get exact, expected value of Subject DN
// we use a config entry
if s.transportSubjectDn != "" {
    claims["tls_client_auth_subject_dn"] = s.transportSubjectDn
} else {
    claims["tls_client_auth_subject_dn"] = s.transportCert.Subject.ToRDNSequence().String()
}
```

The code's own comment admits it avoids parsing the certificate correctly and instead depends on a **manually-supplied config override** (`transport_cert_subject_dn`, wired in from `templates/sbat-conformance-config.jq` / `templates/sbat-create-dcr-config.jq`) — falling back to Go's naive `ToRDNSequence().String()` if the override isn't set, which doesn't escape special characters per RFC 2253 and silently mishandles the `organizationIdentifier` OID (`2.5.4.97`).

Upstream's fix (commits `c30f30c`..`cc00a00`, PR #9 "fix-OID-issues") replaces this with a proper `subjectDN()` function that parses `RawSubject` directly, escapes per RFC 2253, and correctly maps `2.5.4.97` → `organizationIdentifier`. **This is a v3.2-applicable bug fix, unrelated to whether v3.4 is ever adopted** — it's exactly why this fork currently needs the manual `transport_cert_subject_dn` override at all.

Note the DN-escaping fix (item 4) is a legitimate, low-risk, version-independent bug fix — it's not being skipped because it's wrong or risky, purely because of the FAPI-over-OB prioritization stated in the Decision section at the top.

## If this is picked back up later: what would need to be done

If OB priority increases and this work is revisited, here is the concrete scope (full detail in the analysis files referenced above; this is the summary):

### 1. This fork (`SecureApiGateway/conformance-dcr`) — before doing anything, read this

**This fork has diverged substantially from upstream: 74 commits exist on `sbat-master` that don't exist anywhere upstream.** Verified via `git log origin/sbat-master ^upstream/develop ^upstream/master ^upstream/release/v1.3.1` (with `upstream` remote = `https://github.com/OpenBankingUK/conformance-dcr.git`).

- **~45 of those commits are CI/release infrastructure** (Codefresh, GitHub Actions, Docker/Artifact Registry, Slack) — permanent, SBAT-specific, irrelevant to DCR test logic, but **both sides have independently edited `.github/workflows/go.yml`, `.golangci.yml`, `Makefile`, and `Dockerfile`** since diverging, so a blind merge/rebase from upstream will hit conflicts in all four files for reasons that have nothing to do with the DCR changes.
- **~29 commits touch actual Go source** — this fork's real customization layer: `create_software_client_only` mode, config overrides, and extra JWT claims added in `pkg/compliant/auth/signer.go` (`client_id` on updates, `refresh_token` grant type, `authorization_signed_response_alg`) — **the same file** upstream's DN-escaping fix rewrites. There are also **3 test cases this fork added that don't exist upstream** and must be preserved: invalid-signature rejection, redirect_uri-not-in-software_redirect_uris rejection, `alg: none` rejection.

**Recommendation: do not `git merge`/`rebase` from upstream.** Instead, hand-port just the specific files/functions needed:
- Copy as-is (no fork equivalent exists): `pkg/compliant/dcr34.go`, `pkg/compliant/schema/version34.go`, plus its test/testdata files.
- Apply by hand (fork's version already differs): `pkg/compliant/dcr.go` (add the `"3.4"` case alongside this fork's existing `CreateSoftwareClientOnly` branch) and `pkg/compliant/schema/response.go` (register `responseValidator34`).
- Manually reconcile (the one real conflict): `pkg/compliant/auth/signer.go` — take upstream's `subjectDN()`/`escapeRFC2253()`/`oidName()` and rewritten `addTlsClientAuthClaims()` wholesale, then re-apply this fork's own `clientId`/`refresh_token`/`authorization_signed_response_alg` additions on top. Needs a person reading both versions side by side, not a mechanical patch.
- Don't touch the CI/tooling files at all as part of this — irrelevant to DCR behavior, only adds conflict surface.
- After backporting, run `go test ./...` and confirm the 3 fork-only tests above still pass.
- **Verify against a real transport cert** that the new `subjectDN()` output actually starts with `CN=` before enabling v3.4 anywhere — v3.4's schema requires it, v3.2/v3.3 never checked it.
- Record the new sync point once done (this file's "last upstream commit reviewed" table above is the last known point; update it).

### 2. Core (`openig-pyforge`, not in this repo — see `extensions/src/openig/fapi/registration_response_schema.py`)

Pyforge has **no version-selection mechanism at all** — it's a single hardcoded schema mirroring v3.2. If moving Core to v3.4: update `APPLICATION_TYPES` (`"mobile"`→`"native"`, with new tests for both), widen `tls_client_auth_subject_dn` max length and add the `CN=`-pattern check (Core's DN rendering already uses a proper RFC 4514 library, so lower risk than the fork's fix, but still verify against a real cert), and optionally add `jwt-bearer` to `GRANT_TYPES`. If moving Core to v3.3 instead (much cheaper): just add `ciba` to `GRANT_TYPES` and tighten `client_secret` max length from `256` to `36`.

### 3. Decide the target version first

Before doing 1 or 2: decide whether to move to v3.3 (cheap: zero work needed on this fork since it's already there; one small Core change) or v3.4 (bigger: new manifest+schema on this fork, breaking `application_type` rename + DN pattern work on both sides). This is a policy call for whoever owns the OB conformance relationship, not a technical one.
