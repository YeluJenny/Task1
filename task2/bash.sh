#!/bin/bash
csv_file="geoloc-Microsoft.csv"
# Check if the file exists
if [ ! -f "$csv_file" ]; then
    echo "The file $csv_file does not exist."
    exit 1
fi

# Initialize an associative array to hold city counts
declare -A city_counts

# Read the CSV file line by line
while IFS=, read -r ip_range country region city; do
    #validate
    if [ "$country" == "CN" ]; then
        # Increment the city count
        ((city_counts[$city]++))
    fi
done < "$csv_file"
#output
for city in "${!city_counts[@]}"; do
    echo "${city_counts[$city]} $city"
done | sort -k1 -n
