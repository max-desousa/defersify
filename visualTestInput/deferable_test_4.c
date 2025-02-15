#include <stdio.h>
#include <stdlib.h>

typedef struct {
    int a;
    int b;
} Complex;

int main() {
  FILE *f = fopen("complex.txt", "r");
  defer fclose(f);

  if (f == NULL) {
    printf("Error opening file!\n");
    return 1;
  }
  else {
    defer printf("File has been completely read!\n");
    Complex c;
    fscanf(f, "%d %d", &c.a, &c.b);
    printf("Complex number: %d + %di\n", c.a, c.b);
  }

  return 0;
}
