# Stage 1: Build the Next.js application
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package.json and package-lock.json (or yarn.lock) to install dependencies
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Copy the .env.local file for build-time environment variables
# This is crucial for NEXT_PUBLIC_WEBSOCKET_URL to be picked up
COPY .env.local ./

# Remove Next.js build cache to ensure a clean rebuild
RUN rm -rf .next

# Remove node_modules cache
RUN rm -rf node_modules/.cache

# Build the Next.js application
RUN npm run build

# Stage 2: Create a lightweight production image
FROM node:18-alpine

WORKDIR /app

# Set the production environment
ENV NODE_ENV=production

# Copy built assets and necessary files from the builder stage
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./

# Expose the port Next.js runs on
EXPOSE 3000

# Command to run the application
CMD ["npm", "start"]