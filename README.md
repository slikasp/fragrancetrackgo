# FragranceTrackGo

A GO remake of the python application I made previously. It is functioning, but there are still a lot of improvements that I'm working on.

## Motivation

I recently discovered a world of fragrances and wanted to keep track of my collection and also fragrances I have tried, after using notes and Google Sheets I thought it's time I learn a new language and make a proper app for myself.

## Quick Start

<ul> To check how it works go to <still working on hosting it>
<li> Register & login
<li> Search for a fragrance you own and add it to your list with your personal score/comment
<li> View and modify your fragrance list
<li> More functionality is being worked on
</ul>

## Usage

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

## Contributing

Currently I plan to develop it myself and learn more doing it.

## TODOS

<ul> TO DO:
<li> add proper tests for internals (non db/web functions)
<li> make and implement card downloader and add card display (probably in the fragrance manager script)
<li> move local database to remote
<li> add functionality for your fragrance list (better sort, most owned/liked notes/accords, coverage of seasons, suggestions)
<li> find hosting for the web app
<li> add extra security
</ul>


## Notes

Used my own remote database for a list of fragrances, it is managed by. Still need to figure out the best way to keep it up to date and download images.

<ul> Use of AI:
<li> Most of the HTML templates
<li> DB handler interfaces
<li> Some parts of Web handlers
<li> Refactoring multiple files
</ul>