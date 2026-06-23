# Simple Static Site Generator

Don't want to download billions of bytes for just creating a simple website? Or don't want to learn and re-learn a whole bunch of syntax every time you want to customize the theme once in a blue moon? SSSG got you covered.

- Minimal templating language
- Markdown support
- Jupyter notebooks support (with JPEG, PNG and SVG images)
- Minify HTML, CSS and JS files
- Losslessly compressed PNG and JP(E)G images and remove metadata

## Directory structure

There is a specific directory structure that an input directory must follow.

### `pages` directory

All your concrete web pages are stored here. The pages in the directory are also the parse tree roots. Parser goes over all the pages one by one, parsing them and writing them to output directory. All the directories and pages inside this directory conforms to URL schema.

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

Everything inside this is copied to the output directory.

### `template` directory

This is where all the templates reside.

### `config.yml` file

This is where you can store all the global variables that are accessible from `pages` and `templates` directories.

## Templating Language Syntax

The language use tags similar to Jinja's `{% ... %}` syntax. Templates can import templates.

### `{% template <template_name> <props> %}`

The `template_name` is used as template for the page. The contents of the page is replaced with contents of the template. The replaced content is placed in `content` block. Additionally props can be passed to the templates. They follow `x="abc"` syntax. It supports number, strings and list types. To use props, use `prop` tag.

### `{% expand <template_name> <props> %}`

This tag is replaced with the contents of the template.

### `{% prop <variable_name> %}`

You can use passed in props to the template using this tag. Content of the prop is replaced with the value.

### `{% content %}`

This is where the contents of the page are pasted when using `template`.

### `{% global <variable_name> %}`

You can use the variables defined in `config.yml` file using this tag.

### `{% for <variable_name> in <expression> %} ... {% endfor %}` 

You can also loop over expressions that result in list using `for` statement. The variable is scoped within the scope. Variables are accessed using `{% var <variable_name> %}`.

### `{% if <expression> %} ... {% else if <expression> %} ... {% else %} ... {% endif %}`

There is support for conditional statements using if-else statements. "else if" and "else" blocks are optional.

### `{% out_only %}`

Only supported in `ipynb`. If this string is present anywhere in a code cell, the code will not be displayed, only the output

### built-ins

There are a couple of builtin functions for common operations
- `sort(<list>, "asc"|"dec")` - when provided a list, it sorts the list in the given order if the list is of single type

## How to use

Requirements
- C

### Install

