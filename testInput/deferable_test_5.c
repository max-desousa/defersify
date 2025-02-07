#include <stdio.h>

int main()
  {
  FILE *file;
  file = fopen("test_4.c", "r");
  defer printf("First File closed\n");
  defer fclose(file);
  defer printf("Closing first file\n");

  if (file == NULL)
    {
    printf("File not found\n");
    return 1;
    }
  else
    {
    FILE *second_file;
    second_file = fopen("test_3.c", "r");
    defer printf("Second File closed\n");
    defer fclose(file);
    defer printf("Closing second file\n");

    if (second_file == NULL)
      {
      printf("Second file not found\n");
      return 2;
      }
    else
      {
      printf("Both files found\n");
      }
    }

  return 0;
  }

int test()
  {
  defer printf("deferred statement\n");
  printf("Test\n");
  return 0;
  }
