# Simple Go Project (UPDATED AT: 08/01/2025)

This simple project was created by: **Kawinpat Raweepornpisarn**

=============================================================================

## Installation & Setup

1. **Install dependencies (Skip if you are already installed)**:
     go mod tidy

2. **Run the project**:
     Using go run: go run ./main.go
     Or using air (if you want live reloading): air

3. **Test the API**:
     No Access token required: http://localhost:8080/
     Access token required (Bearer): http://localhost:8080/auth

     Example paths:
          http://localhost:8080/signin
          http://localhost:8080/auth/signout

=============================================================================

## Notes

Make sure Go 1.x is installed and properly set up on your machine.

If you're using air, ensure it’s installed. 
You can install it via: go install github.com/cosmtrek/air@latest