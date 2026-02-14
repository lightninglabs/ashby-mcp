# Ashby ATS MCP Server & Claude Code Skill

An MCP server and Claude Code skill that integrates with the
[Ashby ATS API](https://developers.ashbyhq.com/reference/introduction)
for Lightning Labs recruiting. Browse jobs, manage applications, screen
candidates against hiring criteria, view pipeline dashboards, and
annotate candidates.

## MCP Server (Go)

The primary interface is a Go MCP server that exposes 19 Ashby tools
over stdio transport. Works with Claude Desktop, Claude Code, or any
MCP-compatible client.

### Prerequisites

- Go 1.26+
- An Ashby API key with `candidatesRead`, `jobsRead`, and
  `candidatesWrite` permissions
  - Generate one at https://app.ashbyhq.com/admin/api/keys

### Build

```bash
git clone <repo-url> ashby-skills
cd ashby-skills
go build -o ashby-mcp .
```

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "ashby": {
      "command": "/absolute/path/to/ashby-mcp",
      "env": {
        "ASHBY_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

### Claude Code

```bash
claude mcp add ashby /absolute/path/to/ashby-mcp \
  -e ASHBY_API_KEY=your-api-key-here
```

Or add to `.claude/settings.json`:

```json
{
  "mcpServers": {
    "ashby": {
      "type": "stdio",
      "command": "/absolute/path/to/ashby-mcp",
      "env": {
        "ASHBY_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

### Available Tools

#### Read-Only

| Tool | Description |
|------|-------------|
| `list_jobs` | List jobs with optional status filter (Open/Closed/Archived/Draft) |
| `get_job` | Get job details by ID |
| `search_jobs` | Search jobs by title |
| `list_applications` | List applications with filters (jobId, status), cursor pagination |
| `get_application` | Get application details with optional field expansion |
| `list_candidates` | List candidates with pagination |
| `search_candidates` | Search candidates by email or name |
| `get_candidate` | Get candidate details by ID |
| `list_tags` | List all candidate tags |
| `list_candidate_notes` | List notes for a candidate |
| `list_interview_stages` | List interview stages for a job's plan |
| `list_interviews` | List interviews by application |
| `pipeline_dashboard` | Aggregated pipeline stats across all open jobs |
| `screen_candidates` | Score candidates against Lightning Labs hiring criteria |

#### Write

| Tool | Description |
|------|-------------|
| `change_application_stage` | Move application to a different interview stage |
| `create_candidate` | Create a new candidate record |
| `add_candidate_tag` | Add a tag to a candidate |
| `create_candidate_note` | Add an HTML-formatted note to a candidate |
| `create_application` | Create an application linking candidate to job |

### Screening Criteria

The `screen_candidates` tool scores applicants using weighted keyword
matching:

| Category | Weight | Signals |
|----------|--------|---------|
| Bitcoin & Lightning | 3.0 | bitcoin, lnd, taproot, HTLCs, payment channels |
| Go / Golang | 2.5 | golang, gRPC, protobuf, goroutines |
| Systems Languages | 2.0 | Rust, C++, systems programming |
| Distributed Systems | 2.0 | consensus, P2P, fault tolerance |
| Networking & Protocols | 2.0 | protocol design, Tor, noise protocol |
| Cryptography & Security | 2.0 | applied crypto, signatures, ZK proofs |
| Open Source | 1.5 | OSS contributions, maintainership |
| PhD & Research | 1.5 | doctorate, publications |
| Operating Systems | 1.0 | kernel, Linux, embedded |

Tiers: **strong** (>=60%), **moderate** (35-59%), **weak** (15-34%),
**no_signal** (<15%).

## Claude Code Skill (Python, Legacy)

The Python CLI scripts are still available as a Claude Code skill for
backward compatibility.

### Prerequisites

- Python 3.10+
- `requests` library (`pip install requests`)

### Usage

```bash
export ASHBY_API_KEY="your-api-key-here"

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

All communication with Ashby uses their
[REST API](https://developers.ashbyhq.com/reference/introduction):
- **Base URL**: `https://api.ashbyhq.com`
- **Auth**: Basic Auth (API key as username, empty password)
- **Method**: All endpoints are POST with JSON bodies
- **Pagination**: Cursor-based (`moreDataAvailable` + `nextCursor`)

## Limitations

- **No resume content**: The Ashby API provides resume file handles but
  not actual file content. Screening works on application form data,
  custom fields, and metadata only.
- **"Go" keyword ambiguity**: The screener uses "golang", "goroutine",
  "gRPC" etc. as signals rather than bare "Go" which is too common in
  English.
- **Dashboard speed**: Makes N+1 API calls (1 for jobs + 1 per job for
  applications). Fine for ~10 open roles.
