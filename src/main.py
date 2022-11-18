"""Main entrypoint of the program"""

import argparse

if __name__ == "__main__":
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
