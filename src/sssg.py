"""Main entrypoint of the program"""

import re
import argparse
import os
import shutil


def validate_cmd_args(args):
    """Raises exception if the arguments are not what expected"""
    if (not os.path.exists(args.inputdir)) or (not os.path.exists(args.outputdir)):
        raise Exception("Invalid directory paths given")
    if len(os.listdir(args.outputdir)) != 0:
        raise Exception(
            "Output directory not empty, delete everything inside the output directory"
        )


def generate_page(page_path, output_dir):
    with open(page_path, "r") as file:
        contents = "\n".join(file.readlines())
    matches = re.finditer(r"{%(.*?)%}", contents)
    for match in matches:
        print(
            f"'{match.group(1).strip()}' , START_INDEX = {match.start()}, END_INDEX = {match.end()}"
        )


def process_pages(pages_dir, output_dir):
    """Processes html files inside pages directory recursively"""
    files = os.listdir(pages_dir)
    for file in files:
        if file.endswith(".html"):
            file_path = os.path.join(pages_dir, file)
            generate_page(file_path, output_dir)
        elif os.path.isdir(file):
            process_pages(os.path.join(pages_dir, file), output_dir)


def process_public(public_dir, output_dir):
    """Copies everything inside the public directory to the output dir"""
    """
    shutil.copytree(public_dir, output_dir, dirs_exist_ok=True)
    """


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
            process_pages(os.path.join(args.inputdir, "pages"), args.outputdir)
        elif file == "public":
            process_public(os.path.join(args.inputdir, "public"), args.outputdir)
