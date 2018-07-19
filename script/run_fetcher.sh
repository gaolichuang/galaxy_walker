#!/bin/bash

###### production mode ############
# <tag>_<service>_<replica>_<shard>
# sample: crawler_dispatcher_0_0
# bin: /Application/mustard/crawler_dispatcher_0_0/bin/dispatcher_main
# log: /Application/mustard/crawler_dispatcher_0_0/logs/dispatcher.log
##################################
PRODUCT_PREFIX="/Application/mustard"

#### Alpha env default #####
BIN_PATH=../bin
LOG_PATH=../logs
MDATA=../mdata

REPLICA=0
# shard is not used
SHARD=0

# HTTP_PORT_BASE - 50 == SERVICE_PORT
SERVICE_PORT=9100

get_runtime_path() {
  CPATH=$1
  case $CPATH in
    /*) abspath=$CPATH;;
    *)  abspath=$PWD/$CPATH;;
  esac

  if [ -d "$CPATH" ]; then CPATH=$CPATH/.; fi
  abspath=$(cd "$(dirname -- "$CPATH")"; printf %s. "$PWD")
  abspath=${abspath%?}
  abspath=$abspath/${CPATH##*/}
  echo $abspath
}
dpath=`dirname $0`
WORKPATH=`get_runtime_path $dpath`
if [ "${WORKPATH:0:${#PRODUCT_PREFIX}}" == "${PRODUCT_PREFIX}" ];then
  # get from path
  service=`echo $WORKPATH | awk -F"/" '{print $4}'`
  BIN=`echo $service|awk -F"_" '{print $2}'`
  BIN="${BIN}_main"
  BIN_PATH=$WORKPATH
  REPLICA=`echo $service|awk -F"_" '{print $3}'`
  SHARD=`echo $service|awk -F"_" '{print $4}'`
  LOG_PATH="$WORKPATH/../logs"
  MDATA=$PRODUCT_PREFIX
  echo "Product Mode. BIN:$BIN,Replica:$REPLICA,Shard:$SHARD"
else
  # get from script name
  BIN=`basename $0`
  BIN=${BIN#run_}
  BIN=${BIN%.sh}
  BIN="${BIN}_main"
  echo "Alpha Mode. BIN:$BIN,Replica:$REPLICA,Shard:$SHARD"
fi
BIN="$BIN_PATH/$BIN"
# check BIN exist.
if [ ! -f $BIN ];then
  echo "$BIN not exist."
  exit 3
fi
PORT=$(($SERVICE_PORT + $REPLICA))
HPORT=$(($PORT + 50))
LOG_FILE=$LOG_PATH/`basename $BIN`.log
mkdir -p  $LOG_PATH

    #--crawl_handler_chain=FetchHandler;PrepareHandler;DocHandler;StorageHandler;ResponseHandler
CMD="$BIN
    --crawl_handler_chain=DummyRequestProcessor;FetchHandler;PrepareHandler;DocHandler
    --crawl_input_processor=RequestProcessor
    --proxy_conf_file=etc/crawl/fetch_proxys.config
    --channel_buffer_size=10
    --host_load_queue_size=10
    --fetch_connection_number=20
    --crawl_request_port=$PORT
    --conf_path_prefix=$MDATA
    --http_port=$HPORT
    --v=5
    --stdout=true"
checkOnce() {
  pnum=`ps -ef |grep "$BIN"|grep -c $HPORT`
  [ $pnum -eq 1 ]
  return $?
}
check() {
  for (( c=1; c<=15; c++ ))
  do
    sleep 1
    checkOnce
    if [ $? -eq 0 ];then
      return 0
    fi
  done
  return 1
}
start() {
  ulimit -n unlimited
  # clear log...
  rm $LOG_FILE
  # open gctrace.
  export GODEBUG=gctrace=1
  export GOTRACEBACK=crash
  echo $CMD
  nohup $CMD >> $LOG_FILE 2>&1 &
  check
  if [ $? -eq 0 ];then
      echo "Start Finish."
      return 0
  else
      echo "Start Fail."
      return 2
  fi
}
stop() {
  checkOnce
  if [ $? -eq 0 ];then
    pid=`ps -ef |grep "$BIN"|grep $HPORT|awk '{print $2}'`
    kill -9 $pid
    echo "Stop: `basename $BIN` port: $HPORT PID:$pid"
  else
    echo "Process Not exist. `basename $BIN` port: $HPORT"
  fi
}
status() {
  checkOnce
  if [ $? -eq 0 ];then
    echo "It's OK."
    return 0
  else
    echo "It's Gone."
    return 1
  fi
}

startDirect() {
  echo $CMD
  $CMD
}
Usage="Usage:`basename $0` [start|stop|status|stdout]"
if [ $# -ne 1 ];then
  echo $Usage
  exit 1
fi

if [ $1 == "start" ];then
  start
elif [ $1 == "stop" ];then
  stop
elif [ $1 == "status" ];then
  status
elif [ $1 == "stdout" ];then
  startDirect
else
  echo $Usage
  exit 2
fi
