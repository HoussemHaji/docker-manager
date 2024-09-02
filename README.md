This is a Golang application designed for Docker container management.

## Features

- **List Containers:** Quickly view all Docker containers on your system.
- **Start/Stop/Pause/Unpause/Delete Containers:** Take control of your containers with simple commands to manage their lifecycle.
- **Execute Commands:** Seamlessly execute commands within containers, enhancing your control and interaction.
- **Display Logs:**  Quickly view container logs
- **Filter containers:** Filter containers by name or status

## Requirements
- docker
- go


### Run the GO APP
1. Clone the repository
    ```shell
    git clone https://github.com/nicumicle/go-docker.git
    cd go-docker
    ```
2. Run the app
    ```shell
    go run main.go
    ```

### Build from source
1. Clone the repository
    ```shell
    git clone https://github.com/HoussemHaji/docker-manager
    cd docker-manager
    ```
2. Build the app
    ```shell
    go build -o out/docker-manager main.go
    ```
3. Run the application
- Windows:
   ```shell
   .\out\docker-manager.exe
   ```
- Linux / macOS:
   ```shell
   ./out/docker-manager
   ```
