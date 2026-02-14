#!/usr/bin/env python3
"""
Candidate screening tool for Lightning Labs recruiting.

Scores Ashby application data against Lightning Labs hiring criteria
using weighted keyword matching. Claude provides the interpretive layer
on top of the raw scores.

Usage:
    # Pipe from ashby_client.py:
    python ashby_client.py applications --job-id X --status Active | python screen_candidates.py

    # From a JSON file:
    python screen_candidates.py --file /tmp/claude/applications.json

    # Fetch and screen a single application:
    python screen_candidates.py --application-id <id>
"""

import argparse
import json
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

# Allow importing api_call from the sibling ashby_client module.
sys.path.insert(0, str(Path(__file__).resolve().parent))


# =============================================================================
# SCREENING CRITERIA
# =============================================================================

# Weighted keyword categories for Lightning Labs candidate scoring.
# Weights reflect relative importance to the hiring team.
CRITERIA = {
    "bitcoin_lightning": {
        "weight": 3.0,
        "label": "Bitcoin & Lightning",
        "keywords": [
            "bitcoin", "lightning network", "lightning", "lnd", "btcd",
            "taproot", "segwit", "bolt spec", "bolt11", "bolt12",
            "htlc", "payment channel", "state channel",
            "on-chain", "layer 2", "layer-2", "l2 protocol",
            "satoshi", "sats", "utxo",
            "neutrino", "loop", "pool", "taproot assets", "taro",
            "lightning-terminal", "aperture", "lsat", "l402",
            "macaroon", "lnurl", "keysend", "amp",
            "bitcoin core", "bitcoind", "bitcoin improvement proposal",
            "bip", "lightning labs", "blockstream",
            "rgb protocol", "liquid network", "sidechain",
            "nostr", "dlc", "discreet log contract",
        ],
    },
    "golang": {
        "weight": 2.5,
        "label": "Go / Golang",
        "keywords": [
            "golang", "go language", "go programming",
            "goroutine", "goroutines", "go channel", "go channels",
            "grpc", "protobuf", "protocol buffers",
            "go mod", "go module", "go modules",
            "go test", "go vet", "go fmt",
            "cobra", "viper",
        ],
    },
    "systems_languages": {
        "weight": 2.0,
        "label": "Systems Languages",
        "keywords": [
            "rust", "rustlang", "cargo",
            "c++", "cpp", "c programming",
            "systems programming", "low-level", "low level",
            "memory management", "unsafe code",
            "llvm", "compiler", "linker",
        ],
    },
    "distributed_systems": {
        "weight": 2.0,
        "label": "Distributed Systems",
        "keywords": [
            "distributed systems", "distributed computing",
            "consensus", "consensus algorithm",
            "raft", "paxos", "pbft",
            "replication", "fault tolerance", "fault tolerant",
            "cap theorem", "eventual consistency",
            "distributed database", "distributed ledger",
            "etcd", "zookeeper",
            "cluster", "clustering",
            "peer-to-peer", "p2p", "gossip protocol",
            "state machine replication",
            "byzantine", "byzantine fault",
        ],
    },
    "networking": {
        "weight": 2.0,
        "label": "Networking & Protocols",
        "keywords": [
            "protocol engineering", "protocol design",
            "tcp", "udp", "quic",
            "networking", "network programming", "socket",
            "wire protocol", "transport layer",
            "tor", "onion routing", "noise protocol",
            "multiplexing", "muxer",
            "nat traversal", "hole punching",
            "dns", "dhcp", "bgp",
            "packet", "routing", "mesh network",
        ],
    },
    "cryptography": {
        "weight": 2.0,
        "label": "Cryptography & Security",
        "keywords": [
            "cryptography", "cryptographic",
            "encryption", "decryption",
            "elliptic curve", "ecdsa", "ed25519",
            "schnorr", "musig", "musig2",
            "threshold signature", "multisig", "multi-sig",
            "zero knowledge", "zk-snark", "zk-stark", "zkp",
            "hash function", "sha256", "sha-256",
            "digital signature", "signature scheme",
            "commitment scheme", "pedersen commitment",
            "key derivation", "hkdf", "pbkdf",
            "tls", "ssl", "pki",
            "secure multiparty computation", "mpc",
            "homomorphic",
            "security audit", "penetration testing", "pentest",
            "vulnerability", "cve", "bug bounty",
        ],
    },
    "open_source": {
        "weight": 1.5,
        "label": "Open Source",
        "keywords": [
            "open source", "open-source", "oss",
            "github", "gitlab",
            "contributor", "contributions",
            "maintainer", "core maintainer",
            "pull request", "code review",
            "foss", "free software",
            "apache license", "mit license", "bsd license",
        ],
    },
    "phd_research": {
        "weight": 1.5,
        "label": "PhD & Research",
        "keywords": [
            "phd", "ph.d", "ph.d.",
            "doctorate", "doctoral",
            "dissertation", "thesis",
            "research", "researcher",
            "publication", "published paper", "academic paper",
            "conference paper", "journal paper",
            "professor", "postdoc", "post-doc",
            "computer science", "mathematics",
        ],
    },
    "operating_systems": {
        "weight": 1.0,
        "label": "Operating Systems",
        "keywords": [
            "operating system", "operating systems",
            "kernel", "linux kernel", "kernel module",
            "linux", "unix", "posix",
            "embedded", "embedded systems",
            "firmware", "rtos",
            "device driver", "syscall", "system call",
            "memory allocator", "scheduler",
        ],
    },
}


