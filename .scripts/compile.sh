#!/bin/bash

# For each the .schema files, compile them into .go files
for schema in .schema/*.fbs; do
    flatc -o go -g --go-namespace flatdata $schema
    flatc -o python -p $schema
done
