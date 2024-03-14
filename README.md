### [later](https://github.com/rakivo/later) - lightweight, convenient, crossplatform, and self-hosted tool.
##### [later](https://github.com/rakivo/later) allows you to keep all the videos you wanna watch soon without the need to open a lot of unnecessary tabs in your browser.

### To run:
> - Clone and cd the [later repo.](https://github.com/rakivo/later)
```shell
git clone --depth 1 https://github.com/rakivo/later && cd later
```

##### For Linux users:
> You can optionally run ```./scripts/linuxsetpaths.sh```. This script adds the necessary variables(there are only 2 of them) to PATH so that you can run 'later' not only from the project directory but from anywhere after building it. Fun fact: To let your cmd know that you've updated ~/.bashrc file run ```source ~/.bashrc```
##### For Windows users:
> You can optionally try running ```./scripts/windowspaths.bat```. Hovewer, I'm not sure if it will work because I use Arch, btw). (I'm sorry, I'll test it this week).

#### Example of set variables from my ~/.bashrc:
```shell
export PATH=$PATH:/home/rakivo/Coding/later
export LATER_PROJECT_DIR="/home/rakivo/Coding/later/"
```

####  And then finally build and run the project:
```shell
go build -v -ldflags="-s -w" -o later ./src/ # or .\src\ on Windows
./later
```

### Usage:
> Simply paste your link into the input field, click on the submit button and watch it appear on the right side of your screen.

> With a single click on the thumbnail you can open that video in the other window.

### Future plans:
> - Simplify installation
> - Support more platforms

#### References of used dependencies:
> uuid     - https://github.com/google/uuid

> bbolt    - https://github.com/etcd-io/bbolt
