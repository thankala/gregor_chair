# GREGOR Chair

GREGOR Chair is an open-source cyber-physical system designed to orchestrate and manage smart chair components. Built with Go, it leverages a modular architecture to support distributed control, real-time data processing, and seamless integration with IoT devices.
Features

    Modular Design: Easily extend or replace components.

    Real-Time Orchestration: Manage smart chair functions dynamically.

    Go Modules: Utilizes go.mod for dependency management.

    Submodule Integration: Incorporates shared logic via Git submodules.

## Getting Started
## Clone the Repository
``` 
git clone https://github.com/thankala/gregor-chair.git 
```
Navigate to the project directory and install Go dependencies:

cd gregor-chair
./update.sh

## Usage

First bring up the containers:

```
docker compose up -f local-docker-compose.yml up -d
```

Run the application with:


```
./local/local-coordinator.go
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
