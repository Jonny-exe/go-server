# The base go-image
FROM golang:1.14-alpine
 
# Create a directory for the app
RUN mkdir /app
 
# Copy all files from the current directory to the app directory
COPY . /app
 
# Set working directory
WORKDIR /app/httpd

ENV GO_MESSAGES_DIR /app/httpd/
# Run command as described:
# go build will build an executable file named server in the current directory
RUN go build . 
 
# Run the server executable
CMD [ "/app/httpd/httpd" ]

EXPOSE 5000