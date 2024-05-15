# Go Books

## Introduction
This is an attempt to use Google API to creat a simple interface and minimal dependencies to search and download files from personal drive storage(Google Drive). I used a cache file to minimise api calls but included a button to update cache file when needed.

The UI is written in Svelte and embedded into the go code.

 The application can be build by simply running `make` in the main folder and then running the executable.



## Dependencies
Go Packages
```
golang.org/x/oauth2 v0.20.0<br>
google.golang.org/api v0.180.0<br>
```
Svelte
```
"@sveltejs/vite-plugin-svelte": "^3.0.2"
"svelte": "^4.2.12",
"vite": "^5.2.0"
```