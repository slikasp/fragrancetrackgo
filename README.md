# FragranceTrackGo

A GO remake of the python application I made previously, used it to learn a new language, but also something I can use keep track of my tried and owned fragrances and display them when needed.

It is functioning, but there are still a lot of improvements that I'm working on.


## How to use

<ul> To check how it works go to <still working on hosting it>
<li> Register & login
<li> Search for a fragrance you own and add it to your list with your personal score/comment
<li> View and filter your fragrance list
<li> More functionality is being worked on
</ul>


<ol> If you want to run it yourself:

<li> Clone the project:

`git clone https://github.com/slikasp/fragrancetrackgo`

<li> Create appconfig.json file at the root of the project using the _appconfig_template.json file provided:
You will need to add the local database connection for users/scores.
For remote fragrance database you can use the provided read-only connection to my database or create your own (the the dbManagerFragrances project for the schema)

<li> Run goose up migration for local database (schema folder):

`goose postgres <local_db_url> up`

<li> Start the app:

`go run .`

<li> Visit the website that starts

</ol>

## TODOS


<ul> TO DO:
<li> refine the web interface, add functionality
<li> add proper tests for internals (non db/web functions)
<li> move local database to remote
<li> add functionality for your fragrance list (better sort, most owned/liked notes/accords, coverage of seasons, suggestions)
<li> find hosting for the web app
<li> add extra security
</ul>


## Notes

Used my own remote database for a list of fragrances, it is managed by (github.com/slikasp/dbmanagerfrags). Still need to figure out the best way to keep it up to date and download images.

<ul> Use of AI:
<li> Most of the HTML templates
<li> DB handler interfaces
<li> Some parts of Web handlers
<li> Refactoring multiple files
</ul>