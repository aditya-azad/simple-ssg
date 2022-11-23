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


def generate_page(page_path, template_dir, content=None):
    """Expands all the tags present in the html or md page recursively"""
    # TODO: this function can be optimized to just go over
    # the page once for searching and expanding tags
    regex = r"{%(.*?)%}"

    with open(page_path, "r", encoding="utf-8") as file:
        page_contents = "".join(file.readlines())
    if page_path.endswith(".md"):
        page_contents = markdown.markdown(
            page_contents, extensions=["fenced_code", "tables", "footnotes"]
        )

    # fill content
    matches = re.finditer(regex, page_contents)
    for match in matches:
        command, *args = match.group(1).strip().split(" ")
        if command == "fill_content":
            page_contents = (
                page_contents[: match.start()] + content + page_contents[match.end() :]
            )

    # fill vars
    while True:
        found = False
        matches = re.finditer(regex, page_contents)
        for match in matches:
            command, *args = match.group(1).strip().split(" ")
            if command == "fill_global":
                found = True
                if args[0] not in GLOBAL_VARS:
                    error(f"Cannot find '{args[0]}' in the config file")
                page_contents = (
                    page_contents[: match.start()]
                    + GLOBAL_VARS[args[0]]
                    + page_contents[match.end() :]
                )
                break
        if not found:
            break

    # fill templates
    while True:
        found = False
        matches = re.finditer(regex, page_contents)
        for match in matches:
            command, *args = match.group(1).strip().split(" ")
            if command == "fill_template":
                found = True
                template_contents = generate_page(
                    os.path.join(template_dir, args[0] + ".html"),
                    template_dir,
                )
                page_contents = (
                    page_contents[: match.start()]
                    + template_contents
                    + page_contents[match.end() :]
                )
                break
        if not found:
            break

    # expand into templates
    while True:
        matches = re.finditer(regex, page_contents)
        found = False
        for match in matches:
            command, *args = match.group(1).strip().split(" ")
            if command == "use_template":
                found = True
                page_contents = (
                    page_contents[: match.start()] + page_contents[match.end() :]
                )
                page_contents = generate_page(
                    os.path.join(template_dir, args[0] + ".html"),
                    template_dir,
                    page_contents,
                )
            break
        if not found:
            break

    return page_contents


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


def preprocess(input_dir, output_dir):
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
    for file in data["pages"]:
        print(file)
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
    data = preprocess(args.inputdir, args.outputdir)
    # process public
    if data["public"]:
        process_public(os.path.join(data["_input_dir"], "public"), data["_output_dir"])
    # process pages
    process_pages(data)
