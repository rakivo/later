### [later](https://github.com/rakivo/later) - lightweight, convenient, crossplatform, and self-hosted tool.
##### [later](https://github.com/rakivo/later) allows you to keep all the videos you wanna watch soon without the need to open a lot of unnecessary tabs in your browser.

### To run:
> - Clone and cd the [later repo.](https://github.com/rakivo/later)
```shell
git clone --depth 1 https://github.com/rakivo/later && cd later
```
> - You need to obtain free [YouTube Data API v3 api key](https://developers.google.com/youtube/v3), follow this [guide from Google](https://developers.google.com/youtube/v3/getting-started). This is needed because there's no way(AFAIK) to get the title of a YT video without an API key

##### For Linux users:
> Add ```LATER_YOUTUBE_API_KEY``` to your ~/.bashrc and set it to your actual YouTube API key
> You can optionally execute ./linuxsetpaths.sh. This script adds the necessary variables(there are only 2 of them) to PATH so that you can run 'later' not only from the project directory but from anywhere after building it. Fun fact: To let your cmd know that you've updated ~/.bashrc file run ```source ~/.bashrc```
##### For Windows users:
> Add ```LATER_YOUTUBE_API_KEY``` to your PATH variables and set it to your actual YouTube API key
> You can optionally try running windowspaths.bat. Hovewer, I'm not sure if it will work because I use Arch, btw). (I'm sorry, I'll test it this week).

> Example of set variables from my ~/.bashrc:
```shell
export PATH=$PATH:/home/rakivo/Coding/later
export LATER_PROJECT_DIR="/home/rakivo/Coding/later/"
```

> And then finally build and run the project:
```shell
go build -v -ldflags="-s -w" -o later
./later
```

### Usage:
> Simply paste your link into the input field, click on the submit button and watch it appear on the right side

> With a single click on the thumbnail you can open that video in the other window

### Main goal:
> Get rid of using YouTube API

#### References of used dependencies:
> uuid     - https://github.com/google/uuid

> bbolt    - https://github.com/etcd-io/bbolt
