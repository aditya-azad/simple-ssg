#include "scanner.h"

#include <stdlib.h>
#include <string.h>

void error(scanner *s, char error[]) {
  size_t len = strlen(error);
  malloc(len);
  // TODO: finish
  // TODO: switch all fprintf exits to this
}

bool peek(scanner *s, char *out) {
  int ch = fgetc(s->fp);
  if (ch == EOF) {
    return false;
  }
  ungetc(ch, s->fp);
  *out = (char)ch;
  return true;
}

bool advance(scanner *s, char *out) {
  int ch = fgetc(s->fp);
  if (ch == EOF) {
    return false;
  }
  return true;
}

bool match(scanner *s, const char *word) {
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

void add_token(scanner *s, token_type tt) {
  // TODO: finish this
  // TODO: need to add more data to token
}

bool match_keyword(scanner *s) {
  bool matched = false;
  if (match(s, "{%")) {
    matched = true;
    s->tag_open = true;
  } else if (match(s, "%}")) {
    matched = true;
    s->tag_open = false;
  } else if (match(s, "for")) {
    matched = true;
    add_token(s, FOR);
  } else if (match(s, "endfor")) {
    matched = true;
    add_token(s, ENDFOR);
  } else if (match(s, "if")) {
    matched = true;
    add_token(s, IF);
  } else if (match(s, "else if")) {
    matched = true;
    add_token(s, ELSEIF);
  } else if (match(s, "else")) {
    matched = true;
    add_token(s, ELSE);
  } else if (match(s, "expand")) {
    matched = true;
    add_token(s, EXPAND);
  } else if (match(s, "template")) {
    matched = true;
    add_token(s, TEMPLATE);
  } else if (match(s, "prop")) {
    matched = true;
    add_token(s, PROP);
  } else if (match(s, "content")) {
    matched = true;
    add_token(s, CONTENT);
  } else if (match(s, "global")) {
    matched = true;
    add_token(s, GLOBAL);
  }
  return matched;
}

bool match_iden(scanner *s) { 
  // TODO: finish
  return false; 
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

  // TODO: loop
  // TODO: advance till opening tag

  // advance till closing tag
  char ch;
  char next_ch;
  while (advance(s, &ch)) {
    // if we are inside a language block
    if (s->tag_open) {
      switch (ch) {
      case ',':
        add_token(s, COMMA);
        break;
      case '(':
        add_token(s, LEFT_PAREN);
        break;
      case ')':
        add_token(s, RIGHT_PAREN);
        break;
      case '[':
        add_token(s, LEFT_BRACK);
        break;
      case ']':
        add_token(s, RIGHT_BRACK);
        break;
      case '=':
        add_token(s, EQUAL);
        break;
      default:
        // match keywords
        // TODO: need to handle case when indentifiers start with one of the
        // keywords
        // TODO: need to match numbers and strings too
        if (!match_keyword(s)) {
          // match identifiers
          match_iden(s);
        }
      }
    }
  }

  // TODO: advance till end of file

  // close file
  fclose(s->fp);
}
