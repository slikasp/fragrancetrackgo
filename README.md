# FragranceTrackGo

A GO remake of the python application I made previously, mainly just trying out a new language, but hopefully something I can use keep track of my tried and owned fragrances.

## How to use

To use the app just go to <still working on it>, register, login and look for a fragrance, add to your collection.

If you want to run it yourself:

1. Create appconfig.json file at the root of the project using the _appconfig_template.json file provided:
You will need to add the local database connection for users/scores.
For fragrance database you can use the provided read-only connection to my database or create your own (the the dbManagerFragrances project for the schema)

2. Run goose up migration for local database (schema folder):
goose postgres <user_db_url> up

3. Start the app by <doing what?>

## TODOS

- add web interface (simple)

- add actual user authentication

- refine CLI interface so I could use it to manage the app?

- add proper tests

- add extra security

- find hosting for the web app



## Notes

Using my own remote database for a list of fragrances, it is managed by (github.com/slikasp/dbmanagerfrags). Still need to figure out the best way to keep it up to date.