---
name: ashby
description: Interact with Ashby ATS for Lightning Labs recruiting. Use when checking applicant volume, reviewing candidates, screening applicants for relevance, browsing open jobs, viewing the hiring pipeline, or annotating candidates. Trigger phrases include "check ashby", "review applicants", "screen candidates", "hiring pipeline", "applicant volume", "recruiting dashboard", "candidate search", "open roles".
allowed-tools: Read, Bash(python:*)
---

# Ashby ATS — Lightning Labs Recruiting

Browse jobs, manage applications, screen candidates against Lightning Labs
hiring criteria, and view pipeline dashboards.

## Python & Paths

Resolve paths relative to this skill's location. `SKILL_DIR` is the directory
containing this SKILL.md file.

```bash
SKILL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")/.claude/skills/ashby" 2>/dev/null && pwd || echo ".claude/skills/ashby")"
PYTHON="python3"
CLIENT="$SKILL_DIR/scripts/ashby_client.py"
SCREENER="$SKILL_DIR/scripts/screen_candidates.py"
# API key from: $ASHBY_API_KEY (Basic Auth)
# Requires: pip install requests
```

For quick use, set these in your shell:
```bash
export ASHBY_SKILL_DIR="$(git rev-parse --show-toplevel 2>/dev/null || pwd)/.claude/skills/ashby"
alias ashby="python3 $ASHBY_SKILL_DIR/scripts/ashby_client.py"
alias ashby-screen="python3 $ASHBY_SKILL_DIR/scripts/screen_candidates.py"
```

## Commands

### List Open Jobs

```bash
python3 $SKILL_DIR/scripts/ashby_client.py jobs --status Open
```

### Get Job Details

```bash
python3 $SKILL_DIR/scripts/ashby_client.py jobs --id <jobId>
```

### List Applications for a Job

```bash
python3 $SKILL_DIR/scripts/ashby_client.py applications --job-id <jobId> --status Active
python3 $SKILL_DIR/scripts/ashby_client.py applications --job-id <jobId> --status Active --limit 50
```

### Get Application Details (with form submissions)

```bash
python3 $SKILL_DIR/scripts/ashby_client.py applications --id <applicationId> --expand
```

### Search Candidates

```bash
python3 $SKILL_DIR/scripts/ashby_client.py candidates --search-email "user@example.com"
python3 $SKILL_DIR/scripts/ashby_client.py candidates --search-name "Jane Doe"
```

### Get Candidate Details

```bash
python3 $SKILL_DIR/scripts/ashby_client.py candidates --id <candidateId>
```

### Pipeline Dashboard

```bash
python3 $SKILL_DIR/scripts/ashby_client.py dashboard
python3 $SKILL_DIR/scripts/ashby_client.py dashboard --format markdown
```

### Add Note to Candidate

```bash
python3 $SKILL_DIR/scripts/ashby_client.py notes --candidate-id <candidateId> --body "<b>AI Screen</b>: Strong match — Go, Lightning, distributed systems background."
```

### List / Add Tags

```bash
python3 $SKILL_DIR/scripts/ashby_client.py tags
python3 $SKILL_DIR/scripts/ashby_client.py tags --add --candidate-id <candidateId> --tag-id <tagId>
```

### Move Application Stage

```bash
python3 $SKILL_DIR/scripts/ashby_client.py applications --change-stage --id <applicationId> --stage-id <interviewStageId>
```

## Screening Candidates

Screen applications against Lightning Labs hiring criteria using keyword-based
scoring. Claude interprets the raw scores and applies nuanced judgment.

### Screen All Active Applications for a Job

Use `--enrich` to fetch full application details (form submissions, resume
handles) — this gives the screener much more data to work with but is slower
since it fetches each application individually.

```bash
python3 $SKILL_DIR/scripts/ashby_client.py applications --job-id <jobId> --status Active --enrich | python3 $SKILL_DIR/scripts/screen_candidates.py
python3 $SKILL_DIR/scripts/ashby_client.py applications --job-id <jobId> --status Active --enrich --limit 20 | python3 $SKILL_DIR/scripts/screen_candidates.py --format markdown
```