# =============================================================================
# TEXT EXTRACTION
# =============================================================================


def extract_text(application: dict) -> str:
    """Extract all searchable text from an application record.

    Recursively walks the application structure and concatenates all
    string values into a single searchable blob.
    """
    parts = []

    # Candidate info.
    candidate = application.get("candidate", {})
    if candidate:
        parts.append(candidate.get("name", ""))
        email = candidate.get("primaryEmailAddress", {})
        if email:
            parts.append(email.get("value", ""))

    # Custom fields.
    for field in application.get("customFields", []):
        title = field.get("title", "")
        value = field.get("value", "")
        label = field.get("valueLabel", "")
        parts.append(str(title))
        parts.append(str(value))
        if label:
            parts.append(str(label))

    # Application form submissions (expanded).
    for submission in application.get("applicationFormSubmissions", []):
        _extract_form_text(submission, parts)

    # Source info.
    source = application.get("source", {})
    if source:
        parts.append(source.get("title", ""))

    # Resume filename (content not available via API).
    resume = application.get("resumeFileHandle", {})
    if resume:
        parts.append(resume.get("name", ""))

    # Job info.
    job = application.get("job", {})
    if job:
        parts.append(job.get("title", ""))

    return " ".join(filter(None, parts))


def _extract_form_text(obj: Any, parts: list):
    """Recursively extract string values from form submission data."""
    if isinstance(obj, str):
        parts.append(obj)
    elif isinstance(obj, dict):
        for v in obj.values():
            _extract_form_text(v, parts)
    elif isinstance(obj, list):
        for item in obj:
            _extract_form_text(item, parts)


# =============================================================================
# SCORING
# =============================================================================


def score_candidate(text: str) -> dict:
    """Score candidate text against Lightning Labs hiring criteria.

    Returns a breakdown of scores per category and an overall total.
    """
    text_lower = text.lower()
    categories = {}
    total = 0.0
    max_possible = 0.0

    for key, config in CRITERIA.items():
        weight = config["weight"]
        keywords = config["keywords"]

        matched = [kw for kw in keywords if kw.lower() in text_lower]

        # Score: min(matched / 3, 1.0) * weight.
        # Three or more keyword matches in a category earns full score.
        raw = min(len(matched) / 3.0, 1.0)
        category_score = round(raw * weight, 2)

        categories[key] = {
            "label": config["label"],
            "score": category_score,
            "max": weight,
            "matched": matched,
            "match_count": len(matched),
        }

        total += category_score
        max_possible += weight

    pct = round(total / max_possible * 100, 1) if max_possible > 0 else 0.0

    return {
        "total_score": round(total, 2),
        "max_possible": round(max_possible, 2),
        "pct": pct,
        "categories": categories,
    }


def classify_tier(pct: float) -> str:
    """Classify candidate into a tier based on score percentage."""
    if pct >= 60:
        return "strong"
    elif pct >= 35:
        return "moderate"
    elif pct >= 15:
        return "weak"
    else:
        return "no_signal"


# =============================================================================
# SCREENING PIPELINE
# =============================================================================


