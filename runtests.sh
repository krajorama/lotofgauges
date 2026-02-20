#!/bin/bash

set -e

encodings="EncXOR"
# EncXOR2 EncXOR2ST EncXOROptST EncXOROptOtelST"

FEATURE_FLAG="--enable-feature=st-storage"
# FEATURE_FLAG=

for e in ${encodings} ; do
	date
	echo "Build ${e}"
	sed -e 's/const EncodingForFloatST = .*/const EncodingForFloatST = '"${e}"'/g' -i prometheus/tsdb/chunkenc/chunk.go
	pushd prometheus > /dev/null
	make build >& ../buildlog.txt
	popd > /dev/null
	rm -rf data-st
	mkdir data-st
	prometheus/prometheus --storage.tsdb.path=data-st --log.level=error ${FEATURE_FLAG} &
	prompid=$!
	echo "Write only: ${e}"
	sudo /usr/bin/perf stat -p ${prompid} -e task-clock,cpu-clock,faults sleep 1200 2>&1 | grep -E "task-clock|cpu-clock|faults|elapsed"
	/home/krajo/go/github.com/krajorama/lotofgauges/querygauges/querygauges -query 'avg(avg_over_time(example_gauge[10m]))' -max-delay 0ms -min-delay 0ms &
	querypid=$!
	echo "Read/write: ${e}"
	sudo /usr/bin/perf stat -p ${prompid} -e task-clock,cpu-clock,faults sleep 1200 2>&1 | grep -E "task-clock|cpu-clock|faults|elapsed"
	kill ${querypid}
	kill ${prompid}
done

