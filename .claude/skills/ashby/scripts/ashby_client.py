#!/usr/bin/env python3
"""
Ashby ATS API client for Claude Code skills.

Wraps the Ashby REST API (https://api.ashbyhq.com) with a CLI interface
for browsing jobs, managing applications, searching candidates, and
viewing pipeline dashboards.

Usage:
    python ashby_client.py jobs --status Open
    python ashby_client.py jobs --id <jobId>
    python ashby_client.py applications --job-id <jobId> --status Active
    python ashby_client.py applications --id <appId> --expand
    python ashby_client.py candidates --search-email "foo@bar.com"
    python ashby_client.py candidates --id <candidateId>
    python ashby_client.py tags
    python ashby_client.py notes --candidate-id <id> --body "Note text"
    python ashby_client.py dashboard
"""

import argparse
import base64
import json
import os
import sys
import time
from collections import defaultdict
from typing import Any, Optional

try:
    import requests
except ImportError:
    print("Error: requests not installed. Run: pip install requests")
    sys.exit(1)


BASE_URL = "https://api.ashbyhq.com"


# =============================================================================
# CORE API FUNCTIONS
# =============================================================================


def get_api_key() -> str:
    """Get Ashby API key from environment."""
    key = (
        os.environ.get("ASHBY_API_KEY")
        or os.environ.get("ASHBY_KEY")
    )
    if not key:
        print("Error: ASHBY_API_KEY environment variable not set")
        sys.exit(1)
    return key


def api_call(endpoint: str, **params) -> dict:
    """Make an Ashby API call.

    All Ashby endpoints are POST with JSON bodies and Basic Auth.
    """
    key = get_api_key()
    # Basic auth: api_key as username, empty password.
    auth_token = base64.b64encode(f"{key}:".encode()).decode()

    # Strip None-valued params.
    body = {k: v for k, v in params.items() if v is not None}

    url = f"{BASE_URL}/{endpoint}"

    try:
        response = requests.post(
            url,
            json=body,
            headers={
                "Authorization": f"Basic {auth_token}",
                "Content-Type": "application/json",
            },
            timeout=30,
        )
        response.raise_for_status()
        data = response.json()

        if not data.get("success", False):
            error_info = data.get("errorInfo", {})
            msg = error_info.get("message", data.get("errors", ["Unknown error"]))
            print(f"API Error ({endpoint}): {msg}")
            sys.exit(1)

        return data
    except requests.exceptions.HTTPError as e:
        # Handle specific HTTP error codes.
        if e.response is not None:
            status = e.response.status_code
            if status == 401:
                print("Error: Invalid or missing API key (HTTP 401)")
            elif status == 403:
                print("Error: API key lacks required permissions (HTTP 403)")
            else:
                print(f"HTTP Error {status}: {e}")
        else:
            print(f"HTTP Error: {e}")
        sys.exit(1)
    except requests.exceptions.RequestException as e:
        print(f"Request failed: {e}")
        sys.exit(1)


def paginate(endpoint: str, limit: Optional[int] = None, **params) -> list:
    """Paginate through an Ashby list endpoint.

    Accumulates results across pages using cursor-based pagination.
    Respects the optional limit parameter to cap total results.
    """
    all_results = []
    cursor = None

    while True:
        call_params = {**params}
        if cursor:
            call_params["cursor"] = cursor
        # Ashby max per page is 100.
        call_params["limit"] = min(100, limit - len(all_results)) if limit else 100

        data = api_call(endpoint, **call_params)
        results = data.get("results", [])
        all_results.extend(results)

        # Stop if we've hit the requested limit.
        if limit and len(all_results) >= limit:
            all_results = all_results[:limit]
            break

        # Stop if no more data.
        if not data.get("moreDataAvailable", False):
            break

        cursor = data.get("nextCursor")
        if not cursor:
            break

        # Rate-limit politeness.
        time.sleep(0.2)

    return all_results


