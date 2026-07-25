#include "scanner.h"

#include "vector.h"

#include <stdlib.h>
#include <string.h>

void error(scanner *s, char error[]) {
  size_t len = strlen(error);
  malloc(len);
  // TODO: finish
  // TODO: switch all fprintf exits to this
}

bool match(scanner *s, const char *word) {
  // save position so we can rewind on failure
  fpos_t pos;
  if (fgetpos(s->fp, &pos) != 0) {
    return false;
  }

  for (size_t i = 0; word[i] != '\0'; i++) {
    int ch = fgetc(s->fp);
    if (ch == EOF || (unsigned char)ch != (unsigned char)word[i]) {
      fsetpos(s->fp, &pos);
      return false;
    }
  }

  return true;
}

vector *match_string(scanner *s) {
  // save position so we can rewind on failure
  fpos_t pos;
  if (fgetpos(s->fp, &pos) != 0) {
    return NULL;
  }

  // must start with opening quote
  int ch = fgetc(s->fp);
  if (ch != '"') {
    fsetpos(s->fp, &pos);
    return NULL;
  }

  // collect characters into a vector
  vector *vec = malloc(sizeof(vector));
  if (vec == NULL) {
    fsetpos(s->fp, &pos);
    return NULL;
  }
  *vec = vec_make(sizeof(char), 16);

  // read until closing quote
  bool closed = false;
  while ((ch = fgetc(s->fp)) != EOF) {
    if (ch == '"') {
      closed = true;
      break;
    }

    // handle escapes: drop the backslash, keep the next char
    if (ch == '\\') {
      int next = fgetc(s->fp);
      if (next == EOF) {
        break;
      }
      ch = next;
    }

    char c = (char)ch;
    vec_push(vec, &c);
  }

  // unterminated string: rewind and return NULL
  if (!closed) {
    vec_free(vec);
    free(vec);
    fsetpos(s->fp, &pos);
    return NULL;
  }

  return vec;
}

void add_token(scanner *s, token_type tt) {
  // TODO: finish this
  // TODO: need to add more data to token
}

scanner make_scanner(const char path[]) {
  scanner s;
  s.file_path = path;
  s.fp = NULL;
  s.tag_open = false;
  s.tokens = NULL;
  s.errors = NULL;
  s.num_errors = 0;
  return s;
}

void tokenize(scanner *s) {

  // open file
  s->fp = fopen(s->file_path, "rb");
  if (s->fp == NULL) {
    fprintf(stderr, "Cannot open file %s", s->file_path);
    exit(1);
  }

  while (true) {
    // open tag
    if (match(s, "{%")) {
      s->tag_open = true;
    }

    // TODO: need to handle case when indentifiers start with one of the
    // keywords

    // if we are inside code block
    if (s->tag_open) {
      // TODO: skip white space
      if (match(s, "%}"))
        s->tag_open = false;
      else if (match(s, ","))
        add_token(s, COMMA);
      else if (match(s, "("))
        add_token(s, LEFT_PAREN);
      else if (match(s, ")"))
        add_token(s, RIGHT_PAREN);
      else if (match(s, "["))
        add_token(s, LEFT_BRACK);
      else if (match(s, "]"))
        add_token(s, RIGHT_BRACK);
      else if (match(s, "="))
        add_token(s, EQUAL);
      else if (match(s, "for"))
        add_token(s, FOR);
      else if (match(s, "if"))
        add_token(s, IF);
      else if (match(s, "else if"))
        add_token(s, ELSEIF);
      else if (match(s, "else"))
        add_token(s, ELSE);
      else if (match(s, "end"))
        add_token(s, END);
      else if (match(s, "expand"))
        add_token(s, EXPAND);
      else if (match(s, "template"))
        add_token(s, TEMPLATE);
      else if (match(s, "prop"))
        add_token(s, PROP);
      else if (match(s, "content"))
        add_token(s, CONTENT);
      else if (match(s, "global"))
        add_token(s, GLOBAL);
      else if (match(s, "out"))
        add_token(s, OUT_ONLY);
      else if (match(s, "\"")) {
        vector *str = match_string(s);
        if (str != NULL) {
          vec_free(str);
          // TODO: use it
          free(str);
        }
      }
      else if (is_numeric(s))
        match_number(s);
      else
        match_iden(s);
    } else {
      // outside language block is content
      match_content(s);
    }
  }

  // close file
  fclose(s->fp);
}
