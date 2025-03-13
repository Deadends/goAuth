# Project goAuth

This project includes a Go-based authentication system using Google OAuth.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

- Go 1.16 or higher
- Docker
- Make

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Deadends/goAuth
   cd goauth
   ```

2. Create a `.env` file in the `internal/auth` directory with the following content:
   ```properties
   GOOGLE_CLIENT_ID=your-google-client-id
   GOOGLE_CLIENT_SECRET=your-google-client-secret
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Application

We can use the Go make module to run the application, or you can run it manually. For example:

```bash
cd goAuth/cmd/api
go run main.go
```

Check your path before running.

### Frontend Setup

For the frontend, we use React.js with Vite. To spin up the frontend:

```bash
cd goAuth/client
npm install  # Helps to install all required packages
npm run dev
```

### Usage

1. Open your browser and navigate to `http://localhost:5173`.
2. Click on the "Log in with Google" link.
3. After successful authentication, you will be redirected to the frontend with user details.

### Deployment

To deploy the project on a live system, follow these steps:

1. Set up your environment variables for `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET`.
2. Build and run the application using Docker or your preferred method.

### Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

### License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

### Acknowledgments

- [Goth](https://github.com/markbates/goth) for OAuth authentication
- [Chi](https://github.com/go-chi/chi) for the router
- [Gorilla Sessions](https://github.com/gorilla/sessions) for session management
- [Go Blueprint](https://github.com/go-blueprint) for project structure and best practices
- [React.js](https://reactjs.org/) and [Vite](https://vitejs.dev/) for the frontend setup

### Upcoming Releases

We are excited to announce that future releases of `goAuth` will include integration with Keycloak, a powerful open-source identity and access management solution. Stay tuned for future updates and enhancements!
