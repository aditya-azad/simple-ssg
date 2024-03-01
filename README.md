# Under construction

- write tests

- arch
    - public -> compress -> copy
    - templates -> linked list of types + dependencies
    - pages -> linked list of types + dependencies + defines -> compress -> copy
    - check for errors using dag

- global_vars: stores the global variables visible to all files
    - _pages: slugs of files in pages dir (does not recurse)
    - _pages_<subdir>..: files in a sub dir of pages
    - vars defined in config file

# Simple Static Site Generator

Are you sick of downloading a whole lot of build tools and plugins just for creating a simple static website?
Are you sick of learning and re-learning a whole bunch of syntax every time you want to customize the theme once in a blue moon?
SSSG got you covered.

## Features

- Minimal templating language
- Markdown support
- Jupyter notebooks support (with JPEG, PNG and SVG images)
- Minify HTML, CSS and JS files
- Lossless compress PNG and JP(E)G images and remove metadata

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

Everything inside this is copied to the output directory after compression.

### `template` directory

This is where all the templates reside.

### `config.yml` file

This is where you can store all the global variables that are accessible from `pages` and `templates` directories.

## Templating Language Syntax

The language use tags similar to Jinja's `{% ... %}` syntax. Templates themselves have can have other templates within them.

### `{% template <template_name> <props> %}`

The `template_name` is used as template for the page. The contents of the page is replaced with contents of the template. The replaced content is placed where `content` in the template used to be. Additionally, props can be passed to the templates using `<variable>=<value>` syntax (see `var`).

### `{% expand <template_name> <props> %}`

This tag is replaced with the contents of the template.

### `{% content %}`

This is where the contents of the page are pasted when using `template`. This can only appear in template files.

### `{% use <variable_name> %}`

You can use the defined variables passed as props, read from global config file or loop variable.

### `{% for <variable_name> in <directory/global> %}`

- You can also loop over files in a directory. `<variable_name>` inside the loop is accessed using `use`.
- Optionally you can choose to sort using `_sort(<directory/global>,<sort_key>)` where `<sort_key>` is a `var` inside the file to sort by.
- Make sure there are no spaces if you are using sort.
- You can also reverse sort using `_rsort`.
- The `<variable_name>` can be a special `_slug` variable that refers to the relative file path of the page. You can use this to create anchor tags.
- Nesting is not allowed.
- See the example for clear usage.

### `{% endfor %}`

Marks the end of for loop.

### `{% var <variable>=<value> %}`

A variable that can be defined inside a page to be used using `use` within the page or in a for loop somewhere else.

### `{% outonly %}`

Only spported in `.ipynb` files. If this string is present anywhere in a code cell, the code will not be displayed.
