#include "vector.h"

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void increase_cap_if_required(vector *vec) {
  if (vec->num_els > vec->max_els)
    return;
  if (vec->max_els > SIZE_MAX / 2) {
    fprintf(stderr, "Max vector capacity reached");
    exit(EXIT_FAILURE);
  }
  size_t new_max_els = vec->max_els * 2;
  if (new_max_els > SIZE_MAX / vec->el_size) {
    fprintf(stderr, "Max vector capacity reached");
    exit(1);
  }
  size_t new_num_bytes = new_max_els * vec->el_size;
  void *new_data = realloc(vec->data, new_num_bytes);
  if (new_data == NULL) {
    fprintf(stderr, "Failed to reallocate %zu bytes for vector\n",
            new_num_bytes);
    exit(EXIT_FAILURE);
  }
  vec->max_els = new_max_els;
  vec->data = new_data;
}

vector vec_make(size_t el_size, size_t count) {
  vector v;
  size_t mem_req = count * el_size;
  if (count == 0) {
    fprintf(stderr, "Creating vector of zero size is not allowed\n");
    exit(EXIT_FAILURE);
  }
  if (el_size == 0) {
    fprintf(stderr, "Creating vector of element size of 0 is not allowed\n");
    exit(EXIT_FAILURE);
  }
  v.data = malloc(mem_req);
  if (v.data == NULL) {
    fprintf(stderr, "Failed to allocate %zu bytes for vector\n", mem_req);
    exit(EXIT_FAILURE);
  }
  v.max_els = count;
  v.num_els = 0;
  v.el_size = el_size;
  return v;
}

void vec_push(vector *vec, void *element) {
  increase_cap_if_required(vec);
  size_t offset = vec->num_els * vec->el_size;
  memcpy((char *)vec->data + offset, element, vec->el_size);
  vec->num_els++;
}

void vec_get(vector *vec, size_t pos, void *element) {
  if (pos >= vec->num_els) {
    fprintf(stderr, "Out of bound vector element access\n");
    exit(EXIT_FAILURE);
  }
  size_t offset = pos * vec->el_size;
  memcpy((char *)vec->data + offset, element, vec->el_size);
}

void vec_pop(vector *vec, void *element) {
  vec_get(vec, vec->num_els - 1, element);
  vec->num_els--;
}

void vec_free(vector *vec) { free(vec->data); }
