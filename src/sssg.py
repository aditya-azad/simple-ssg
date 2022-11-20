"""Main entrypoint of the program"""

from collections import namedtuple
import re
import argparse
import os
import shutil


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
    """Expands all the tags present in the page recursively"""
    # TODO: this function can be optimized to just go over the page once for searching and expanding tags
    with open(page_path, "r", encoding="utf-8") as file:
        page_contents = "".join(file.readlines())
    regex = r"{%(.*?)%}"

    # fill content
    matches = re.finditer(regex, page_contents)
    for match in matches:
        command, *_ = match.group(1).strip().split(" ")
        if command == "fill_content":
            page_contents = (
                page_contents[: match.start()] + content + page_contents[match.end() :]
            )

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
                    page_contents[match.start() :],
                )
            break
        if not found:
            break

    return page_contents


def process_pages(base_pages_dir, pages_dir, output_dir, template_dir):
    """Processes html files inside pages directory recursively"""
    files = os.listdir(pages_dir)
    for file in files:
        if file.endswith(".html"):
            file_path = os.path.join(pages_dir, file)
            print(f"Processing: {file_path}")
            contents = generate_page(file_path, template_dir)
            os.makedirs(
                os.path.join(output_dir, leftover_path(pages_dir, base_pages_dir)),
                exist_ok=True,
            )
            with open(
                os.path.join(
                    output_dir, leftover_path(pages_dir, base_pages_dir), file
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