def format_output(data: Any, fmt: str = "json") -> str:
    """Format output as JSON."""
    if fmt == "compact":
        return json.dumps(data, separators=(",", ":"), default=str)
    return json.dumps(data, indent=2, default=str)


# =============================================================================
# JOBS SUBCOMMAND
# =============================================================================


def cmd_jobs(args):
    """Handle jobs subcommand."""
    if args.id:
        # Get single job.
        data = api_call("job.info", jobId=args.id)
        print(format_output(data.get("results", {}), args.format))
        return

    # List all jobs, then filter client-side by status.
    # The job.list API does not support filtering by job status (Open/Closed/Archived)
    # directly — its status param uses application statuses.
    results = paginate("job.list", limit=None)

    if args.status:
        results = [j for j in results if j.get("status") == args.status]

    if args.limit:
        results = results[:args.limit]

    output = {"jobs": results, "total": len(results)}
    print(format_output(output, args.format))


# =============================================================================
# APPLICATIONS SUBCOMMAND
# =============================================================================


def cmd_applications(args):
    """Handle applications subcommand."""
    # Change stage operation.
    if args.change_stage:
        if not args.id or not args.stage_id:
            print("Error: --id and --stage-id required with --change-stage")
            sys.exit(1)
        change_params = {
            "applicationId": args.id,
            "interviewStageId": args.stage_id,
        }
        data = api_call("application.changeStage", **change_params)
        print(format_output(data.get("results", {}), args.format))
        return

    # Get single application.
    if args.id:
        info_params = {"applicationId": args.id}
        if args.expand:
            info_params["expand"] = [
                "applicationFormSubmissions",
                "openings",
                "referrals",
            ]
        data = api_call("application.info", **info_params)
        print(format_output(data.get("results", {}), args.format))
        return

    # List applications with optional filters.
    params = {}
    if args.job_id:
        params["jobId"] = args.job_id
    if args.status:
        params["status"] = args.status

    results = paginate("application.list", limit=args.limit, **params)

    # Enrich mode: fetch full details for each application (slower but gives
    # form submissions, resume handles, and referrals for screening).
    if args.enrich:
        enriched = []
        for app in results:
            app_id = app.get("id")
            if not app_id:
                continue
            detail = api_call(
                "application.info",
                applicationId=app_id,
                expand=["applicationFormSubmissions", "openings", "referrals"],
            )
            enriched.append(detail.get("results", app))
            time.sleep(0.2)
        results = enriched

    output = {"applications": results, "total": len(results)}
    print(format_output(output, args.format))


# =============================================================================
# CANDIDATES SUBCOMMAND
# =============================================================================


def cmd_candidates(args):
    """Handle candidates subcommand."""
    # Get single candidate.
    if args.id:
        data = api_call("candidate.info", candidateId=args.id)
        print(format_output(data.get("results", {}), args.format))
        return

    # Search by email or name.
    if args.search_email or args.search_name:
        search_params = {}
        if args.search_email:
            search_params["email"] = args.search_email
        if args.search_name:
            search_params["name"] = args.search_name
        data = api_call("candidate.search", **search_params)
        results = data.get("results", [])
        output = {"candidates": results, "total": len(results)}
        print(format_output(output, args.format))
        return

    # List all candidates.
    results = paginate("candidate.list", limit=args.limit)
    output = {"candidates": results, "total": len(results)}
    print(format_output(output, args.format))


# =============================================================================
# TAGS SUBCOMMAND
# =============================================================================


def cmd_tags(args):
    """Handle tags subcommand."""
    # Add tag to candidate.
    if args.add:
        if not args.candidate_id or not args.tag_id:
            print("Error: --candidate-id and --tag-id required with --add")
            sys.exit(1)
        data = api_call(
            "candidate.addTag",
            candidateId=args.candidate_id,
            tagId=args.tag_id,
        )
        print(format_output(data.get("results", {}), args.format))
        return

    # List all tags.
    results = paginate("candidateTag.list")
    output = {"tags": results, "total": len(results)}
    print(format_output(output, args.format))


