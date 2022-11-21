# Simple Static Site Generator

A simple static site generator written with the following features:

- Minimal templating language.
- Bloat free websites (but that is mostly up to you).
- Markdown support

## Requirements

- Python 3.10 (possibly 3.8 >= (not tested))

## TODO

- templating language
  - for loops
  - arguments to templates
- blogging example
- auto serve during development
- minify files
- compress images
- config file and variables for use inside templates
- improve performance

## Directory structure

There is a specific directory structure that all input directories must follow.

### `pages` directory

All your concrete web pages are stored here. The pages in the directory are also the parse tree roots. Parser goes over all the pages one by one, parsing them and putting writing them to output directory. All the directories and pages inside this directory, conforms to URL schema.

For example,

``` text
input_dir
|--pages
   |--posts
   |  |--first-post.html
   |--index.html
   |--about.html
...
```

This will result in URLs:

- `yourwebsite.domain/index.html`
- `yourwebsite.domain/about.html`
- `yourwebsite.domain/posts/first-post.html`.

### `public` directory

- everything inside this is copied as it is to the output directory.

### `template` directory

- this is where all the templates for pages (and other templates templates) reside.

## Templating Language Syntax

The language use tags similar to Nunjucks' `{% ... %}` syntax. There are currently following tags implemented. Templates themselves have can have other templates within them.

### `{% use_template <template_name> %}`

The `<template_name>` is used as template for the page. Basically, the contents of the page is is replaced with contents of the template. The replaced content is placed where `{% fill_content %}` in the template used to be.

### `{% fill_template <template_name> %}`

This tag is replaced with the contents of the template.

### `{% fill_content %}`

This is where the contents of the page are pasted when using `{% use_template <template_name> %}`.

## How to use

``` text
python ./src/sssg.py -i <input_directory> -o <output directory>
```

Make sure output directory is empty. There is an example provided in `example` directory. Use that directory as the `input_directory` in the above command to see how it looks.

## Development Environment

Assuming you are using python 3. Run the following command inside the output directory and go to `127.0.0.1:8000` in your browser. You have to refresh the page though.

``` text
python -m http.server 8000 --bind 127.0.0.1
```
