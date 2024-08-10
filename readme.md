# URL Shortener

A simple url shortener server developed on Golang using [Gin](https://gin-gonic.com/). It uses in-memory data storage. 

### Steps to run
- Clone the repo.
- In the terminal, go the directory where the repo has been cloned and run following command to download the dependencies:
  ```bash
  go mod download 
  ```
- To run the application, run the following command:
  ```bash
  go run main.go 
  ```

### Endpoints

- POST `/shorten`: Create a short url for the given long url.
  - Sample Payload:
    ```json
    {
        "long_url": "https://www.example.com/some/very/long/url",
        "custom_alias": "myalias",  // optional
        "ttl_seconds": 300 //optional
    }
    ```
  - Sample Response:
    - Status: 201 Created
      ```json
      {
          "short_url": "myalias"
      }
      ```
    - Status: 409 Conflict
      ```json
      "alias already exists"
      ```
  **Note:** Default value for `custom_alias` is 6 character long random string and for `ttl_seconds` is 120.

- GET `/:alias`: Redirects to the long url for the given alias.
  - Sample Response:
    - Status: 307 TemporaryRedirect
      Redirects to long url
    - Status: 404 Not Found
      ```json
      "alias does not exist or has expired"
      ```

- GET `/analytics/:alias`: Returns the analytics of alias. Analytics includes the alias, long url, access counts and last 10 access times.
  - Sample Response:
    - Status: 200 OK
      ```json
      {
          "alias": "myalias",
          "long_url": "https://www.example.com/some/very/long/url",
          "ttl_seconds": 300,
          "access_count": 1,
          "access_times": [
              "2024-08-10 16:39:54.547296 +0000 UTC"
          ]
      }
      ```
    - Status: 404 Not Found
      ```json
      "alias does not exist or has expired"
      ```

- PUT `/:alias`: Updates the long url and TTLS of the given alias. Both parameters are optional.
  - Sample Payload: 
    ```json
    {
      "custom_alias": "newalias", // optional
      "ttl_seconds": 90 // optional
    }
     ```
  - Sample Response:
    - Status: 200 OK
      ```json
      "Successfully updated"
      ```
    - Status: 404 Not Found
      ```json
      "alias does not exist or has expired"
      ```

- DELETE `/:alias`: Deletes the given alias.
  - Sample Response:
    - Status: 204 No Content
    - Status: 404 Not Found
      ```json
      "alias does not exist or has expired"
      ```

If the alias is not deleted manually then will automatically get deleted automatically after TTLS has expired.