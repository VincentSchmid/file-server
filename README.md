Basic File Upload and Serving in Go
===================================

This repository contains a simple Go application that allows uploading files to a server. The application is containerized using Docker, providing an isolated and reproducible environment.

Getting Started
---------------

Follow the steps below to build and run the Go application using Docker:

1.  Clone the repository:
    ```shell
    git clone git@github.com:VincentSchmid/file-server.git
    cd file-server
    ```

2.  Build the Docker image:
    ```shell
    docker build -t file-server .
    ```
3. Generate an api-key:
    ```shell
    API_KEY=$(openssl rand -base64 32)
    ```

4.  Run the Docker container:
    ```shell
    docker run -p 8080:8080 -e API_KEY=$API_KEY -v $(pwd)/uploads:/app/uploads file-server
    ```

    This command runs a Docker container from the `file-server` image. It exposes port 8080 on the host machine and maps it to port 8080 inside the container.

    The `-v` option mounts the `uploads` directory on the host machine to the `/app/uploads` directory inside the container. This allows uploaded files to be stored on the host machine.

5.  Access the application:  

    Open your web browser and visit `http://localhost:8080/files` to access the upload page. You can use tools like cURL or web browsers to upload files to the server.

Upload and File access
----------------------

files can be uploaded using the api key header:
```shell
curl -X POST -H "X-API-Key: $API_KEY" -F "file=@test.txt" http://localhost:8080/upload
```

and the files can be accessed browsed:  
`http://localhost:8080/files`

and downloaded using:  
`http://localhost:8080/files/test.txt`

Additional Notes
----------------

*   `.env`: (Optional) Environment file to store the API key. Not required if passing the API key as an environment variable directly.
*   Uploaded files will be stored in the `uploads` directory on the host machine due to the volume mount configuration `-v $(pwd)/uploads:/app/uploads`. Adjust the path as needed.
