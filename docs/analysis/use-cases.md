# Use Cases

## Use Case: HR Creates and Assigns Assessment

**Actors:** HR Specialist, System, AI Service

**Trigger:** HR needs to screen candidates for a Go developer position

**Preconditions:**
- HR is authenticated and has necessary permissions
- Competency matrix is configured in the system

**Main Success Scenario:**
1. HR navigates to assessment creation page
2. System displays competency matrix with levels
3. HR selects required competencies (Go, Algorithms, Databases)
4. HR sets difficulty level (Middle) and duration (60 minutes)
5. System generates assessment using AI with 20 questions
6. HR imports candidate list (names and emails)
7. System sends invitation emails with unique test links
8. Candidates receive notifications and begin assessments

**Postconditions:**
- Assessment is created and available for candidates
- HR can monitor progress in real-time
- Results are automatically compiled when candidates complete

**Alternative Flows:**
- AI service unavailable: System uses predefined question bank
- Candidate email invalid: System flags and allows correction

## Use Case: Candidate Takes Assessment

**Actors:** Candidate, System, Code Execution Service

**Trigger:** Candidate receives assessment invitation and clicks link

**Preconditions:**
- Candidate has valid assessment link
- Assessment is active and not expired
- Candidate has not completed assessment before

**Main Success Scenario:**
1. Candidate opens assessment link
2. System verifies access and displays instructions
3. Candidate starts assessment timer
4. System presents questions one by one with code editor
5. Candidate writes code and runs test cases
6. System executes code in secure container and provides feedback
7. Candidate completes all questions and submits assessment
8. System calculates scores and generates report
9. Candidate receives immediate results with detailed feedback

**Postconditions:**
- Assessment marked as completed
- Results available to candidate and HR
- Candidate receives skill analysis report

**Alternative Flows:**
- Time expires: System auto-submits with current progress
- Technical issues: System allows resume from last save point
- Code execution failure: System provides error details and allows retry

## Use Case: Technical Expert Reviews AI-Generated Questions

**Actors:** Technical Expert, System, AI Service

**Trigger:** AI generates batch of new questions for review

**Preconditions:**
- Technical expert is authenticated with appropriate permissions
- AI service has generated new questions
- Questions are in "pending review" status

**Main Success Scenario:**
1. System notifies technical expert of pending questions
2. Expert reviews question content, code examples, and test cases
3. Expert approves, rejects, or edits questions
4. System updates question status and adds to question bank if approved
5. Approved questions become available for assessments
6. System logs expert feedback for AI training

**Postconditions:**
- Questions are validated and available for use
- AI model receives feedback for improvement
- Question quality metrics are updated

**Alternative Flows:**
- Question requires major changes: Expert sends back to AI with detailed feedback
- Multiple experts disagree: System escalates to senior expert for final decision

## Use Case: Admin Configures System Integration

**Actors:** Admin, System, External Services

**Trigger:** Company wants to integrate with existing HR systems

**Preconditions:**
- Admin has system configuration permissions
- External system APIs are available and documented

**Main Success Scenario:**
1. Admin navigates to integration settings
2. System displays available integration options (ATS, SSO, etc.)
3. Admin configures API endpoints and authentication
4. System tests connection to external services
5. Admin enables data synchronization settings
6. System begins syncing candidate and assessment data

**Postconditions:**
- Integration is active and functional
- Data flows between systems automatically
- Admin receives integration status reports

**Alternative Flows:**
- Authentication fails: System provides detailed error and troubleshooting steps
- API version mismatch: System suggests compatible versions or workarounds

## Use Case: AI Generates Adaptive Assessment

**Actors:** AI Service, System, Question Bank

**Trigger:** HR requests assessment for specific role and level

**Preconditions:**
- Competency matrix is defined and populated
- Question bank has sufficient questions for required competencies
- AI service is available and responsive

**Main Success Scenario:**
1. System receives assessment parameters (role, level, competencies)
2. AI analyzes competency requirements and difficulty level
3. AI retrieves relevant questions from question bank using RAG
4. AI generates new questions to fill gaps using competency context
5. System validates question distribution and time estimates
6. Assessment is compiled and made available for assignment

**Postconditions:**
- Custom assessment is created and ready for use
- Questions are balanced across competencies and difficulty
- Assessment meets time and quality requirements

**Alternative Flows:**
- Insufficient base questions: AI generates more questions to meet requirements
- Quality validation fails: System flags for technical expert review