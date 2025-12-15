# Build stage
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code (includes pre-generated types)
COPY . .

# Accept build arguments for Vite environment variables
ARG VITE_API_URL
ARG VITE_WS_URL

# Make them available as environment variables during build
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_WS_URL=$VITE_WS_URL

# Types are already generated and committed - no need to regenerate in Docker
# If you need to regenerate, run `npm run generate` locally before building

# Build the application (skip generate - types are pre-committed)
RUN npm run build:docker

# Runtime stage - serve with nginx
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
