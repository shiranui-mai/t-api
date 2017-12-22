include ../common.mk

BIN=api
OBJ_PATH=/tmp/api/

SOURCES=$(wildcard *.cc)
GEN_OBJ=$(patsubst %.cc, %.o, $(SOURCES))
GEN_OBJS=$(addprefix $(OBJ_PATH), $(GEN_OBJ))

$(BIN): $(GEN_OBJS)
	$(CC) $(CC_FLAGS) $(CC_LIBS) $< -o $@

$(GEN_OBJS): $(SOURCES) 
	@echo $(GEN_OBJS)
	$(CC) -c $(CC_FLAGS) $(CC_LIBS) $< -o $@
	#mv $@ $(OBJ_PATH)

.PHONY: clean
clean:
	rm -rf $(GEN_OBJS) $(BIN)
