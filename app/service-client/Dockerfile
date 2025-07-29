# Stage 1: Build the application
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files and install all dependencies
COPY package*.json ./
RUN npm i

# Copy the rest of the application source code
COPY . .

CMD npm run dev
