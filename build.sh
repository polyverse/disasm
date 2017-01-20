cc -c src/i386-dis.c -o obj/i386-dis.o -Iinc
cc -c src/dis-init.c -o obj/dis-init.o -Iinc
cc -c src/dis-buf.c -o obj/dis-buf.o -Iinc
cc -c DisAsm.c -o obj/DisAsm.o -Iinc
ar -rcs lib/libDisAsm.a obj/*.o

cc -c test.c -o test.o -Iinc
cc -o test test.o -Llib -lDisAsm
