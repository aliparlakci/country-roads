# Country Roads

# Server
is a REST API developed on Go using MongoDB and Redis.
- [x] Ride CRUD operations
    - [x] Tests
- [x] Location CRUD operations
    - [x] Tests
- [ ] Authentication Flow
    - [x] Signup
    - [ ] (Passwordless) Login
- [ ] Detailed Logging

Start the server with:
```bash
cd server
go get .
go run .
```

## List of design mistakes
1. Validator interface makes queries against database. This logic might be moved to controllers.
2. UserValidator has side effects. It modifies the given email.


# PWA (Progressive Web App)
is a web application developed on React
- [ ] Ride listing page with filtering
- [ ] Ride details page (for sharing post links)
- [ ] Create ride page
- [ ] User profile page
    - [ ] Edit user details

Start the application with:
```bash
cd pwa
yarn
yarn start
```