### Screen a Single Application

```bash
python3 $SKILL_DIR/scripts/screen_candidates.py --application-id <applicationId>
```

### Screen from File

```bash
python3 $SKILL_DIR/scripts/screen_candidates.py --file /tmp/claude/applications.json
```

### Filter by Minimum Tier

```bash
python3 $SKILL_DIR/scripts/ashby_client.py applications --job-id <jobId> --status Active --enrich | python3 $SKILL_DIR/scripts/screen_candidates.py --min-tier moderate
```

### Screening Output

Candidates are scored (0–100%) and classified into tiers:
- **strong** (≥60%): Highly relevant background — prioritize for review
- **moderate** (35–59%): Some relevant experience — worth a closer look
- **weak** (15–34%): Minimal signals — lower priority
- **no_signal** (<15%): No relevant keywords detected

Each candidate gets a per-category breakdown showing matched keywords.

## Lightning Labs Screening Criteria

The screener scores candidates across these weighted categories:

| Category | Weight | What We Look For |
|---|---|---|
| Bitcoin & Lightning | 3.0 | Bitcoin protocol, Lightning Network, lnd, taproot, payment channels |
| Go / Golang | 2.5 | Go proficiency, gRPC, protobuf, goroutines |
| Systems Languages | 2.0 | Rust, C++, systems programming |
| Distributed Systems | 2.0 | Consensus, fault tolerance, P2P, gossip protocols |
| Networking & Protocols | 2.0 | Protocol design, TCP/UDP, Tor, noise protocol |
| Cryptography & Security | 2.0 | Applied crypto, signatures, ZK proofs, security audits |
| Open Source | 1.5 | OSS contributions, maintainership, code review |
| PhD & Research | 1.5 | Doctorate, publications, academic research |
| Operating Systems | 1.0 | Kernel, Linux, embedded, POSIX |

**Note**: Resume content is not available through the Ashby API. Screening works
on application form data, custom fields, and metadata. Flag moderate-tier
candidates for manual resume review when the form data is sparse.

## Active Roles (Lightning Labs)

- Senior Security Engineer
- Lightning Protocol Engineer
- Lightning Infrastructure Engineer
- Assets Protocol Engineer
- Platform Engineer
- Technical Product Manager
- Business Development Strategist
- Lightning Developer Evangelist
- Senior Engineering Manager

## Workflow Examples

### Weekly Applicant Volume Check

```bash
# Get the dashboard overview.
python3 $SKILL_DIR/scripts/ashby_client.py dashboard --format markdown

# Drill into a specific role.
python3 $SKILL_DIR/scripts/ashby_client.py applications --job-id <jobId> --status Active --limit 20
```

### Screen New Candidates for Protocol Engineer

```bash
# Find the job ID.
python3 $SKILL_DIR/scripts/ashby_client.py jobs --status Open

# Screen all active applications (with enriched data for better scoring).
python3 $SKILL_DIR/scripts/ashby_client.py applications --job-id <jobId> --status Active --enrich | python3 $SKILL_DIR/scripts/screen_candidates.py --format markdown

# Get details on a strong candidate.
python3 $SKILL_DIR/scripts/ashby_client.py applications --id <applicationId> --expand
```

### Annotate Strong Candidates

```bash
# Screen and identify strong matches.
python3 $SKILL_DIR/scripts/ashby_client.py applications --job-id <jobId> --status Active --enrich | python3 $SKILL_DIR/scripts/screen_candidates.py --min-tier strong

# Add a note with screening summary.
python3 $SKILL_DIR/scripts/ashby_client.py notes --candidate-id <candidateId> --body "<b>AI Screen</b>: Strong match (78%). Bitcoin/Lightning expertise, Go proficiency, distributed systems background."

# Tag the candidate.
python3 $SKILL_DIR/scripts/ashby_client.py tags  # List available tags first.
python3 $SKILL_DIR/scripts/ashby_client.py tags --add --candidate-id <candidateId> --tag-id <tagId>
```

## Output Format

All commands produce JSON by default. Use `--format compact` for minimal output
or `--format markdown` for human-readable tables (dashboard and screener).
