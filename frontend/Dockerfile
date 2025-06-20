# Stage 1: Build the Next.js application
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package.json and package-lock.json (or yarn.lock) to install dependencies
COPY frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY frontend/ .

# Copy environment files for build-time variables
COPY .env.frontend ./.env.local

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
COPY --from=builder /app/.env.local ./

# Expose the port Next.js runs on
EXPOSE 3000

# Command to run the application
CMD ["npm", "start"]