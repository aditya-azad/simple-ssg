#include <stdio.h>
#include <stdlib.h>

#include "fs.h"

typedef struct {
  char *input_directory;
  char *output_directory;
} Args;

Args parse_args(int argc, char **argv) {
  Args args;
  // check number of args
  if (argc != 3) {
    printf("Required input directory and output directory as arguments!\n");
    exit(1);
  }
  // check input dir exists
  if (!dir_exists(argv[1])) {
    printf("Input directory does not exist!\n");
    exit(1);
  }
  args.input_directory = argv[1];
  // output directory create if not exists
  if (!dir_create(argv[2], true)) {
    printf("Cannot create output directory!\n");
    exit(1);
  }
  if (!dir_empty(argv[2])) {
    printf("Output directory is not empty!\n");
    exit(1);
  }
  args.output_directory = argv[2];
  return args;
}

int main(int argc, char **argv) {
  Args args = parse_args(argc, argv);
  return 0;
}
