#!/bin/bash -e

DIR=$(cd `dirname $0` && pwd -P)

out_file=$(mktemp)

(
	echo "---"
	echo "# DO NOT MODIFY - AUTO GENERATED"
	echo
	drone jsonnet --stream --stdout --source ${DIR}/ci.jsonnet | tail -n +2
) >$out_file

mv $out_file ${DIR}/ci.yml
