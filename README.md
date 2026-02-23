# FragranceTrackGo

A GO remake of the python application I made previously, mainly just trying out a new language, but hopefully something I can use keep track of my tried and owned fragrances.

## How to use

To check how it works go to <still working on hosting it>

If you want to run it yourself:

0. Clone the project
`git clone https://github.com/slikasp/fragrancetrackgo`

1. Create appconfig.json file at the root of the project using the _appconfig_template.json file provided:
You will need to add the local database connection for users/scores.
For remote fragrance database you can use the provided read-only connection to my database or create your own (the the dbManagerFragrances project for the schema)

2. Run goose up migration for local database (schema folder):
`goose postgres <local_db_url> up`

3. Start the app
`go run .`

4. Visit the website that starts

## TODOS

- refine the web interface, add functionality

- add proper tests

- add extra security

- find hosting for the web app



## Notes

Using my own remote database for a list of fragrances, it is managed by (github.com/slikasp/dbmanagerfrags). Still need to figure out the best way to keep it up to date.