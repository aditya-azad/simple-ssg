"""Main entrypoint of the program"""

import re
import argparse
import json
import os
import sys
import shutil
import markdown
import yaml
import minify_html_onepass
import operator
from functools import reduce
from PIL import Image
from rich.console import Console
from rich.theme import Theme

console = Console(theme=Theme(inherit=False))

md_extensions = ["fenced_code", "tables", "footnotes", "codehilite"]


def get_file_contents(file_path):
    with open(file_path, "r", encoding="utf-8") as f:
        return "".join(f.readlines())


def error(msg):
    """Print error message and exit"""
    console.print(f"[red]> {msg}[/red]")
    sys.exit(1)


def message(msg):
    """Print a message to the screen"""
    console.print(f"[green]>[/green] {msg}")


def leftover_path(shorter_path, longer_path):
    """Essentially subtract shorter path from longer_path and return the difference"""
    i = 0
    while i < len(shorter_path) and longer_path[i] == shorter_path[i]:
        i += 1
    if i < len(longer_path) and (longer_path[i] == "/" or longer_path[i] == "\\"):
        return longer_path[i + 1:]
    return longer_path[i:]


def parse_arguments():
    """Parse and return comman dline arguments"""
    parser = argparse.ArgumentParser(
        prog="sssg", description="Simple Static Site Generator"
    )
    parser.add_argument(
        "-i", "--inputdir",
        help="root directory of your website's files", required=True
    )
    parser.add_argument(
        "-o",
        "--outputdir",
        help="output directory of your compiled website",
        required=True,
    )
    args = parser.parse_args()
    if (not os.path.exists(args.inputdir)) or (not os.path.exists(args.outputdir)):
        error("Invalid directory paths given")
    if len(os.listdir(args.outputdir)) != 0:
        error(
            "Output directory not empty, delete everything inside the output directory"
        )
    return args


def get_tree(input_dir):
    """Walks the directory and returns the syntax tree assuming all files are
    compatible with templating engine"""
    tree = {}
    input_root = ""

    for root, _, files in os.walk(input_dir):
        if not input_root:
            input_root = root
        for name in files:
            file_path = os.path.join(root, name)
            leftover = leftover_path(input_root, file_path)
            if file_path.endswith(".ipynb"):
                tree[leftover] = get_ipynb_tree(file_path)
            else:
                tree[leftover] = get_text_tree(get_file_contents(file_path))
    return tree


def get_ipynb_tree(file_path):
    """convert ipynb to markdown"""
    out = ""
    regex = r"{%(.*?)%}"
    with open(file_path, "r", encoding="utf-8") as f:
        data = json.load(f)
    lang = data["metadata"]["language_info"]["name"]
    for cell in data["cells"]:
        if cell["cell_type"] == "markdown":
            for text in cell["source"]:
                out += text
        elif cell["cell_type"] == "code":
            code = ""
            code_output = ""
            for text in cell["source"]:
                match = re.search(regex, text, flags=re.S)
                if match and match.group(1).strip() == "out_only":
                    code = ""
                    break
                else:
                    code += text
            if code:
                out += f"```{lang}\n{code.strip()}\n```\n\n"
            for output in cell["outputs"]:
                code_output += "".join(output["text"])
            if code_output:
                out += f"```output\n{code_output.strip()}\n```\n\n"
        # new cell = new paragaraph
        if not out.endswith("\n"):
            out += "\n\n"
        elif out.endswith("\n") and not out.endswith("\n\n"):
            out += "\n"
    return get_text_tree(out)


def get_text_tree(page_contents):
    """convert text types to tree representation"""
    regex = r"{%(.*?)%}"
    tree = []
    matches = [match for match in re.finditer(regex, page_contents, flags=re.S)]
    curr_ptr = 0
    match_ptr = 0
    while match_ptr < len(matches):
        # parse content
        content = page_contents[curr_ptr: matches[match_ptr].start()]
        tag = matches[match_ptr].group(1).strip().split(" ")
        if content:
            tree.append(("_content", content))
        tree.append(tuple(tag))
        curr_ptr = matches[match_ptr].end()
        match_ptr += 1
    if curr_ptr < len(page_contents):
        tree.append(("_content", page_contents[curr_ptr:]))
    return tree


