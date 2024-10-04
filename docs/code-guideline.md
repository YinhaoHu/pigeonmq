# Code Guideline

In this article, I write about the code guideline I have followed when I was coding PigeonMQ.

## Naming

- **Keys on etcd**: Use forward slashes as separators of a key and use hyphen to seperate the words. For example, `/pigeonmq/topics/my-topic/partitions/0`.

## Project Layout

This project layout is inspired by the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).
