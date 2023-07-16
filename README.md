# pcd [![Build Status](https://travis-ci.org/kvannotten/pcd.svg?branch=master)](https://travis-ci.org/kvannotten/pcd)

## Philosophy

Downloading and listening to podcasts should be simple. It doesn't require massively complex interfaces that eat all your memory and CPU. 

Pcd is a simple CLI tool that does nothing more than downloading your favorite podcasts. It doesn't run in the background, it doesn't eat all your memory and/or cpu. Everything that needs to be done is your responsibility. 

## Why?

I wanted to be able to download my favorite podcasts in a simple way, and on the CLI. I stumbled upon a few utilities like `marrie`. It inspired me to make a version that doesn't need all those annoying python dependencies. Also I wanted to be able to access podcasts that are behind some http authentication method.

## Installation

### Package managers

Pcd is available on the Arch User Repository (AUR). Use your favorite AUR helper (yay, paru, etc.) to install pcd.

### Binary releases

You can download the latest prebuilt binary from the [releases tab](https://github.com/kvannotten/pcd/releases).

### Building from source

Make sure you have the latest Go compiler installed. You can do so on Arch Linux-based systems using `sudo pacman -S go`.
```
git clone https://github.com/kvannotten/pcd
cd pcd
go build -o pcd cmd/pcd/main.go
sudo mv -f pcd /usr/local/bin
```

## Usage

- You will need to create a configuration file under ~/.config/pcd.yml that has the following options: 
```
---
podcasts:
  - id: 1
    name: biggest_problem
    path: /some/path/to/biggest_problem
    feed: http://feeds.feedburner.com/TheBiggestProblemInTheUniverse
    filenameTemplate: "{{ .rand }}_{{ .title }}{{ .ext }}"
  - id: 2
    name: some_other
    path: /your/podcast/path/to/some_other
    feed:  http://feeds.example.com/SomeOther.rss
    username: foo
    password: bar1234
```
- You have to "sync" the feeds: `pcd sync`
- (Optionally) List the episodes of a podcast: `pcd ls 1` or `pcd ls biggest_problem`
- Download the first episode of `biggest_problem`: `pcd d 1 1` or `pcd d biggest_problem 1`

### Filename template

The `filenameTemplate` configuration entry is a per podcast configuration that allows you to
customize the name of the file. You can add static strings to be shared by all files. Furthermore
the following variables are pushed into the template for your usage:
* `title`: the title of the podcast episode (provided by podcast)
* `name`: the filename parsed from the url (this usually includes the extension)
* `date`: the date provided by the podcast, note that this is unparsed and provided as is. Podcasts use very different formats so there is no uniformity here.
* `current_date`: the current date (when you download it)
* `rand`: a string of 8 random characters
* `ext`: the extension (including the prefix dot)
* `episode_id`: the relative, generated id of the episode. Please note that this is not necessarily idempotent. It depends on the management of the RSS feed.

The `filenameTemplate` is optional. It will default to: `{{ .name }}`

## Support

Community support can be had via the matrix channel: https://matrix.to/#/#pcd:kristof.tech

## Contributions

Contributions are welcome, as long as they are in line with the philosophy of keeping it simple and to the point. No features that are out of the scope of this application will be accepted.
