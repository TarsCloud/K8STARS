#!/bin/bash

from_version=":latest"
to_version=":v1.1.0"

if [[ "$OSTYPE" == "darwin"* ]]; then
    alias sed="sed -i ''"
else
    alias sed="sed -i"
fi

for f in `ls yaml/*.yaml`; do
    sed "s/$from_version/$to_version/g" $f
done

