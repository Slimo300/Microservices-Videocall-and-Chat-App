FROM node:alpine as builder
# Set the working directory to /app inside the container
WORKDIR /app
# Copy app files
COPY . .
# Install dependencies (npm ci makes sure the exact versions in the lockfile gets installed)
RUN npm install --omit=dev
RUN npm ci 
# Build the app
RUN npm run build

# Bundle static assets with nginx
FROM nginx:1.21.0-alpine

# Nginx config
RUN rm -rf /etc/nginx/conf.d
COPY nginx-conf /etc/nginx

# Copy built assets from `builder` image
COPY --from=builder /app/build /usr/share/nginx/html

# Expose port
EXPOSE 80

# Copy .env file and shell script to container
WORKDIR /usr/share/nginx/html
COPY ./env.sh .
COPY .env .

# Add bash
RUN apk add --no-cache bash

# Making env.sh executable
RUN chmod +x env.sh

# Run env.sh and start nginx
CMD ["/bin/bash", "-c", "/usr/share/nginx/html/env.sh && nginx -g \"daemon off;\""]