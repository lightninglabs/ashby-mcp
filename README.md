# Ashby ATS Skill for Claude Code

A Claude Code skill that integrates with the [Ashby ATS API](https://developers.ashbyhq.com/reference/introduction) for Lightning Labs recruiting. Browse jobs, manage applications, screen candidates against hiring criteria, view pipeline dashboards, and annotate candidates — all from Claude Code.

## Setup

### Prerequisites

- Python 3.10+
- `requests` library (`pip install requests`)
- An Ashby API key with `candidatesRead`, `jobsRead`, and `candidatesWrite` permissions
  - Generate one at https://app.ashbyhq.com/admin/api/keys

### Installation

1. Clone this repo:
   ```bash
   git clone <repo-url> ashby-skills
   cd ashby-skills
   ```

2. Install the Python dependency:
   ```bash
   pip install requests
   ```

3. Set your Ashby API key:
   ```bash
   export ASHBY_API_KEY="your-api-key-here"
   ```
   Add this to your shell profile (`~/.zshrc`, `~/.bashrc`, etc.) to persist it.

4. Open Claude Code from the repo directory:
   ```bash
   cd ashby-skills
   claude
   ```

   The skill auto-discovers from `.claude/skills/ashby/SKILL.md`. You can now
   ask Claude things like "check the hiring pipeline" or "screen candidates
   for the protocol engineer role".

### Claude Code Agent Setup

If you are a Claude Code agent (not a human), the skill is automatically
available when your working directory is this repo. Use these paths:

```bash
SKILL_DIR=".claude/skills/ashby"
python3 $SKILL_DIR/scripts/ashby_client.py jobs --status Open
python3 $SKILL_DIR/scripts/ashby_client.py dashboard --format markdown
python3 $SKILL_DIR/scripts/ashby_client.py applications --job-id <id> --status Active --enrich \
  | python3 $SKILL_DIR/scripts/screen_candidates.py --format markdown
```

The `ASHBY_API_KEY` environment variable must be set.

## What's Included

```
.claude/skills/ashby/
  SKILL.md                    # Claude Code skill definition (auto-discovered)
  scripts/
    ashby_client.py           # CLI wrapper around the Ashby REST API
    screen_candidates.py      # Candidate screening/scoring tool
```

### ashby_client.py

A CLI tool with these subcommands:

| Command | Description |
|---------|-------------|
| `jobs` | List/view open positions (`--status Open`, `--id`) |
| `applications` | List/view/move applications (`--job-id`, `--status`, `--enrich`, `--expand`) |
| `candidates` | List/search/view candidates (`--search-email`, `--search-name`, `--id`) |
| `tags` | List tags or add tags to candidates (`--add`, `--candidate-id`, `--tag-id`) |
| `notes` | Add notes to candidates (`--candidate-id`, `--body`) |
| `dashboard` | Aggregated pipeline view across all open jobs |

All commands output JSON by default. Use `--format markdown` for tables or `--format compact` for minimal output.

### screen_candidates.py

Scores candidates against Lightning Labs hiring criteria using weighted keyword matching:

| Category | Weight | What We Look For |
|----------|--------|------------------|
| Bitcoin & Lightning | 3.0 | bitcoin, lnd, taproot, HTLCs, payment channels |
| Go / Golang | 2.5 | golang, gRPC, protobuf, goroutines |
| Systems Languages | 2.0 | Rust, C++, systems programming |
| Distributed Systems | 2.0 | consensus, P2P, fault tolerance |
| Networking & Protocols | 2.0 | protocol design, Tor, noise protocol |
| Cryptography & Security | 2.0 | applied crypto, signatures, ZK proofs |
| Open Source | 1.5 | OSS contributions, maintainership |
| PhD & Research | 1.5 | doctorate, publications |
| Operating Systems | 1.0 | kernel, Linux, embedded |

Candidates are classified into tiers: **strong** (>=60%), **moderate** (35-59%), **weak** (15-34%), **no_signal** (<15%).

Input modes:
- Pipe from ashby_client: `ashby_client.py applications ... | screen_candidates.py`
- From file: `screen_candidates.py --file applications.json`
- Single application: `screen_candidates.py --application-id <id>`

## Quick Start Examples

```bash
# See all open roles
python3 .claude/skills/ashby/scripts/ashby_client.py jobs --status Open

# Pipeline dashboard
python3 .claude/skills/ashby/scripts/ashby_client.py dashboard --format markdown

# Screen active candidates for a specific job
python3 .claude/skills/ashby/scripts/ashby_client.py applications \
  --job-id <jobId> --status Active --enrich \
  | python3 .claude/skills/ashby/scripts/screen_candidates.py --format markdown

# Search for a candidate
python3 .claude/skills/ashby/scripts/ashby_client.py candidates --search-name "Jane Doe"

# Add a screening note
python3 .claude/skills/ashby/scripts/ashby_client.py notes \
  --candidate-id <id> \
  --body "<b>AI Screen</b>: Strong match (72%). Bitcoin/Lightning, distributed systems."
```

## API Reference

All communication with Ashby uses their [REST API](https://developers.ashbyhq.com/reference/introduction):
- **Base URL**: `https://api.ashbyhq.com`
- **Auth**: Basic Auth (API key as username, empty password)
- **Method**: All endpoints are POST with JSON bodies
- **Pagination**: Cursor-based (`moreDataAvailable` + `nextCursor`)

## Limitations

- **No resume content**: The Ashby API provides resume file handles but not actual file content. Screening works on application form data, custom fields, and metadata only.
- **"Go" keyword ambiguity**: The screener uses "golang", "goroutine", "gRPC" etc. as signals rather than bare "Go" which is too common in English.
- **Dashboard speed**: Makes N+1 API calls (1 for jobs + 1 per job for applications). Fine for ~10 open roles.
