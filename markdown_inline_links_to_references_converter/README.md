# Markdown: convert inline links to references

I wanted to cleanup my `*.md` files and use reference-style links instead of inline urls.

## Example

### Input

```markdown
# My markdown file

I have an inline link to [Google](https://www.google.com) and [Facebook](https://www.facebook.com).

Last line of the file.
```

### Output

```markdown
# My markdown file

I have an inline link to [Google][1] and [Facebook][2].

Last line of the file.

[1]: https://www.google.com
[2]: https://www.facebook.com
```

## Usage

Clone this repo and run the following command:

`PATH_TO_YOUR_FILES` is the path to the directory containing your `*.md` files. It will parse the directory recursively. It could be also path to a single file.

```bash
docker compose run --rm -v "$(pwd)"/PATH_TO_YOUR_FILES:/input convert
```

## Known issues and potential improvements

- not everything is right if you run the script multiple times on the same content
- path handling should be better - if you're passing a directory, it will parse all files in the directory recursively
- it would be nice to have a flag to specify the output directory
- it wouldn't hurt to make sure that someone is really expecting changes in the existing files
- would be cool to run this without having to clone the repo

## Feedback

If you have any feedback, please reach out to me on Twitter or email - you can find the links on my GitHub profile [@lubieniebieski](https://https://github.com/lubieniebieski/). Probably a lot of things could be done better, so I'm open to suggestions.
