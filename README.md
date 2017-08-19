# Purpose

This application serves to generate graphs based on transaction values read from a CSV file.  After
reading the CSV[s], the transaction logs are loaded into a sqlite3 database.

# CSV Format

A single line has the following form:

`status,,date,,business name,category,amount`

with the blank fields being reserved for future use.

