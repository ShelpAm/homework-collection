# Homework Collection

## Quick start

1.   Install [Golang](https://go.dev/)
2.   `git clone https://github.com/ShelpAm/homework-collection`
3.   `cd homework-collection`
4.   `go get .`
5.   `go run src/*.go`

## Nginx

The server by default listens on port 8080. If you want it to listen on port 80 (allowing visiting it without explicitly specifying port number), you can install `nginx` and apply port forwarding via it.

Here is the way to use it.

1.  Add the following to the http section of `/etc/nginx/nginx.conf`. It configures the nginx to do reverse proxy to port 8080.
    ```nginx.conf
            server {
                    listen 80;
                    listen [::]:80;
                    server_name  homework_collection;
    
                    client_max_body_size 20M;
    
                    location / {
                      proxy_pass http://127.0.0.1:8080;
                      proxy_set_header Host $host;
                      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                      proxy_set_header X-Real-IP $remote_addr;
                    }
            }
    ```
2.  Remove the line (or comment out) in `/etc/nginx/nginx.conf`: `include /etc/nginx/sites-enabled/*;`