def generate_data(input_dir, output_dir):
    """Reads and returns a convenient structure for processing files"""
    data = {"_input_dir": input_dir, "_output_dir": output_dir}
    data["_globals"] = None
    data["public"] = False

    # config file read
    config_file_path = os.path.join(input_dir, "config.yml")
    if os.path.exists(config_file_path):
        with open(config_file_path, "r", encoding="utf-8") as file:
            data["_globals"] = yaml.safe_load(file)

    # directories read
    for file in os.listdir(input_dir):
        file_path = os.path.join(input_dir, file)
        if os.path.isdir(file_path):
            if file == "public":
                data["public"] = True
            elif file == "templates":
                data["templates"] = get_tree(os.path.join(input_dir, file))
            elif file == "pages":
                data["pages"] = get_tree(os.path.join(input_dir, file))
    return data


def process_public(public_dir, output_dir):
    """Copies everything inside the public directory to the output dir"""
    message("Copying over public files...")
    shutil.copytree(public_dir, output_dir, dirs_exist_ok=True)


def process_pages(data):
    """Parse pages and templates directory to create final files"""
    variables = {}

    def def_processor(file, contents, reject=False):
        new_contents = []
        for item in contents:
            if item[0] == "def":
                if reject:
                    error(f"You are not allowed to use defs inside '{file}'")
                variables.setdefault(file, {})
                value = " ".join(item[2:]).strip()
                var_name = item[1]
                if not value:
                    error(
                        f"You must provide a value for '{var_name}' in '{file}'")
                variables[file][var_name] = value
                # add slug speciala variabel
                if not variables[file].get("_slug"):
                    slug = file.replace("\\", "/")
                    if not slug.startswith("/"):
                        slug = f"/{slug}"
                    variables[file]["_slug"] = slug
            else:
                new_contents.append(item)
        return new_contents

    def use_processor(file, contents, reject=False):
        new_contents = []
        for item in contents:
            if item[0] == "use":
                if reject:
                    error(f"You are not allowed to use use inside '{file}'")
                try:
                    new_contents.append(("_content", variables[file][item[1]]))
                except Exception:
                    error(
                        f"Error processing variable '{item[1]}' for '{file}'")
            else:
                new_contents.append(item)
        return new_contents

    def globals_processor(file, contents):
        new_contents = []
        for item in contents:
            if item[0] == "global":
                try:
                    new_contents.append(
                        ("_content", data["_globals"][item[1]]))
                except Exception:
                    error(
                        f"Error processing global variable '{item[1]}' for '{file}'")
            else:
                new_contents.append(item)
        return new_contents

    def expand_processor(file, contents):
        new_contents = []
        for item in contents:
            if item[0] == "expand":
                try:
                    new_contents += data["templates"][item[1]]
                except Exception:
                    error(f"Error expanding template '{item[1]}' for '{file}'")
            else:
                new_contents.append(item)
        return new_contents

    def template_processor(file, contents, ignore_props=False):
        new_contents = []
        template = None
        for item in contents:
            if item[0] == "template":
                if not template:
                    template = item
                else:
                    error(
                        f"You cannot have more than one template declarations inside '{file}'"
                    )
            else:
                new_contents.append(item)

        if template:
            # expand template
            try:
                template_data = data["templates"][template[1]]
            except Exception:
                if file.endswith(".html"):
                    file = file[:-5]
                error(
                    f"Cannot find template '{template[1]}' for file '{file}'")
            template_props = set(template[2:])
            for i, item in enumerate(template_data):
                if item[0] == "content":
                    new_contents = template_data[:i] + new_contents
                    if i + 1 < len(template_data):
                        new_contents += template_data[i + 1:]
                    break
            # expand props
            if not ignore_props:
                prop_expanded_contents = []
                for i, item in enumerate(new_contents):
                    if item[0] == "prop":
                        if item[1] not in template_props:
                            error(
                                f"Prop '{item[1]}' not passed to the page '{template[1]}'"
                            )
                        else:
                            prop_expanded_contents.append(
                                ("_content", variables[file][item[1]])
                            )
                    else:
                        prop_expanded_contents.append(item)
                new_contents = prop_expanded_contents
            # recurse to see if there is more tempaltes
            new_contents = template_processor(file, new_contents, ignore_props)
        return new_contents

    def content_processor(file, contents):
        new_contents = []
        current_contents = ""
        for item in contents:
            if item[0] == "_content":
                try:
                    current_contents += item[1]
                except Exception:
                    error(f"Error merging content '{item[1]}' for '{file}'")
            else:
                new_contents.append(("_content", current_contents))
                current_contents = ""
        if current_contents:
            new_contents.append(("_content", current_contents))
        return new_contents


    def process_markdown(file, contents):
        new_contents = []
        file = file[:-3] + ".html"
        for item in contents:
            if item[0] == "_content":
                try:
                    new_contents.append(
                        ("_content", markdown.markdown(item[1], extensions=md_extensions))
                    )
                except Exception:
                    error(f"Error error converting markdown '{item[1]}' for '{file}'")
            else:
                new_contents.append(item)
        return file, new_contents


    def process_ipynb(file, contents):
        new_contents = []
        file = file[:-6] + ".html"
        for item in contents:
            if item[0] == "_content":
                try:
                    new_contents.append(
                        ("_content", markdown.markdown(item[1], extensions=md_extensions))
                    )
                except Exception:
                    error(f"Error error converting markdown '{item[1]}' for '{file}'")
            else:
                new_contents.append(item)
        return file, new_contents


    def for_processor(file, contents):
        new_contents = []
        loop_var = ""
        loop_vars = []
        for item in contents:
            if item[0] == "for":
                sort_key = None
                reversed = False
                # thing
                loop_var = item[1].strip()
                # in
                if item[2] != "in":
                    error(f"Syntax error in for loop for '{file}'")
                # variable
                if item[3].startswith("_"):
                    if item[3].find("sort(") != -1:
                        opening = item[3].find("(")
                        closing = item[3].rfind(")")
                        starting = item[3].find(",")
                        sort_key = item[3][starting + 1:closing].strip()
                        if item[3].find("rsort(") != -1:
                            reversed = True
                        root = item[3][opening + 1:starting].strip()
                    else:
                        root = item[3].strip()
                    for i, v in variables.items():
                        if i.startswith(root):
                            loop_vars.append(v)
                # content
                content = " ".join(item[4:])
                # loop
                # parse out the contents
                regex = r"\{\$(.*?)\$\}"
                parsed_content = []
                curr = 0
                for match in re.finditer(regex, content):
                    var = match.group(1).strip()
                    parsed_content.append(
                        ("_content", content[curr: match.start()]))
                    var_tree = var.split(".")
                    if var.startswith(loop_var):
                        parsed_content.append(("loop_var", var_tree[1:]))
                    elif var.startswith("this"):
                        if len(var_tree) != 2:
                            error(
                                f"Invalid use of 'this' in '{file}'. Need to have exactly one '.'"
                            )
                        parsed_content.append(("use", var_tree[1]))
                    elif var.startswith("_global"):
                        if len(var_tree) != 2:
                            error(
                                f"Invalid use of '_global' in '{file}'. Need to have exactly one '.'"
                            )
                        parsed_content.append(("global", var_tree[1]))
                    curr = match.end()
                if curr < len(content):
                    parsed_content.append(("_content", content[curr:]))
                # put in the contents
                if sort_key:
                    loop_vars = sorted(
                        loop_vars,
                        key=lambda x: x[sort_key],
                        reverse=reversed)
                for v in loop_vars:
                    for i in parsed_content:
                        if i[0] == "loop_var":
                            try:
                                new_contents.append(
                                    ("_content", reduce(
                                        operator.getitem, i[1], v))
                                )
                            except Exception:
                                error(
                                    f"Cannot find a key in the loop variable of for loop in '{file}'"
                                )
                        else:
                            new_contents.append(i)
                if not loop_vars:
                    new_contents += parsed_content
            else:
                new_contents.append(item)
        return new_contents

    # markdown
    for file, contents in data["pages"].copy().items():
        if file.endswith(".md"):
            del data["pages"][file]
            file, contents = process_markdown(file, contents)
            data["pages"][file] = contents
    for file, contents in data["templates"].copy().items():
        if file.endswith(".md"):
            error("You cannot have markdown templates")

    # ipynb
    for file, contents in data["pages"].copy().items():
        if file.endswith(".ipynb"):
            del data["pages"][file]
            file, contents = process_ipynb(file, contents)
            data["pages"][file] = contents
    for file, contents in data["templates"].copy().items():
        if file.endswith(".md"):
            error("You cannot have markdown templates")

    # def
    for file, contents in data["pages"].items():
        data["pages"][file] = def_processor(file, contents)
    for file, contents in data["templates"].items():
        def_processor(file, contents, True)

    # for loops
    for file, contents in data["pages"].items():
        data["pages"][file] = for_processor(file, contents)

    # globals
    for file, contents in data["pages"].items():
        data["pages"][file] = globals_processor(file, contents)
    for file, contents in data["templates"].items():
        data["templates"][file] = globals_processor(file, contents)

    # use
    for file, contents in data["pages"].items():
        data["pages"][file] = use_processor(file, contents)
    for file, contents in data["templates"].items():
        use_processor(file, contents, True)

    # expand
    for file, contents in data["pages"].items():
        data["pages"][file] = expand_processor(file, contents)
    for file, contents in data["templates"].items():
        data["templates"][file] = expand_processor(file, contents)

    # template
    for file, contents in data["templates"].items():
        data["templates"][file] = template_processor(file, contents, True)
    for file, contents in data["pages"].items():
        data["pages"][file] = template_processor(file, contents)

    # _content
    for file, contents in data["pages"].items():
        data["pages"][file] = content_processor(file, contents)


