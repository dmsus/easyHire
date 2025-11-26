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