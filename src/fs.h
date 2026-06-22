#ifndef FS_H
#define FS_H

#include <stdbool.h>

bool dir_exists(const char *path);
bool dir_create(const char *path, bool exists_ok);
bool dir_empty(const char *path);

#endif
