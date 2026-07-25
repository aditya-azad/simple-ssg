#ifndef VECTOR_H
#define VECTOR_H

#include <stdbool.h>
#include <stddef.h>

typedef struct {
  size_t num_els;
  size_t max_els;
  size_t el_size;
  void *data;
} vector;

vector make_vec(size_t el_size, size_t count);
void push(vector *vec, void *element);
void get(vector *vec, size_t pos, void *element);
void pop(vector *vec, void *element);
void free_vec(vector *vec);

#endif
