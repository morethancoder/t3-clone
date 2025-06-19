# Gunter

A simple project setup using PocketBase and Go. Follow the steps below to get the project up and running.

---

## Getting Started

### Prerequisites
- [Go](https://go.dev/) installed on your system
- [PocketBase](https://pocketbase.io/) installed and running
- Google OAuth provider details (Client ID and Client Secret)

---

### Installation

1. **Clone the Repository**  
   Clone the project repository to your local machine:
   ```bash
   git clone <repository-url>
   cd gunter
   ```

2. **Install Dependencies**  
   Run the following command to install all required dependencies:
   ```bash
   go mod tidy
   ```

3. **Set Up PocketBase Schema**  
   - Open your PocketBase admin panel.
   - Load the schema from the `pb_schema.json` file into PocketBase.  
     (This file contains the necessary collections and fields for the project.)

4. **Run PocketBase**  
   You can start the PocketBase server using the provided `Makefile`. Simply run:
   ```bash
   make db
   ```

5. **Configure Google OAuth**  
   - Add your Google OAuth provider details (Client ID and Client Secret) in the PocketBase admin panel under the "Auth Providers" section.

---

### Running the Project

Once the setup is complete, everything should work as expected. Start your PocketBase server using `make db` and run the project.

---

## Notes

- Ensure that your PocketBase server is running before interacting with the project.
- If you encounter any issues, double-check the schema and OAuth provider configuration.

---

Enjoy building with **Gunter**! 🎉


## Docker

### How to Use
1. Build the Docker Image
```bash
cd /path/to/your/t3-clone
docker build -t morethancoder/t3-clone .
```
2. Run with Default Settings
```bash
docker run -p 8080:8080 -p 8090:8090 morethancoder/t3-clone
```
3. Run with Custom Environment Variables
```bash
docker run -p 8080:8080 -p 8090:8090   -e OPENROUTER_API_KEY=your_actual_api_key   -e ENV=production   morethancoder/t3-clone
```
4. Run with Persistent Database
```bash
docker run -p 8080:8080 -p 8090:8090   -v $(pwd)/pb_data:/root/pb_data   -e OPENROUTER_API_KEY=your_actual_api_key   morethancoder/t3-clone
```
5. Run in Background (Detached Mode)
```bash
docker run -d -p 8080:8080 -p 8090:8090   --name t3-clone-app   -e OPENROUTER_API_KEY=your_actual_api_key   morethancoder/t3-clone
```

#### Access Your Application
```

Main app: http://localhost:8080
Database server: http://localhost:8090
```
#### Stop the Container
```bash
docker stop t3-clone-app
docker rm t3-clone-app
```
