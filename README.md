# fillin [![Travis Build Status](https://travis-ci.org/itchyny/fillin.svg?branch=master)](https://travis-ci.org/itchyny/fillin)
### fill-in your command line
A command line tool to improve your cli life.

## Motivation
We rely on shell history in our terminal operation.
We search from our shell history and execute commands dozens times in a day.

However, shell history sometimes contains some authorization token that we don't care while searching the commands.
Some incremental fuzzy searchers have troubles when there are many random tokens in the shell history.
Yeah, I know that I should not type a authorization token directly in the command line, but it's much easier than creating some shell script snippets.

Another hint to implement `fillin` is that programmers execute same commands switching servers.
We do not just login with `ssh {{hostname}}`, we also connect to the database with `psql -h {{hostname}} {{dbname}} -U {{username}}` and to Redis server with `redis-cli -h {{hostname}} -p {{port}}`.
We switch the host argument from the localhost (you may omit this), staging and production servers.

The main idea is that splitting the command history and the template variable history.
With this `fillin` command line tool, you can

- make your commands reusable and it will make incremental shell history searching easy.
- fill in template variables interactively and their history will be stored locally.

## Installation
### Homebrew
```sh
 $ brew install itchyny/fillin/fillin
```

### Download binary from GitHub Releases
[Releases・itchyny/fillin - GitHub](https://github.com/itchyny/fillin/releases)

### Build from source
```sh
 $ go get -u github.com/itchyny/fillin
```

## Usage
The interface of the `fillin` command is very simple.
Prepend `fillin` to the command and create template variables with `{{...}}`.
So the hello world for the `fillin` command is as follows.
```sh
 $ fillin echo {{message}}
message: Hello, world!        # you type here
Hello, world!                 # fillin executes: echo Hello, world!
```
The value of `message` variable is stored locally.
You can use the recently used value with the upwards key (this may be replaced with more rich interface in the future but I'm not sure).

The `{{message}}` is called as a template part of the command.
As the identifier, you can use alphabets, numbers, underscore and hyphen.
Thus `{{sample-id}}`, `{{SAMPLE_ID}}`, `{{X01}}` and `{{FOO_example-identifier0123}}` are all valid template parts.

One of the important features of `fillin` is scope grouping.
Let's look into more practical example.
When you connect to PostgreSQL server, you can use:
```sh
 $ fillin psql -h {{psql:hostname}} {{psql:dbname}} -U {{psql:username}}
[psql] hostname: example.com
[psql] dbname: example-db
[psql] username: example-user
```
What's the benefit of `psql:` prefix?
You'll notice the answer when you execute the command again:
```sh
 $ fillin psql -h {{psql:hostname}} {{psql:dbname}} -U {{psql:username}}
[psql] hostname, dbname, username: example.com, example-db, example-user   # you can select the most recently used entry with the upwards key
```
The identifiers with the same scope name (`psql` scope here) can be selected as pairs.
You can input individual values to create a new pair after skipping the multi input prompt.
```sh
 $ fillin psql -h {{psql:hostname}} {{psql:dbname}} -U {{psql:username}}
[psql] hostname, dbname, username:             # just type enter to input values for each identifiers
[psql] hostname: example.org
[psql] dbname: example-org-db
[psql] username: example-org-user
```

The scope grouping behaviour is useful with some authorization keys.
```sh
 $ fillin curl {{example-api:base-url}}/api/example -H 'Api-Key: {{example-api:api-key}}'
[example-api] base-url, api-key: example.com, apikeyabcde012345
```
The `base-url` and `api-key` are stored as tuples so you can easily switch local, staging and production environment authorization.
Without the grouping behaviour, variable history searching will lead you to an unmatched pair of `base-url` and `api-key`.

In order to have the benefit of this grouping behaviour, it's strongly recommended to prepend the scope name.
The `psql:` prefix on connecting to PostgreSQL database server, `redis:` prefix for Redis server are useful best practice in my opinion.

## Problem with pipe and redirect
The terminal interface of `fillin` is currently have problem with pipe and redirect.
For example, the following command will get stuck the terminal interface.
```sh
 $ fillin echo {{message}} | jq .
^M^M^C
```
This is because the interface of `fillin` is rely on the standard output.
Instead of connecting the output of `fillin` to another command, pass the pipe character as an argument.
```sh
 $ fillin echo {{message}} \| jq .
message: {}
{}
 $ # or
 $ fillin echo {{message}} '|' jq .
message: {}
{}
```
Same problem occurs with redirect so please escape `>`.
```sh
 $ fillin echo {{message}} \> /tmp/message
 $ # or
 $ fillin echo {{message}} '>' /tmp/message
```

## Disclaimer
This command line tool is in its early developing stage.
The user interface may be changed without any announcement.

This tool is not an encryption tool.
The command saves the inputted values in a JSON file with no encryption.

## Bug Tracker
Report bug at [Issues・itchyny/fillin - GitHub](https://github.com/itchyny/fillin/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
