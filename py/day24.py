import time
import heapq
import math
from dataclasses import dataclass, field
from collections import deque, defaultdict, Counter
from typing import *

from scipy.signal import convolve2d

import numpy as np
import itertools
import re

import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    instructions = list()
    with open(f'../inputs/day24/{infile}', 'r') as f:
        for i, line in enumerate(f.readlines()):
            instr = line.strip().split(' ')
            instructions.append(instr)

    biggest = 92967699949891
    print(f'1: {biggest}', z(list(map(int, str(biggest)))))

    smallest = 92967699949891
    print(f'2: {smallest}', z(list(map(int, str(smallest)))))


"""
Constraints:
C09: s3 - s4 = 3
C15: S7 - S6 = 3
C17: S8 - S5 = 2
C20: S9 - S10 = 5
C23: S11 - S12 = 1
C25: S13 - S2 = 7
C27: S1 - S14 = 8

01: 9
02: 2/1
03: 9
04: 6
05: 7
06: 6
07: 9
08: 9
09: 9
10: 4
11: 9
12: 8
13: 9/8
14: 1
biggest = 92967699949891

01: 9
02: 1
03: 4
04: 1
05: 1
06: 1
07: 4
08: 3
09: 6
10: 1
11: 2
12: 1
13: 8
14: 1
smallest: 91411143612181
"""


def execute(instructions, input):
    assert all(v != 0 for v in input)
    a = {k: 0 for k in 'wxyz'}
    idx = 0
    for n, i in enumerate(instructions):
        try:
            if i[0] == 'inp':
                a[i[1]] = input[idx]
                idx += 1
                continue
            b = a[i[2]] if i[2] in 'wxyz' else int(i[2])
            if i[0] == 'add':
                a[i[1]] += b
            elif i[0] == 'mul':
                a[i[1]] *= b
            elif i[0] == 'div':
                assert b != 0
                sign = int(np.sign(a[i[1]] / b))
                a[i[1]] = sign * (abs(a[i[1]]) // abs(b))
            elif i[0] == 'mod':
                assert a[i[1]] >= 0
                assert b > 0
                a[i[1]] %= b
            elif i[0] == 'eql':
                a[i[1]] = 1 if a[i[1]] == b else 0
        except Exception as e:
            print(n, i, a)
            raise e

    return a


def div(a, b):
    assert b != 0
    sign = int(np.sign(a / b))
    return sign*(abs(a)//abs(b))


@dataclass
class ZVals:
    x: Optional[int]
    z: int
    s: int
    t: int


ZKEYS = {
    27: ZVals(-11,  25, 14, 1),
    25: ZVals(0,  23, 13, 12),
    23: ZVals(-1,  21, 12, 7),
    21: ZVals(None, 20, 11, 0),
    20: ZVals(-11,  18, 10, 4),
    18: ZVals(None, 17,  9, 6),
    17: ZVals(-12,  15,  8, 9),
    15: ZVals(-4,  12,  7, 9),
    12: ZVals(None, 11,  6, 7),
    11: ZVals(None,  9,  5, 14),
    9:  ZVals(-4,   6,  4, 6),
    6:  ZVals(None,  4,  3, 1),
}


def z(input, zkey=27):
    zk = ZKEYS[zkey]
    if zk.z == 4:
        def zfunc():
            return (input[0] + 3)*26 + input[1] + 7
    else:
        def zfunc():
            return z(input, zk.z)

    if zk.x is None:
        return zfunc()*26 + input[zk.s - 1] + zk.t

    zval = zfunc()
    tmp = div(zval, 26)
    x = zval % 26 + zk.x
    if x == input[zk.s - 1]:
        return tmp
    else:
        return tmp*26 + input[zk.s - 1] + zk.t


if __name__ == '__main__':
    fire.Fire(main)
