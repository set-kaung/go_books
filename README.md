# Go Books

## Introduction
This is an attempt to use Google API to creat a simple interface and minimal dependencies to search and download ebooks files (.pdf and .epub) from personal drive storage(Google Drive). ~~I used a cache file to minimise api calls but included a button to update cache file when needed.~~ Implements a cache file as well as sqlite database to reduces api calls.

 Also a way to find duplicate books in drive.

The UI is written in Svelte and embedded into the go code.

The application can be build by simply running `make` in the main folder and then running the executable.


## Dependencies
Go Packages
```
golang.org/x/oauth2 v0.20.0
google.golang.org/api v0.180.0
modernc.org/sqlite v1.32.0
```
Svelte
```
"@sveltejs/vite-plugin-svelte": "^3.0.2"
"svelte": "^4.2.12",
"vite": "^5.2.0"
```