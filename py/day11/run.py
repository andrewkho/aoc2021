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

    lines = []
    with open(f'../inputs/day11/{infile}', 'r') as f:
        for line in f.readlines():
            lines.append(list(int(x) for x in line.strip()))
    input = np.array(lines)
    print(input)

    def step():
        dq = deque()
        for i in range(input.shape[0]):
            for j in range(input.shape[1]):
                input[i, j] += 1
                if input[i, j] > 9:
                    dq.append((i, j))

        flashed = set()
        while len(dq) > 0:
            node = dq.popleft()
            if node in flashed:
                continue
            flashed.add(node)
            for di, dj in itertools.product([-1, 0, 1], [-1, 0, 1]):
                i, j = node
                if not (0 <= i+di < input.shape[0]):
                    continue
                if not (0 <= j+dj < input.shape[1]):
                    continue
                if (di, dj) == (0, 0):
                    continue

                input[i+di, j+dj] += 1
                if input[i+di, j+dj] > 9:
                    dq.append((i+di, j+dj))

        for i, j in flashed:
            input[i, j] = 0

        return len(flashed)

    p1_total = 0
    for i in range(100):
        p1_total += step()

    print(input)
    print(f'1: {p1_total}')

    s = 100
    while True:
        s += 1
        if step() == input.shape[0]*input.shape[1]:
            break

    print(input)
    print(f'2: {s}')


if __name__ == '__main__':
    fire.Fire(main)
