#include <stdio.h>
#include "usedByC.h"
int main(int argc, char **argv) {
  GoInt x = 12;
  GoInt y = 23;
  printf("About to call a Go function!\n");
  // 导入usedByC.h之后，就可以使用其中实现的函数
  PrintMessage();

  //从Go函数中获取一个整数值
  GoInt p = Multiply(x,y);
  // 使用(int)p把它转换成C语言的整数
  printf("Product: %d\n",(int)p);
  printf("It worked!\n");
  return 0;
}

// gcc -o willUseGo willUseGo.c ./usedByC.o
// ./willUseGo