# =============================================================================
# NOTES SUBCOMMAND
# =============================================================================


def cmd_notes(args):
    """Handle notes subcommand."""
    if not args.candidate_id:
        print("Error: --candidate-id required")
        sys.exit(1)
    if not args.body:
        print("Error: --body required")
        sys.exit(1)

    note_params = {
        "candidateId": args.candidate_id,
        "note": args.body,
        "sendNotifications": False,
    }

    data = api_call("candidate.createNote", **note_params)
    print(format_output(data.get("results", {}), args.format))


# =============================================================================
# DASHBOARD SUBCOMMAND
# =============================================================================


def cmd_dashboard(args):
    """Handle dashboard subcommand — aggregated pipeline view."""
    # Get all jobs then filter to open ones client-side.
    all_jobs = paginate("job.list")
    jobs = [j for j in all_jobs if j.get("status") == "Open"]

    job_summaries = []
    grand_total = 0
    grand_active = 0

    for job in jobs:
        job_id = job.get("id")
        job_title = job.get("title", "Unknown")

        # Get all applications for this job.
        apps = paginate("application.list", jobId=job_id)

        # Count by status and stage.
        status_counts = defaultdict(int)
        stage_counts = defaultdict(int)

        for app in apps:
            status = app.get("status", "Unknown")
            status_counts[status] += 1

            stage = app.get("currentInterviewStage", {})
            stage_title = stage.get("title", "Unknown") if stage else "Unknown"
            if status == "Active":
                stage_counts[stage_title] += 1

        total = len(apps)
        grand_total += total
        active = status_counts.get("Active", 0)
        grand_active += active

        job_summaries.append({
            "id": job_id,
            "title": job_title,
            "total_applications": total,
            "active": active,
            "archived": status_counts.get("Archived", 0),
            "hired": status_counts.get("Hired", 0),
            "lead": status_counts.get("Lead", 0),
            "by_stage": dict(stage_counts),
        })

        # Rate-limit politeness between jobs.
        time.sleep(0.2)

    # Sort by total applications descending.
    job_summaries.sort(key=lambda j: j["total_applications"], reverse=True)

    output = {
        "jobs": job_summaries,
        "totals": {
            "jobs": len(jobs),
            "applications": grand_total,
            "active": grand_active,
        },
    }

    if args.format == "markdown":
        _print_dashboard_markdown(output)
    else:
        print(format_output(output, args.format))


def _print_dashboard_markdown(data: dict):
    """Print dashboard in markdown table format."""
    totals = data["totals"]
    print(f"# Ashby Pipeline Dashboard")
    print(f"\n**{totals['jobs']} open jobs** | "
          f"**{totals['applications']} total applications** | "
          f"**{totals['active']} active**\n")

    print("| Job | Active | Archived | Hired | Total |")
    print("|-----|--------|----------|-------|-------|")
    for job in data["jobs"]:
        print(f"| {job['title']} | {job['active']} | "
              f"{job['archived']} | {job['hired']} | "
              f"{job['total_applications']} |")

    print()
    for job in data["jobs"]:
        if job["by_stage"]:
            print(f"### {job['title']}")
            for stage, count in sorted(
                job["by_stage"].items(),
                key=lambda x: x[1],
                reverse=True,
            ):
                print(f"- {stage}: {count}")
            print()


# =============================================================================
# MAIN
# =============================================================================


