ar x libassimp.a

ld -r *.o /usr/lib/gcc/x86_64-linux-gnu/4.8/libstdc++.a /usr/lib/x86_64-linux-gnu/libm.a /usr/lib/x86_64-linux-gnu/libz.a -o ../../assimp_linux_amd64.syso

rm -rf *.o
