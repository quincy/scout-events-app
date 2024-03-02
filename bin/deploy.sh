#!/bin/bash

# create the app in fly.io
fly launch --no-deploy --dockerfile Dockerfile

# deploy it
fly deploy --strategy canary
