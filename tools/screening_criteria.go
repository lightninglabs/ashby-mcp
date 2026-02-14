package tools

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// Criterion defines a weighted keyword category used for
// candidate screening.
type Criterion struct {
	// Weight is the relative importance of this category.
	Weight float64

	// Label is the human-readable category name.
	Label string

	// Keywords is the list of terms to match against.
	Keywords []string
}

// criteria maps category keys to their screening configuration.
// Weights and keywords are identical to the Python
// screen_candidates.py implementation.
var criteria = map[string]Criterion{
	"bitcoin_lightning": {
		Weight: 3.0,
		Label:  "Bitcoin & Lightning",
		Keywords: []string{
			"bitcoin", "lightning network", "lightning",
			"lnd", "btcd",
			"taproot", "segwit", "bolt spec", "bolt11",
			"bolt12",
			"htlc", "payment channel", "state channel",
			"on-chain", "layer 2", "layer-2",
			"l2 protocol",
			"satoshi", "sats", "utxo",
			"neutrino", "loop", "pool",
			"taproot assets", "taro",
			"lightning-terminal", "aperture", "lsat",
			"l402",
			"macaroon", "lnurl", "keysend", "amp",
			"bitcoin core", "bitcoind",
			"bitcoin improvement proposal",
			"bip", "lightning labs", "blockstream",
			"rgb protocol", "liquid network", "sidechain",
			"nostr", "dlc", "discreet log contract",
		},
	},
	"golang": {
		Weight: 2.5,
		Label:  "Go / Golang",
		Keywords: []string{
			"golang", "go language", "go programming",
			"goroutine", "goroutines",
			"go channel", "go channels",
			"grpc", "protobuf", "protocol buffers",
			"go mod", "go module", "go modules",
			"go test", "go vet", "go fmt",
			"cobra", "viper",
		},
	},
	"systems_languages": {
		Weight: 2.0,
		Label:  "Systems Languages",
		Keywords: []string{
			"rust", "rustlang", "cargo",
			"c++", "cpp", "c programming",
			"systems programming", "low-level",
			"low level",
			"memory management", "unsafe code",
			"llvm", "compiler", "linker",
		},
	},
	"distributed_systems": {
		Weight: 2.0,
		Label:  "Distributed Systems",
		Keywords: []string{
			"distributed systems",
			"distributed computing",
			"consensus", "consensus algorithm",
			"raft", "paxos", "pbft",
			"replication", "fault tolerance",
			"fault tolerant",
			"cap theorem", "eventual consistency",
			"distributed database",
			"distributed ledger",
			"etcd", "zookeeper",
			"cluster", "clustering",
			"peer-to-peer", "p2p", "gossip protocol",
			"state machine replication",
			"byzantine", "byzantine fault",
		},
	},
	"networking": {
		Weight: 2.0,
		Label:  "Networking & Protocols",
		Keywords: []string{
			"protocol engineering", "protocol design",
			"tcp", "udp", "quic",
			"networking", "network programming", "socket",
			"wire protocol", "transport layer",
			"tor", "onion routing", "noise protocol",
			"multiplexing", "muxer",
			"nat traversal", "hole punching",
			"dns", "dhcp", "bgp",
			"packet", "routing", "mesh network",
		},
	},
	"cryptography": {
		Weight: 2.0,
		Label:  "Cryptography & Security",
		Keywords: []string{
			"cryptography", "cryptographic",
			"encryption", "decryption",
			"elliptic curve", "ecdsa", "ed25519",
			"schnorr", "musig", "musig2",
			"threshold signature", "multisig",
			"multi-sig",
			"zero knowledge", "zk-snark", "zk-stark",
			"zkp",
			"hash function", "sha256", "sha-256",
			"digital signature", "signature scheme",
			"commitment scheme", "pedersen commitment",
			"key derivation", "hkdf", "pbkdf",
			"tls", "ssl", "pki",
			"secure multiparty computation", "mpc",
			"homomorphic",
			"security audit", "penetration testing",
			"pentest",
			"vulnerability", "cve", "bug bounty",
		},
	},
	"open_source": {
		Weight: 1.5,
		Label:  "Open Source",
		Keywords: []string{
			"open source", "open-source", "oss",
			"github", "gitlab",
			"contributor", "contributions",
			"maintainer", "core maintainer",
			"pull request", "code review",
			"foss", "free software",
			"apache license", "mit license",
			"bsd license",
		},
	},
	"phd_research": {
		Weight: 1.5,
		Label:  "PhD & Research",
		Keywords: []string{
			"phd", "ph.d", "ph.d.",
			"doctorate", "doctoral",
			"dissertation", "thesis",
			"research", "researcher",
			"publication", "published paper",
			"academic paper",
			"conference paper", "journal paper",
			"professor", "postdoc", "post-doc",
			"computer science", "mathematics",
		},
	},
	"operating_systems": {
		Weight: 1.0,
		Label:  "Operating Systems",
		Keywords: []string{
			"operating system", "operating systems",
			"kernel", "linux kernel", "kernel module",
			"linux", "unix", "posix",
			"embedded", "embedded systems",
			"firmware", "rtos",
			"device driver", "syscall", "system call",
			"memory allocator", "scheduler",
		},
	},
}

// CategoryScore holds the scoring breakdown for a single
// category.
type CategoryScore struct {
	// Label is the human-readable category name.
	Label string `json:"label"`

	// Score is the weighted score earned.
	Score float64 `json:"score"`

	// Max is the maximum possible score (the weight).
	Max float64 `json:"max"`

	// Matched lists the keywords that matched.
	Matched []string `json:"matched"`

	// MatchCount is the number of matched keywords.
	MatchCount int `json:"matchCount"`
}

