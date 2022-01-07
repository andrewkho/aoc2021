import math
from dataclasses import dataclass
from collections import deque, defaultdict, Counter
from typing import *
import time

import numpy as np
import itertools
import re

import fire


def main(infile: str):
    print('hi!')

    instr = dict()
    base = ''
    with open(infile, 'r') as f:
        for i, line in enumerate(f.readlines()):
            if i == 0:
                base = line.strip()
            elif i == 1:
                continue
            else:
                l, r = line.strip().split(' -> ')
                instr[l] = r
    print(base)
    print(len(instr))

    c = defaultdict(lambda: 0, Counter(base))
    matches = defaultdict(lambda: 0)
    for i in range(len(base)-1):
        ch = base[i:i+2]
        if ch in instr:
            matches[ch] += 1

    def step(matches, c):
        new_matches = defaultdict(lambda: 0)
        for ch, n in matches.items():
            chl = ch[0] + instr[ch]
            chr = instr[ch] + ch[1]
            if chl in instr:
                new_matches[chl] += n
            if chr in instr:
                new_matches[chr] += n
            c[instr[ch]] += n

        return new_matches, c

    t0 = time.time()
    for i in range(10):
        matches, c = step(matches, c)

    mn = min(c.values())
    mx = max(c.values())
    print('1:', mn, mx, mx-mn, time.time() - t0, 's')
    
    t1 = time.time()
    for i in range(30):
        matches, c = step(matches, c)

    mn = min(c.values())
    mx = max(c.values())
    print('2:', mn, mx, mx-mn, time.time() - t1, 's')
    print('2: sum(values)', sum(c.values())/1024/1024/1024/1024, 'TB')


if __name__ == '__main__':
    fire.Fire(main)
