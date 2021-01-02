# Using Open Emoji API service provider example

* Check whether a string contains an emoji.
* Get all available emojis

# How to build

```
make
```

# How to run

```
./openemoji -accessKey=XXXXX
```
You can get access key [here](https://emoji-api.com) after signing up 

# Example output
```sh
String: Hello world!, contains emoji: false
String: Hello world 🤗, contains emoji: true
Emojis count: 1357
```