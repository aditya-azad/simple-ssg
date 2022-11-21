"""Main entrypoint of the program"""

import re
import argparse
import os
import shutil
import markdown
import yaml

GLOBAL_VARS = {}


def leftover_path(longer_path, shorter_path):
    """Returns the path in the longer path starting from the first different character"""
    i = 0
    while i < len(shorter_path) and longer_path[i] == shorter_path[i]:
        i += 1
    if i < len(longer_path) and (longer_path[i] == "/" or longer_path[i] == "\\"):
        return longer_path[i + 1 :]
    return longer_path[i:]


def validate_cmd_args(args):
    """Raises exception if the arguments are not what expected"""
    if (not os.path.exists(args.inputdir)) or (not os.path.exists(args.outputdir)):
        raise Exception("Invalid directory paths given")
    if len(os.listdir(args.outputdir)) != 0:
        raise Exception(
            "Output directory not empty, delete everything inside the output directory"
        )


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
            if command == "fill_var":
                found = True
                if args[0] not in GLOBAL_VARS:
                    raise Exception(f"Cannot find '{args[0]}' in the config file")
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


def process_pages(base_pages_dir, pages_dir, output_dir, template_dir):
    """Processes html and md files inside pages directory recursively"""
    files = os.listdir(pages_dir)
    for file in files:
        if file.endswith(".html") or file.endswith(".md"):
            file_path = os.path.join(pages_dir, file)
            print(f"Processing: {file_path}")
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


def process_public(public_dir, output_dir):
    """Copies everything inside the public directory to the output dir"""
    print("Copying over public files...")
    shutil.copytree(public_dir, output_dir, dirs_exist_ok=True)


if __name__ == "__main__":
    # parse and validate arguments
    parser = argparse.ArgumentParser(
        prog="Simple SSG",
        description="Simple Static Site Generator",
        epilog="Text at the bottom of help",
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
    validate_cmd_args(args)
    # read config
    config_file_path = os.path.join(args.inputdir, "config.yml")
    if os.path.exists(config_file_path):
        with open(config_file_path) as file:
            GLOBAL_VARS = yaml.safe_load(file)
    # run it
    files = os.listdir(args.inputdir)
    for file in files:
        if file == "pages":
            process_pages(
                os.path.join(args.inputdir, "pages"),
                os.path.join(args.inputdir, "pages"),
                args.outputdir,
                os.path.join(args.inputdir, "templates"),
            )
        elif file == "public":
            process_public(os.path.join(args.inputdir, "public"), args.outputdir)
