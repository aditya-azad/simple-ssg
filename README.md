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

This is where all the templates for pages (and other templates) reside.

### `config.yml` file

This is where you can store all the global variables that are accessable from `pages` and `templates` directories.

## Templating Language Syntax

The language use tags similar to Jinja's `{% ... %}` syntax. Templates themselves have can have other templates within them.

### `{% template <template_name> <props> %}`

The `template_name` is used as template for the page. Basically, the contents of the page is replaced with contents of the template. The replaced content is placed where `content` in the template used to be. Additionally props can be passed to the templates. They need to be variables therefore you need to first define them using `def`. To use those props in the template file, you have to use `prop` tag.

### `{% expand <template_name> %}`

This tag is replaced with the contents of the template.

### `{% content %}`

This is where the contents of the page are pasted when using `template`.

### `{% global <variable_name> %}`

You can use the variables defined in `config.yml` file using this tag.

### `{% def <variable_name> <value> %}`

You can define variables inside the files of `pages` directory. These variables can be used in the page using `use`. The order of definition does not matter since def statements are processed before use statements. `value` can be space separated. Basically everything after `variable_name` is part of `value`.

### `{% use <variable_name> %}`

You can use the defined variables using this (see `def`).

### `{% prop <variable_name> %}`

You can fill the place with the variable passed to the template. Any other variable use will give you an error.

### `{% for <variable_name> in <directory> <content> %}`

You can also loop over files in a directory. You can have content on a separate line for readability. `<variable_name>` inside the loop is accessed using `{$<variable_name>$}`. You can access the `def`s inside those pages in the directory using the `.` operator. Optionally you can choose to sort using `_sort(<directory>,<sort_key>)` where `<sort_key>` is the `def` inside the file to sort by. Make sure there are no spaces if you are using sort. You can also reverse sort using `_rsort`. The `<variable_name>` has a special `_slug` variable that refers to the relative filepath of the page. You can use this to create anchor tags. Nesting is not allowed. See the example for clear usage.

### `{% out_only %}`

Only spported in `ipynb`. If this string is present anywhere in a code cell, the code will not be displayed.

## How to use

Requirements
- Python 3.10
- Pyinstaller (if creating executable)

### Install

The executable will be created in `dist` folder.

``` text
python -m venv venv
./venv/Scripts/activate
pip install -r requirements.txt
pyinstaller ./src/sssg.py
```

Add the `sssg` folder to path to use the executable from anywhere on your system.

### From source

``` text
python -m venv venv
./venv/Scripts/activate
pip install -r requirements.txt
python ./src/sssg.py -i <input_directory> -o <output_directory>
```

Make sure output directory is empty.

## Development Environment

Run the following command inside the output directory and go to `127.0.0.1:8000` in your browser. You have to refresh the page though.

``` text
python -m http.server 8000 --bind 127.0.0.1
```
