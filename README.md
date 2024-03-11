### The [later](https://github.com/rakivo/later) is lightweight, convenient and self-hosted tool.
##### The [later](https://github.com/rakivo/later) allows you to keep all the videos you wanna watch soon without the need to open a lot of unnecessary tabs in your browser.

### To run:
> - Clone and cd the [later repo.](https://github.com/rakivo/later)
```shell
git clone --depth 1 https://github.com/rakivo/later && cd later
```
> - Create a .env file with the following content:
```txt
YOUTUBE_API_KEY = "PUT YOUR API KEY BETWEEN THIS QUOTATION MARKS"
```
> - You need to obtain free [YouTube Data API v3 api key](https://developers.google.com/youtube/v3), follow this [guide from Google](https://developers.google.com/youtube/v3/getting-started). This is needed because there's no way(AFAIK) to get the title of a YT video without an API key
> - After that just build & run the project
```shell
go build -v -ldflags="-s -w" -o later
./later
```

### Usage:
> Simply paste your link into the input field, click on the submit button and watch it appear on the right side

> With a single click on the thumbnail you can open that video in the other window

### Future Plans:
> 1.1 proper frontend

> 1.2  proper readme

> 2.1 get rid of using youtube api

> 2.2 get rid of using gin

> 3.1 gif preview

#### References of used dependencies:
> gin      - https://github.com/gin-gonic/gin

> uuid     - https://github.com/google/uuid

> bbolt    - https://github.com/etcd-io/bbolt

> godotenv - https://github.com/joho/godotenv
