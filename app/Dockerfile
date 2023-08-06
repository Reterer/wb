##
## STEP 1 - BUILD
##

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.20-alpine AS build

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY . .

# download Go modules and dependencies
RUN go mod download

# compile application
RUN go build -o app main.go

##
## STEP 2 - DEPLOY
##
FROM alpine:3.18

WORKDIR /app

COPY ./static ./static
COPY --from=build ./app .

EXPOSE 8080
ENTRYPOINT [ "./app" ]

