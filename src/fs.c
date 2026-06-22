#include "fs.h"

#include <dirent.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>

bool dir_exists(const char *path) {
  struct stat st;
  return stat(path, &st) == 0 && S_ISDIR(st.st_mode);
}

bool dir_create(const char *path, bool exists_ok) {
  if (mkdir(path, 0777) == 0) {
    return true;
  }
  if (errno == EEXIST) {
    if (exists_ok)
      return true;
    return false;
  }
  return false;
}

bool dir_empty(const char *path) {
  DIR *dir = opendir(path);
  // error opening directory
  if (dir == NULL) {
    printf("Cannot open directory %s", path);
    exit(1);
  }
  errno = 0;
  struct dirent *entry;
  // empty or not
  while ((entry = readdir(dir)) != NULL) {
    if (strcmp(entry->d_name, ".") != 0 && strcmp(entry->d_name, "..") != 0) {
      closedir(dir);
      return false;
    }
  }
  // error while reading
  if (errno != 0) {
    printf("Cannot read directory %s", path);
    exit(1);
  }
  closedir(dir);
  // empty
  return true;
}
