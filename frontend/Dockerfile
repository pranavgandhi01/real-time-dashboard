# Stage 1: Build the Next.js application
FROM node:18-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
# Using package*.json ensures both package.json and package-lock.json are copied
COPY package*.json ./

# Install dependencies
# This installs 'next' and all other dependencies into node_modules
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the application (generates .next folder)
RUN npm run build

# Stage 2: Create a production-ready image
FROM node:18-alpine

# Set the working directory
WORKDIR /app

# Copy the installed node_modules from the builder stage
# This ensures 'next' and its dependencies are available for runtime
COPY --from=builder /app/node_modules ./node_modules

# Copy the built Next.js application files
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
# Copy package.json to ensure npm start works correctly, even though node_modules is present
COPY --from=builder /app/package.json ./

# Expose the port Next.js listens on
EXPOSE 3000

# Set the command to start the server
# 'npm start' will now correctly use the 'next' binary from the copied node_modules
CMD ["npm", "start"]