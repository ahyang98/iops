# Using official python runtime base image
# FROM python:3.9-slim
# FROM debian:latest
FROM ubuntu:20.04

# add curl for healthcheck
# RUN apt-get update \
#     && apt-get install -y --no-install-recommends python3 curl

RUN apt update \
    && apt install -y software-properties-common \
    && add-apt-repository ppa:deadsnakes/ppa \
    && apt install -y python3.9 \
    && apt install -y python3-pip

# Set the application directory
WORKDIR /app

# Install our requirements.txt
COPY requirements.txt /app/requirements.txt
RUN pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple \
    && pip install -r requirements.txt

# Copy our code from the current folder to /app inside the container
COPY . .

# Make port 80 available for links and/or publish
EXPOSE 80

# Define our command to be run when launching the container
CMD ["gunicorn", "app:app", "-b", "0.0.0.0:80", "--log-file", "-", "--access-logfile", "-", "--workers", "4", "--keep-alive", "0"]
