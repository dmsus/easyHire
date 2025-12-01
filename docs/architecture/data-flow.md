# Data Flow Diagrams

## Assessment Creation Flow

┌──────────┐ 1. Create Assessment ┌──────────┐
│ HR │ ─────────────────────────► │ Frontend │
└──────────┘ └──────────┘
│
│ 2. API Request
▼
┌──────────┐
│ Backend │
│ API │
└──────────┘
│
┌────────────────────────────┼─────────────────────────┐
│ 3. Generate Questions │ 4. Store Assessment │
▼ ▼ │
┌──────────────┐ ┌────────────────┐ │
│ AI Service │ │ Database │ │
│ │ │ │ │
└──────────────┘ └────────────────┘ │
│ │ │
│ 5. Return Questions │ 6. Return Assessment ID │
▼ ▼ │
┌──────────────┐ ┌────────────────┐ │
│ Questions │ │ Assessment │ │
│ with IDs │ │ Record │ │
└──────────────┘ └────────────────┘ │
│ │ │
└─────────────┬──────────────┘ │
│ 7. Send Invitations │
▼ │
┌──────────────┐ │
│ Email Service│ │
│ │ │
└──────────────┘ │
│ │
▼ │
┌──────────────┐ │
│ Candidates │ │
│ Notified │ │
└──────────────┘ │
│
8. Return Success Response │
▼ │
┌──────────────┐ │
│ HR │◄──────────────────────────────┘
│ Dashboard │
└──────────────┘


## Candidate Assessment Flow

┌─────────────┐ 1. Click Link ┌─────────────┐
│ Candidate │ ──────────────────► │ Frontend │
└─────────────┘ └─────────────┘
│
│ 2. Verify & Start
▼
┌─────────────┐
│ Backend API │
└─────────────┘
│
┌──────────────────────────────┼──────────────────────────┐
│ 3. Get Questions │ 4. Start Timer │
▼ ▼ │
┌─────────────┐ ┌─────────────┐ │
│ Database │ │ Redis Cache │ │
│ │ │ │ │
└─────────────┘ └─────────────┘ │
│ │ │
│ 5. Return Questions │ 6. Return Session │
▼ ▼ │
┌─────────────┐ ┌─────────────┐ │
│ Questions │ │ Session │ │
│ Displayed │ │ Started │ │
└─────────────┘ └─────────────┘ │
│ │ │
│ 7. Answer Questions │ 8. Auto-save Progress │
▼ ▼ │
┌─────────────┐ ┌─────────────┐ │
│ Code Editor │ │ Database │ │
│ & UI │ │ │ │
└─────────────┘ └─────────────┘ │
│ │ │
│ 9. Submit Code │ │
▼ │ │
┌─────────────┐ │ │
│ Code Exec. │ │ │
│ Service │ │ │
└─────────────┘ │ │
│ │ │
│ 10. Return Results │ │
▼ │ │
┌─────────────┐ │ │
│ Immediate │ │ │
│ Feedback │ │ │
└─────────────┘ │ │
│ │ │
│ 11. Complete Assessment │ │
▼ ▼ │
┌─────────────┐ ┌─────────────┐ │
│ Scoring │ │ Finalize │ │
│ Engine │ │ Session │ │
└─────────────┘ └─────────────┘ │
│ │ │
│ 12. Calculate Score │ 13. Store Results │
▼ ▼ │
┌─────────────┐ ┌─────────────┐ │
│ Fibonacci │ │ Database │ │
│ Scoring │ │ │ │
└─────────────┘ └─────────────┘ │
│ │ │
│ 14. Generate Report │ 15. Send Notifications │
▼ ▼ │
┌─────────────┐ ┌─────────────┐ │
│ Results │ │ Email & │ │
│ Displayed │ │ Webhooks │ │
└─────────────┘ └─────────────┘ │
│ │
│ 16. Notify HR │
▼ │
┌─────────────┐ │
│ HR │◄────────────────┘
│ Dashboard │
└─────────────┘

## AI Question Generation Flow (RAG)

┌─────────────┐ 1. Request Question ┌─────────────┐
│ Backend API │ ──────────────────────► │ AI Service │
└─────────────┘ └─────────────┘
│
│ 2. Retrieve Similar
▼
┌─────────────┐
│ Vector DB │
│ (pgvector) │
└─────────────┘
│
│ 3. Return Context
▼
┌─────────────┐
│ RAG Engine │
│ │
└─────────────┘
│
│ 4. Build Prompt
▼
┌─────────────┐
│ Prompt │
│ Engine │
└─────────────┘
│
│ 5. Call AI Model
▼
┌─────────────┐
│ Gemini │
│ API │
└─────────────┘
│
│ 6. Return Generated
▼
┌─────────────┐
│ Question │
│ Generator │
└─────────────┘
│
│ 7. Validate & Score
▼
┌─────────────┐
│ Quality │
│ Scorer │
└─────────────┘
│
│ 8. Store Embedding
▼
┌─────────────┐
│ Vector DB │
│ │
└─────────────┘
│
│ 9. Return Question
▼
┌─────────────┐
│ Backend API │◄──┐
└─────────────┘ │
│ │
│ 10. Queue │
▼ │
┌─────────────┐ │
│ Expert │───┘
│ Review │
└─────────────┘


## Data Retention & Archiving

┌─────────────────────────────────────────────────┐
│ Active System │
│ ┌─────────────┐ Daily ┌─────────────┐ │
│ │ Hot Data │ ─────────────► │ Warm Data │ │
│ │ (<30d) │ Move Old │ (30-365d) │ │
│ └─────────────┘ └─────────────┘ │
│ Monthly │
│ Move │
│ ▼ │
│ ┌─────────────┐ │
│ │ Cold Data │ │
│ │ (>1 year) │ │
│ └─────────────┘ │
│ │ │
│ After 3 years │
│ ▼ │
│ ┌─────────────┐ │
│ │ Archived │ │
│ │ Data │ │
│ │ (Anonymized)│ │
│ └─────────────┘ │
└─────────────────────────────────────────────────┘

