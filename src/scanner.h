#ifndef SCANNER_H
#define SCANNER_H

#include <stdbool.h>
#include <stdio.h>

#include "vector.h"

typedef enum {
  // characters
  COMMA,
  LEFT_PAREN,
  RIGHT_PAREN,
  LEFT_BRACK,
  RIGHT_BRACK,
  EQUAL,
  // keywords
  FOR,
  IF,
  ELSEIF,
  ELSE,
  END,
  EXPAND,
  TEMPLATE,
  PROP,
  CONTENT,
  GLOBAL,
  OUT_ONLY,
  // literals
  STRING,
  NUMBER,
  // identifier
  IDEN,
  // content
  HTML,
  CODE,
  MD
} token_type;

typedef struct {
  token_type type;
  int line;
} token;

typedef struct {
  token *tokens;
  FILE *fp;
  const char *file_path;
  bool tag_open;
  char **errors;
  size_t num_errors;
} scanner;

scanner make_scanner(const char path[]);
vector *match_string(scanner *s);
void tokenize(scanner *s);

#endif
