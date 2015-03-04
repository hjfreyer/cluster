#!/usr/bin/env python

import sys
import os.path
import os
import json
import base64

if len(sys.argv) < 3:
    print('bad')
    sys.exit(1)

name = sys.argv[1]
path = sys.argv[2]

secrets = dict()

for fp in os.listdir(path):
    fullpath = os.path.join(path, fp)
    if not os.path.isfile(fullpath):
        raise Error('Is a directory: %s' % fullpath)

    with open(fullpath) as f:
        secrets[fp] = f.read()

output = dict(
    apiVersion='v1beta1',
    kind='Secret',
    id=name,
    data=dict())

for k, v in secrets.iteritems():
    output['data'][k] = base64.b64encode(v)

json.dump(output, sys.stdout)
