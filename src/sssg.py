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
                if file not in variables:
                    variables[file] = {}
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
    # template
    # markdown

    console.print(data)
    console.print(variables)

    """
    if file.endswith(".html") or file.endswith(".md"):
        file_path = os.path.join(pages_dir, file)
        contents = generate_page(file_path, template_dir)
        os.makedirs(
            os.path.join(output_dir, leftover_path(pages_dir, base_pages_dir)),
            exist_ok=True,
        )
        final_file_name = file
        if file.endswith(".md"):
            final_file_name = final_file_name[:-3] + ".html"
        with open(
            os.path.join(
                output_dir,
                leftover_path(pages_dir, base_pages_dir),
                final_file_name,
            ),
            "w",
            encoding="utf-8",
        ) as file:
            file.write(contents)
    elif os.path.isdir(os.path.join(pages_dir, file)):
        process_pages(
            base_pages_dir, os.path.join(pages_dir, file), output_dir, template_dir
        )
    """


if __name__ == "__main__":
    args = parse_arguments()
    data = generate_data(args.inputdir, args.outputdir)
    # process public
    if data["public"]:
        process_public(os.path.join(data["_input_dir"], "public"), data["_output_dir"])
    # process pages
    process_pages(data)
