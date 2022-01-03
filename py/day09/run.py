import math
from dataclasses import dataclass
from collections import deque
from typing import *
import numpy as np
import itertools
import re

import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    input = []
    with open(infile, 'r') as f:
        for line in f.readlines():
            input.append([int(c) for c in line.strip()])

    # print(input)
    s = 0
    for i, line in enumerate(input):
        for j, n in enumerate(line):
            is_bottom = True
            for di, dj in [(-1, 0), (1, 0), (0, -1), (0, 1)]:
                if i+di < 0 or i+di >= len(input):
                    continue
                if j+dj < 0 or j+dj >= len(line):
                    continue
                if n >= input[i+di][j+dj]:
                    is_bottom = False
                    break
            if is_bottom:
                s += n+1

    print(f'1: {s}')

    seen = set()

    def find_region(i, j):
        dq = deque()
        dq.append((i, j))
        region = set()
        while len(dq) > 0:
            i, j = dq.popleft()
            if (i, j) in seen:
                continue

            region.add((i, j))
            seen.add((i, j))
            for di, dj in [(-1, 0), (1, 0), (0, -1), (0, 1)]:
                if (i+di, j+dj) in seen:
                    continue
                if not (0 <= i+di < len(input) and 0 <= j+dj < len(input[0])):
                    continue
                if input[i+di][j+dj] == 9:
                    continue
                dq.append((i+di, j+dj))

        return region

    regions = []
    for i, j in itertools.product(range(len(input)), range(len(input[0]))):
        if input[i][j] != 9 and (i, j) not in seen:
            regions.append(find_region(i, j))

    lengths = sorted(len(r) for r in regions)
    print(lengths[-3:])
    print(f'2: {np.prod(lengths[-3:])}')


if __name__ == '__main__':
    fire.Fire(main)
