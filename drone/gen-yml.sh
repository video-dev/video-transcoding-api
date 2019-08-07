#!/bin/bash -e

DIR=$(cd `dirname $0` && pwd -P)

out_file=$(mktemp)

(
	echo "# DO NOT MODIFY - AUTO GENERATED"
	echo
	drone jsonnet --stream --stdout --source ${DIR}/ci.jsonnet
) >$out_file

mv $out_file ${DIR}/ci.yml
