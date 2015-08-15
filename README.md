# pcd

## Philosophy

Downloading and listening to podcasts should be simple. It doesn't require massively complex interfaces that eat all your memory and CPU. 

Pcd is a simple CLI tool that does nothing more than downloading your favorite podcasts. It doesn't run in the background, it doesn't eat all your memory and/or cpu. Everything that needs to be done is your responsability. 

## Why?

I wanted to be able to download my favorite podcasts in a simple way, and on the CLI. I stumbled upon a few utilities like `marrie`. It inspired me to make a version that doesn't need all those annoying python dependencies. Also I wanted to be able to access podcasts that are behind some http authentication method.

## Usage

- You will need to create a configuration file under ~/.pcd that has the following options: 
```
  ---
  commands:
    player: mplayer
  podcasts:
  - id: 1
    name: biggest_problem
    path: /home/kristof/pods/biggest_problem
    feed: http://feeds.feedburner.com/TheBiggestProblemInTheUniverse
    username:
    password:
  - id: 2
    name: some_other
    path: /home/kristof/pods/some_other
    feed:  http://feeds.example.com/SomeOther.rss
```
- You have to "sync" the feeds: `pcd s`
- (Optionally) List the episodes of a podcast: `pcd l 1` or `pcd l biggest_problem`
- Download the latest episode of `biggest_problem`: `pcd d 1` or `pcd d biggest_problem`
- Play the latest episode with your favorite player by using: `pcd p 1` or `pcd p biggest_problem`

