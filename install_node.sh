#!/bin/bash

# Install Node.js and npm
curl -fsSL https://deb.nodesource.com/setup_16.x | bash -
apt-get install -y nodejs

# Verify installation
node -v
npm -v