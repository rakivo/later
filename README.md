## The [later](https://github.com/rakivo/later) is lightweight, convenient and self-hosted tool that allows you to keep all the videos you wanna watch soon without the need to open a lot of unnecessary tabs in your browser.

### To run:
> - Clone and cd the [later repo](https://github.com/rakivo/later)
```shell
git clone --depth 1 https://github.com/rakivo/later && cd later
```
> - Create a .env file with the following content::
```txt
YOUTUBE_API_KEY = "PUT YOUR API KEY BETWEEN THIS QUOTATION MARKS"
```
> - You need to obtain free [YouTube Data API v3 api key](https://developers.google.com/youtube/v3), follow this [guide from Google](https://developers.google.com/youtube/v3/getting-started). This is needed because there's no way(AFAIK) to get the title of a YT video without an API key

### Usage:
> Simply paste your link into the input field and watch it appear on the right side, with a single click on thumbnail you can that video in the other window

### Future Plans:
> I plan to expand support beyond just [YouTube](https://www.youtube.com/) videos to include other platforms and content types.
