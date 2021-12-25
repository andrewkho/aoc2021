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

    input = list()
    with open(f'../inputs/day25/{infile}', 'r') as f:
        for i, line in enumerate(f.readlines()):
            input.append(list(line.strip()))

    N, M = len(input), len(input[0])

    def step(input):
        nbr = [l.copy() for l in input]
        # >
        moves = 0
        for i, j in itertools.product(range(N), range(M)):
            nxt = j+1 if j+1 < M else 0
            if input[i][j] == '>' and input[i][nxt] == '.':
                moves += 1
                nbr[i][nxt] = '>'
                nbr[i][j] = '.'
        # v
        nbd = [l.copy() for l in nbr]
        for j, i in itertools.product(range(M), range(N)):
            nxt = i+1 if i+1 < N else 0
            if nbr[i][j] == 'v' and nbr[nxt][j] == '.':
                moves += 1
                nbd[nxt][j] = 'v'
                nbd[i][j] = '.'

        return nbd, moves

    print('\n'.join(''.join(l) for l in input))
    print()
    t0 = time.time()
    steps = 0
    while True:
        input, moves = step(input)
        steps += 1
        if moves == 0:
            break

    print(f'1: {steps}', time.time() - t0, 's')


if __name__ == '__main__':
    fire.Fire(main)
