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

    img = []
    with open(f'../inputs/day20/{infile}', 'r') as f:
        for i, line in enumerate(f.readlines()):
            if i == 0:
                code = [1 if x == '#' else 0 for x in line.strip()]
            elif i == 1:
                pass
            else:
                img.append([1 if x == '#' else 0 for x in line.strip()])

    t0 = time.time()
    code = np.array(code, dtype=int)
    img = np.array(img, dtype=int)
    print(len(code), code)
    print(img.shape)
    for i in range(2):
        if code[0] == 1:
            fill = i%2
        else:
            fill = 0
        img = conv(img, code, fill)
        display(img)
        print()

    display(img)

    print('1:', img.sum(), time.time() - t0, 's')

    t1 = time.time()
    for i in range(48):
        if code[0] == 1:
            fill = i%2
        else:
            fill = 0
        img = conv(img, code, fill)

    print('2:', img.sum(), time.time() - t1, 's')


KERNEL = np.array([
    [1, 2, 4],
    [8, 16, 32],
    [64, 128, 256],
])


def conv(img, code, fill):
    return code[convolve2d(img, KERNEL, fillvalue=fill)]


def display(img):
    N, M = img.shape
    for i in range(N):
        line = []
        for j in range(M):
            line.append('#' if img[i, j] == 1 else '.')
        print(''.join(line))


if __name__ == '__main__':
    fire.Fire(main)
