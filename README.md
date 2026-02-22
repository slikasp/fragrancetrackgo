# FragranceTrackGo

A GO remake of the python application I made previously, mainly just trying out a new language, but hopefully something I can use keep track of my tried and owned fragrances.

## How to use

To use the app just go to <still working on it>, register, login and look for a fragrance, add to your collection.

If you want to run it yourself:

1. Add config in appconfig.json at the root of the project with these contents:
{
    "user_db_url":"<your_local_database_url>",
    "fragrance_db_url":"postgresql://ftg_readonly.rqjghqqmmmzenutmzveq:menulis@aws-1-eu-central-1.pooler.supabase.com:5432/postgres",
}

2. Run goose <add exact command> Up migration for local database

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