def screen_applications(applications: list) -> dict:
    """Screen a list of application records and return ranked results."""
    screened = []

    for app in applications:
        text = extract_text(app)
        score = score_candidate(text)
        tier = classify_tier(score["pct"])

        candidate = app.get("candidate", {})
        job = app.get("job", {})
        stage = app.get("currentInterviewStage", {})

        screened.append({
            "candidate_id": candidate.get("id", ""),
            "candidate_name": candidate.get("name", "Unknown"),
            "application_id": app.get("id", ""),
            "job_title": job.get("title", "Unknown"),
            "stage": stage.get("title", "Unknown") if stage else "Unknown",
            "status": app.get("status", "Unknown"),
            "tier": tier,
            "score": score,
        })

    # Sort by score descending.
    screened.sort(key=lambda x: x["score"]["total_score"], reverse=True)

    # Summary counts.
    tier_counts = {"strong": 0, "moderate": 0, "weak": 0, "no_signal": 0}
    for s in screened:
        tier_counts[s["tier"]] += 1

    return {
        "screened_at": datetime.now(timezone.utc).isoformat(),
        "total_candidates": len(screened),
        "summary": tier_counts,
        "candidates": screened,
    }


# =============================================================================
# OUTPUT FORMATTING
# =============================================================================


def print_markdown(results: dict):
    """Print screening results as markdown."""
    summary = results["summary"]
    total = results["total_candidates"]

    print(f"# Candidate Screening Results")
    print(f"\n**{total} candidates screened** | "
          f"Strong: {summary['strong']} | "
          f"Moderate: {summary['moderate']} | "
          f"Weak: {summary['weak']} | "
          f"No Signal: {summary['no_signal']}\n")

    print("| Tier | Candidate | Job | Stage | Score | Top Matches |")
    print("|------|-----------|-----|-------|-------|-------------|")

    for c in results["candidates"]:
        # Collect top matched keywords across categories.
        top_matches = []
        for cat_data in c["score"]["categories"].values():
            if cat_data["matched"]:
                top_matches.extend(cat_data["matched"][:2])

        matches_str = ", ".join(top_matches[:6]) if top_matches else "-"
        tier_emoji = {
            "strong": "+++",
            "moderate": "++",
            "weak": "+",
            "no_signal": "-",
        }[c["tier"]]

        print(f"| {tier_emoji} {c['tier']} | {c['candidate_name']} | "
              f"{c['job_title']} | {c['stage']} | "
              f"{c['score']['pct']}% | {matches_str} |")


# =============================================================================
# MAIN
# =============================================================================


def main():
    parser = argparse.ArgumentParser(
        description="Screen Ashby candidates against Lightning Labs hiring criteria",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "--file",
        help="Read application JSON from file (output of ashby_client.py applications)",
    )
    parser.add_argument(
        "--application-id",
        help="Fetch and screen a single application by ID (calls Ashby API)",
    )
    parser.add_argument(
        "--min-tier",
        choices=["strong", "moderate", "weak", "no_signal"],
        help="Only show candidates at or above this tier",
    )
    parser.add_argument(
        "--format", "-f",
        choices=["json", "compact", "markdown"],
        default="json",
        help="Output format (default: json)",
    )

    args = parser.parse_args()

    # Determine input source.
    if args.application_id:
        # Fetch single application from API.
        from ashby_client import api_call
        data = api_call(
            "application.info",
            applicationId=args.application_id,
            expand=["applicationFormSubmissions", "openings", "referrals"],
        )
        applications = [data.get("results", {})]
    elif args.file:
        with open(args.file, "r") as f:
            raw = json.load(f)
        # Handle both {"applications": [...]} and bare [...] formats.
        if isinstance(raw, list):
            applications = raw
        else:
            applications = raw.get("applications", [raw])
    else:
        # Read from stdin.
        raw = json.load(sys.stdin)
        if isinstance(raw, list):
            applications = raw
        else:
            applications = raw.get("applications", [raw])

    # Screen candidates.
    results = screen_applications(applications)

    # Filter by minimum tier if requested.
    if args.min_tier:
        tier_order = {"strong": 3, "moderate": 2, "weak": 1, "no_signal": 0}
        min_level = tier_order[args.min_tier]
        results["candidates"] = [
            c for c in results["candidates"]
            if tier_order[c["tier"]] >= min_level
        ]
        results["total_candidates"] = len(results["candidates"])

    # Output.
    if args.format == "markdown":
        print_markdown(results)
    elif args.format == "compact":
        print(json.dumps(results, separators=(",", ":"), default=str))
    else:
        print(json.dumps(results, indent=2, default=str))


if __name__ == "__main__":
    main()
