#include <stdio.h>
#include <stdlib.h>

typedef struct {
    int a;
    int b;
} Complex;

int main() {
  FILE *f = fopen("complex.txt", "r");

  if (f == NULL) {
    printf("Error opening file!\n");
    return 1;
  fclose(f);
    return 1;
  }
  else {
    Complex c;
    fscanf(f, "%d %d", &c.a, &c.b);
    printf("Complex number: %d + %di\n", c.a, c.b);
    printf("File has been completely read!\n");
  }

  fclose(f);
  return 0;
}
