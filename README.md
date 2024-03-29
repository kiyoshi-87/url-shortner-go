
# URL Shortner

As the name suggests, this is an URL shortner made using Go and Redis as a database.


## Getting Started

You will need Docker and Docker-Compose installed on your system. To use this backend API to generate shortened URLs follow these steps:

1. Install Docker and Docker-Compose on your system. Make sure the service is running.

2. Clone this project repository:

```bash
  git clone https://github.com/kiyoshi-87/url-shortner-go.git
```
3. Navigate to the root directory, where the docker-compose file is present.

4. Go into the API directory and make a .env file with the following environment variables:
      - __DB_ADDRS="db:6379"__ (set this as the value as redis runs on port 6379 of the system)
      - __DB_PASSWD__ = "" (set a password of the database if you want or leave it an empty string)
      - __APP_PORT__ = ":4000" (Set the port at which the api would run on)
      - __DOMAIN__ = "localhost:4000" (This would be replaced with the public ip of the cloud server on which the api would run)
      - __API_QUOTA__ =   (any numerical value, user can make only a given number of API calls in a given period of time)

5. Then go back to the root directory and use the following command: 

```bash
  go mod tidy
```
This would install all the dependencies that are required by the project.

6. After that everything is set and we can start the application using docker compose:

```bash
  docker-compose up -d    
```
![image](https://github.com/kiyoshi-87/url-shortner-go/assets/90674598/84d8235c-e5b8-464e-b1f1-c263ad0f7334)

    
