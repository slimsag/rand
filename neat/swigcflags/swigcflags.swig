// #go LDFLAGS: -lX11

%{
#include <X11/Xlib.h>
%}

typedef void* Display;
Display *XOpenDisplay(char* display_name);