def main():
    # Shared parent parser for --format flag so it works in any position.
    format_parser = argparse.ArgumentParser(add_help=False)
    format_parser.add_argument(
        "--format", "-f",
        choices=["json", "compact", "markdown"],
        default="json",
        help="Output format (default: json)",
    )

    parser = argparse.ArgumentParser(
        description="Ashby ATS API client",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        parents=[format_parser],
    )

    subparsers = parser.add_subparsers(dest="command", help="Command category")

    # -------------------------------------------------------------------------
    # Jobs subcommand.
    # -------------------------------------------------------------------------
    jobs_parser = subparsers.add_parser(
        "jobs", help="List and view jobs", parents=[format_parser],
    )
    jobs_parser.add_argument("--id", help="Get a specific job by ID")
    jobs_parser.add_argument(
        "--status",
        choices=["Open", "Closed", "Archived", "Draft"],
        help="Filter by job status",
    )
    jobs_parser.add_argument("--limit", type=int, help="Max results to return")
    jobs_parser.set_defaults(func=cmd_jobs)

    # -------------------------------------------------------------------------
    # Applications subcommand.
    # -------------------------------------------------------------------------
    apps_parser = subparsers.add_parser(
        "applications", help="List, view, and manage applications",
        parents=[format_parser],
    )
    apps_parser.add_argument("--id", help="Get a specific application by ID")
    apps_parser.add_argument("--job-id", help="Filter applications by job ID")
    apps_parser.add_argument(
        "--status",
        choices=["Active", "Hired", "Archived", "Lead"],
        help="Filter by application status",
    )
    apps_parser.add_argument(
        "--expand",
        action="store_true",
        help="Expand application form submissions, openings, and referrals",
    )
    apps_parser.add_argument("--limit", type=int, help="Max results to return")
    apps_parser.add_argument(
        "--enrich",
        action="store_true",
        help="Fetch full details for each application (slower, but includes form submissions for screening)",
    )
    apps_parser.add_argument(
        "--change-stage",
        action="store_true",
        help="Move application to a different stage (requires --id and --stage-id)",
    )
    apps_parser.add_argument(
        "--stage-id", help="Target interview stage ID (for --change-stage)",
    )
    apps_parser.set_defaults(func=cmd_applications)

    # -------------------------------------------------------------------------
    # Candidates subcommand.
    # -------------------------------------------------------------------------
    cand_parser = subparsers.add_parser(
        "candidates", help="List, search, and view candidates",
        parents=[format_parser],
    )
    cand_parser.add_argument("--id", help="Get a specific candidate by ID")
    cand_parser.add_argument("--search-email", help="Search candidates by email")
    cand_parser.add_argument("--search-name", help="Search candidates by name")
    cand_parser.add_argument("--limit", type=int, help="Max results to return")
    cand_parser.set_defaults(func=cmd_candidates)

    # -------------------------------------------------------------------------
    # Tags subcommand.
    # -------------------------------------------------------------------------
    tags_parser = subparsers.add_parser(
        "tags", help="List and manage candidate tags", parents=[format_parser],
    )
    tags_parser.add_argument(
        "--add",
        action="store_true",
        help="Add tag to candidate (requires --candidate-id and --tag-id)",
    )
    tags_parser.add_argument("--candidate-id", help="Candidate ID for tag operations")
    tags_parser.add_argument("--tag-id", help="Tag ID to add")
    tags_parser.set_defaults(func=cmd_tags)

    # -------------------------------------------------------------------------
    # Notes subcommand.
    # -------------------------------------------------------------------------
    notes_parser = subparsers.add_parser(
        "notes", help="Add notes to candidates", parents=[format_parser],
    )
    notes_parser.add_argument(
        "--candidate-id", required=True, help="Candidate ID to add note to",
    )
    notes_parser.add_argument(
        "--body", required=True, help="Note content (HTML supported: b, i, u, a, ul, ol, li, code, pre)",
    )
    notes_parser.set_defaults(func=cmd_notes)

    # -------------------------------------------------------------------------
    # Dashboard subcommand.
    # -------------------------------------------------------------------------
    dash_parser = subparsers.add_parser(
        "dashboard", help="Aggregated pipeline view across all open jobs",
        parents=[format_parser],
    )
    dash_parser.set_defaults(func=cmd_dashboard)

    # Parse and execute.
    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        sys.exit(1)

    args.func(args)


if __name__ == "__main__":
    main()