// ScoreResult holds the complete screening score for a
// candidate.
type ScoreResult struct {
	// TotalScore is the sum of weighted category scores.
	TotalScore float64 `json:"totalScore"`

	// MaxPossible is the sum of all category weights.
	MaxPossible float64 `json:"maxPossible"`

	// Pct is the percentage score (0-100).
	Pct float64 `json:"pct"`

	// Categories maps category keys to their breakdowns.
	Categories map[string]CategoryScore `json:"categories"`
}

// ScoreCandidate scores the given text against the Lightning
// Labs screening criteria. The scoring formula is identical to
// the Python implementation: min(matchCount/3, 1.0) * weight.
func ScoreCandidate(text string) ScoreResult {
	textLower := strings.ToLower(text)

	categories := make(
		map[string]CategoryScore, len(criteria),
	)
	var total, maxPossible float64

	for key, crit := range criteria {
		var matched []string
		for _, kw := range crit.Keywords {
			// Keywords are already lowercase in the
			// criteria definitions.
			if strings.Contains(textLower, kw) {
				matched = append(matched, kw)
			}
		}

		// Score: min(matched/threshold, 1.0) * weight.
		// Three or more matches in a category earns the
		// full weight.
		raw := math.Min(
			float64(len(matched))/matchesForFullScore,
			1.0,
		)
		score := math.Round(
			raw*crit.Weight*100,
		) / 100

		categories[key] = CategoryScore{
			Label:      crit.Label,
			Score:      score,
			Max:        crit.Weight,
			Matched:    matched,
			MatchCount: len(matched),
		}

		total += score
		maxPossible += crit.Weight
	}

	var pct float64
	if maxPossible > 0 {
		pct = math.Round(
			total/maxPossible*1000,
		) / 10
	}

	return ScoreResult{
		TotalScore:  math.Round(total*100) / 100,
		MaxPossible: math.Round(maxPossible*100) / 100,
		Pct:         pct,
		Categories:  categories,
	}
}

const (
	// matchesForFullScore is the number of keyword matches
	// needed in a category to earn the full category weight.
	matchesForFullScore = 3.0

	// strongTierThreshold is the minimum percentage score for
	// the "strong" tier.
	strongTierThreshold = 60.0

	// moderateTierThreshold is the minimum percentage score
	// for the "moderate" tier.
	moderateTierThreshold = 35.0

	// weakTierThreshold is the minimum percentage score for
	// the "weak" tier.
	weakTierThreshold = 15.0
)

// ClassifyTier returns a tier classification based on the
// percentage score: strong (>=60), moderate (>=35), weak (>=15),
// or no_signal (<15).
func ClassifyTier(pct float64) string {
	switch {
	case pct >= strongTierThreshold:
		return "strong"
	case pct >= moderateTierThreshold:
		return "moderate"
	case pct >= weakTierThreshold:
		return "weak"
	default:
		return "no_signal"
	}
}

// ExtractText extracts all searchable text from an application
// record represented as a raw JSON map. It recursively walks the
// structure and concatenates all string values.
func ExtractText(app map[string]any) string {
	var parts []string

	// Candidate info.
	if cand, ok := app["candidate"].(map[string]any); ok {
		if name, ok := cand["name"].(string); ok {
			parts = append(parts, name)
		}

		if email, ok := cand["primaryEmailAddress"].(map[string]any); ok {
			if v, ok := email["value"].(string); ok {
				parts = append(parts, v)
			}
		}
	}

	// Custom fields.
	if fields, ok := app["customFields"].([]any); ok {
		for _, f := range fields {
			fm, ok := f.(map[string]any)
			if !ok {
				continue
			}

			if t, ok := fm["title"].(string); ok {
				parts = append(parts, t)
			}
			parts = append(
				parts,
				fmt.Sprintf("%v", fm["value"]),
			)
			if l, ok := fm["valueLabel"].(string); ok {
				parts = append(parts, l)
			}
		}
	}

	// Application form submissions.
	if subs, ok := app["applicationFormSubmissions"].([]any); ok {
		for _, sub := range subs {
			extractFormText(sub, &parts)
		}
	}

	// Source info.
	if src, ok := app["source"].(map[string]any); ok {
		if t, ok := src["title"].(string); ok {
			parts = append(parts, t)
		}
	}

	// Resume filename.
	if resume, ok := app["resumeFileHandle"].(map[string]any); ok {
		if n, ok := resume["name"].(string); ok {
			parts = append(parts, n)
		}
	}

	// Job info.
	if job, ok := app["job"].(map[string]any); ok {
		if t, ok := job["title"].(string); ok {
			parts = append(parts, t)
		}
	}

	return strings.Join(parts, " ")
}

// extractFormText recursively extracts string values from form
// submission data.
func extractFormText(obj any, parts *[]string) {
	switch v := obj.(type) {
	case string:
		*parts = append(*parts, v)

	case map[string]any:
		for _, val := range v {
			extractFormText(val, parts)
		}

	case []any:
		for _, item := range v {
			extractFormText(item, parts)
		}
	}
}

// ExtractTextFromJSON is a convenience helper that unmarshals
// raw JSON into a map and extracts searchable text.
func ExtractTextFromJSON(raw json.RawMessage) string {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}

	return ExtractText(m)
}
