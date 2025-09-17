# --- BUILD UI ---
FROM node:22-alpine as nuxt-build
WORKDIR /app

# Copy package.json and package-lock.json
COPY ./ui/package*.json ./

# Install dependencies
RUN npm ci

# Copy entire project
COPY ./ui ./

# Build the project
RUN npm run generate



# --- BUILD API ---
FROM golang:1.25-alpine AS go-build
WORKDIR /app

# Copy go.mod and go.sum
COPY ./api/go.mod ./api/go.sum ./
RUN go mod download && go mod verify

# Copy entire project
COPY ./api ./

# Build the project
RUN go build -v -o app ./cmd/api



# --- RUN ---
FROM alpine:latest AS run
WORKDIR /root

# Copy ui static files
COPY --from=nuxt-build /app/.output/public ./dist

# Copy api app
COPY --from=go-build /app/app ./
COPY --from=go-build /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Entrypoint
CMD ["./app"]
