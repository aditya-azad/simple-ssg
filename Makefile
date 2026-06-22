.PHONY: all clean run

CC ?= cc
CFLAGS ?= -Wall -Wextra -std=c11 -g
SRCS := src/main.c src/fs.c

all: ssg

ssg: $(SRCS)
	$(CC) $(CFLAGS) $(SRCS) -o sssg

run: sssg
	./sssg example out

clean:
	rm -f sssg
	rm -rf out
