# pcd [![Build Status](https://travis-ci.org/kvannotten/pcd.svg?branch=master)](https://travis-ci.org/kvannotten/pcd)

## Philosophy

Downloading and listening to podcasts should be simple. It doesn't require massively complex interfaces that eat all your memory and CPU. 

Pcd is a simple CLI tool that does nothing more than downloading your favorite podcasts. It doesn't run in the background, it doesn't eat all your memory and/or cpu. Everything that needs to be done is your responsability. 

## Why?

I wanted to be able to download my favorite podcasts in a simple way, and on the CLI. I stumbled upon a few utilities like `marrie`. It inspired me to make a version that doesn't need all those annoying python dependencies. Also I wanted to be able to access podcasts that are behind some http authentication method.

## Usage

- You will need to create a configuration file under ~/.config/pcd.yml that has the following options: 
```
---
podcasts:
  - id: 1
    name: biggest_problem
    path: /some/path/to/biggest_problem
    feed: http://feeds.feedburner.com/TheBiggestProblemInTheUniverse
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

## Contributions

Contributions are welcome, as long as they are in line with the philosophy of keeping it simple and to the point. No features that are out of the scope of this application will be accepted.
