# Ashby API Coverage Gaps

Total endpoints in Ashby API: 150
Currently implemented: ~18

## Currently Implemented

| MCP Tool | Ashby Endpoint |
|----------|---------------|
| list_applications | application.list |
| get_application | application.info |
| change_application_stage | application.changeStage |
| create_application | application.create |
| list_candidates | candidate.list |
| get_candidate | candidate.info |
| search_candidates | candidate.search |
| create_candidate | candidate.create |
| add_candidate_tag | candidate.addTag |
| create_candidate_note | candidate.createNote |
| list_candidate_notes | candidate.listNotes |
| get_file_url | file.info |
| list_jobs | job.list |
| get_job | job.info |
| search_jobs | job.search |
| list_interview_stages | interviewStage.list |
| list_interviews | interview.list |
| list_tags | candidateTag.list |
| pipeline_dashboard | custom |
| screen_candidates | custom |

---

## Priority 1: High Value for Recruiting Workflow

These unblock common recruiting actions we currently cannot perform.

### Application

| Endpoint | Why Needed | Notes |
|----------|-----------|-------|
| application.changeSource | Set/update application source | Needs source.list too |
| application.update | Update application fields | e.g. credit, rejection |
| application.listHistory | Audit trail for application | |
| application.transfer | Move candidate between jobs | |

### Candidate

| Endpoint | Why Needed | Notes |
|----------|-----------|-------|
| candidate.update | Update candidate name/email/phone | |
| candidate.uploadResume | Upload resume to candidate | Needs file upload |
| candidate.uploadFile | Upload arbitrary file | Needs file upload |

### Source / Archive

| Endpoint | Why Needed | Notes |
|----------|-----------|-------|
| source.list | Required for application.changeSource | Simple list |
| archiveReason.list | Required to archive/reject candidates | Simple list |
| closeReason.list | Required to close jobs | Simple list |

### Job

| Endpoint | Why Needed | Notes |
|----------|-----------|-------|
| job.setStatus | Open/close a job | |
| job.update | Update job details | |

### Interview Plan / Stage

| Endpoint | Why Needed | Notes |
|----------|-----------|-------|
| interviewPlan.list | List all interview plans | |
| interviewStage.info | Details on a specific stage | |

### User

| Endpoint | Why Needed | Notes |
|----------|-----------|-------|
| user.list | List team members | Useful for assigning |
| user.info | Look up a specific user | |
| user.search | Search users by name/email | |

### Job Posting

| Endpoint | Why Needed | Notes |
|----------|-----------|-------|
| jobPosting.list | List public job postings | |
| jobPosting.info | Details on a posting | |

### Opening

| Endpoint | Why Needed | Notes |
|----------|-----------|-------|
| opening.list | List openings | |
| opening.info | Opening details | |
| opening.search | Search openings | |

### Department / Location

| Endpoint | Why Needed | Notes |
|----------|-----------|-------|
| department.list | List departments | Needed for job creation |
| department.info | Department details | |
| location.list | List office locations | Needed for job creation |
| location.info | Location details | |

### Hiring Team

| Endpoint | Why Needed | Notes |
|----------|-----------|-------|
| hiringTeam.addMember | Add hiring team member | |
| hiringTeam.removeMember | Remove hiring team member | |
| hiringTeamRole.list | List available roles | |

---

## Priority 2: Useful but Not Blocking

### Offer

| Endpoint | Why Needed |
|----------|-----------|
| offer.list | View offers for candidates |
| offer.info | Offer details |

### Custom Fields

| Endpoint | Why Needed |
|----------|-----------|
| customField.list | See available custom fields |
| customField.setValue | Set a custom field on a record |

### Interview

| Endpoint | Why Needed |
|----------|-----------|
| interview.info | Details on a specific interview |
| interviewSchedule.list | View scheduled interviews |

### Feedback Form

| Endpoint | Why Needed |
|----------|-----------|
| feedbackFormDefinition.list | List available feedback forms |
| feedbackFormDefinition.info | Form details |

---

## Priority 3: Low Value / Out of Scope

These are primarily admin, partner integrations, or rarely needed:

- **Assessments** (assessment.*) -- partner integrations
- **Webhooks** (webhook.*) -- admin setup, not recruiting workflow
- **Surveys** (surveyFormDefinition.*, surveyRequest.*, surveySubmission.*) -- low usage
- **Interviewer pools** (interviewerPool.*) -- admin
- **Location admin** (location.create, location.archive, etc.) -- admin
- **Department admin** (department.create, department.archive, etc.) -- admin
- **Reports** (report.generate, report.synchronous) -- complex, analytics-focused
- **Approval** (approval.list, approvalDefinition.update) -- workflow admin
- **Brand** (brand.list) -- not recruiting
- **Job templates** (jobTemplate.list) -- admin
- **Source tracking links** (sourceTrackingLink.list) -- marketing
- **Referral forms** (referralForm.info) -- niche
- **Application forms** (applicationForm.submit) -- public-facing
- **candidate.anonymize** -- GDPR, niche

---

## Implementation Order

1. **source.list** + **application.changeSource** -- add source management
2. **archiveReason.list** -- needed to reject/archive candidates
3. **candidate.update** -- update candidate fields
4. **user.list** / **user.search** / **user.info** -- team member lookup
5. **department.list** / **location.list** -- supporting lookups
6. **jobPosting.list** / **jobPosting.info** -- job posting visibility
7. **opening.list** / **opening.info** / **opening.search** -- openings
8. **interviewPlan.list** / **interviewStage.info** -- interview pipeline
9. **job.setStatus** / **job.update** -- job management
10. **offer.list** / **offer.info** -- offer tracking
