#!/bin/bash

echo "Generating sitemap..."

cd ./src/frontend || exit
node generate-sitemap.js
cd ../../
