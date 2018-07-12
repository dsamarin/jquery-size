# jquery-size

This package uses the [GitHub API](https://developer.github.com/v3/repos/#list-tags) to list the [jQuery releases](https://github.com/jquery/jquery/releases), downloads them using the [jQuery CDN](https://code.jquery.com/), and records their sizes into a
[CSV file](https://en.wikipedia.org/wiki/Comma-separated_values).

![image](https://user-images.githubusercontent.com/110829/42649607-010a8e14-85bf-11e8-9130-94a6d7309675.png)

> Disclaimer: it’s not the file size, it’s how you use it. Sure, jQuery has gotten bigger over time, but every new release patches bugs and/or introduces new features. This [project] aims to demonstrate the importance of HTTP compression and minification.
>
> Also note that jQuery 1.8 allows you to create custom builds containing only the modules you need, in case file size is an issue.
> 
> &mdash; [Matthias Bynens](https://mathiasbynens.be/demo/jquery-size)


## Usage

```sh
jquery-size > jquery-size-data.csv
```
