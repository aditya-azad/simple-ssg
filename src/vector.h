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

vector vec_make(size_t el_size, size_t count);
void vec_push(vector *vec, void *element);
void vec_get(vector *vec, size_t pos, void *element);
void vec_pop(vector *vec, void *element);
void vec_free(vector *vec);

#endif
