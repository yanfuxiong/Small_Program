
CC = gcc

CXX = g++

DEBUG =
CXX_FLAGS = $(DEBUG) -std=c++17 -Wall -g 

AllDIRS := $(shell ls -R | grep '^\./.*:$$' | awk '{gsub(":","");print}') .

CSRCS = $(foreach n,$(AllDIRS) , $(wildcard $(n)/*.c)) 

CXXSRCS = $(foreach n,$(AllDIRS) , $(wildcard $(n)/*.cpp)) 

OBJS := $(CSRCS:.c=.o) $(CXXSRCS:.cpp=.o)

SRC_PATH = .

MV = mv -f  
RM = rm -rf  
LN = ln -sf  

TARGET = multi_download

INC_PATH = -I$(SRC_PATH) \
			-I$(SRC_PATH)/include \
			-I$(SRC_PATH)/common

LIB_PATH = -lcurl -lpthread

			


%.o : %.cpp

	$(CXX) $(CXX_FLAGS) $< -c -o $@ $(INC_PATH) $(LIB_PATH);

%.o : %.c

	$(CC) $(CXX_FLAGS) $< -c -o  $@ $(INC_PATH) $(LIB_PATH);


$(TARGET) : $(OBJS)

	$(CXX) $(CXX_FLAGS) $^ -o $@ $(INC_PATH) $(LIB_PATH);

	@echo Create $(TARGET) ok , build success


.PHONY:clean


clean:

	@ $(RM) $(OBJS) $(TARGET)
	@echo clean all...



