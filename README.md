# Go Waitlist API

Simple API for double opt-in waitlist signups using Go, Gin, GORM, mariadb and [Resend](https://resend.com) for email sending.

Provides the following endpoints:
- `POST /v1/waitlist` with JSON body: `{"email": "mail@example.com", "additionalInformation": {"firstName": "John"}}` (additionalInformation is optional)
- `POST /v1/waitlist/opt-in` with JSON body: `{"email": "mail@example.com", "confirmationToken": "[token provided in email]"}`
- `POST /v1/waitlist/opt-out` with JSON body: `{"email": "mail@example.com"}`

## Setup

1. Clone the repository
2. Copy the .env.sample (`cp .env.sample .env`) and adjust its contents (please take a look at their usages, a more detailed documentation will be provided soon)
3. Make sure docker (with compose) is installed
4. Build the Go app and image using `docker compose build`
5. Run the app and DB using `docker compose up -d` (leave out -d if you don't want it to run in background, e.g. for debugging)

## Road Map

- More detailed documentation on setup and environment variables
- Export of registrations to JSON and an authorized admin API endpoint
- Support second email after confirmation
- Support manual email sending without Resend, using SMTP connection
