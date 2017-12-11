#ifndef __API_INCLUDE__
#define __API_INCLUDE__

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <assert.h>

#include <sys/types.h>
#include <sys/epoll.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#include <list>
#include <string>
#include <vector>

using namespace std;

#define __CHECK_SUCC__(res) assert(res == 0);

#endif
