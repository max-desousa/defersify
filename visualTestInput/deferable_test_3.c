#include <stdio.h>

int main() {
    printf("Hello, World!\n");
    defer {
        printf("Goodbye, World!\n");
    } 
    return 0;
}
