#!/bin/bash

echo "Welcome to Akami"
echo "Note: This tool is still in beta version."
echo "Update: Low CPU usage, more powerful."
echo

# Input validation for target URL
read -p "Enter target URL: " target_url
if [[ -z "$target_url" ]]; then
    echo "Error: Target URL cannot be empty."
    exit 1
fi

# Input validation for number of workers
read -p "Enter number of workers: " workers
if [[ -z "$workers" || ! "$workers" =~ ^[0-9]+$ ]]; then
    echo "Error: Number of workers must be a valid number."
    exit 1
fi

# Input validation for duration
read -p "Enter duration (in seconds): " duration
if [[ -z "$duration" || ! "$duration" =~ ^[0-9]+$ ]]; then
    echo "Error: Duration must be a valid number."
    exit 1
fi

echo

# Run the Go program with the provided parameters
go run main.go "$target_url" "$workers" "$duration"