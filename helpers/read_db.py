#!/bin/env python3

import sqlite3
import base64
import os
import sys
from collections import namedtuple

Record = namedtuple('Record', 'id date sourceIP user password clientVersion')

if __name__ == '__main__':
    if len(sys.argv) > 2:
        print(f"Usage: {sys.argv[0]} <SQLITE3 FILE>")
        os.exit(1)
    elif len(sys.argv) == 1:
        db_name = 'mydb.sqlite'
    else:
        db_name = sys.argv[1]
    conn = sqlite3.connect(db_name)
    c = conn.cursor()

    for row in c.execute("select * from sshconnections ORDER BY date"):
        r = Record(*row)
        user = base64.b64decode(r.user.encode())
        password = base64.b64decode(r.password.encode())
        clientVersion = base64.b64decode(r.clientVersion.encode())
        print("{:>5} | {} | {:<20} | {:<30} | {:<40} | {:} ".format(
            r.id,
            r.date,
            r.sourceIP,
            user.decode(),
            password.decode(),
            clientVersion.decode(),
        ))
