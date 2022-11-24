"""Main entrypoint of the program"""

import re
import argparse
import os
import sys
import shutil
import markdown
import yaml
from rich.console import Console
from rich.theme import Theme

console = Console(theme=Theme(inherit=False))


def error(message):
    """Print error message and exit"""
    console.print(message)
    sys.exit(1)


def parse_arguments():
    """Parse and return comman dline arguments"""
    parser = argparse.ArgumentParser(
        prog="Simple SSG", description="Simple Static Site Generator"
    )
    parser.add_argument(
        "-i", "--inputdir", help="root directory of your website's files", required=True
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
    """Walks the directory and returns the syntax tree assuming all files are compatible with templating engine"""
    tree = {}
    regex = r"{%(.*?)%}"
    input_root = ""

    def leftover_path(shorter_path, longer_path):
        i = 0
        while i < len(shorter_path) and longer_path[i] == shorter_path[i]:
            i += 1
        if i < len(longer_path) and (longer_path[i] == "/" or longer_path[i] == "\\"):
            return longer_path[i + 1 :]
        return longer_path[i:]

    for root, _, files in os.walk(input_dir):
        if not input_root:
            input_root = root
        for name in files:
            file_path = os.path.join(root, name)
            leftover = leftover_path(input_root, file_path)
            tree[leftover] = []
            with open(file_path, "r", encoding="utf-8") as f:
                page_contents = "".join(f.readlines())
            matches = [match for match in re.finditer(regex, page_contents)]
            curr_ptr = 0
            match_ptr = 0
            while match_ptr < len(matches):
                # parse content
                content = page_contents[curr_ptr : matches[match_ptr].start()]
                tag = matches[match_ptr].group(1).strip().split(" ")
                if content:
                    tree[leftover].append(("_content", content))
                tree[leftover].append(tuple(tag))
                curr_ptr = matches[match_ptr].end()
                match_ptr += 1
            if curr_ptr < len(page_contents):
                tree[leftover].append(("_content", page_contents[curr_ptr:]))
    return tree


def generate_data(input_dir, output_dir):
    """Reads and returns a convenient structure for processing files"""
    d = {"_input_dir": input_dir, "_output_dir": output_dir}
    d["_globals"] = None
    d["public"] = False

    # config file read
    config_file_path = os.path.join(args.inputdir, "config.yml")
    if os.path.exists(config_file_path):
        with open(config_file_path) as file:
            d["_globals"] = yaml.safe_load(file)

    # directories read
    for file in os.listdir(input_dir):
        file_path = os.path.join(input_dir, file)
        if os.path.isdir(file_path):
            if file == "public":
                d["public"] = True
            elif file == "templates":
                d["templates"] = get_tree(os.path.join(input_dir, file))
            elif file == "pages":
                d["pages"] = get_tree(os.path.join(input_dir, file))
    return d


def process_public(public_dir, output_dir):
    """Copies everything inside the public directory to the output dir"""
    console.print("Copying over public files...")
    shutil.copytree(public_dir, output_dir, dirs_exist_ok=True)


def process_pages(data):
    """Parse pages directory to create final files"""
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
                    error(f"You must provide a value for '{var_name}' in '{file}'")
                variables[file][var_name] = value
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
                    error(f"Error processing variable '{item[1]}' for '{file}'")
            else:
                new_contents.append(item)
        return new_contents

    def globals_processor(file, contents, reject=False):
        new_contents = []
        for item in contents:
            if item[0] == "global":
                if reject:
                    error(f"You are not allowed to use global inside '{file}'")
                try:
                    new_contents.append(("_content", data["_globals"][item[1]]))
                except Exception:
                    error(f"Error processing global variable '{item[1]}' for '{file}'")
            else:
                new_contents.append(item)
        return new_contents

    def expand_processor(file, contents, reject=False):
        new_contents = []
        for item in contents:
            if item[0] == "expand":
                if reject:
                    error(f"You are not allowed to use expand inside '{file}'")
                try:
                    new_contents += data["templates"][item[1]]
                except Exception:
                    error(f"Error expanding template '{item[1]}' for '{file}'")
            else:
                new_contents.append(item)
        return new_contents

    def template_processor(file, contents, reject=False):
        new_contents = []
        template = None
        for item in contents:
            if item[0] == "template":
                if reject:
                    error(f"You are not allowed to use template inside '{file}'")
                if not template:
                    template = item
                else:
                    error(
                        f"You cannot have more than one template declarations inside '{file}'"
                    )
            else:
                new_contents.append(item)

        if template:
            template_data = data["templates"][template[1]]
            for i, item in enumerate(template_data):
                if item[0] == "content":
                    new_contents = template_data[:i] + new_contents
                    if i + 1 < len(template_data):
                        new_contents += template_data[i + 1 :]
                    break
        return new_contents

    def content_processor(file, contents, reject=False):
        new_contents = []
        current_contents = ""
        for i, item in enumerate(contents):
            if item[0] == "_content":
                if reject:
                    error(f"You are not allowed to use _content inside '{file}'")
                try:
                    current_contents += item[1]
                except Exception:
                    error(f"Error expanding template '{item[1]}' for '{file}'")
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
                        (
                            "_content",
                            markdown.markdown(
                                item[1],
                                extensions=["fenced_code", "tables", "footnotes"],
                            ),
                        )
                    )
                except Exception:
                    error(f"Error expanding template '{item[1]}' for '{file}'")
            else:
                new_contents.append(item)
        return file, new_contents

    # markdown
    for file, contents in data["pages"].copy().items():
        if file.endswith(".md"):
            del data["pages"][file]
            file, contents = process_markdown(file, contents)
            data["pages"][file] = contents
    for file, contents in data["templates"].copy().items():
        if file.endswith(".md"):
            _, contents = process_markdown(file, contents)
            data["templates"][file] = contents

    # def
    for file, contents in data["pages"].items():
        data["pages"][file] = def_processor(file, contents)
    for file, contents in data["templates"].items():
        def_processor(file, contents, True)

    # use
    for file, contents in data["pages"].items():
        data["pages"][file] = use_processor(file, contents)
    for file, contents in data["templates"].items():
        use_processor(file, contents, True)

    # globals
    for file, contents in data["pages"].items():
        data["pages"][file] = globals_processor(file, contents)
    for file, contents in data["templates"].items():
        data["templates"][file] = globals_processor(file, contents)

    # expand
    for file, contents in data["pages"].items():
        data["pages"][file] = expand_processor(file, contents)
    for file, contents in data["templates"].items():
        data["templates"][file] = expand_processor(file, contents)

    # template
    for file, contents in data["pages"].items():
        data["pages"][file] = template_processor(file, contents)
    for file, contents in data["templates"].items():
        data["templates"][file] = template_processor(file, contents)

    # _content
    for file, contents in data["pages"].items():
        data["pages"][file] = content_processor(file, contents)

    # write
    for file_path, contents in data["pages"].items():
        console.print(f"Writing: {file_path}")
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


if __name__ == "__main__":
    args = parse_arguments()
    data = generate_data(args.inputdir, args.outputdir)
    # process public
    if data["public"]:
        process_public(os.path.join(data["_input_dir"], "public"), data["_output_dir"])
    # process pages
    process_pages(data)
