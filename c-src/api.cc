#include "api.h"
#include "ini.h"
#include "server.h"

#include <fcntl.h>

#include <map>
#include <fstream>

#include <sys/wait.h>

#include <thread>
#include <functional>

// Global conf
// Parser g_conf;

#define MAX_CONNS 999999

void set_nonblock(int fd)
{
    int fl = fcntl(fd, F_GETFL);
    fcntl(fd, F_SETFL, fl | O_NONBLOCK);
}

void add_2_epoll(int ep_fd, int fd, struct epoll_event* ev)
{
    epoll_ctl(ep_fd, EPOLL_CTL_ADD, fd, ev);
}

int load_conf(const char* cp)
{
    fstream fs(cp);
    INI::Parser conf(fs);
    cout << conf.top()("test")["test_key"] << endl;
    cout << conf.top()["test1"] << endl;
    return 0;
}

void work_func(server* s, int n)
{
    cout << "[work_func] id: " << this_thread::get_id() << ". n: " << n << endl;

    work& w = s->works[n];

    while (true) {
        struct epoll_event events[1024]; 
        int num = epoll_wait(w.ep_fd, events, 1024, 1000);
        for (int i = 0; i < num; ++i) {
            printf("[work_func] Receive event\n");
            struct epoll_event _ev = events[i];
            if (_ev.events & EPOLLIN) {
                // Have new data 
                char buf[1024] = "";
                int nread = read(_ev.data.fd, buf, 1023);
                if (nread == 0) {// close {
                    epoll_ctl(s->ep_fd, EPOLL_CTL_DEL, _ev.data.fd, NULL);
                    close(_ev.data.fd);
                    continue;
                }

                printf("[work_func] buf: %s\n", buf);
            } 
        }
    }
}


void init_thread(server* s)
{
    s->works.resize(s->thread_num);

    for (int i = 0; i < s->thread_num; ++i) {
        s->works[i].ep_fd = epoll_create(MAX_CONNS);
        s->works[i].t = thread(work_func, s, i);
        //thread(work_func, s, i).detach(); 
    }
}

server* init_server(const char* cp)
{
    server* s = (server*)malloc(sizeof(server));
    memset(s, 0, sizeof(server));
    assert(s);

    s->conf = cp;
    s->port = stoi(s->get_conf_value("core", "port"));
    s->thread_num = stoi(s->get_conf_value("core", "thread_num"));

    if ((s->fd = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
        assert("socket error");
    }

    // Multi thread
    init_thread(s);

    int opt = 1;
    setsockopt(s->fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    // server add
    struct sockaddr_in s_addr;
    memset(&s_addr, 0, sizeof(s_addr));
    s_addr.sin_family = AF_INET;
    s_addr.sin_addr.s_addr = INADDR_ANY;
    s_addr.sin_port = htons(s->port);

    if (bind(s->fd, (struct sockaddr*)&s_addr, sizeof(struct sockaddr)) < 0) {
        assert("bind failure!");
    }

    listen(s->fd, 1024);

    // new epoll
    if ((s->ep_fd = epoll_create(MAX_CONNS)) < 0) {
        assert("epoll_create error");
    }

    struct epoll_event _ev;
    _ev.events = EPOLLIN|EPOLLET|EPOLLERR|EPOLLHUP;
    // _ev.events = EPOLLIN;
    _ev.data.fd = s->fd;
    // _ev->data.ptr = s;
    add_2_epoll(s->ep_fd, s->fd, &_ev);

    return s;
}

int wait_server(server* s)
{
    while (true) {
        struct epoll_event events[1024]; 
        int num = epoll_wait(s->ep_fd, events, 1024, -1);
        if (num <= 0) continue;
        for (int i = 0; i < num; ++i) {
            printf("[wait_server] Receive event.i: %d \n", i);
            struct epoll_event _ev = events[i];
            printf("[wait_server] _ev.data.fd: %d, s->fd: %d, s->ep_fd: %d, events: %d\n", _ev.data.fd, s->fd, s->ep_fd, _ev.events);
            if (_ev.data.fd == s->fd) {
                // Have new connection
                int cli_fd = accept(s->fd, NULL, NULL);
                struct epoll_event cli_ev;
                cli_ev.data.fd = cli_fd;
                // cli_ev.data.ptr = s;
                cli_ev.events = EPOLLIN|EPOLLET;
                set_nonblock(cli_fd);
                add_2_epoll(s->works[cli_fd%s->thread_num].ep_fd, cli_fd, &cli_ev);
            } 
        }
    }
    return 0;
}

int daemon(const char* cp)
{
    // Load conf
    if (load_conf(cp) != 0) {
        printf("Load conf error\n");
        exit(-1);
    }

    // Loop of child process
    while (true) {
        pid_t pid = fork();
        if (pid == -1) {
            printf("fork error!\n");
            exit(-1);
        }
        if (pid == 0) {
            // Child Process
            printf("I am Child process: %d, My father: %d\n", getpid() ,getppid()); 

            server* s = init_server(cp);
            return wait_server(s);
        } else {
            printf("I am Father process: %d, My Child: %d\n", getpid(), pid); 
        }
        wait(NULL);
        return 0;
    }

    return 0;
}

int main(int argc, const char* argv[])
{
    const char* conf_path = "./api.ini";
    if (argc > 1) conf_path = argv[1]; 

    pid_t pid = fork();
    if (pid == -1) {
        printf("fork error!\n");
        exit(-1);
    }
    // Child Process
    if (pid == 0)
        daemon(conf_path);

    return 0;
}
