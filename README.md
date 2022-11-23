# Simple Static Site Generator

- Minimal templating language
- Markdown support

``` text
TODO
- templating language
    - for loops
    - if statements
    - template props
    - associated vars
- blogging example
- auto serve during development
- minify files
- compress images
- improve performance
```

## Requirements

- Python 3.10 (possibly 3.8 >= (not tested))

## Directory structure

There is a specific directory structure that all input directories must follow.

### `pages` directory

All your concrete web pages are stored here. The pages in the directory are also the parse tree roots. Parser goes over all the pages one by one, parsing them and putting writing them to output directory. All the directories and pages inside this directory, conforms to URL schema.

For example,

``` text
input_dir
|--pages
   |--posts
   |  |--first-post.md
   |--index.html
   |--about.html
...
```

This will result in URLs:

- `yourwebsite.domain/index.html`
- `yourwebsite.domain/about.html`
- `yourwebsite.domain/posts/first-post.html`.

### `public` directory

Everything inside this is copied as it is to the output directory.

### `template` directory

This is where all the templates for pages (and other templates templates) reside.

### `config.yml` file

This is where you can store all the global variables that is accessable from `pages` and `templates` directories.

## Templating Language Syntax

The language use tags similar to Jinja's `{% ... %}` syntax. There are currently following tags implemented. Templates themselves have can have other templates within them.

### `{% template <template_name> %}`

The `<template_name>` is used as template for the page. Basically, the contents of the page is is replaced with contents of the template. The replaced content is placed where `{% content %}` in the template used to be.

### `{% expand <template_name> %}`

This tag is replaced with the contents of the template.

### `{% content %}`

This is where the contents of the page are pasted when using `{% template <template_name> %}`.

### `{% global <variable_name> %}`

You can use the variables defined in `config.yml` file using this tag.

## How to use

### Install

Use pyinstaller for creating an executable. After you run this command, you will get your executable in `dist` folder.

``` text
pyinstaller ./src/sssg.py
```

Copy the `sssg` folder to anywhere on your system and add it to path to use the executable from anywhere on your system.

### From source

``` text
pip install -r requirements.txt
python ./src/sssg.py -i <input_directory> -o <output directory>
```

Make sure output directory is empty. There is an example provided in `example` directory. Use that directory as the `input_directory` in the above command to see how it looks.

## Development Environment

Assuming you are using python 3. Run the following command inside the output directory and go to `127.0.0.1:8000` in your browser. You have to refresh the page though.

``` text
python -m http.server 8000 --bind 127.0.0.1
```
