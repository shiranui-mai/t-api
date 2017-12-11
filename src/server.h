#include "api.h"

#include <thread>
#include <functional>

#ifndef __SERVER_H__
#define __SERVER_H__

typedef struct work {
    thread t;

    int ep_fd;
} work;

typedef struct server {
    
    int fd;
    int port;
    int ep_fd;

    const char* conf;

    int thread_num;
    
    vector<work>  works;

    const string& get_conf_value(const char* sec, const char* key)
    {
        fstream fs(conf);
        return INI::Parser(fs).top()(sec)[key];
    }
} server;

#endif