def write_files(data):
    """Write the files of pages into the output directory
    assuming there is only one _content block for all pages"""
    for file_path, contents in data["pages"].items():
        message(f"Writing: {file_path}")
        contents = contents[0][1]
        if ("/" in file_path) or ("\\" in file_path):
            head, tail = os.path.split(file_path)
            base_path = os.path.join(data["_output_dir"], head)
            file_path = tail
            os.makedirs(base_path, exist_ok=True)
        else:
            base_path = data["_output_dir"]
        with open(os.path.join(base_path, file_path), "w", encoding="utf-8") as file:
            file.write(contents)


def minify(output_dir):
    """Minify HTML, CSS, JS and image files present in output directory"""
    message("Minifiying files...")
    for root, _, files in os.walk(output_dir):
        for file in files:
            # html, css, js
            file_path = os.path.join(root, file)
            if (
                file_path.endswith(".html")
                or file_path.endswith(".css")
                or file_path.endswith(".js")
            ):
                with open(file_path, "r", encoding="utf-8") as fl:
                    contents = "".join(fl.readlines())
                with open(file_path, "w", encoding="utf-8") as fl:
                    fl.write(minify_html_onepass.minify(
                        contents, minify_js=True))
            # images
            if (
                file_path.lower().endswith(".png")
                or file_path.lower().endswith(".jpg")
                or file_path.lower().endswith(".jpeg")
            ):
                pic = Image.open(file_path)
                # remove metadata
                stripped = Image.new(pic.mode, pic.size)
                stripped.putdata(pic.getdata())
                if 'P' in pic.mode:
                    stripped.putpalette(pic.getpalette())  # type: ignore
                stripped.save(file_path, optimized=True, quality=95)


def run():
    """Run the app"""
    console.print("[bold cyan]Simple SSG[/bold cyan]")
    args = parse_arguments()
    data = generate_data(args.inputdir, args.outputdir)
    # process public
    if data["public"]:
        process_public(os.path.join(
            data["_input_dir"], "public"), data["_output_dir"])
    # process pages
    process_pages(data)
    # write files
    write_files(data)
    # minify
    minify(data["_output_dir"])


if __name__ == "__main__":
    run()
