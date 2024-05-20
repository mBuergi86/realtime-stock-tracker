# Use the official Node.js image as the base image
FROM node:alpine

# Create a new directory for the app
RUN mkdir -p /usr/src/app/node_modules && chown -R node:node /usr/src/app

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy the package.json and package-lock.json files to the container
COPY package*.json ./

# Allow the container to run as a non-root user
USER node

# Install the dependencies
RUN npm install

# Copy the source code to the container
COPY --chown=node:node . . 

# Expose the port on which the server will run
EXPOSE 3000

# Set the environment variables
ENV MONGODB_URL="mongodb://localhost:27017/"
ENV MONGODB_DB="stockmarket"
ENV MONGODB_COLLECTION="stocks"
ENV NODE_ENV="development"

# Start the server when the container starts
CMD ["node", "app.js"]
