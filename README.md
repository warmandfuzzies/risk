#Risky Plumbers

##Question

Risk application:
1. Listens on port 8080 for standard HTTP traffic (not HTTPS)
2. Can return a list of
Risk given a GET to /v1/risks on the endpoint; and
3. Can create a new Risk given a POST to /v1/risks on the endpoint; and
4. Can retrieve an individual Risk given a GET to /v1/risks/<id> .

A Risk should consist of:
A Risk ID in the form of a UUID.
This ID should be auto-generated on creation.
A state value as a string
This can be one of [open, closed, accepted, investigating]
MUST be present for all Risks.
A Risk title as a string.
A Risk description as a string.
Arctic Wolf

Data transfer should be done in JSON.
The storage of risks can be done in memory; no database creation or usage is
required.
The endpoints should use standard HTTP response status codes. E.g.:
200 OK for a successful GET to the /v1/risks/ or /v1/risks/<id> endpoints
500 Internal Server Error for problem with the internal server; etc.
You may use any Golang framework or library to assist in the creation of this app.
Please include:
Instructions via a README.md or README.adoc file including how to run your
service and any tests you may have written.
In a For-Reviewers.md file, include any notes, thoughts, etc. which you
would
like the interviewers to be aware of.
